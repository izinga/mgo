package mgo

import (
	"reflect"

	"github.com/izinga/mgo/bson"
)

// isEmptyFilter reports whether a query selects the whole collection.
//
// Callers build filters in many shapes - nil, bson.M, bson.D, a plain map - and
// frequently reduce one to "everything" at runtime, e.g.
//
//	query := bson.M{"project": projectID}
//	if projectID == "" { query = bson.M{} }
//
// A static check cannot see that, which is exactly how the 2026-09 production
// incident was missed in review for six years.
func isEmptyFilter(q interface{}) bool {
	if q == nil {
		return true
	}
	switch v := q.(type) {
	case bson.M:
		return len(v) == 0
	case bson.D:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	}
	rv := reflect.ValueOf(q)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len() == 0
	}
	return false
}
