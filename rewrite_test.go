package elasticquery

import (
	"testing"

	"github.com/mtlkn/json"
)

func TestRewrite(t *testing.T) {
	t.Run("test rewrite", func(t *testing.T) {
		s := parseRewrite(nil)
		if s != "" {
			t.Fail()
		}

		jv := &json.Value{}
		s = parseRewrite(jv)
		if s != "" {
			t.Fail()
		}

		jo := json.New().Add("rewrite", REWRITE_SCORING_BOOLEAN)
		jv = jo.Properties[0].Value
		s = parseRewrite(jv)
		if s != REWRITE_SCORING_BOOLEAN {
			t.Fail()
		}

		jo = json.New().Add("rewrite", "xyz")
		jv = jo.Properties[0].Value
		s = parseRewrite(jv)
		if s != "" {
			t.Fail()
		}

		jo = json.New()
		appendRewriteJSON(jo, REWRITE_TOP_TERMS_N)
		if s, _ := jo.GetString("fuzzy_rewrite"); s != REWRITE_TOP_TERMS_N {
			t.Fail()
		}

		jo = json.New()
		appendRewriteJSON(jo, "")
		if len(jo.Properties) > 0 {
			t.Fail()
		}

		jo = json.New()
		appendRewriteJSON(jo, "xyz")
		if len(jo.Properties) > 0 {
			t.Fail()
		}
	})
}
