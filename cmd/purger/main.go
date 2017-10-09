package main

import (
	"context"
	"github.com/yinyin/go-maildock/cmd"
	"log"
	"time"
)

func main() {
	cfg, err := cmd.LoadConfigurationWithFlags()
	if nil != err {
		log.Printf("cannot load configuration: %v", err)
		return
	}
	if cfg.PurgeDays < 1 {
		log.Printf("configurated purge-days less than 1: %v", cfg.PurgeDays)
		return
	}
	purgeDuration := time.Duration(cfg.PurgeDays*24) * time.Hour
	dbconn, err := cfg.Database.Config.OpenConnection()
	if nil != err {
		log.Printf("failed on connection to database: %v", err)
		return
	}
	defer dbconn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()
	if err = dbconn.PurgeMail(ctx, purgeDuration); nil != err {
		log.Printf("failed on purging mail: %v", err)
		return
	}
}
