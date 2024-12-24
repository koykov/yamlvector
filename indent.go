package yamlvector

type indent uint8

const (
	indentEqual indent = iota
	indentUp
	indentDown
)

// indent direction/width check
func (vec *Vector) indentDW(p []byte, offset, n int) (dir indent, d int) {
	i := offset
	_ = p[n-1]
	for p[i] == ' ' {
		i++
	}
	d = i - offset
	switch {
	case d < vec.indw:
		dir = indentUp
	case d > vec.indw:
		dir = indentDown
	}

	return
}
