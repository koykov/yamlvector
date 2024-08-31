package yamlvector

func ensureTrueBin(src []byte, poffset *int) (ok bool) {
	bin1 := binUnsafe(src, *poffset, 1)
	bin2 := binUnsafe(src, *poffset, 2)
	bin3 := binUnsafe(src, *poffset, 3)
	bin4 := binUnsafe(src, *poffset, 4)
	_ = binBoolTrue[10]
	switch {
	case bin1 == binBoolTrue[0] || bin1 == binBoolTrue[1]:
		*poffset += 1
		ok = true
	case bin2 == binBoolTrue[2] || bin2 == binBoolTrue[3] || bin2 == binBoolTrue[4]:
		*poffset += 2
		ok = true
	case bin3 == binBoolTrue[5] || bin3 == binBoolTrue[6] || bin3 == binBoolTrue[7]:
		*poffset += 3
		ok = true
	case bin4 == binBoolTrue[8] || bin4 == binBoolTrue[9] || bin4 == binBoolTrue[10]:
		*poffset += 4
		ok = true
	}
	return
}

func ensureFalseBin(src []byte, poffset *int) (ok bool) {
	bin2 := binUnsafe(src, *poffset, 2)
	bin3 := binUnsafe(src, *poffset, 3)
	bin5 := binUnsafe(src, *poffset, 4)
	_ = binBoolFalse[8]
	switch {
	case bin2 == binBoolFalse[0] || bin2 == binBoolFalse[1] || bin2 == binBoolFalse[2]:
		*poffset += 2
		ok = true
	case bin3 == binBoolFalse[3] || bin3 == binBoolFalse[4] || bin3 == binBoolFalse[5]:
		*poffset += 3
		ok = true
	case bin5 == binBoolFalse[6] || bin5 == binBoolFalse[7] || bin5 == binBoolFalse[8]:
		*poffset += 5
		ok = true
	}
	return
}

var (
	bBoolTrue = [11][]byte{
		[]byte("y"),
		[]byte("Y"),
		[]byte("on"),
		[]byte("On"),
		[]byte("ON"),
		[]byte("yes"),
		[]byte("Yes"),
		[]byte("YES"),
		[]byte("true"),
		[]byte("True"),
		[]byte("TRUE"),
	}
	binBoolTrue = [11]uint64{}
	bBoolFalse  = [9][]byte{
		[]byte("no"),
		[]byte("No"),
		[]byte("NO"),
		[]byte("off"),
		[]byte("Off"),
		[]byte("OFF"),
		[]byte("false"),
		[]byte("False"),
		[]byte("FALSE"),
	}
	binBoolFalse = [9]uint64{}
)

func init() {
	for i := 0; i < len(bBoolTrue); i++ {
		binBoolTrue[i] = binSafe(bBoolTrue[i], 0, len(bBoolTrue[i]))
	}
	for i := 0; i < len(bBoolFalse); i++ {
		binBoolFalse[i] = binSafe(bBoolFalse[i], 0, len(bBoolFalse[i]))
	}
}
