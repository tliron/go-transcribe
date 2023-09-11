package transcribe

import (
	"io"
)

const DIRECTORY_WRITE_PERMISSIONS = 0700
const FILE_WRITE_PERMISSIONS = 0600

// Convenience function. If output is an empty string will be identical
// to [Transcriber.Print], otherwise will call [Transcriber.WriteToFile]
// on the output, ignoring the writer argument.
func (self *Transcriber) WriteOrPrint(value any, writer io.Writer, output string, format string) error {
	if output == "" {
		return self.Print(value, writer, format)
	} else {
		return self.WriteToFile(value, output, format)
	}
}
