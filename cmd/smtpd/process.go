package main

import (
	"net"
	"context"
	"github.com/siebenmann/smtpd"
	"os"
	"log"
	"net/mail"
	"github.com/yinyin/go-maildock/database"
)

var smtpConfig smtpd.Config = smtpd.Config{
	SftName: "go-maildock",
}

func init() {
	if hostname, err := os.Hostname(); nil != err {
		log.Printf("cannot have local host name: %v", err)
	} else {
		smtpConfig.LocalName = hostname
	}

}

type smtpProcessor struct {
	conn * smtpd.Conn
	heloName string
	fromAddress *mail.Address
	toAddresses []*mail.Address
}

func newSMTPProcessor(conn net.Conn) (p * smtpProcessor) {
	smtpConn := smtpd.NewConn(conn, smtpConfig, nil)
	return &smtpProcessor{
		conn: smtpConn,
		toAddresses: make([]*mail.Address, 0),
	}
}

func (p * smtpProcessor) handleCommand(evt * smtpd.EventInfo) {
	switch evt.Cmd {
	case smtpd.EHLO, smtpd.HELO:
		p.heloName = evt.Arg
		p.toAddresses = make([]*mail.Address, 0)
	case smtpd.MAILFROM:
		if addr, err := mail.ParseAddress(evt.Arg); nil != err {
			log.Printf("failed on parsing MAILFROM %v: %v", evt.Arg, err)
			p.conn.Reject()
		} else {
			if nil != p.fromAddress {
				log.Printf("replacing mail from %v with %v", p.fromAddress, evt.Arg)
			}
			p.fromAddress = addr
		}
	case smtpd.RCPTTO:
		if addr, err := mail.ParseAddress(evt.Arg); nil != err {
			log.Printf("failed on parsing RCPTTO %v: %v", evt.Arg, err)
			p.conn.Reject()
		} else {
			p.toAddresses = append(p.toAddresses, addr)
		}
	case smtpd.DATA:
		if 0 == len(p.toAddresses) {
			p.conn.RejectMsg("delivery target is required")
		}
	}
}

func (p * smtpProcessor) handleData(ctx context.Context, evt * smtpd.EventInfo, dbconn database.Connection) {
	data := evt.Arg
	if err := dbconn.AppendMail(ctx, p.fromAddress, p.toAddresses, data); nil != err {
		log.Printf("failed on put mail record into database: %v", err)
		p.conn.Reject()
	}
}


func (p * smtpProcessor) Process(ctx context.Context, dbconn database.Connection) {
	for {
		evt := p.conn.Next()
		switch evt.What {
		case smtpd.COMMAND:
			p.handleCommand(&evt)
		case smtpd.GOTDATA:
			p.handleData(ctx, &evt, dbconn)
		case smtpd.DONE, smtpd.ABORT:
			break
		}
	}
}

func processSMTPConnection(ctx context.Context, conn net.Conn, dbconn database.Connection) {
	processor := newSMTPProcessor(conn)
	processor.Process(ctx, dbconn)
}
