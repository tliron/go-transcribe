package transcribe

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func Print(value any, format string, writer io.Writer, strict bool, pretty bool, base64 bool, reflector *ard.Reflector) error {
	// Special handling for strings (ignore format)
	if s, ok := value.(string); ok {
		if _, err := io.WriteString(writer, s); err == nil {
			if pretty {
				return util.WriteNewline(writer)
			} else {
				return nil
			}
		} else {
			return err
		}
	}

	// Special handling for etree (ignore format)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return PrintXMLDocument(xmlDocument, writer, pretty)
	}

	switch format {
	case "yaml", "":
		return PrintYAML(value, writer, strict, pretty, reflector)

	case "json":
		return PrintJSON(value, writer, pretty)

	case "xjson":
		return PrintXJSON(value, writer, pretty, reflector)

	case "xml":
		return PrintXML(value, writer, pretty, reflector)

	case "cbor":
		return PrintCBOR(value, writer, base64)

	case "messagepack":
		return PrintMessagePack(value, writer, base64)

	case "go":
		return PrintGo(value, writer, pretty)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func PrintYAML(value any, writer io.Writer, strict bool, pretty bool, reflector *ard.Reflector) error {
	if pretty && terminal.Colorize {
		if code, err := StringifyYAML(value, terminal.Indent, strict, reflector); err == nil {
			return ColorizeYAML(code, writer)
		} else {
			return err
		}
	} else {
		return WriteYAML(value, writer, "  ", strict, reflector)
	}
}

func PrintJSON(value any, writer io.Writer, pretty bool) error {
	if pretty {
		if terminal.Colorize {
			formatter := NewJSONColorFormatter(terminal.IndentSpaces)
			if bytes, err := formatter.Marshal(value); err == nil {
				if _, err := writer.Write(bytes); err == nil {
					return util.WriteNewline(writer)
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			return WriteJSON(value, writer, terminal.Indent)
		}
	} else {
		return WriteJSON(value, writer, "")
	}
}

func PrintXJSON(value any, writer io.Writer, pretty bool, reflector *ard.Reflector) error {
	if value_, err := ard.PrepareForEncodingXJSON(value, reflector); err == nil {
		return PrintJSON(value_, writer, pretty)
	} else {
		return err
	}
}

func PrintXML(value any, writer io.Writer, pretty bool, reflector *ard.Reflector) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteXML(value, writer, indent, reflector); err == nil {
		if pretty {
			return util.WriteNewline(writer)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func PrintXMLDocument(xmlDocument *etree.Document, writer io.Writer, pretty bool) error {
	if pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(writer)
	return err
}

func PrintCBOR(value any, writer io.Writer, base64 bool) error {
	return WriteCBOR(value, writer, base64)
}

func PrintMessagePack(value any, writer io.Writer, base64 bool) error {
	return WriteMessagePack(value, writer, base64)
}

func PrintGo(value any, writer io.Writer, pretty bool) error {
	if pretty {
		return WriteGo(value, writer, terminal.Indent)
	} else {
		return WriteGo(value, writer, "")
	}
}
