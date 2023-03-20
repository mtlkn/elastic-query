package elasticquery

import (
	"github.com/mtlkn/json"
	"golang.org/x/exp/constraints"
)

type TermsQuery struct {
	Field   string
	Values  []any
	Boost   float64
	Uncased bool // "case_insensitive"; term query (single value) only
}

func Term[T constraints.Ordered | ~bool](field string, value T) *TermsQuery {
	return &TermsQuery{
		Field:  field,
		Values: []any{value},
	}
}

func Terms[T constraints.Ordered](field string, values []T) *TermsQuery {
	var vs []any
	for _, v := range values {
		vs = append(vs, any(v))
	}

	return &TermsQuery{
		Field:  field,
		Values: vs,
	}
}

func (q *TermsQuery) SetBoost(v float64) *TermsQuery {
	q.Boost = v
	return q
}

func (q *TermsQuery) CaseInsensitive() *TermsQuery {
	q.Uncased = true
	return q
}

func (q *TermsQuery) JSON() *json.Object {
	if q.Field == "" || len(q.Values) == 0 {
		return nil
	}

	if len(q.Values) == 1 {
		jo := json.New().Add("value", q.Values[0])

		if q.Boost > 0 {
			jo.Add("boost", q.Boost)
		}

		if q.Uncased {
			jo.Add("case_insensitive", true)
		}

		return json.New().Add("term", json.New().Add(q.Field, jo))
	}

	jo := json.New().Add(q.Field, q.Values)

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	return json.New().Add("terms", jo)
}

func parseTerms(jo *json.Object) (*TermsQuery, bool) {
	if jo == nil || len(jo.Properties) == 0 {
		return nil, false
	}

	var (
		field   string
		values  []any
		boost   float64
		uncased bool
	)

	for _, jp := range jo.Properties {
		if jp.Value.Type == json.OBJECT {
			field = jp.Name

			o, _ := jp.Value.GetObject()
			p, ok := o.GetProperty("value")
			if !ok {
				return nil, false
			}

			vs, ok := parseTermValues(p)
			if !ok {
				return nil, false
			}
			values = vs

			boost, _ = o.GetFloat("boost")
			uncased, _ = o.GetBool("case_insensitive")

			break
		}

		if jp.Name != "boost" && jp.Name != "case_insensitive" {
			field = jp.Name

			vs, ok := parseTermValues(jp)
			if !ok {
				return nil, false
			}
			values = vs

			boost, _ = jo.GetFloat("boost")
			uncased, _ = jo.GetBool("case_insensitive")
		}
	}

	if field == "" || len(values) == 0 {
		return nil, false
	}

	return &TermsQuery{
		Field:   field,
		Values:  values,
		Boost:   boost,
		Uncased: uncased,
	}, true
}

func parseTermValues(jp *json.Property) ([]any, bool) {
	ja, ok := jp.Value.GetArray()
	if ok {
		var vs []any

		for _, jv := range ja.Values {
			v, err := jv.GetValue()
			if err != nil || v == nil {
				return nil, false
			}

			vs = append(vs, v)
		}

		return vs, true
	}

	v, err := jp.Value.GetValue()
	if err != nil {
		return nil, false
	}

	return []any{v}, true
}
