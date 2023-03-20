package elasticquery

import "github.com/mtlkn/json"

const (
	REWRITE_CONSTANT_SCORE            = "constant_score"
	REWRITE_CONSTANT_SCORE_BOOLEAN    = "constant_score_boolean"
	REWRITE_SCORING_BOOLEAN           = "scoring_boolean"
	REWRITE_TOP_TERMS_BLENDED_FREQS_N = "top_terms_blended_freqs_N"
	REWRITE_TOP_TERMS_BOOST_N         = "top_terms_boost_N"
	REWRITE_TOP_TERMS_N               = "top_terms_N"
)

func parseRewrite(jv *json.Value) string {
	if jv == nil {
		return ""
	}

	s, _ := jv.GetString()
	if s == "" {
		return ""
	}

	for _, v := range validRewrites() {
		if v == s {
			return s
		}
	}

	return ""
}

func appendRewriteJSON(parent *json.Object, name, value string) {
	if value == "" {
		return
	}

	for _, v := range validRewrites() {
		if v == value {
			parent.Add(name, value)
			return
		}
	}
}

func validRewrites() []string {
	return []string{
		REWRITE_CONSTANT_SCORE,
		REWRITE_CONSTANT_SCORE_BOOLEAN,
		REWRITE_SCORING_BOOLEAN,
		REWRITE_TOP_TERMS_BLENDED_FREQS_N,
		REWRITE_TOP_TERMS_BOOST_N,
		REWRITE_TOP_TERMS_N,
	}
}
