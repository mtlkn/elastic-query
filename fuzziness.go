package elasticquery

import (
	"fmt"

	"github.com/mtlkn/json"
)

type Fuzziness struct {
	Value          any
	PrefixLength   int    // prefix_length
	Transpositions bool   // fuzzy_transpositions
	Rewrite        string // fuzzy_rewrite
}

func FuzzyLevenshtein(value int) *Fuzziness {
	return &Fuzziness{
		Value:          value,
		Transpositions: true,
	}
}

func FuzzyAuto(low, high int) *Fuzziness {
	value := "AUTO"
	if high > 0 {
		value = fmt.Sprintf("AUTO:%d,%d", low, high)
	}

	return &Fuzziness{
		Value:          value,
		Transpositions: true,
	}
}

func (fuzzy *Fuzziness) SetPrefixLength(v int) *Fuzziness {
	fuzzy.PrefixLength = v
	return fuzzy
}

func (fuzzy *Fuzziness) SetTranspositions(v bool) *Fuzziness {
	fuzzy.Transpositions = v
	return fuzzy
}

func (fuzzy *Fuzziness) SetRewrite(v string) *Fuzziness {
	fuzzy.Rewrite = v
	return fuzzy
}

func (fuzzy *Fuzziness) appendJSON(parent *json.Object, prefix string) {
	if fuzzy == nil || fuzzy.Value == nil || parent == nil {
		return
	}

	parent.Add("fuzziness", fuzzy.Value)

	if fuzzy.PrefixLength > 0 {
		parent.Add("prefix_length", fuzzy.PrefixLength)
	}

	if !fuzzy.Transpositions {
		parent.Add(prefix+"transpositions", false)
	}

	appendRewriteJSON(parent, prefix+"rewrite", fuzzy.Rewrite)
}

func parseFuzziness(jo *json.Object) *Fuzziness {
	if jo == nil || len(jo.Properties) == 0 {
		return nil
	}

	fuzzy := &Fuzziness{
		Transpositions: true,
	}

	for _, jp := range jo.Properties {
		switch jp.Name {
		case "fuzziness":
			fuzzy.Value, _ = jp.Value.GetValue()
		case "prefix_length":
			fuzzy.PrefixLength, _ = jp.Value.GetInt()
		case "transpositions", "fuzzy_transpositions":
			b, ok := jp.Value.GetBool()
			if ok && !b {
				fuzzy.Transpositions = false
			}
		case "rewrite", "fuzzy_rewrite":
			fuzzy.Rewrite = parseRewrite(jp.Value)
		}
	}

	return fuzzy
}
