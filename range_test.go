package elasticquery

import (
	"fmt"
	"testing"

	"github.com/mtlkn/json"
)

func TestRange(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{"age": { "lt": 30 }}`)
		q, ok := parseRange(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "age" || q.Right.Value != 30 || q.Right.IsExact {
			t.Fail()
		}
		if q.Left != nil {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = LT("age", 30)
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "lte": 30 }}`)
		q, ok = parseRange(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "age" || q.Right.Value != 30 || !q.Right.IsExact {
			t.Fail()
		}
		if q.Left != nil {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = LTE("age", 30)
		r = q.JSON().String()
		if l != r {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "gt": 30 }}`)
		q, ok = parseRange(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "age" || q.Left.Value != 30 || q.Left.IsExact {
			t.Fail()
		}
		if q.Right != nil {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = GT("age", 30)
		r = q.JSON().String()
		if l != r {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "gte": 30 }}`)
		q, ok = parseRange(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "age" || q.Left.Value != 30 || !q.Left.IsExact {
			t.Fail()
		}
		if q.Right != nil {
			t.Fail()
		}
		l = q.JSON().String()
		fmt.Println(l)

		q = GTE("age", 30)
		r = q.JSON().String()
		if l != r {
			t.Fail()
		}

		q = Range("age", 20, 30, true)
		q.From(RangeValue(21, false)).To(RangeValue(29, false))
		fmt.Println(q.JSON().String())
		if q.Field != "age" || q.Left.Value != 21 || q.Left.IsExact || q.Right.Value != 29 || q.Right.IsExact {
			t.Fail()
		}
	})

	t.Run("complex", func(t *testing.T) {
		jo, _ := json.ParseObjectString(`{
			"date": {
			  "gte": "now-30d",
			  "lte": "now",
			  "boost": 1.2,
			  "format": "yyyy-MM-dd",
			  "time_zone": "+01:00",
			  "relation": "CONTAINS"
			}
		  }`)

		q, ok := parseRange(jo)
		if !ok {
			t.Fail()
			return
		}
		if q.Field != "date" || q.Left.Value != "now-30d" || !q.Left.IsExact || q.Right.Value != "now" || !q.Right.IsExact {
			t.Fail()
		}
		if q.Format != "yyyy-MM-dd" || q.TimeZone != "+01:00" || q.Relation != "CONTAINS" || q.Boost != 1.2 {
			t.Fail()
		}
		l := q.JSON().String()
		fmt.Println(l)

		q = Range("date", "now-30d", "now", true)
		q.SetBoost(1.2).SetFormat("yyyy-MM-dd").SetRelation("CONTAINS").SetTimeZone("+01:00")
		r := q.JSON().String()
		if l != r {
			t.Fail()
		}
	})

	t.Run("errors", func(t *testing.T) {
		q := &RangeQuery{}
		if q.JSON() != nil {
			t.Fail()
		}

		if _, ok := parseRange(json.New()); ok {
			t.Fail()
		}

		jo, _ := json.ParseObjectString(`{"message": {}}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"pi": 3.14}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "lt": now }}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "lte": now }}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "gt": now }}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}

		jo, _ = json.ParseObjectString(`{"age": { "gte": now }}`)
		if _, ok := parseRange(jo); ok {
			t.Fail()
		}
	})
}
