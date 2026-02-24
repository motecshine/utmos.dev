package wpml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	nbioxml "github.com/nbio/xml"
)

// XMLSerializer provides XML marshaling and unmarshaling for WPML documents.
type XMLSerializer struct {
	indent bool
}

// NewXMLSerializer creates a new XMLSerializer with the given indentation setting.
func NewXMLSerializer(indent bool) *XMLSerializer {
	return &XMLSerializer{
		indent: indent,
	}
}

// Marshal marshals a Document to XML bytes, prepending the XML declaration header.
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

// Unmarshal unmarshals XML bytes into a Document using namespace-aware parsing.
func (s *XMLSerializer) Unmarshal(data []byte, doc *Document) error {
	if err := nbioxml.Unmarshal(data, doc); err != nil {
		return fmt.Errorf(ErrUnmarshalDocument, err)
	}
	return nil
}

// MarshalToWriter marshals a Document and writes the XML bytes to the given io.Writer.
func (s *XMLSerializer) MarshalToWriter(doc *Document, writer io.Writer) error {
	data, err := s.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

// UnmarshalFromReader reads all data from the given io.Reader and unmarshals it into a Document.
func (s *XMLSerializer) UnmarshalFromReader(reader io.Reader, doc *Document) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf(ErrReadData, err)
	}

	return s.Unmarshal(data, doc)
}

// FormatXML formats an XML string with proper indentation and an XML declaration header.
func FormatXML(xmlString string) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlString)))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
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

// ValidateXML validates that the given bytes contain well-formed XML.
func ValidateXML(data []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		_, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf(ErrInvalidXML, err)
		}
	}

	return nil
}

func writeStartElement(buf *bytes.Buffer, name string, attrs []xml.Attr) {
	buf.WriteString("<")
	buf.WriteString(name)
	for _, attr := range attrs {
		buf.WriteString(" ")
		buf.WriteString(attr.Name.Local)
		buf.WriteString(`="`)
		buf.WriteString(attr.Value)
		buf.WriteString(`"`)
	}
	buf.WriteString(">")
}

// GetXMLElements extracts all XML elements matching the given element name from the XML data.
func GetXMLElements(data []byte, elementName string) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var elements []string
	var depth int
	var capturing bool
	var element bytes.Buffer

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
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
				writeStartElement(&element, t.Name.Local, t.Attr)
			} else if capturing {
				depth++
				writeStartElement(&element, t.Name.Local, t.Attr)
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

// xmlNamespacer is implemented by document types that have KML/WPML namespace attributes.
type xmlNamespacer interface {
	setDefaultNamespaces()
}

func marshalDocumentWithDefaults[T xmlNamespacer](doc T, errMsg string) ([]byte, error) {
	doc.setDefaultNamespaces()

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	data, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	buf.Write(data)
	return buf.Bytes(), nil
}

// MarshalTemplate marshals a TemplateDocument to XML bytes with default namespaces and XML declaration.
func MarshalTemplate(template *TemplateDocument) ([]byte, error) {
	return marshalDocumentWithDefaults(template, ErrMarshalTemplateDocument)
}

// MarshalWaylines marshals a WaylinesDocument to XML bytes with default namespaces and XML declaration.
func MarshalWaylines(waylines *WaylinesDocument) ([]byte, error) {
	return marshalDocumentWithDefaults(waylines, ErrMarshalWaylinesDocument)
}

func unmarshalDocument[T any](data []byte, errMsg string) (*T, error) {
	var doc T
	if err := nbioxml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &doc, nil
}

// UnmarshalTemplate unmarshals XML bytes into a TemplateDocument.
func UnmarshalTemplate(data []byte) (*TemplateDocument, error) {
	return unmarshalDocument[TemplateDocument](data, ErrUnmarshalTemplateDocument)
}

// UnmarshalWaylines unmarshals XML bytes into a WaylinesDocument.
func UnmarshalWaylines(data []byte) (*WaylinesDocument, error) {
	return unmarshalDocument[WaylinesDocument](data, ErrUnmarshalWaylinesDocument)
}
