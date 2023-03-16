package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestMatch(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{"message": "this is a test"}`)
		q, ok := parseMatch("match", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "this is a test" {
			t.Fail()
		}
		if q.IsPhrase || q.IsPhrasePrefix {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Match("message", "this is a test")
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}
	})

	t.Run("complex query", func(t *testing.T) {
		s := `{
			"message": {
			  "query": "this is a test",
			  "analyzer": "test",
			  "boost": 1.2,
			  "operator": "and",
			  "minimum_should_match": 2,
			  "auto_generate_synonyms_phrase_query": false,
			  "zero_terms_query": "all",
			  "fuzziness": "AUTO",
			  "max_expansions": 60,
			  "prefix_length": 1,
			  "fuzzy_transpositions": false,
			  "fuzzy_rewrite": "top_terms_N",
			  "lenient": true
			}
		  }`
		jo, _ := json.ParseObjectString(s)
		q, ok := parseMatch("match", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "this is a test" {
			t.Fail()
		}
		if q.Boost != 1.2 || !q.IsAnd || !q.IsZeroTerms || !q.Lenient || q.IsSynonyms || q.MaxExpansions != 60 {
			t.Fail()
		}
		if q.Fuzziness == nil || q.Fuzziness.Value != "AUTO" || q.Fuzziness.PrefixLength != 1 || q.Fuzziness.Transpositions || q.Fuzziness.Rewrite != "top_terms_N" {
			t.Fail()
		}
		if q.IsPhrase || q.IsPhrasePrefix {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Match("message", "this is a test").SetBoost(1.2).SetAnd().SetZeroTerms()
		q.SetAnalyzer("test").SetMinMatch(2).SetSynonyms(false)
		q.SetMaxExpansions(60).SetLenient()
		q.SetFuzziness(FuzzyAuto(0, 0).SetPrefixLength(1).SetTranspositions(false).SetRewrite(REWRITE_TOP_TERMS_N))

		r := q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("exotic match queries", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{"message": "this is a test"}`)
		q, ok := parseMatch("match_phrase", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "this is a test" {
			t.Fail()
		}
		if !q.IsPhrase {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = MatchPhrase("message", "this is a test")
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}

		s := `{
			"message": {
			  "query": "this is a test",
			  "slop": 2
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		q, ok = parseMatch("match_phrase", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "this is a test" {
			t.Fail()
		}
		if !q.IsPhrase || q.Slop != 2 {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = MatchPhrase("message", "this is a test").SetSlop(2)

		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		s = `{
			"message": {
			  "query": "quick brown f",
			  "slop": 2
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		q, ok = parseMatch("match_phrase_prefix", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "quick brown f" {
			t.Fail()
		}
		if !q.IsPhrasePrefix || q.Slop != 2 {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = MatchPhrasePrefix("message", "quick brown f").SetSlop(2)

		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		s = `{
			"message": {
				"query": "quick brown f",
				"analyzer": "keywords"
			}
		  }`
		jo, _ = json.ParseObjectString(s)
		q, ok = parseMatch("match_bool_prefix", jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Query != "quick brown f" {
			t.Fail()
		}
		if !q.IsBoolPrefix || q.Analyzer != "keywords" {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = MatchBoolPrefix("message", "quick brown f").SetAnalyzer("keywords")

		r = q.JSON().String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := Match("text", "")
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseMatch("match", json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"message": {}}`)
		if _, ok := parseMatch("match", jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"pi": 3.14}`)
		if _, ok := parseMatch("match", jo); ok {
			t.Fail()
		}
	})
}
