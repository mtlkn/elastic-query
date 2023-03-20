package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestExists(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{ "field": "message", "boost": 1.2 }`)
		q, ok := parseExists(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "message" {
			t.Fail()
		}

		l := q.JSON().String()
		fmt.Println(l)

		q = Exists("message").SetBoost(1.2)
		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := Exists("")
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseExists(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"xyz": "message}`)
		if _, ok := parseExists(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"field": 3.14}`)
		if _, ok := parseExists(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"field": ""}`)
		if _, ok := parseExists(jo); ok {
			t.Fail()
		}
	})
}
