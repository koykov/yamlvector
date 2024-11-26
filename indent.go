package yamlvector

type indent uint8

const (
	indentEqual indent = iota
	indentUp
	indentDown
)

// indent direction/width check
func (vec *Vector) indentDW(p []byte, offset, n int) (indent, int, int) {
	dir := indentEqual
	i := offset
	_ = p[n-1]
	for p[i] == ' ' {
		i++
	}
	d := i - offset
	switch {
	case d < vec.indent:
		dir = indentUp
	case d > vec.indent:
		dir = indentDown
	}
	offset = i

	return dir, offset, d
}
