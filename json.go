package transcribe

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func (self *Transcriber) WriteJSON(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	if self.ForTerminal && terminal.Colorize {
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
	}

	self = self.fixIndentForTerminal()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", self.Indent)
	return encoder.Encode(value)
}

func (self *Transcriber) StringifyJSON(value any) (string, error) {
	var writer strings.Builder
	self = self.Clone()
	self.Writer = &writer

	if err := self.WriteJSON(value); err == nil {
		s := writer.String()
		if self.Indent == "" {
			// json.Encoder adds a "\n", unlike json.Marshal
			s = strings.TrimRight(s, "\n")
		}
		return s, nil
	} else {
		return "", err
	}
}

// Utils

func NewJSONColorFormatter(indent int) *prettyjson.Formatter {
	return &prettyjson.Formatter{
		KeyColor:        color.New(color.FgGreen),
		StringColor:     color.New(color.FgBlue),
		BoolColor:       color.New(color.FgCyan),
		NumberColor:     color.New(color.FgMagenta),
		NullColor:       color.New(color.FgCyan),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          indent,
		Newline:         "\n",
	}
}
