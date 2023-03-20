package elasticquery

import (
	"github.com/mtlkn/json"
)

type IDsQuery struct {
	Values []string
	Boost  float64
}

func ID(id string) *IDsQuery {
	return &IDsQuery{
		Values: []string{id},
	}
}

func IDs(ids []string) *IDsQuery {
	return &IDsQuery{
		Values: ids,
	}
}

func (q *IDsQuery) SetBoost(v float64) *IDsQuery {
	q.Boost = v
	return q
}

func (q *IDsQuery) JSON() *json.Object {
	if len(q.Values) == 0 {
		return nil
	}

	jo := json.New().Add("values", q.Values)

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	return json.New().Add("ids", jo)
}

func parseIDs(jo *json.Object) (*IDsQuery, bool) {
	if jo == nil || len(jo.Properties) == 0 {
		return nil, false
	}

	q := &IDsQuery{}

	for _, jp := range jo.Properties {
		switch jp.Name {
		case "values":
			q.Values, _ = getJSONStrings(jp.Value)
		case "boost":
			q.Boost, _ = jp.Value.GetFloat()
		}
	}

	return q, len(q.Values) > 0
}
