package transcribe

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
)

//
// Transcriber
//

type Transcriber struct {
	File        string    // if not empty will supersede Writer
	Writer      io.Writer // if nil then os.Stdout will be used
	Format      string    // "yaml", "json", "xjson", "xml", "cbor", "messagepack", or "go"
	ForTerminal bool

	Indent    string         // used by YAML, JSON, XML, and Go
	Strict    bool           // used by YAML
	Base64    bool           // used by CBOR and MessagePack
	InPlace   bool           // used by XJSON and XML
	Reflector *ard.Reflector // used by XJSON and XML
}

func NewTranscriber() *Transcriber {
	return new(Transcriber)
}

func (self *Transcriber) Clone() *Transcriber {
	return &Transcriber{
		File:        self.File,
		Writer:      self.Writer,
		Format:      self.Format,
		ForTerminal: self.ForTerminal,
		Indent:      self.Indent,
		Strict:      self.Strict,
		Base64:      self.Base64,
		InPlace:     self.InPlace,
		Reflector:   self.Reflector,
	}
}

// Writes the value to a writer according to [Transcriber.Format]. If the
// format is any empty string will default to YAML.
//
// If [Transcriber.File] is not empty will create the file and write to it.
// Sub-directories in the path will be created if they don't exist. Otherwise
// will write to [Transcriber.Writer], and if that is nil will write to stdout.
//
// When [Transcriber.ForTerminal] is true will optimize for terminals, such
// as ensuring a newline at the end and colorization if supported by the format
// and the terminal. In this mode [Transcriber.Indent] is ignored and
// [terminal.Indent] will be used instead.
//
// If value is a string then will ignore the format and use
// [Transcriber.WriteString]. If value is a [*etree.Document] will ignore the
// format and use [Transcriber.WriteXMLDocument].
//
// Note special handling for YAML when value is []any. In this case, instead
// of writing a YAML seq, it will treat the value as a list of documents
// separated by "---".
func (self *Transcriber) Write(value any) error {
	if self.File != "" {
		if file, err := OpenFileForWrite(self.File); err == nil {
			defer file.Close()

			self = self.Clone()
			self.Writer = file
		} else {
			return err
		}
	}

	// Special handling for bare strings (format is ignored)
	if string_, ok := value.(string); ok {
		return self.WriteString(string_)
	}

	// Special handling for XML etree document (format is ignored)
	if xmlDocument, ok := value.(*etree.Document); ok {
		return self.WriteXMLDocument(xmlDocument)
	}

	switch self.Format {
	case "yaml", "":
		return self.WriteYAML(value)

	case "json":
		return self.WriteJSON(value)

	case "xjson":
		return self.WriteXJSON(value)

	case "xml":
		return self.WriteXML(value)

	case "cbor":
		return self.WriteCBOR(value)

	case "messagepack":
		return self.WriteMessagePack(value)

	case "go":
		return self.WriteGo(value)

	default:
		return fmt.Errorf("unsupported format: %q", self.Format)
	}
}

// Converts the value to a string according to [Transcriber.Format]. If the
// format is any empty string will default to YAML.
//
// The binary formats (CBOR and MessagePack) will be converted to base64,
// ensuring that the returned value is always a valid string.
func (self *Transcriber) Stringify(value any) (string, error) {
	switch self.Format {
	case "yaml", "":
		return self.StringifyYAML(value)

	case "json":
		return self.StringifyJSON(value)

	case "xjson":
		return self.StringifyXJSON(value)

	case "xml":
		return self.StringifyXML(value)

	case "cbor":
		return self.StringifyCBOR(value)

	case "messagepack":
		return self.StringifyMessagePack(value)

	case "go":
		return self.StringifyGo(value)

	default:
		return "", fmt.Errorf("unsupported format: %q", self.Format)
	}
}

// Utils

func (self *Transcriber) fixIndentForTerminal() *Transcriber {
	if self.ForTerminal && (self.Indent != terminal.Indent) {
		self = self.Clone()
		self.Indent = terminal.Indent
		return self
	} else {
		return self
	}
}
