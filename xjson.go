package transcribe

import (
	"github.com/tliron/go-ard"
)

// If inPlace is false then the function is non-destructive:
// the written data structure is a [ard.ValidCopy] of the value
// argument. Otherwise, the value may be changed during
// preparing it for writing.
func (self *Transcriber) WriteXJSON(value any) error {
	if value_, err := ard.PrepareForEncodingXJSON(value, self.InPlace, self.Reflector); err == nil {
		return self.WriteJSON(value_)
	} else {
		return err
	}
}

func (self *Transcriber) StringifyXJSON(value any) (string, error) {
	if value_, err := ard.PrepareForEncodingXJSON(value, self.InPlace, self.Reflector); err == nil {
		return self.StringifyJSON(value_)
	} else {
		return "", err
	}
}
