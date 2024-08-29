package yamlvector

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

type stage struct {
	key string

	origin, json []byte
}

type multiStage struct {
	key string
	buf []stage
}

var (
	stages      []stage
	multiStages []multiStage
	bNl         = []byte("\n")
)

func init() {
	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" && !strings.Contains(filepath.Base(path), ".json") {
			st := stage{}
			st.key = strings.Replace(filepath.Base(path), ".yaml", "", 1)
			st.origin, _ = os.ReadFile(path)
			if st.json, _ = os.ReadFile(strings.Replace(path, ".yaml", ".json", 1)); len(st.json) > 0 {
				st.json = bytealg.Trim(st.json, bNl)
			}
			stages = append(stages, st)
			return nil
		}

		if info.IsDir() && path != "testdata" {
			mstg := multiStage{key: filepath.Base(path)}
			_ = filepath.Walk(path, func(path1 string, info1 os.FileInfo, err1 error) error {
				if filepath.Ext(path1) == ".yaml" && !strings.Contains(filepath.Base(path1), ".json") {
					st := stage{}
					st.key = strings.Replace(filepath.Base(path1), ".yaml", "", 1)
					st.origin, _ = os.ReadFile(path1)
					if st.json, _ = os.ReadFile(strings.Replace(path1, ".yaml", ".json", 1)); len(st.json) > 0 {
						st.json = bytealg.Trim(st.json, bNl)
					}
					mstg.buf = append(mstg.buf, st)
					return nil
				}
				return nil
			})
			multiStages = append(multiStages, mstg)
		}
		return nil
	})
}

func getStage(key string) (st *stage) {
	for i := 0; i < len(stages); i++ {
		st1 := &stages[i]
		if st1.key == key {
			st = st1
		}
	}
	return st
}

func getStageMulti(key string) (st *multiStage) {
	for i := 0; i < len(multiStages); i++ {
		st1 := &multiStages[i]
		if st1.key == key {
			st = st1
		}
	}
	return st
}

func getTBName(tb testing.TB) string {
	key := tb.Name()
	return key[strings.Index(key, "/")+1:]
}

func assertParse(tb testing.TB, dst *Vector, err error, errOffset int) *Vector {
	dst, _ = assertParseStage(tb, dst, err, errOffset)
	return dst
}

func assertParseStage(tb testing.TB, dst *Vector, err error, errOffset int) (*Vector, *stage) {
	key := getTBName(tb)
	st := getStage(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			if !errors.Is(err1, err) || dst.ErrorOffset() != errOffset {
				tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			}
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
		}
	}
	return dst, st
}

func assertParseMulti(tb testing.TB, dst *Vector, buf *bytes.Buffer, err error, errOffset int) *Vector {
	key := getTBName(tb)
	mst := getStageMulti(key)
	if mst == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	for i := 0; i < len(mst.buf); i++ {
		st := &mst.buf[i]
		err1 := dst.ParseCopy(st.origin)
		if err1 != nil {
			if err != nil {
				if !errors.Is(err1, err) || dst.ErrorOffset() != errOffset {
					tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
				}
			} else {
				tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
			}
		}
		root := dst.RootTop()
		buf.Reset()
		_ = root.Beautify(buf)
		if fmt1 := buf.Bytes(); !bytes.Equal(fmt1, st.json) {
			tb.Fatalf("node mismatch, need '%s'\ngot '%s'", string(st.json), string(fmt1))
		}
	}
	return dst
}

func assertType(tb testing.TB, vec *Vector, path string, typ vector.Type) {
	if typ1 := vec.Dot(path).Type(); typ1 != typ {
		tb.Error("type mismatch, need", typ, "got", typ1)
	}
}

func assertBool(tb testing.TB, vec *Vector, path string, val bool) {
	if val1 := vec.Dot(path).Bool(); val1 != val {
		tb.Error("type mismatch, need", val, "got", val1)
	}
}
