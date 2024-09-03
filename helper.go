package yamlvector

import (
	"io"

	"github.com/koykov/vector"
)

const (
	flagEscapedString = 0
	flagFoldBlock     = 1
)

type Helper struct{}

var (
	helper = Helper{}
)

func (h Helper) Indirect(p *vector.Byteptr) []byte {
	return p.RawBytes()
}

func (h Helper) Beautify(w io.Writer, node *vector.Node) error {
	_, _ = w, node
	return nil
}

func (h Helper) Marshal(w io.Writer, node *vector.Node) error {
	_, _ = w, node
	return nil
}
