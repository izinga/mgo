package mgo

import (
	"testing"

	"github.com/izinga/mgo/bson"
)

func TestIsEmptyFilter(t *testing.T) {
	empty := []interface{}{
		nil,
		bson.M{},
		bson.D{},
		map[string]interface{}{},
		// the shape that caused the incident: a filter reduced at runtime
		func() interface{} { q := bson.M{"project": "x"}; q = bson.M{}; return q }(),
	}
	for i, q := range empty {
		if !isEmptyFilter(q) {
			t.Errorf("empty[%d] %#v: want empty, got not-empty", i, q)
		}
	}
	nonEmpty := []interface{}{
		bson.M{"project": "abc"},
		bson.D{{Name: "project", Value: "abc"}},
		map[string]interface{}{"a": 1},
	}
	for i, q := range nonEmpty {
		if isEmptyFilter(q) {
			t.Errorf("nonEmpty[%d] %#v: want not-empty, got empty", i, q)
		}
	}
}
