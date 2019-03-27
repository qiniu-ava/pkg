package mongo

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeStorage struct {
	Storage
}

func (s *fakeStorage) fakeCollection() (*fakeCollection, error) {
	c := Collection{s.DB.C("fake")}
	ss := c.CopySession()
	defer ss.CloseSession()

	if e := ss.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true}); e != nil {
		return nil, errors.Wrap(e, "init fake collection with index failed")
	}

	return &fakeCollection{c}, nil
}

type fakeCollection struct {
	Collection
}

func (c *fakeCollection) Set(t *fakeType) error {
	ss := c.CopySession()
	defer ss.CloseSession()

	return ss.Insert(t)
}

func (c *fakeCollection) Get(id string) (*fakeType, error) {
	ss := c.CopySession()
	defer ss.CloseSession()

	out := &fakeType{}
	if e := ss.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(out); e != nil {
		return nil, e
	}
	return out, nil
}

type fakeType struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

func TestFakeCollection(t *testing.T) {
	s, closer, e := NewTestStorage()
	require.NoError(t, e)
	defer closer()

	ms := &fakeStorage{*s}
	coll, e := ms.fakeCollection()
	require.NoError(t, e)

	ft := &fakeType{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	sameName := &fakeType{
		ID:   bson.NewObjectId(),
		Name: "test",
	}
	{
		_, e := coll.Get(ft.ID.Hex())
		assert.Error(t, e)
	}
	{
		e := coll.Set(ft)
		assert.NoError(t, e)
	}
	{
		out, e := coll.Get(ft.ID.Hex())
		assert.NoError(t, e)
		assert.True(t, reflect.DeepEqual(ft, out))
	}
	{
		e := coll.Set(sameName)
		assert.Error(t, e)
	}
}
