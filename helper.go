package yamlvector

import (
	"io"

	"github.com/koykov/vector"
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
