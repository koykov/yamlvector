package yamlvector

func skipTillNlGeneric(b []byte) int {
	if len(b) == 0 {
		return -1
	}
	_ = b[len(b)-1]
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' || b[i] == '\r' {
			return i
		}
	}
	return -1
}
