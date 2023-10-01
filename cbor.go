package transcribe

import (
	"io"
	"os"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/util"
)

// When [Transcriber.Base64] is true will first convert to base64 and
// will also add a trailing newline if [Self.ForTerminal] is true.
func (self *Transcriber) WriteCBOR(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	if self.Base64 {
		if value_, err := self.StringifyCBOR(value); err == nil {
			if _, err = io.WriteString(writer, value_); err == nil {
				if self.ForTerminal {
					return util.WriteNewline(writer)
				} else {
					return nil
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		encoder := cbor.NewEncoder(writer)
		return encoder.Encode(value)
	}
}

// Will always use base64.
func (self *Transcriber) StringifyCBOR(value any) (string, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		return util.ToBase64(bytes), nil
	} else {
		return "", err
	}
}
