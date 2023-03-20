package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestFuzzy(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{"name": "ki"}`)
		q, ok := parseFuzzy(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ki" {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Fuzzy("name", "ki")
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}
	})

	t.Run("complex query", func(t *testing.T) {
		s := `{
			"name": {
			  "value": "ki",
			  "boost": 1.2,
			  "fuzziness": "AUTO",
			  "max_expansions": 60,
			  "prefix_length": 1,
			  "transpositions": false,
			  "rewrite": "top_terms_N"
			}
		  }`
		jo, _ := json.ParseObjectString(s)
		q, ok := parseFuzzy(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Value != "ki" || q.Field != "name" {
			t.Fail()
		}
		if q.Boost != 1.2 || q.Fuzziness == nil || q.Fuzziness.Value != "AUTO" || q.Fuzziness.PrefixLength != 1 || q.Fuzziness.Transpositions || q.Fuzziness.Rewrite != "top_terms_N" {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Fuzzy("name", "ki").SetBoost(1.2).SetMaxExpansions(60)
		q.SetFuzziness(FuzzyAuto(0, 0).SetPrefixLength(1).SetTranspositions(false).SetRewrite(REWRITE_TOP_TERMS_N))

		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := Fuzzy("name", "")
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseFuzzy(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"name": {}}`)
		if _, ok := parseFuzzy(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"pi": 3.14}`)
		if _, ok := parseFuzzy(jo); ok {
			t.Fail()
		}
	})
}
