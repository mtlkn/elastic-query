package elasticquery

import "github.com/mtlkn/json"

type BoolQuery struct {
	Must     []Query
	Filter   []Query
	Should   []Query
	MustNot  []Query
	Boost    float64
	MinMatch any // minimum_should_match
}

func (q *BoolQuery) JSON() *json.Object {
	var buf []*json.Object
	jo := json.New()

	buf = append(buf, appendBoolJSON(jo, "must", q.Must)...)
	buf = append(buf, appendBoolJSON(jo, "filter", q.Filter)...)
	buf = append(buf, appendBoolJSON(jo, "must_not", q.MustNot)...)
	buf = append(buf, appendBoolJSON(jo, "should", q.Should)...)

	if len(buf) == 0 {
		return nil
	}

	many := len(buf) > 1 || len(q.MustNot) > 0

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
		many = true
	}

	if q.MinMatch != nil && len(q.Should) > 0 && many {
		jo.Add("minimum_should_match", q.MinMatch)
	}

	if !many {
		return buf[0]
	}

	return json.New().Add("bool", jo)
}

func (q *BoolQuery) And(qs ...Query) *BoolQuery {
	for _, v := range qs {
		switch v.(type) {
		case *MatchQuery:
			q.Must = append(q.Must, v)
		case *TermsQuery, *RangeQuery, *WildcardQuery, *ExistsQuery, *FuzzyQuery, *IDsQuery:
			q.Filter = append(q.Filter, v)
		}
	}

	return q
}

func (q *BoolQuery) Or(qs ...Query) *BoolQuery {
	q.Should = append(q.Should, qs...)
	return q
}

func (q *BoolQuery) Not(qs ...Query) *BoolQuery {
	q.MustNot = append(q.MustNot, qs...)
	return q
}

func And(qs ...Query) *BoolQuery {
	q := &BoolQuery{}

	for _, v := range qs {
		switch v.(type) {
		case *MatchQuery:
			q.Must = append(q.Must, v)
		case *TermsQuery, *RangeQuery, *WildcardQuery, *ExistsQuery:
			q.Filter = append(q.Filter, v)
		}
	}

	return q
}

func Or(qs ...Query) *BoolQuery {
	return &BoolQuery{
		Should: qs,
	}
}

func Not(qs ...Query) *BoolQuery {
	return &BoolQuery{
		MustNot: qs,
	}
}

func parseBool(jo *json.Object) (*BoolQuery, bool) {
	if jo == nil || len(jo.Properties) == 0 {
		return nil, false
	}

	q := new(BoolQuery)
	ok := true

	for _, jp := range jo.Properties {
		switch jp.Name {
		case "must":
			q.Must, ok = parseBoolQueries(jp)
			if !ok {
				return nil, false
			}
		case "filter":
			q.Filter, ok = parseBoolQueries(jp)
			if !ok {
				return nil, false
			}
		case "should":
			q.Should, ok = parseBoolQueries(jp)
			if !ok {
				return nil, false
			}
		case "must_not":
			q.MustNot, ok = parseBoolQueries(jp)
			if !ok {
				return nil, false
			}
		case "boost":
			q.Boost, _ = jp.Value.GetFloat()
		case "minimum_should_match":
			q.MinMatch, _ = jp.Value.GetValue()
		}
	}

	return q, true
}

func parseBoolQueries(jp *json.Property) ([]Query, bool) {
	jo, ok := jp.Value.GetObject()
	if ok {
		q, ok := Parse(jo)
		if !ok {
			return nil, false
		}
		return []Query{q}, true
	}

	ja, ok := jp.Value.GetArray()
	if !ok {
		return nil, false
	}

	var qs []Query

	for _, jv := range ja.Values {
		jo, ok := jv.GetObject()
		if !ok {
			return nil, false
		}

		q, ok := Parse(jo)
		if !ok {
			return nil, false
		}

		qs = append(qs, q)
	}

	return qs, true
}

func appendBoolJSON(parent *json.Object, name string, qs []Query) []*json.Object {
	if len(qs) == 0 {
		return []*json.Object{}
	}

	if len(qs) == 1 {
		jo := qs[0].JSON()
		parent.Add(name, jo)
		return []*json.Object{jo}
	}

	var buf []*json.Object
	for _, q := range qs {
		jo := q.JSON()
		buf = append(buf, jo)
	}
	parent.Add(name, json.NewArray(buf))
	return buf
}
