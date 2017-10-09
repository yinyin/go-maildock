package database

import (
	"net/mail"
	"context"
	maildock "github.com/yinyin/go-maildock"
	"time"
)

type Connection interface {
	AppendMail(ctx context.Context, fromAddress *mail.Address, toAddresses []*mail.Address, mailBody string) (err error)
	PurgeMail(ctx context.Context, d time.Duration) (err error)
	SearchForRecipient(ctx context.Context, recipientAddress string) (mailRecords []*maildock.MailRecord, err error)
	Close() (error)
}

type Configuration interface {
	OpenConnection() (Connection, error)
}