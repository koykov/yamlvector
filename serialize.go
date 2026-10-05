package yamlvector

import (
	"errors"
	"io"

	"github.com/koykov/vector"
)

var (
	// Byte constants.
	btSpace = []byte(` `)
	btQuote = []byte(`"`)
	btComma = []byte(`,`)
	btColon = []byte(`:`)
	btNl    = []byte("\n")
	btTab   = []byte("\t")
	btArrO  = []byte(`[`)
	btArrC  = []byte(`]`)
	btArrE  = []byte(`[]`)
	btObjO  = []byte(`{`)
	btObjC  = []byte(`}`)
	btObjE  = []byte(`{}`)
	btNull  = []byte(`null`)

	ErrUnsupportedType = errors.New("unsupported type")
)

func serialize(w io.Writer, node *vector.Node, depth int, indent bool) (err error) {
	switch node.Type() {
	case vector.TypeNull:
		_, err = w.Write(btNull)
	case vector.TypeNumber, vector.TypeBool:
		_, err = w.Write(node.ForceBytes())
	case vector.TypeString:
		_, err = w.Write(btQuote)
		_, err = w.Write(node.RawBytes())
		_, err = w.Write(btQuote)
	case vector.TypeAlias:
		if child := node.FirstChild(); child != nil && child != node {
			err = serialize(w, child, depth, indent)
		} else {
			_, err = w.Write(btNull)
		}
	case vector.TypeArray:
		if node.Limit() == 0 {
			_, err = w.Write(btArrE)
		} else {
			_, err = w.Write(btArrO)
			if indent {
				_, err = w.Write(btNl)
			}
			node.Each(func(idx int, node *vector.Node) {
				if idx > 0 {
					_, err = w.Write(btComma)
					if indent {
						_, err = w.Write(btNl)
					}
				}
				if indent {
					writePad(w, node.Depth())
				}
				err = serialize(w, node, depth+1, indent)
			})
			if indent {
				_, err = w.Write(btNl)
				writePad(w, node.Depth())
			}
			_, err = w.Write(btArrC)
		}
	case vector.TypeObject:
		if node.Limit() == 0 {
			_, err = w.Write(btObjE)
		} else {
			_, err = w.Write(btObjO)
			if indent {
				_, err = w.Write(btNl)
			}
			node.Each(func(idx int, node *vector.Node) {
				if idx > 0 {
					_, err = w.Write(btComma)
					if indent {
						_, err = w.Write(btNl)
					}
				}
				if indent {
					writePad(w, node.Depth())
				}
				_, err = w.Write(btQuote)
				_, err = w.Write(node.KeyBytes())
				_, err = w.Write(btQuote)
				_, err = w.Write(btColon)
				if indent {
					_, err = w.Write(btSpace)
				}
				err = serialize(w, node, depth+1, indent)
			})
			if indent {
				_, err = w.Write(btNl)
				writePad(w, node.Depth())
			}
			_, err = w.Write(btObjC)
		}
	default:
		err = ErrUnsupportedType
	}
	return
}

// Write number of tabs to w.
func writePad(w io.Writer, cnt int) {
	for i := 0; i < cnt; i++ {
		_, _ = w.Write(btTab)
	}
}
