package yamlvector

import (
	"fmt"
)

type ttoken uint8

const (
	tokenError ttoken = iota
	tokenEOF
	tokenIndent
	tokenDash     // -
	tokenColon    // :
	tokenComma    // ,
	tokenLBracket // [
	tokenRBracket // ]
	tokenLBrace   // {
	tokenRBrace   // }
	tokenString
	tokenNumber
	tokenBool
	tokenNull
	tokenComment
	tokenAnchor    // &
	tokenAlias     // *
	tokenTag       // !!
	tokenDirective // %YAML, %TAG
)

func (tt ttoken) String() string {
	switch tt {
	case tokenError:
		return "Error"
	case tokenEOF:
		return "EOF"
	case tokenIndent:
		return "Indent"
	case tokenDash:
		return "Dash"
	case tokenColon:
		return "Colon"
	case tokenComma:
		return "Comma"
	case tokenLBracket:
		return "LBracket"
	case tokenRBracket:
		return "RBracket"
	case tokenLBrace:
		return "LBrace"
	case tokenRBrace:
		return "RBrace"
	case tokenString:
		return "String"
	case tokenNumber:
		return "Number"
	case tokenBool:
		return "Bool"
	case tokenNull:
		return "Null"
	case tokenComment:
		return "Comment"
	case tokenAnchor:
		return "Anchor"
	case tokenAlias:
		return "Alias"
	case tokenTag:
		return "Tag"
	case tokenDirective:
		return "Directive"
	default:
		return "Unknown"
	}
}

type token struct {
	typ    ttoken
	lo, hi uint64
}

func (t *token) String() string {
	return fmt.Sprintf("%s (%d:%d)", t.typ.String(), t.lo, t.hi)
}

func (t *token) setlo(lo uint64) *token {
	t.lo = lo
	return t
}

func (t *token) sethi(hi uint64) *token {
	t.hi = hi
	return t
}
