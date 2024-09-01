package yamlvector

func ensureDigit(c byte) bool {
	// todo check infinity
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E' || c == '.'
}
