package transcribe

import (
	"os"

	"github.com/kortschak/utter"
	"github.com/tliron/kutil/terminal"
)

func (self *Transcriber) WriteGo(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	var indent string
	if self.ForTerminal {
		indent = terminal.Indent
	} else {
		indent = self.Indent
	}

	NewUtterConfig(indent).Fdump(writer, value)
	return nil
}

func (self *Transcriber) StringifyGo(value any) (string, error) {
	return NewUtterConfig(self.Indent).Sdump(value), nil
}

// Utils

func NewUtterConfig(indent string) *utter.ConfigState {
	var config = utter.NewDefaultConfig()
	config.Indent = indent
	config.SortKeys = true
	config.CommentPointers = true
	return config
}
