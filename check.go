package yamlvector

import "unicode"

// todo use table approach

// isAnchorChar checks rune allow in anchor/alias
func isAnchorChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-'
}

// isTagChar checks rune allow in tag
func isTagChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '/' || r == '#' || r == ';' || r == ':' || r == '?' || r == '@' || r == '&' || r == '=' || r == '+' || r == '$' || r == ',' || r == '~' || r == '%' || r == '*' || r == '[' || r == ']'
}

// isDirectiveChar checks rune allow in directive
func isDirectiveChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '/' || r == '#' || r == ';' || r == ':' || r == '?' || r == '@' || r == '&' || r == '=' || r == '+' || r == '$' || r == ',' || r == '~' || r == '%' || r == '*' || r == '[' || r == ']'
}
