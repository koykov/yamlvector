package yamlvector

type indent uint8

const (
	indentEqual indent = iota
	indentUp
	indentDown
)

func (vec *Vector) indentDW(p []byte, offset int) (indent, int) {
	dir := indentEqual
	i := offset
	for p[i] == ' ' {
		i++
	}
	d := i - offset
	switch {
	case d == vec.indent:
		return dir, offset
	case d < vec.indent:
		dir = indentUp
	case d > vec.indent:
		dir = indentDown
	}

	return dir, offset
}
