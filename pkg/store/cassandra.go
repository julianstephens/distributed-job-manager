package store

import (
	"fmt"
	"sync"

	"github.com/gocql/gocql"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/scylladb/gocqlx/v3"
)

var (
	Session gocqlx.Session
	once    sync.Once
	conf    = config.GetConfig()
)

type DBSession struct {
	Client *gocqlx.Session
}

func GetDB(ks string) (*DBSession, error) {
	var err error
	once.Do(func() {
		err = initStore(ks)
	})
	if err != nil {
		return nil, err
	}
	return &DBSession{
		Client: &Session,
	}, nil
}

func initStore(ks string) error {
	var err error
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", conf.Cassandra.Host, conf.Cassandra.Port))
	cluster.Keyspace = ks

	if Session, err = gocqlx.WrapSession(cluster.CreateSession()); err != nil {
		return err
	}
	return nil
}
