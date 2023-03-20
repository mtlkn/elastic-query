package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestWildcards(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{"message": "ky*i?"}`)
		q, ok := parseWildcard("wildcard", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ky*i?" {
			t.Fail()
		}
		if q.IsPrefix {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Wildcard("message", "ky*i?")
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"message": "ky"}`)
		q, ok = parseWildcard("prefix", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ky" {
			t.Fail()
		}
		if !q.IsPrefix {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = Prefix("message", "ky")
		r = q.JSON().String()
		if l != r {
			t.Fail()
		}
	})

	t.Run("complex query", func(t *testing.T) {
		s := `{
			"message": {
			  "value": "ky*i?",
			  "boost": 1.2,
			  "case_insensitive": true,
			  "rewrite": "top_terms_N"
			}
		  }`
		jo, _ := json.ParseObjectString(s)
		q, ok := parseWildcard("wildcard", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ky*i?" {
			t.Fail()
		}
		if q.Boost != 1.2 || !q.Uncased || q.Rewrite != "top_terms_N" {
			t.Fail()
		}
		if q.IsPrefix {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Wildcard("message", "ky*i?").SetBoost(1.2).CaseInsensitive().SetRewrite(REWRITE_TOP_TERMS_N)
		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		s = `{
			"message": {
			  "wildcard": "ky*i?",
			  "boost": 1.2
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		q, ok = parseWildcard("wildcard", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ky*i?" {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)
		q = Wildcard("message", "ky*i?").SetBoost(1.2)
		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		s = `{
			"message": {
			  "value": "ky",
			  "boost": 1.2,
			  "case_insensitive": true,
			  "rewrite": "top_terms_N"
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		q, ok = parseWildcard("prefix", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ky" {
			t.Fail()
		}
		if q.Boost != 1.2 || !q.Uncased || q.Rewrite != "top_terms_N" {
			t.Fail()
		}
		if !q.IsPrefix {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)
		q = Prefix("message", "ky").SetBoost(1.2).CaseInsensitive().SetRewrite(REWRITE_TOP_TERMS_N)
		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := Wildcard("text", "")
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseWildcard("wildcard", json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"message": {}}`)
		if _, ok := parseWildcard("wildcard", jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"pi": 3.14}`)
		if _, ok := parseWildcard("wildcard", jo); ok {
			t.Fail()
		}

		s := `{
			"message": {
			  "wildcard": "ky",
			  "boost": 1.2,
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		if _, ok := parseWildcard("prefix", jo); ok {
			t.Fail()
		}
	})
}
