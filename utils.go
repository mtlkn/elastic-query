package elasticquery

import "github.com/mtlkn/json"

func getJSONStrings(jv *json.Value) ([]string, bool) {
	ja, ok := jv.GetArray()
	if ok {
		return ja.GetStrings()
	}

	s, ok := jv.GetString()
	if !ok {
		return nil, false
	}

	return []string{s}, true
}
