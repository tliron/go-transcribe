package transcribe

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
	"gopkg.in/yaml.v3"
)

func (self *Transcriber) Write(value any, writer io.Writer, format string) error {
	// Special handling for bare strings (format is ignored)
	if s, ok := value.(string); ok {
		_, err := io.WriteString(writer, s)
		return err
	}

	// Special handling for XML etree document (format is ignored)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return self.WriteXMLDocument(xmlDocument, writer)
	}

	switch format {
	case "yaml", "":
		return self.WriteYAML(value, writer)

	case "json":
		return self.WriteJSON(value, writer)

	case "xjson":
		return self.WriteXJSON(value, writer)

	case "xml":
		return self.WriteXML(value, writer)

	case "cbor":
		return self.WriteCBOR(value, writer)

	case "messagepack":
		return self.WriteMessagePack(value, writer)

	case "go":
		return self.WriteGo(value, writer)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

// Opens the output file for write and calls [Transcriber.Write] on it.
func (self *Transcriber) WriteToFile(value any, output string, format string) error {
	if f, err := OpenFileForWrite(output); err == nil {
		defer f.Close()
		return self.Write(value, f, format)
	} else {
		return err
	}
}

func (self *Transcriber) WriteYAML(value any, writer io.Writer) error {
	if self.Strict {
		var err error
		if value, err = ard.ToYAMLDocumentNode(value, true, self.Reflector); err != nil {
			return err
		}
	}

	encoder := yaml.NewEncoder(writer)

	encoder.SetIndent(len(self.Indent)) // This might not work as expected for tabs!
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

func (self *Transcriber) WriteJSON(value any, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", self.Indent)
	return encoder.Encode(value)
}

func (self *Transcriber) WriteXJSON(value any, writer io.Writer) error {
	if value_, err := ard.PrepareForEncodingXJSON(value, self.Reflector); err == nil {
		return self.WriteJSON(value_, writer)
	} else {
		return err
	}
}

func (self *Transcriber) WriteXML(value any, writer io.Writer) error {
	value, err := ard.PrepareForEncodingXML(value, self.Reflector)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", self.Indent)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	if self.Indent == "" {
		// When there's no indent the XML encoder does not emit a final newline
		// (We want it for consistency with YAML and JSON)
		if _, err := io.WriteString(writer, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (self *Transcriber) WriteXMLDocument(xmlDocument *etree.Document, writer io.Writer) error {
	xmlDocument.Indent(len(self.Indent))
	_, err := xmlDocument.WriteTo(writer)
	return err
}

func (self *Transcriber) WriteCBOR(value any, writer io.Writer) error {
	if self.Base64 {
		if value_, err := self.StringifyCBOR(value); err == nil {
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

func (self *Transcriber) WriteMessagePack(value any, writer io.Writer) error {
	// MessagePack encoder has problems with map[any]any
	if value_, err := ard.ValidCopyMapsToStringMaps(value, nil); err == nil {
		if self.Base64 {
			if value__, err := self.StringifyMessagePack(value_); err == nil {
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

func (self *Transcriber) WriteGo(value any, writer io.Writer) error {
	self.NewUtterConfig().Fdump(writer, value)
	return nil
}

func OpenFileForWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), DIRECTORY_WRITE_PERMISSIONS); err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FILE_WRITE_PERMISSIONS)
	} else {
		return nil, err
	}
}
