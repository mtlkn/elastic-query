package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestFuzziness(t *testing.T) {
	t.Run("test fuzziness", func(t *testing.T) {
		fz := FuzzyLevenshtein(2)
		if fz.Value != 2 {
			t.Fail()
		}

		fz = FuzzyAuto(2, 4)
		if fz.Value != "AUTO:2,4" {
			t.Fail()
		}

		fz.appendJSON(nil, "")

		jo := json.New()
		fz.SetPrefixLength(1).SetRewrite(REWRITE_TOP_TERMS_BLENDED_FREQS_N).SetTranspositions(false)
		fz.appendJSON(jo, "fuzzy_")
		l := jo.String()
		fmt.Println(l)

		fz = parseFuzziness(jo)
		if fz == nil {
			t.Fail()
			return
		}

		jo = json.New()
		fz.appendJSON(jo, "fuzzy_")
		r := jo.String()
		if l != r {
			fmt.Println(r)
			t.Fail()
		}

		if parseFuzziness(nil) != nil {
			t.Fail()
		}
	})
}
