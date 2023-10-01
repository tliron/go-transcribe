package transcribe

import (
	"io"
	"os"
	"strings"

	"github.com/tliron/kutil/util"
)

// Converts the value to string (using [util.ToString]) and
// writes it, adding a newline if the string doesn't already
// have it.
func (self *Transcriber) WriteString(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	string_ := util.ToString(value)
	if _, err := io.WriteString(writer, string_); err == nil {
		if self.ForTerminal && !strings.HasSuffix(string_, "\n") {
			return util.WriteNewline(writer)
		}
		return nil
	} else {
		return err
	}
}
