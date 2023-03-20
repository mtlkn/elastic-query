package elasticquery

import (
	"github.com/mtlkn/json"
)

type Query interface {
	JSON() *json.Object
}

func Parse(jo *json.Object) (Query, bool) {
	if jo == nil || len(jo.Properties) != 1 {
		return nil, false
	}

	jp := jo.Properties[0]

	o, ok := jp.Value.GetObject()
	if !ok {
		return nil, false
	}

	switch jp.Name {
	case "bool":
		return parseBool(o)
	case "match", "match_phrase", "match_phrase_prefix", "match_bool_prefix":
		return parseMatch(jp.Name, o)
	case "term", "terms":
		return parseTerms(o)
	case "prefix", "wildcard":
		return parseWildcard(jp.Name, o)
	case "range":
		return parseRange(o)
	case "exists":
		return parseExists(o)
	case "fuzzy":
		return parseFuzzy(o)
	case "ids":
		return parseIDs(o)
	}

	return nil, false
}
