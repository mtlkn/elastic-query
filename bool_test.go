package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestBool(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		q := And(Match("name", "YM"), Term("age", 26))
		if len(q.Must) != 1 || len(q.Filter) != 1 || len(q.Should) != 0 || len(q.MustNot) != 0 {
			t.Fail()
		}
		fmt.Println(q.JSON().String())

		q.And(LT("date", "now")).Or(Match("title", "test")).Not(Exists("status"))
		if len(q.Must) != 1 || len(q.Filter) != 2 || len(q.Should) != 1 || len(q.MustNot) != 1 {
			t.Fail()
		}
		fmt.Println(q.JSON().String())

		q = Or(Match("title", "test")).And(Match("summary", "joe biden"))
		if len(q.Must) != 1 || len(q.Filter) != 0 || len(q.Should) != 1 || len(q.MustNot) != 0 {
			t.Fail()
		}
		fmt.Println(q.JSON().String())

		q = Not(Term("id", 666)).And(Match("summary", "joe biden"))
		if len(q.Must) != 1 || len(q.Filter) != 0 || len(q.Should) != 0 || len(q.MustNot) != 1 {
			t.Fail()
		}
		fmt.Println(q.JSON().String())
	})

	t.Run("errors", func(t *testing.T) {
		if _, ok := parseBool(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"must": "abc"}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"filter": "abc"}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"must_not": "abc"}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"should": "abc"}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"filter":{"xyz":"abc"}}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"filter":[{"xyz":"abc"}]}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"filter":[123]}`)
		if _, ok := parseBool(jo); ok {
			t.Fail()
		}

		q := &BoolQuery{}
		if q.JSON() != nil {
			t.Fail()
		}

		q = And(Term("age", 26))
		if q.JSON().String() != Term("age", 26).JSON().String() {
			t.Fail()
		}

	})
}
