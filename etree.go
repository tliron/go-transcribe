package transcribe

import (
	"os"

	"github.com/beevik/etree"
	"github.com/tliron/kutil/terminal"
)

// Will change the indentation of the document.
func (self *Transcriber) WriteXMLDocument(xmlDocument *etree.Document) error {
	writer := self.Writer
	if writer == nil {
		writer = os.Stdout
	}

	if self.ForTerminal {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		// This might not work as expected for tabs!
		xmlDocument.Indent(len(self.Indent))
	}

	_, err := xmlDocument.WriteTo(writer)
	return err
}
