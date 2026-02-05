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

func CreateKmz(mission *WPMLMission, kmzPath string) error {
	buffer, err := CreateKmzBuffer(mission)
	if err != nil {
		return fmt.Errorf(ErrGenerateKMZBuffer, err)
	}

	dir := filepath.Dir(kmzPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf(ErrCreateDirectory, err)
	}

	if err := os.WriteFile(kmzPath, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf(ErrWriteKMZFile, err)
	}

	return nil
}

func CreateKmzBuffer(mission *WPMLMission) (*bytes.Buffer, error) {
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
		zipWriter.Close()
		return nil, fmt.Errorf(ErrCreateTemplateEntry, err)
	}
	if _, writeErr := templateWriter.Write(templateData); writeErr != nil {
		zipWriter.Close()
		return nil, fmt.Errorf(ErrWriteTemplate, writeErr)
	}

	waylinesWriter, err := zipWriter.Create("wpmz/waylines.wpml")
	if err != nil {
		zipWriter.Close()
		return nil, fmt.Errorf(ErrCreateWaylinesEntry, err)
	}
	if _, err := waylinesWriter.Write(waylinesData); err != nil {
		zipWriter.Close()
		return nil, fmt.Errorf(ErrWriteWaylines, err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf(ErrCloseZIPWriter, err)
	}

	return buffer, nil
}

func CreateKmzBufferFromWaylines(waylines *Waylines) (*bytes.Buffer, error) {
	mission, err := ConvertWaylinesToWPMLMission(waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertWaylines, err)
	}

	return CreateKmzBuffer(mission)
}

func GetKmzInfo(mission *WPMLMission) (map[string]interface{}, error) {
	buffer, err := CreateKmzBuffer(mission)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if err != nil {
		return nil, fmt.Errorf(ErrParseZIP, err)
	}

	info := map[string]interface{}{
		"total_size": buffer.Len(),
		"files":      make([]map[string]interface{}, 0),
	}

	for _, file := range zipReader.File {
		fileInfo := map[string]interface{}{
			"name":              file.Name,
			"compressed_size":   file.CompressedSize64,
			"uncompressed_size": file.UncompressedSize64,
			"method":            file.Method,
		}
		info["files"] = append(info["files"].([]map[string]interface{}), fileInfo)
	}

	return info, nil
}

func ParseKMZBuffer(buffer []byte) (*WPMLMission, error) {
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

	return &WPMLMission{
		Template:  template,
		Waylines:  waylines,
		Resources: resources,
	}, nil
}

func ParseKMZFile(filePath string) (*WPMLMission, error) {
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取KMZfilefailure: %w", err)
	}
	return ParseKMZBuffer(buffer)
}

func GenerateKMZJSON(mission *WPMLMission, fileName string) (string, error) {
	if mission == nil {
		return "", fmt.Errorf("mission不能为空")
	}

	result := map[string]interface{}{
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
	defer reader.Close()
	return io.ReadAll(reader)
}
