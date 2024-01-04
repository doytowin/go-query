package goooqo

import (
	"database/sql"
	"github.com/doytowin/goooqo/rdb"
	. "github.com/doytowin/goooqo/test"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Export Interface", func(t *testing.T) {
		var db *sql.DB
		tm := rdb.NewTransactionManager(db)
		dataAccess := rdb.NewTxDataAccess[UserEntity](tm, func() UserEntity { return UserEntity{} })

		BuildRestService[UserEntity, UserQuery](
			"/user/", dataAccess,
			func() UserEntity { return UserEntity{} },
			func() UserQuery { return UserQuery{} },
		)
	})
}