package main

import (
	"log"
	"context"
	"github.com/yinyin/go-maildock/cmd"
	"github.com/yinyin/go-maildock/database"
	"net"
	"time"
)

type smtpRunner struct {
	ctx context.Context
	dbcfg database.Configuration
	connections <-chan net.Conn
	listeners []net.Listener
}

func (u *smtpRunner) processConnection(netConn net.Conn) {
	defer netConn.Close()
	dbConn, err := u.dbcfg.OpenConnection()
	if nil != err {
		log.Printf("failed on connecting to database: %v", err)
		return
	}
	defer dbConn.Close()
	ctxProcess, _ := context.WithTimeout(u.ctx, time.Minute*3)
	processSMTPConnection(ctxProcess, netConn, dbConn)
}

func (u *smtpRunner) Run() {
	defer u.closeListeners()
	for {
		select {
		case <-u.ctx.Done():
			log.Print("leaving connection handling loop")
			return
		case conn := <- u.connections:
			if nil == conn {
				continue
			}
			go u.processConnection(conn)
		}
	}
}

func (u *smtpRunner) closeListeners() {
	for _, l := range u.listeners {
		l.Close()
	}
}

func setupRunnerWithConfiguration(ctx context.Context) (runner *smtpRunner, err error) {
	cfg, err := cmd.LoadConfigurationWithFlags()
	if nil != err {
		log.Printf("cannot load configuration: %v", err)
		return nil, err
	}
	connections, listeners, err := startListeners(cfg.SMTPListenOn)
	if nil != err {
		log.Printf("cannot start listeners: %v", err)
		return nil,err
	}
	runner = &smtpRunner{
		ctx: ctx,
		dbcfg: cfg.Database.Config,
		connections: connections,
		listeners: listeners,
	}
	return runner, nil
}
