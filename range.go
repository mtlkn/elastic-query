package elasticquery

import (
	"github.com/mtlkn/json"
	"golang.org/x/exp/constraints"
)

const (
	RELATION_INTERSECTS = "INTERSECTS"
	RELATION_CONTAINS   = "CONTAINS"
	RELATION_WITHIN     = "WITHIN"
)

type RangeQueryValue struct {
	Value   any
	IsExact bool
}

func RangeValue[T constraints.Ordered](value T, exact bool) *RangeQueryValue {
	return &RangeQueryValue{
		Value:   value,
		IsExact: exact,
	}
}

type RangeQuery struct {
	Field    string
	Left     *RangeQueryValue
	Right    *RangeQueryValue
	Boost    float64
	Format   string
	Relation string
	TimeZone string
}

func Range[T constraints.Ordered](field string, from, to T, exact bool) *RangeQuery {
	return &RangeQuery{
		Field: field,
		Left:  RangeValue(from, exact),
		Right: RangeValue(to, exact),
	}
}

func GT[T constraints.Ordered](field string, value T) *RangeQuery {
	return &RangeQuery{
		Field: field,
		Left:  RangeValue(value, false),
	}
}

func GTE[T constraints.Ordered](field string, value T) *RangeQuery {
	return &RangeQuery{
		Field: field,
		Left:  RangeValue(value, true),
	}
}

func LT[T constraints.Ordered](field string, value T) *RangeQuery {
	return &RangeQuery{
		Field: field,
		Right: RangeValue(value, false),
	}
}

func LTE[T constraints.Ordered](field string, value T) *RangeQuery {
	return &RangeQuery{
		Field: field,
		Right: RangeValue(value, true),
	}
}

func (q *RangeQuery) SetBoost(v float64) *RangeQuery {
	q.Boost = v
	return q
}

func (q *RangeQuery) SetFormat(v string) *RangeQuery {
	q.Format = v
	return q
}

func (q *RangeQuery) SetRelation(v string) *RangeQuery {
	q.Relation = v
	return q
}

func (q *RangeQuery) SetTimeZone(v string) *RangeQuery {
	q.TimeZone = v
	return q
}

func (q *RangeQuery) From(value *RangeQueryValue) *RangeQuery {
	q.Left = value
	return q
}

func (q *RangeQuery) To(value *RangeQueryValue) *RangeQuery {
	q.Right = value
	return q
}

func (q *RangeQuery) JSON() *json.Object {
	if q.Field == "" || q.Left == nil && q.Right == nil {
		return nil
	}

	jo := json.New()

	if q.Left != nil {
		if q.Left.IsExact {
			jo.Add("gte", q.Left.Value)
		} else {
			jo.Add("gt", q.Left.Value)
		}
	}

	if q.Right != nil {
		if q.Right.IsExact {
			jo.Add("lte", q.Right.Value)
		} else {
			jo.Add("lt", q.Right.Value)
		}
	}

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	if q.Format != "" {
		jo.Add("format", q.Format)
	}

	if q.TimeZone != "" {
		jo.Add("time_zone", q.TimeZone)
	}

	if q.Relation == RELATION_CONTAINS || q.Relation == RELATION_WITHIN {
		jo.Add("relation", q.Relation)
	}

	return json.New().Add("range", json.New().Add(q.Field, jo))
}

func parseRange(jo *json.Object) (*RangeQuery, bool) {
	if jo == nil || len(jo.Properties) != 1 {
		return nil, false
	}

	jp := jo.Properties[0]

	q := &RangeQuery{
		Field: jp.Name,
	}

	o, ok := jp.Value.GetObject()
	if !ok || len(o.Properties) == 0 {
		return nil, false
	}

	for _, p := range o.Properties {
		switch p.Name {
		case "gt":
			v, err := p.Value.GetValue()
			if err != nil {
				return nil, false
			}
			q.Left = &RangeQueryValue{
				Value: v,
			}
		case "gte":
			v, err := p.Value.GetValue()
			if err != nil {
				return nil, false
			}
			q.Left = &RangeQueryValue{
				Value:   v,
				IsExact: true,
			}
		case "lt":
			v, err := p.Value.GetValue()
			if err != nil {
				return nil, false
			}
			q.Right = &RangeQueryValue{
				Value: v,
			}
		case "lte":
			v, err := p.Value.GetValue()
			if err != nil {
				return nil, false
			}
			q.Right = &RangeQueryValue{
				Value:   v,
				IsExact: true,
			}
		case "boost":
			q.Boost, _ = p.Value.GetFloat()
		case "format":
			q.Format, _ = p.Value.GetString()
		case "time_zone":
			q.TimeZone, _ = p.Value.GetString()
		case "relation":
			q.Relation, _ = p.Value.GetString()
		}
	}

	return q, true
}
