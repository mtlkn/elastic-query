package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestTerms(t *testing.T) {
	t.Run("terms", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{ "user.id": "kimchy" }`)
		q, ok := parseTerms(jo)
		if !ok || q == nil || len(q.Values) != 1 || q.Values[0] != "kimchy" {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Term("user.id", "kimchy")
		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "user.id": { "value": "kimchy", "boost": 2.0, "case_insensitive": true } }`)
		q, ok = parseTerms(jo)
		if !ok || q == nil || len(q.Values) != 1 || q.Values[0] != "kimchy" {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = Term("user.id", "kimchy").SetBoost(2).CaseInsensitive()
		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "user.id": [ "kimchy", "chymki" ], "boost": 2.0, "case_insensitive": true }`)
		q, ok = parseTerms(jo)
		if !ok || q == nil || len(q.Values) != 2 || q.Values[0] != "kimchy" {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = Terms("user.id", []string{"kimchy", "chymki"}).SetBoost(2).CaseInsensitive()
		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := Terms("name", []string{})
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseTerms(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{ "user.id": kimchy }`)
		if _, ok := parseTerms(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "user.id": { "value": kimchy } }`)
		if _, ok := parseTerms(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "user.id": { "query": 1 } }`)
		if _, ok := parseTerms(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "ids": [ kimchy ] }`)
		if _, ok := parseTerms(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "boost": 1.2 }`)
		if _, ok := parseTerms(jo); ok {
			t.Fail()
		}
	})
}
