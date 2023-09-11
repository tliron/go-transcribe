package transcribe

import (
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

// When Pretty is true, takes into account [terminal.Colorize] and uses
// [terminal.Indent] or [terminal.IndentSpaces], overriding our Indent.
//
// If value is a [string] will print it as is, ignoring the format argument.
//
// If value is a [*etree.Document] will use [Transcriber.PrintXMLDocument],
// ignoring the format argument.
func (self *Transcriber) Print(value any, writer io.Writer, format string) error {
	if string_, ok := value.(string); ok {
		return self.PrintString(string_, writer)
	}

	if xmlDocument, ok := value.(*etree.Document); ok {
		return self.PrintXMLDocument(xmlDocument, writer)
	}

	switch format {
	case "yaml", "":
		return self.PrintYAML(value, writer)

	case "json":
		return self.PrintJSON(value, writer)

	case "xjson":
		return self.PrintXJSON(value, writer)

	case "xml":
		return self.PrintXML(value, writer)

	case "cbor":
		return self.PrintCBOR(value, writer)

	case "messagepack":
		return self.PrintMessagePack(value, writer)

	case "go":
		return self.PrintGo(value, writer)

	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func (self *Transcriber) PrintString(value any, writer io.Writer) error {
	string_ := util.ToString(value)
	if _, err := io.WriteString(writer, string_); err == nil {
		if self.Pretty && !strings.HasSuffix(string_, "\n") {
			return util.WriteNewline(writer)
		}
		return nil
	} else {
		return err
	}
}

func (self *Transcriber) PrintYAML(value any, writer io.Writer) error {
	if self.Pretty {
		self = self.WithIndent(terminal.Indent)
		if terminal.Colorize {
			if code, err := self.StringifyYAML(value); err == nil {
				return ColorizeYAML(code, writer)
			} else {
				return err
			}
		}
	}

	return self.WriteYAML(value, writer)
}

func (self *Transcriber) PrintJSON(value any, writer io.Writer) error {
	if self.Pretty {
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
			return self.WithIndent(terminal.Indent).WriteJSON(value, writer)
		}
	} else {
		return self.WriteJSON(value, writer)
	}
}

func (self *Transcriber) PrintXJSON(value any, writer io.Writer) error {
	if value_, err := ard.PrepareForEncodingXJSON(value, self.Reflector); err == nil {
		return self.PrintJSON(value_, writer)
	} else {
		return err
	}
}

func (self *Transcriber) PrintXML(value any, writer io.Writer) error {
	if self.Pretty {
		self = self.WithIndent(terminal.Indent)
	}

	if err := self.WriteXML(value, writer); err == nil {
		if self.Pretty {
			return util.WriteNewline(writer)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func (self *Transcriber) PrintXMLDocument(xmlDocument *etree.Document, writer io.Writer) error {
	if self.Pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(len(self.Indent))
	}

	_, err := xmlDocument.WriteTo(writer)
	return err
}

func (self *Transcriber) PrintCBOR(value any, writer io.Writer) error {
	return self.WriteCBOR(value, writer)
}

func (self *Transcriber) PrintMessagePack(value any, writer io.Writer) error {
	return self.WriteMessagePack(value, writer)
}

func (self *Transcriber) PrintGo(value any, writer io.Writer) error {
	if self.Pretty {
		self = self.WithIndent(terminal.Indent)
	}

	return self.WriteGo(value, writer)
}
