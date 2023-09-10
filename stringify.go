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

func Stringify(value any, format string, indent string, strict bool, reflector *ard.Reflector) (string, error) {
	switch format {
	case "yaml", "":
		return StringifyYAML(value, indent, strict, reflector)

	case "json":
		return StringifyJSON(value, indent)

	case "xjson":
		return StringifyXJSON(value, indent, reflector)

	case "xml":
		return StringifyXML(value, indent, reflector)

	case "cbor":
		return StringifyCBOR(value)

	case "messagepack":
		return StringifyMessagePack(value)

	case "go":
		return StringifyGo(value, indent)

	default:
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func StringifyYAML(value any, indent string, strict bool, reflector *ard.Reflector) (string, error) {
	var writer strings.Builder
	if err := WriteYAML(value, &writer, indent, strict, reflector); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func StringifyJSON(value any, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteJSON(value, &writer, indent); err == nil {
		s := writer.String()
		if indent == "" {
			// json.Encoder adds a "\n", unlike json.Marshal
			s = strings.TrimRight(s, "\n")
		}
		return s, nil
	} else {
		return "", err
	}
}

func StringifyXJSON(value any, indent string, reflector *ard.Reflector) (string, error) {
	if value_, err := ard.PrepareForEncodingXJSON(value, reflector); err == nil {
		return StringifyJSON(value_, indent)
	} else {
		return "", err
	}
}

func StringifyXML(value any, indent string, reflector *ard.Reflector) (string, error) {
	var writer strings.Builder
	if err := WriteXML(value, &writer, indent, reflector); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

// To Base64
func StringifyCBOR(value any) (string, error) {
	if bytes, err := cbor.Marshal(value); err == nil {
		return util.ToBase64(bytes), nil
	} else {
		return "", err
	}
}

// To Base64
func StringifyMessagePack(value any) (string, error) {
	var buffer bytes.Buffer
	encoder := util.NewMessagePackEncoder(&buffer)
	if err := encoder.Encode(value); err == nil {
		return util.ToBase64(buffer.Bytes()), nil
	} else {
		return "", err
	}
}

func StringifyGo(value any, indent string) (string, error) {
	return NewUtterConfig(indent).Sdump(value), nil
}

func NewUtterConfig(indent string) *utter.ConfigState {
	var config = utter.NewDefaultConfig()
	config.Indent = indent
	config.SortKeys = true
	config.CommentPointers = true
	return config
}
