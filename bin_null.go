package yamlvector

func ensureNullBin(src []byte, poffset *int) (ok bool) {
	bin1 := binUnsafe(src, *poffset, 1)
	bin4 := binUnsafe(src, *poffset, 4)
	_ = binNull[3]
	switch {
	case bin1 == binNull[0]:
		*poffset += 1
		ok = true
	case bin4 == binNull[1] || bin4 == binNull[2] || bin4 == binNull[3]:
		*poffset += 4
		ok = true
	}
	return
}

var (
	bNull = [4][]byte{
		[]byte("~"),
		[]byte("null"),
		[]byte("None"),
		[]byte("NONE"),
	}
	binNull [4]uint64
)

func init() {
	for i := 0; i < len(bNull); i++ {
		binNull[i] = binSafe(bNull[i], 0, len(bNull[i]))
	}
}
