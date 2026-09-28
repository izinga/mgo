package mgo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/izinga/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/mgocompat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// withDriver switches the package to the official-driver path for one test.
func withDriver(t *testing.T) {
	old := UseMongoDriver
	UseMongoDriver = true
	t.Cleanup(func() { UseMongoDriver = old })
}

// A session with no driver client used to dereference nil and panic; every
// operation must now return ErrNoDriverClient.
func TestNoDriverClientIsAnError(t *testing.T) {
	withDriver(t)
	c := NewDummySession("no_client").DB("no_client").C("c")
	checks := map[string]func() error{
		"Insert":    func() error { return c.Insert(bson.M{"a": 1}) },
		"Update":    func() error { return c.Update(bson.M{"a": 1}, bson.M{"$set": bson.M{"a": 2}}) },
		"UpdateId":  func() error { return c.UpdateId(bson.NewObjectId(), bson.M{"$set": bson.M{"a": 2}}) },
		"UpdateAll": func() error { _, err := c.UpdateAll(bson.M{}, bson.M{"$set": bson.M{"a": 2}}); return err },
		"Remove":    func() error { return c.Remove(bson.M{"a": 1}) },
		"RemoveAll": func() error { _, err := c.RemoveAll(bson.M{"a": 1}); return err },
		"One":       func() error { var m bson.M; return c.Find(bson.M{"a": 1}).One(&m) },
		"All":       func() error { var m []bson.M; return c.Find(bson.M{"a": 1}).All(&m) },
		"Count":     func() error { _, err := c.Find(bson.M{"a": 1}).Count(); return err },
		"FindId":    func() error { var m bson.M; return c.FindId(bson.NewObjectId()).One(&m) },
	}
	for name, fn := range checks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked: %v", name, r)
				}
			}()
			if err := fn(); !errors.Is(err, ErrNoDriverClient) {
				t.Errorf("%s: %v, want ErrNoDriverClient", name, err)
			}
		}()
	}
}

// FindId must work on the driver path (it passed this package's bson.D, which
// the driver cannot encode).
func TestFindIdOnDriverPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetRegistry(mgocompat.Registry))
	if err != nil || client.Ping(ctx, nil) != nil {
		t.Skip("no MongoDB on 127.0.0.1:27017")
	}
	defer client.Disconnect(context.Background())
	const dbName = "mgo_driverpath_test"
	defer client.Database(dbName).Drop(context.Background())
	withDriver(t)

	s := NewDummySession(dbName)
	s.SetDriverDatabase(dbName)
	s.SetMongoDriverClient(client)
	c := s.DB(dbName).C("c")
	id := bson.NewObjectId()
	if err := c.Insert(bson.M{"_id": id, "name": "x"}); err != nil {
		t.Fatal(err)
	}
	var got bson.M
	if err := c.FindId(id).One(&got); err != nil || got["name"] != "x" {
		t.Fatalf("FindId: %v %v", got, err)
	}
}
