package mongo

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

// Storage is a mongo storage backend for all DAO
type Storage struct {
	// Session holds the root session connecting to the mongo server
	Session *mgo.Session

	// all collections are int the same db instance
	DB *mgo.Database
}

type Config struct {
	Host        string
	DBName      string
	MaxPoolSize int
}

// New Storage instance and a done channel which will be closed after all the sessions closed
func New(ctx context.Context, cfg *Config) (*Storage, <-chan struct{}, error) {
	if cfg.Host == "" {
		return nil, nil, errors.New("empty mongodb host")
	}
	if cfg.DBName == "" {
		return nil, nil, errors.New("empty mongodb name")
	}
	if cfg.MaxPoolSize == 0 {
		cfg.MaxPoolSize = 200
	}

	ss, e := mgo.Dial(cfg.Host)
	if e != nil {
		return nil, nil, errors.Wrap(e, "create mongo session failed")
	}
	db := ss.DB(cfg.DBName)

	if cfg.MaxPoolSize > 0 {
		ss.SetPoolLimit(cfg.MaxPoolSize)
	}

	s := &Storage{
		Session: ss,
		DB:      db,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		ss.Close()
		close(done)
	}()

	return s, done, nil
}

// Collection adds session supports to the mgo.Collection
type Collection struct {
	*mgo.Collection
}

func (c *Collection) CopySession() *Collection {
	db := c.Database
	return &Collection{db.Session.Copy().DB(db.Name).C(c.Name)}
}

func (c *Collection) CloseSession() {
	c.Database.Session.Close()
}
