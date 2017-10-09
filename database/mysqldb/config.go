package mysqldb

import (
	"github.com/yinyin/go-maildock/database"
	"math/rand"
)

type mysqlConfiguration struct {
	DSN []string	`yaml:"dsn"`
}

func NewEmptyConfiguration() (r database.Configuration) {
	return &mysqlConfiguration {}
}

func (c * mysqlConfiguration) OpenConnection() (conn database.Connection, err error) {
	l := len(c.DSN)
	shuffleDSN := make([]string, l)
	copy(shuffleDSN, c.DSN)
	for idx := 1; idx < l; idx++ {
		t := rand.Intn(l)
		if t != idx {
			shuffleDSN[t], shuffleDSN[idx] = shuffleDSN[idx], shuffleDSN[t]
		}
	}
	return newMySQLConnection(shuffleDSN)
}