package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestIDs(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{ "values": ["abc", "xyz" ], "boost": 1.2 }`)
		q, ok := parseIDs(jo)
		if !ok {
			t.Fail()
			return
		}
		if len(q.Values) != 2 || q.Values[1] != "xyz" {
			t.Fail()
		}

		l := q.JSON().String()
		fmt.Println(l)

		q = IDs([]string{"abc", "xyz"}).SetBoost(1.2)
		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{ "values": "abc" }`)
		q, ok = parseIDs(jo)
		if !ok {
			t.Fail()
			return
		}
		if len(q.Values) != 1 || q.Values[0] != "abc" {
			t.Fail()
		}

		l = q.JSON().String()
		fmt.Println(l)

		q = ID("abc")
		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := IDs(nil)
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseIDs(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"xyz": "message}`)
		if _, ok := parseExists(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"values": 3.14}`)
		if _, ok := parseExists(jo); ok {
			t.Fail()
		}
	})
}
