package yamlvector

import (
	"math"

	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

const bcamod = 10000

func tokenHash(src []byte, offset *int) (hsum uint64, eol bool) {
	lim, n := 8, len(src)
	if rest := n - *offset; rest < 8 {
		lim = rest
	}
	_ = tokend[math.MaxUint8-1]
	var i int
	for i = 0; i < lim && !tokend[i]; i++ {
		hsum = hsum | uint64(src[*offset])<<(i*8)
		*offset++
	}
	eol = tokend[i] || *offset == n
	return
}

func ensureNullOrBool(src []byte, offset *int, typ *vector.Type, b *bool) bool {
	hsum, eol := tokenHash(src, offset)
	if !eol {
		return false
	}
	idx := hsum % bcamod
	_ = bca[bcamod-1]
	if tuple := bca[idx]; tuple != nil {
		*typ = tuple.t
		*b = tuple.b
		return true
	}
	return false
}

type bct struct {
	t vector.Type
	b bool
}

var (
	bcRegistry = map[string]bct{
		"~":     {vector.TypeNull, false},
		"null":  {vector.TypeNull, false},
		"None":  {vector.TypeNull, false},
		"NONE":  {vector.TypeNull, false},
		"y":     {vector.TypeBool, true},
		"Y":     {vector.TypeBool, true},
		"on":    {vector.TypeBool, true},
		"On":    {vector.TypeBool, true},
		"ON":    {vector.TypeBool, true},
		"yes":   {vector.TypeBool, true},
		"Yes":   {vector.TypeBool, true},
		"YES":   {vector.TypeBool, true},
		"true":  {vector.TypeBool, true},
		"True":  {vector.TypeBool, true},
		"TRUE":  {vector.TypeBool, true},
		"n":     {vector.TypeBool, false},
		"N":     {vector.TypeBool, false},
		"no":    {vector.TypeBool, false},
		"No":    {vector.TypeBool, false},
		"NO":    {vector.TypeBool, false},
		"off":   {vector.TypeBool, false},
		"Off":   {vector.TypeBool, false},
		"OFF":   {vector.TypeBool, false},
		"false": {vector.TypeBool, false},
		"False": {vector.TypeBool, false},
		"FALSE": {vector.TypeBool, false},
	}
	bca    [bcamod]*bct
	tokend [math.MaxUint8]bool
)

func init() {
	for k, v := range bcRegistry {
		cpy := v
		x := binSafe(byteconv.S2B(k), 0, len(k))
		i := x % bcamod
		bca[i] = &cpy
	}

	tokend[' '] = true
	tokend[','] = true
	tokend['\n'] = true
	tokend['\r'] = true
	tokend['\t'] = true
	tokend[']'] = true
	tokend['['] = true
	tokend['}'] = true
	tokend['{'] = true
}
