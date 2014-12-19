package outputs

import (
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"io"
	"log"
	"strings"
)

var newLine = []byte("\n")

func New(t typed.Typed) core.Output {
	switch strings.ToLower(t.String("type")) {
	case "file":
		return NewFile(t)
	case "stdout":
		return NewStdout(t)
	case "stderr":
		return NewStderr(t)
	default:
		log.Fatalf("invalid output type %v", string(t.MustBytes("")))
		return nil
	}
}

func writeTo(results *core.Results, writer io.Writer, newline bool) {
	writer.Write(results.Serialized())
	if newline {
		writer.Write(newLine)
	}
}
