package wpml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	nbioxml "github.com/nbio/xml"
)

type XMLSerializer struct {
	indent bool
}

func NewXMLSerializer(indent bool) *XMLSerializer {
	return &XMLSerializer{
		indent: indent,
	}
}

func (s *XMLSerializer) Marshal(doc *Document) ([]byte, error) {

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	var data []byte
	var err error

	if s.indent {
		data, err = xml.MarshalIndent(doc, "", "  ")
	} else {
		data, err = xml.Marshal(doc)
	}

	if err != nil {
		return nil, fmt.Errorf(ErrMarshalDocument, err)
	}

	buf.Write(data)
	return buf.Bytes(), nil
}

func (s *XMLSerializer) Unmarshal(data []byte, doc *Document) error {
	if err := nbioxml.Unmarshal([]byte(data), doc); err != nil {
		return fmt.Errorf(ErrUnmarshalDocument, err)
	}
	return nil
}

func (s *XMLSerializer) MarshalToWriter(doc *Document, writer io.Writer) error {
	data, err := s.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (s *XMLSerializer) UnmarshalFromReader(reader io.Reader, doc *Document) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf(ErrReadData, err)
	}

	return s.Unmarshal(data, doc)
}

func FormatXML(xmlString string) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlString)))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf(ErrParseXML, err)
		}

		if err := encoder.EncodeToken(token); err != nil {
			return "", fmt.Errorf(ErrEncodeToken, err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", fmt.Errorf(ErrFlushEncoder, err)
	}

	return buf.String(), nil
}

func ValidateXML(data []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		_, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf(ErrInvalidXML, err)
		}
	}

	return nil
}

func GetXMLElements(data []byte, elementName string) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var elements []string
	var depth int
	var capturing bool
	var element bytes.Buffer

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf(ErrParseXML, err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == elementName || t.Name.Space+":"+t.Name.Local == elementName {
				capturing = true
				depth = 1
				element.Reset()
				element.WriteString("<")
				element.WriteString(t.Name.Local)
				for _, attr := range t.Attr {
					element.WriteString(" ")
					element.WriteString(attr.Name.Local)
					element.WriteString(`="`)
					element.WriteString(attr.Value)
					element.WriteString(`"`)
				}
				element.WriteString(">")
			} else if capturing {
				depth++
				element.WriteString("<")
				element.WriteString(t.Name.Local)
				for _, attr := range t.Attr {
					element.WriteString(" ")
					element.WriteString(attr.Name.Local)
					element.WriteString(`="`)
					element.WriteString(attr.Value)
					element.WriteString(`"`)
				}
				element.WriteString(">")
			}
		case xml.EndElement:
			if capturing {
				element.WriteString("</")
				element.WriteString(t.Name.Local)
				element.WriteString(">")
				depth--
				if depth == 0 {
					elements = append(elements, element.String())
					capturing = false
				}
			}
		case xml.CharData:
			if capturing {
				element.Write(t)
			}
		}
	}

	return elements, nil
}

func MarshalTemplate(template *TemplateDocument) ([]byte, error) {

	if template.XMLNS == "" {
		template.XMLNS = "http://www.opengis.net/kml/2.2"
	}
	if template.WPMLNS == "" {
		template.WPMLNS = "http://www.dji.com/wpmz/1.0.6"
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	data, err := xml.MarshalIndent(template, "", "  ")
	if err != nil {
		return nil, fmt.Errorf(ErrMarshalTemplateDocument, err)
	}

	buf.Write(data)
	return buf.Bytes(), nil
}

func MarshalWaylines(waylines *WaylinesDocument) ([]byte, error) {

	if waylines.XMLNS == "" {
		waylines.XMLNS = "http://www.opengis.net/kml/2.2"
	}
	if waylines.WPMLNS == "" {
		waylines.WPMLNS = "http://www.dji.com/wpmz/1.0.6"
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	data, err := xml.MarshalIndent(waylines, "", "  ")
	if err != nil {
		return nil, fmt.Errorf(ErrMarshalWaylinesDocument, err)
	}

	buf.Write(data)
	return buf.Bytes(), nil
}

func UnmarshalTemplate(data []byte) (*TemplateDocument, error) {
	var template TemplateDocument
	if err := nbioxml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf(ErrUnmarshalTemplateDocument, err)
	}
	return &template, nil
}

func UnmarshalWaylines(data []byte) (*WaylinesDocument, error) {
	var waylines WaylinesDocument
	if err := nbioxml.Unmarshal(data, &waylines); err != nil {
		return nil, fmt.Errorf(ErrUnmarshalWaylinesDocument, err)
	}
	return &waylines, nil
}
