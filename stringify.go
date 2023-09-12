package transcribe

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/kortschak/utter"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
)

func (self *Transcriber) Stringify(value any, format string) (string, error) {
	switch format {
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
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func (self *Transcriber) StringifyYAML(value any) (string, error) {
	var writer strings.Builder
	if err := self.WriteYAML(value, &writer); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func (self *Transcriber) StringifyJSON(value any) (string, error) {
	var writer strings.Builder
	if err := self.WriteJSON(value, &writer); err == nil {
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

func (self *Transcriber) StringifyXJSON(value any) (string, error) {
	if value_, err := ard.PrepareForEncodingXJSON(value, self.Reflector); err == nil {
		return self.StringifyJSON(value_)
	} else {
		return "", err
	}
}

func (self *Transcriber) StringifyXML(value any) (string, error) {
	var writer strings.Builder
	if err := self.WriteXML(value, &writer); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

// Note: will always use base64.
func (self *Transcriber) StringifyCBOR(value any) (string, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		return util.ToBase64(bytes), nil
	} else {
		return "", err
	}
}

// Note: will always use base64.
func (self *Transcriber) StringifyMessagePack(value any) (string, error) {
	var buffer bytes.Buffer
	encoder := ard.NewMessagePackEncoder(&buffer)
	if err := encoder.Encode(value); err == nil {
		return util.ToBase64(buffer.Bytes()), nil
	} else {
		return "", err
	}
}

func (self *Transcriber) StringifyGo(value any) (string, error) {
	return self.NewUtterConfig().Sdump(value), nil
}

func (self *Transcriber) NewUtterConfig() *utter.ConfigState {
	var config = utter.NewDefaultConfig()
	config.Indent = self.Indent
	config.SortKeys = true
	config.CommentPointers = true
	return config
}
