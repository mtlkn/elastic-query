package elasticquery

import (
	"strings"

	"github.com/mtlkn/json"
)

type ExistsQuery struct {
	Field string
	Boost float64
}

func Exists(field string) *ExistsQuery {
	return &ExistsQuery{
		Field: strings.TrimSpace(field),
	}
}

func (q *ExistsQuery) SetBoost(v float64) *ExistsQuery {
	q.Boost = v
	return q
}

func (q *ExistsQuery) JSON() *json.Object {
	if q.Field == "" {
		return nil
	}

	jo := json.New().Add("field", q.Field)

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	return json.New().Add("exists", jo)
}

func parseExists(jo *json.Object) (*ExistsQuery, bool) {
	if jo == nil || len(jo.Properties) == 0 {
		return nil, false
	}

	q := &ExistsQuery{}

	for _, jp := range jo.Properties {
		switch jp.Name {
		case "field":
			q.Field, _ = jp.Value.GetString()
		case "boost":
			q.Boost, _ = jp.Value.GetFloat()
		}
	}

	return q, q.Field != ""
}
