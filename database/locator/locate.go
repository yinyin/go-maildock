package locator

import (
	"fmt"
	"github.com/yinyin/go-maildock/database"
	"github.com/yinyin/go-maildock/database/mysqldb"
)

func NewEmptyConfiguration(databaseType string) (cfg database.Configuration, err error) {
	cfg = nil
	switch databaseType {
	case "mysql":
		cfg = mysqldb.NewEmptyConfiguration()
	}
	if nil == cfg {
		err = fmt.Errorf("unknown database type: %v", databaseType)
		return nil, err
	}
	return cfg, nil
}
