package transcribe

import (
	"encoding/xml"
	"io"
	"os"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
)

// If inPlace is false then the function is non-destructive:
// the written data structure is a [ard.ValidCopy] of the value
// argument. Otherwise, the value may be changed during
// preparing it for writing.
func (self *Transcriber) WriteXML(value any) error {
	value, err := ard.PrepareForEncodingXML(value, self.InPlace, self.Reflector)
	if err != nil {
		return err
	}

	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	self = self.fixIndentForTerminal()

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", self.Indent)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	if self.ForTerminal || (self.Indent == "") {
		// When there's no indent the XML encoder does not emit a final newline
		// (We want it for consistency with YAML and JSON)
		return util.WriteNewline(writer)
	} else {
		return nil
	}
}

func (self *Transcriber) StringifyXML(value any) (string, error) {
	var writer strings.Builder
	self = self.Clone()
	self.Writer = &writer

	if err := self.WriteXML(value); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}
