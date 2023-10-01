package transcribe

import (
	"io"
	"os"
	"strings"

	yamllexer "github.com/goccy/go-yaml/lexer"
	yamlprinter "github.com/goccy/go-yaml/printer"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"gopkg.in/yaml.v3"
)

// Note special handling when value is []any. In this case, instead of writing
// a YAML seq, it will treat the value as a list of documents separated by
// "---".
func (self *Transcriber) WriteYAML(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	if self.ForTerminal && terminal.Colorize {
		// Unfortunately we need to stringify before colorizing
		self = self.Clone()
		self.ForTerminal = false
		self.Indent = terminal.Indent

		if code, err := self.StringifyYAML(value); err == nil {
			return ColorizeYAML(code, writer)
		} else {
			return err
		}
	}

	self = self.fixIndentForTerminal()

	if self.Strict {
		var err error
		if value, err = ard.ToYAMLDocumentNode(value, true, self.Reflector); err != nil {
			return err
		}
	}

	encoder := yaml.NewEncoder(writer)

	// This might not work as expected for tabs!
	// BUG: currently does not allow an indent value of 1, see: https://github.com/go-yaml/yaml/issues/501
	encoder.SetIndent(len(self.Indent))

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

// Note special handling when value is []any. In this case, instead of writing
// a YAML seq, it will treat the value as a list of documents separated by
// "---".
func (self *Transcriber) StringifyYAML(value any) (string, error) {
	var writer strings.Builder
	self = self.Clone()
	self.Writer = &writer

	if err := self.WriteYAML(value); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func ColorizeYAML(code string, writer io.Writer) error {
	tokens := yamllexer.Tokenize(code)
	if _, err := io.WriteString(writer, YAMLColorPrinter.PrintTokens(tokens)); err == nil {
		return util.WriteNewline(writer)
	} else {
		return err
	}
}

// Utils

var YAMLColorPrinter = yamlprinter.Printer{
	String: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.BlueCode,
			Suffix: terminal.ResetCode,
		}
	},
	Number: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.MagentaCode,
			Suffix: terminal.ResetCode,
		}
	},
	Bool: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.CyanCode,
			Suffix: terminal.ResetCode,
		}
	},

	MapKey: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.GreenCode,
			Suffix: terminal.ResetCode,
		}
	},

	Anchor: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.RedCode,
			Suffix: terminal.ResetCode,
		}
	},
	Alias: func() *yamlprinter.Property {
		return &yamlprinter.Property{
			Prefix: terminal.YellowCode,
			Suffix: terminal.ResetCode,
		}
	},
}
