package mongo

import (
	"io/ioutil"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/dbtest"
	"github.com/pkg/errors"
)

// Storage is a mongo storage backend for all DAO
type Storage struct {
	// Session holds the root session connecting to the mongo server
	Session *mgo.Session

	// all collections are int the same db instance
	DB *mgo.Database

	Config *Config
}

type Config struct {
	Host        string `json:"host,omitempty"`
	DBName      string `json:"db_name,omitempty"`
	MaxPoolSize int    `json:"max_pool_size,omitempty"`
}

// New Storage instance and a done channel which will be closed after all the sessions closed
func New(cfg *Config) (*Storage, func(), error) {
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
		Config:  cfg,
	}

	return s, ss.Close, nil
}

func NewTestStorage() (*Storage, func(), error) {
	var server dbtest.DBServer
	const testDBName = "test_db"
	tempDir, e := ioutil.TempDir("", "mgo_test")
	if e != nil {
		return nil, nil, e
	}
	server.SetPath(tempDir)

	ss := server.Session()
	hosts := ss.LiveServers()
	if len(hosts) == 0 {
		return nil, nil, errors.New("no living servers")
	}

	db := ss.DB(testDBName)

	closer := func() {
		ss.Close()
		server.Stop()
		server.Wipe()
		os.Remove(tempDir)
	}

	return &Storage{
		Session: ss,
		DB:      db,
		Config: &Config{
			DBName: testDBName,
			Host:   hosts[0],
		},
	}, closer, nil
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
