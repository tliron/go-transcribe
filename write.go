package transcribe

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
	"gopkg.in/yaml.v3"
)

func Write(value any, format string, indent string, strict bool, writer io.Writer, base64 bool, reflector *ard.Reflector) error {
	// Special handling for bare strings (format is ignored)
	if s, ok := value.(string); ok {
		_, err := io.WriteString(writer, s)
		return err
	}

	// Special handling for XML etree document (format is ignored)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return WriteXMLDocument(xmlDocument, writer, indent)
	}

	switch format {
	case "yaml", "":
		return WriteYAML(value, writer, indent, strict, reflector)

	case "json":
		return WriteJSON(value, writer, indent)

	case "xjson":
		return WriteXJSON(value, writer, indent, reflector)

	case "xml":
		return WriteXML(value, writer, indent, reflector)

	case "cbor":
		return WriteCBOR(value, writer, base64)

	case "messagepack":
		return WriteMessagePack(value, writer, base64)

	case "go":
		return WriteGo(value, writer, indent)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func WriteYAML(value any, writer io.Writer, indent string, strict bool, reflector *ard.Reflector) error {
	if strict {
		var err error
		if value, err = ard.ToYAMLDocumentNode(value, true, reflector); err != nil {
			return err
		}
	}

	encoder := yaml.NewEncoder(writer)

	encoder.SetIndent(len(indent)) // This might not work as expected for tabs!
	// BUG: currently does not allow an indent value of 1, see: https://github.com/go-yaml/yaml/issues/501

	if slice, ok := value.([]any); !ok {
		return encoder.Encode(value)
	} else {
		// YAML separates each entry with "---"
		// (In JSON the slice would be written as an array)
		for _, data_ := range slice {
			if err := encoder.Encode(data_); err != nil {
				return err
			}
		}
		return nil
	}
}

func WriteJSON(value any, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(value)
}

func WriteXJSON(value any, writer io.Writer, indent string, reflector *ard.Reflector) error {
	if value_, err := ard.PrepareForEncodingXJSON(value, reflector); err == nil {
		return WriteJSON(value_, writer, indent)
	} else {
		return err
	}
}

func WriteXML(value any, writer io.Writer, indent string, reflector *ard.Reflector) error {
	value, err := ard.PrepareForEncodingXML(value, reflector)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", indent)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	if indent == "" {
		// When there's no indent the XML encoder does not emit a final newline
		// (We want it for consistency with YAML and JSON)
		if _, err := io.WriteString(writer, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func WriteXMLDocument(xmlDocument *etree.Document, writer io.Writer, indent string) error {
	xmlDocument.Indent(len(indent))
	_, err := xmlDocument.WriteTo(writer)
	return err
}

func WriteCBOR(value any, writer io.Writer, base64 bool) error {
	if base64 {
		if value_, err := StringifyCBOR(value); err == nil {
			_, err = io.WriteString(writer, value_)
			return err
		} else {
			return err
		}
	} else {
		encoder := cbor.NewEncoder(writer)
		return encoder.Encode(value)
	}
}

func WriteMessagePack(value any, writer io.Writer, base64 bool) error {
	// MessagePack encoder has problems with map[any]any
	if value_, err := ard.ValidCopyMapsToStringMaps(value, nil); err == nil {
		if base64 {
			if value__, err := StringifyMessagePack(value_); err == nil {
				_, err = io.WriteString(writer, value__)
				return err
			} else {
				return err
			}
		} else {
			encoder := util.NewMessagePackEncoder(writer)
			return encoder.Encode(value_)
		}
	} else {
		return err
	}
}

func WriteGo(value any, writer io.Writer, indent string) error {
	NewUtterConfig(indent).Fdump(writer, value)
	return nil
}
