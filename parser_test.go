package yamlvector

import (
	"testing"

	"github.com/koykov/vector"
)

func TestParser(t *testing.T) {
	t.Run("go-yaml/v3", func(t *testing.T) {
		// go-yaml tests as reference.
		stages := getStageMulti("tests_v3")
		if stages == nil {
			t.Fatal("stage tests_v3 not found")
		}
		skips := map[string]bool{"tests_v3/59": true}
		vec := NewVector()
		for i := range stages.buf {
			st := stages.buf[i]
			if skips[st.key] {
				continue
			}
			t.Run(st.key, func(t *testing.T) {
				vec.Reset()
				if err := vec.ParseCopy(st.origin); err != nil {
					t.Fatalf("parse error: %v at %d; src=%q", err, vec.ErrorOffset(), string(st.origin))
				}
			})
		}
	})

	t.Run("object", func(t *testing.T) {
		vec := NewVector()
		t.Run("multi_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeObject)
			assertString(t, vec, "name", "John Smith")
			assertString(t, vec, "desc", "foobar")
			assertNumber(t, vec, "age", 33)
		})
		t.Run("single_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeObject)
			assertString(t, vec, "name", "John Smith")
			assertNumber(t, vec, "age", 33)
		})
		t.Run("mixed", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeArray)
			assertNodeLen(t, vec.Root(), 2)
			assertNodeString(t, vec.Root().At(0).Look("name"), "John Smith")
			assertNodeNumber(t, vec.Root().At(0).Look("age"), 33)
			assertNodeString(t, vec.Root().At(1).Look("name"), "Mary Smith")
			assertNodeNumber(t, vec.Root().At(1).Look("age"), 27)
		})
		t.Run("alias", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeObject)

			aliases := vec.Root().Look("aliases")
			assertNodeType(t, aliases, vector.TypeArray)
			assertNodeLen(t, aliases, 2)
			assertNodeString(t, aliases.At(0).Look("datakey"), "dataval 1")
			assertNodeString(t, aliases.At(0).Look("moredata"), "morevals 1")
			assertNodeString(t, aliases.At(1).Look("datakey"), "dataval 2")
			assertNodeString(t, aliases.At(1).Look("moredata"), "morevals 2")

			config := vec.Root().Look("config")
			assertNodeType(t, config, vector.TypeArray)
			assertNodeLen(t, config, 1)
			assertNodeString(t, config.At(0).Look("datakey"), "dataval 1")
			assertNodeString(t, config.At(0).Look("moredata"), "morevals 1")
		})
	})

	t.Run("array", func(t *testing.T) {
		vec := NewVector()
		t.Run("single_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeArray)
			assertNodeLen(t, vec.Root(), 4)
			for i, want := range []string{"milk", "bread", "eggs", "juice"} {
				assertNodeString(t, vec.Root().At(i), want)
			}
		})
		t.Run("multi_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeArray)
			assertNodeLen(t, vec.Root(), 3)
			for i, want := range []string{"Casablanca", "Spellbound", "Notorious"} {
				assertNodeString(t, vec.Root().At(i), want)
			}
		})
		t.Run("mixed", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			assertType(t, vec, "", vector.TypeObject)

			men := vec.Root().Look("men")
			assertNodeType(t, men, vector.TypeArray)
			assertNodeLen(t, men, 2)
			assertNodeString(t, men.At(0), "John Smith")
			assertNodeString(t, men.At(1), "Bill Jones")

			women := vec.Root().Look("women")
			assertNodeType(t, women, vector.TypeArray)
			assertNodeLen(t, women, 2)
			assertNodeString(t, women.At(0), "Mary Smith")
			assertNodeString(t, women.At(1), "Susan Williams")
		})
	})

	t.Run("mixed", func(t *testing.T) {
		vec := NewVector()
		methods := []string{"newUri", "deleteUri", "randomUri"}

		check := func(t *testing.T) {
			assertType(t, vec, "", vector.TypeObject)
			bindings := vec.Root().Look("bindings")
			assertNodeType(t, bindings, vector.TypeArray)
			assertNodeLen(t, bindings, len(methods))
			for i, method := range methods {
				assertNodeString(t, bindings.At(i).Look("ircEvent"), "PRIVMSG")
				assertNodeString(t, bindings.At(i).Look("method"), method)
			}
		}

		t.Run("single_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			check(t)
		})
		t.Run("multi_line", func(t *testing.T) {
			vec = assertParse(t, vec, nil, 0)
			check(t)
		})
	})
}
