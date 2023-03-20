package elasticquery

import "github.com/mtlkn/json"

type FuzzyQuery struct {
	Field         string
	Value         string
	Boost         float64
	Fuzziness     *Fuzziness
	MaxExpansions int // max_expansions
}

func Fuzzy(field, value string) *FuzzyQuery {
	return &FuzzyQuery{
		Field: field,
		Value: value,
	}
}

func (q *FuzzyQuery) SetBoost(v float64) *FuzzyQuery {
	q.Boost = v
	return q
}

func (q *FuzzyQuery) SetFuzziness(v *Fuzziness) *FuzzyQuery {
	q.Fuzziness = v
	return q
}

func (q *FuzzyQuery) SetMaxExpansions(v int) *FuzzyQuery {
	q.MaxExpansions = v
	return q
}

func (q *FuzzyQuery) JSON() *json.Object {
	if q.Field == "" || q.Value == "" {
		return nil
	}

	jo := json.New().Add("value", q.Value)

	if q.Boost > 0 {
		jo.Add("boost", q.Boost)
	}

	if q.MaxExpansions > 0 {
		jo.Add("max_expansions", q.MaxExpansions)
	}

	q.Fuzziness.appendJSON(jo, "")

	return json.New().Add("fuzzy", json.New().Add(q.Field, jo))
}

func parseFuzzy(jo *json.Object) (*FuzzyQuery, bool) {
	if jo == nil || len(jo.Properties) != 1 {
		return nil, false
	}

	jp := jo.Properties[0]

	q := &FuzzyQuery{
		Field: jp.Name,
	}

	if jp.Value.Type == json.STRING {
		q.Value, _ = jp.Value.GetString()
		return q, true
	}

	var fz *json.Object

	if jp.Value.Type == json.OBJECT {
		o, ok := jp.Value.GetObject()
		if !ok || o == nil || len(o.Properties) == 0 {
			return nil, false
		}

		for _, p := range o.Properties {
			switch p.Name {
			case "value":
				q.Value, _ = p.Value.GetString()
			case "boost":
				q.Boost, _ = p.Value.GetFloat()
			case "max_expansions":
				q.MaxExpansions, _ = p.Value.GetInt()
			case "fuzziness", "prefix_length", "transpositions", "rewrite":
				if fz == nil {
					fz = json.New()
				}
				fz.Properties = append(fz.Properties, p)
			}
		}

		q.Fuzziness = parseFuzziness(fz)

		return q, true
	}

	return nil, false
}
