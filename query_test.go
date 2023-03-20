package elasticquery

import (
	"fmt"
	"os"
	"testing"

	"github.com/mtlkn/json"
)

func TestQuery(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		bs, _ := os.ReadFile("testdata/parse.json")
		jo, _ := json.ParseObject(bs)

		for _, jp := range jo.Properties {
			fmt.Println(jp.Name)

			ja, _ := jp.Value.GetArray()

			l, _ := ja.Values[0].GetObject()
			r, _ := ja.Values[1].GetObject()
			if len(r.Properties) == 0 {
				r = l
			}

			q, ok := Parse(l)
			if !ok {
				fmt.Println(l.String())
				t.Fail()
			}

			s := q.JSON().String()
			fmt.Println(s)

			if s != r.String() {
				fmt.Println(r.String())
				t.Fail()
			}

			fmt.Println()
		}
	})

	t.Run("errors", func(t *testing.T) {
		if _, ok := Parse(json.New()); ok {
			t.Fail()
		}

		if _, ok := Parse(json.New().Add("query", "match")); ok {
			t.Fail()
		}

		if _, ok := Parse(json.New().Add("query", json.New())); ok {
			t.Fail()
		}
	})
}
