package wpml

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateKmz creates a KMZ file at the specified path from a Mission.
func CreateKmz(mission *Mission, kmzPath string) error {
	buffer, err := CreateKmzBuffer(mission)
	if err != nil {
		return fmt.Errorf(ErrGenerateKMZBuffer, err)
	}

	dir := filepath.Dir(kmzPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf(ErrCreateDirectory, err)
	}

	if err := os.WriteFile(kmzPath, buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf(ErrWriteKMZFile, err)
	}

	return nil
}

// CreateKmzBuffer creates a KMZ file as an in-memory buffer from a Mission.
func CreateKmzBuffer(mission *Mission) (*bytes.Buffer, error) {
	if mission == nil {
		return nil, ErrMissionCannotBeEmpty
	}
	if mission.Template == nil {
		return nil, ErrMissionTemplateCannotBeEmpty
	}
	if mission.Waylines == nil {
		return nil, ErrMissionWaylinesCannotBeEmpty
	}
	templateData, err := MarshalTemplate(mission.Template)
	if err != nil {
		return nil, fmt.Errorf(ErrSerializeTemplate, err)
	}
	waylinesData, err := MarshalWaylines(mission.Waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrSerializeWaylines, err)
	}
	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)
	templateWriter, err := zipWriter.Create("wpmz/template.kml")
	if err != nil {
		_ = zipWriter.Close()
		return nil, fmt.Errorf(ErrCreateTemplateEntry, err)
	}
	if _, writeErr := templateWriter.Write(templateData); writeErr != nil {
		_ = zipWriter.Close()
		return nil, fmt.Errorf(ErrWriteTemplate, writeErr)
	}

	waylinesWriter, err := zipWriter.Create("wpmz/waylines.wpml")
	if err != nil {
		_ = zipWriter.Close()
		return nil, fmt.Errorf(ErrCreateWaylinesEntry, err)
	}
	if _, err := waylinesWriter.Write(waylinesData); err != nil {
		_ = zipWriter.Close()
		return nil, fmt.Errorf(ErrWriteWaylines, err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf(ErrCloseZIPWriter, err)
	}

	return buffer, nil
}

// CreateKmzBufferFromWaylines creates a KMZ buffer by converting a Waylines schema to a Mission first.
func CreateKmzBufferFromWaylines(waylines *Waylines) (*bytes.Buffer, error) {
	mission, err := ConvertWaylinesToMission(waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertWaylines, err)
	}

	return CreateKmzBuffer(mission)
}

// KmzFileInfo holds metadata for a single file inside a KMZ archive.
type KmzFileInfo struct {
	Name             string `json:"name"`
	CompressedSize   uint64 `json:"compressed_size"`
	UncompressedSize uint64 `json:"uncompressed_size"`
	Method           uint16 `json:"method"`
}

// KmzInfo holds metadata about a KMZ archive.
type KmzInfo struct {
	TotalSize int           `json:"total_size"`
	Files     []KmzFileInfo `json:"files"`
}

// GetKmzInfo returns metadata about the KMZ file generated from a Mission, including file sizes.
func GetKmzInfo(mission *Mission) (*KmzInfo, error) {
	buffer, err := CreateKmzBuffer(mission)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if err != nil {
		return nil, fmt.Errorf(ErrParseZIP, err)
	}

	info := &KmzInfo{
		TotalSize: buffer.Len(),
		Files:     make([]KmzFileInfo, 0, len(zipReader.File)),
	}

	for _, file := range zipReader.File {
		info.Files = append(info.Files, KmzFileInfo{
			Name:             file.Name,
			CompressedSize:   file.CompressedSize64,
			UncompressedSize: file.UncompressedSize64,
			Method:           file.Method,
		})
	}

	return info, nil
}

// ParseKMZBuffer parses a KMZ file from a byte buffer and returns a Mission.
func ParseKMZBuffer(buffer []byte) (*Mission, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	if err != nil {
		return nil, fmt.Errorf(ErrParseZIPFile, err)
	}

	var templateData, waylinesData []byte
	resources := make(map[string][]byte)

	for _, file := range zipReader.File {
		switch {
		case strings.HasSuffix(file.Name, "template.kml"):
			templateData, err = readZipFile(file)
			if err != nil {
				return nil, fmt.Errorf(ErrReadTemplate, err)
			}
		case strings.HasSuffix(file.Name, "waylines.wpml"):
			waylinesData, err = readZipFile(file)
			if err != nil {
				return nil, fmt.Errorf(ErrReadWaylinesWPML, err)
			}
		case strings.HasPrefix(file.Name, "res/"):
			resData, readErr := readZipFile(file)
			if readErr == nil {
				resources[file.Name] = resData
			}
		}
	}

	if templateData == nil || waylinesData == nil {
		return nil, ErrKMZFormatIncorrect
	}

	template, err := UnmarshalTemplate(templateData)
	if err != nil {
		return nil, fmt.Errorf(ErrParseTemplateKML, err)
	}

	waylines, err := UnmarshalWaylines(waylinesData)
	if err != nil {
		return nil, fmt.Errorf(ErrParseWaylinesWPML, err)
	}

	return &Mission{
		Template:  template,
		Waylines:  waylines,
		Resources: resources,
	}, nil
}

// ParseKMZFile reads and parses a KMZ file from the filesystem and returns a Mission.
func ParseKMZFile(filePath string) (*Mission, error) {
	buffer, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("读取KMZfilefailure: %w", err)
	}
	return ParseKMZBuffer(buffer)
}

// GenerateKMZJSON generates a JSON representation of the Mission with the given filename and creation timestamp.
func GenerateKMZJSON(mission *Mission, fileName string) (string, error) {
	if mission == nil {
		return "", fmt.Errorf("mission不能为空")
	}

	result := map[string]any{
		"file_name":  fileName,
		"created_at": time.Now().Format(time.RFC3339),
		"template":   nil,
		"waylines":   nil,
	}

	if mission.Template != nil {
		result["template"] = mission.Template
	}

	if mission.Waylines != nil {
		result["waylines"] = mission.Waylines
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serializationJSONfailure: %w", err)
	}

	return string(jsonData), nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()
	return io.ReadAll(reader)
}
