package transcribe

import (
	"github.com/tliron/go-ard"
)

//
// Transcriber
//

type Transcriber struct {
	Indent    string
	Strict    bool
	Pretty    bool // only applies to Print; when true then Indent is ignored
	Base64    bool
	Reflector *ard.Reflector
}

func NewTranscriber() *Transcriber {
	return new(Transcriber)
}

func (self *Transcriber) WithIndent(indent string) *Transcriber {
	return &Transcriber{
		Indent:    indent,
		Strict:    self.Strict,
		Pretty:    self.Pretty,
		Base64:    self.Base64,
		Reflector: self.Reflector,
	}
}
