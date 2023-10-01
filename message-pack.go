package transcribe

import (
	"bytes"
	"io"
	"os"

	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
)

// When [Transcriber.Base64] is true will first convert to base64 and
// will also add a trailing newline if [Self.ForTerminal] is true.
func (self *Transcriber) WriteMessagePack(value any) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	// MessagePack encoder has problems with map[any]any
	if value_, err := ard.ValidCopyMapsToStringMaps(value, nil); err == nil {
		if self.Base64 {
			if value__, err := self.StringifyMessagePack(value_); err == nil {
				if _, err = io.WriteString(writer, value__); err == nil {
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
			encoder := ard.NewMessagePackEncoder(writer)
			return encoder.Encode(value_)
		}
	} else {
		return err
	}
}

// Will always use base64.
func (self *Transcriber) StringifyMessagePack(value any) (string, error) {
	var buffer bytes.Buffer
	encoder := ard.NewMessagePackEncoder(&buffer)
	if err := encoder.Encode(value); err == nil {
		return util.ToBase64(buffer.Bytes()), nil
	} else {
		return "", err
	}
}
