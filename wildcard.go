package elasticquery

import (
	"github.com/mtlkn/json"
)

type WildcardQuery struct {
	Field    string
	Value    string
	IsPrefix bool // prefix query
	Boost    float64
	Uncased  bool // "case_insensitive"
	Rewrite  string
}

func Wildcard(field, value string) *WildcardQuery {
	return &WildcardQuery{
		Field: field,
		Value: value,
	}
}

func Prefix(field, value string) *WildcardQuery {
	return &WildcardQuery{
		Field:    field,
		Value:    value,
		IsPrefix: true,
	}
}

func (q *WildcardQuery) SetBoost(v float64) *WildcardQuery {
	q.Boost = v
	return q
}

func (q *WildcardQuery) CaseInsensitive() *WildcardQuery {
	q.Uncased = true
	return q
}

func (q *WildcardQuery) SetRewrite(v string) *WildcardQuery {
	q.Rewrite = v
	return q
}

func (q *WildcardQuery) JSON() *json.Object {
	if q.Field == "" || q.Value == "" {
		return nil
	}

	var name = "wildcard"
	if q.IsPrefix {
		name = "prefix"
	}

	jo := json.New().Add("value", q.Value)

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	if q.Uncased {
		jo.Add("case_insensitive", true)
	}

	appendRewriteJSON(jo, "rewrite", q.Rewrite)

	return json.New().Add(name, json.New().Add(q.Field, jo))
}

func parseWildcard(name string, jo *json.Object) (*WildcardQuery, bool) {
	if jo == nil || len(jo.Properties) != 1 {
		return nil, false
	}

	jp := jo.Properties[0]

	q := &WildcardQuery{
		Field:    jp.Name,
		IsPrefix: name == "prefix",
	}

	if jp.Value.Type == json.STRING {
		q.Value, _ = jp.Value.GetString()
		return q, true
	}

	if jp.Value.Type == json.OBJECT {
		jo, ok := jp.Value.GetObject()
		if !ok || jo == nil || len(jo.Properties) == 0 {
			return nil, false
		}

		for _, p := range jo.Properties {
			switch p.Name {
			case "value":
				q.Value, _ = p.Value.GetString()
			case "wildcard":
				if q.IsPrefix {
					return nil, false
				}
				q.Value, _ = p.Value.GetString()
			case "boost":
				q.Boost, _ = p.Value.GetFloat()
			case "rewrite":
				q.Rewrite = parseRewrite(p.Value)
			case "case_insensitive":
				q.Uncased, _ = p.Value.GetBool()
			}
		}

		return q, true
	}

	return nil, false
}
