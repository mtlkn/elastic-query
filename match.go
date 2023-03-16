package elasticquery

import (
	"strings"

	"github.com/mtlkn/json"
)

type MatchQuery struct {
	Field          string
	Query          string
	IsPhrase       bool // match_phrase
	IsPhrasePrefix bool // match_phrase_prefix
	IsBoolPrefix   bool // match_bool_prefix
	Analyzer       string
	Boost          float64
	IsAnd          bool // operator: "and"
	Slop           int  // for phrases only
	IsZeroTerms    bool // zero_terms_query
	IsSynonyms     bool // auto_generate_synonyms_phrase_query
	MaxExpansions  int  // max_expansions
	Fuzziness      *Fuzziness
	Lenient        bool
	MinMatch       any // minimum_should_match
}

func Match(field, query string) *MatchQuery {
	return &MatchQuery{
		Field:      strings.TrimSpace(field),
		Query:      strings.TrimSpace(query),
		IsSynonyms: true,
	}
}

func MatchPhrase(field, query string) *MatchQuery {
	return &MatchQuery{
		Field:      field,
		Query:      query,
		IsSynonyms: true,
		IsPhrase:   true,
	}
}

func MatchPhrasePrefix(field, query string) *MatchQuery {
	return &MatchQuery{
		Field:          field,
		Query:          query,
		IsSynonyms:     true,
		IsPhrase:       true,
		IsPhrasePrefix: true,
	}
}

func MatchBoolPrefix(field, query string) *MatchQuery {
	return &MatchQuery{
		Field:        field,
		Query:        query,
		IsSynonyms:   true,
		IsBoolPrefix: true,
	}
}

func (q *MatchQuery) SetAnd() *MatchQuery {
	q.IsAnd = true
	return q
}

func (q *MatchQuery) SetBoost(v float64) *MatchQuery {
	q.Boost = v
	return q
}

func (q *MatchQuery) SetAnalyzer(v string) *MatchQuery {
	q.Analyzer = v
	return q
}

func (q *MatchQuery) SetMinMatch(v any) *MatchQuery {
	q.MinMatch = v
	return q
}

func (q *MatchQuery) SetFuzziness(v *Fuzziness) *MatchQuery {
	q.Fuzziness = v
	return q
}

func (q *MatchQuery) SetSlop(v int) *MatchQuery {
	q.Slop = v
	return q
}

func (q *MatchQuery) SetMaxExpansions(v int) *MatchQuery {
	q.MaxExpansions = v
	return q
}

func (q *MatchQuery) SetZeroTerms() *MatchQuery {
	q.IsZeroTerms = true
	return q
}

func (q *MatchQuery) SetLenient() *MatchQuery {
	q.Lenient = true
	return q
}

func (q *MatchQuery) SetSynonyms(v bool) *MatchQuery {
	q.IsSynonyms = v
	return q
}

func (q *MatchQuery) JSON() *json.Object {
	if q.Field == "" || q.Query == "" {
		return nil
	}

	var name = "match"
	if q.IsPhrase {
		if q.IsPhrasePrefix {
			name = "match_phrase_prefix"
		} else {
			name = "match_phrase"
		}
	} else if q.IsBoolPrefix {
		name = "match_bool_prefix"
	}

	query := json.New()

	and := q.IsAnd && !q.IsPhrase
	if and || q.Boost > 0 || (q.Slop > 0 && q.IsPhrase) {
		jo := json.New().Add("query", q.Query)

		if q.Boost > 0 {
			jo.Add("boost", q.Boost)
		}

		if q.Analyzer != "" {
			jo.Add("analyzer", q.Analyzer)
		}

		if q.MinMatch != nil {
			jo.Add("minimum_should_match", q.MinMatch)
		}

		if and {
			jo.Add("operator", "and")
		}

		if q.Slop > 0 && q.IsPhrase {
			jo.Add("slop", q.Slop)
		}

		if q.IsZeroTerms {
			jo.Add("zero_terms_query", "all")
		}

		if !q.IsSynonyms {
			jo.Add("auto_generate_synonyms_phrase_query", false)
		}

		if q.MaxExpansions > 0 {
			jo.Add("max_expansions", q.MaxExpansions)
		}

		q.Fuzziness.appendJSON(jo)

		if q.Lenient {
			jo.Add("lenient", true)
		}

		query.Add(q.Field, jo)
	} else {
		query.Add(q.Field, q.Query)
	}

	return json.New().Add(name, query)
}

func parseMatch(name string, jo *json.Object) (*MatchQuery, bool) {
	if jo == nil || len(jo.Properties) != 1 {
		return nil, false
	}

	var phrase, phrasePrefix, boolPrefix bool

	if name == "match_phrase" {
		phrase = true
	} else if name == "match_phrase_prefix" {
		phrase = true
		phrasePrefix = true
	} else if name == "match_bool_prefix" {
		boolPrefix = true
	}

	jp := jo.Properties[0]

	q := &MatchQuery{
		Field:          jp.Name,
		IsPhrase:       phrase,
		IsPhrasePrefix: phrasePrefix,
		IsBoolPrefix:   boolPrefix,
		IsSynonyms:     true,
	}

	if jp.Value.Type == json.STRING {
		q.Query, _ = jp.Value.GetString()
		return q, true
	}

	var fz *json.Object

	if jp.Value.Type == json.OBJECT {
		jo, ok := jp.Value.GetObject()
		if !ok || jo == nil || len(jo.Properties) == 0 {
			return nil, false
		}

		for _, p := range jo.Properties {
			switch p.Name {
			case "query":
				q.Query, _ = p.Value.GetString()
			case "boost":
				q.Boost, _ = p.Value.GetFloat()
			case "analyzer":
				q.Analyzer, _ = p.Value.GetString()
			case "minimum_should_match":
				q.MinMatch, _ = p.Value.GetValue()
			case "operator":
				s, _ := p.Value.GetString()
				q.IsAnd = s == "and"
			case "slop":
				q.Slop, _ = p.Value.GetInt()
			case "zero_terms_query":
				s, _ := p.Value.GetString()
				q.IsZeroTerms = s == "all"
			case "auto_generate_synonyms_phrase_query":
				b, ok := p.Value.GetBool()
				if ok && !b {
					q.IsSynonyms = false
				}
			case "max_expansions":
				q.MaxExpansions, _ = p.Value.GetInt()
			case "fuzziness", "prefix_length", "fuzzy_transpositions", "fuzzy_rewrite":
				if fz == nil {
					fz = json.New()
				}
				fz.Properties = append(fz.Properties, p)
			case "lenient":
				q.Lenient, _ = p.Value.GetBool()
			}
		}

		q.Fuzziness = parseFuzziness(fz)

		return q, true
	}

	return nil, false
}
