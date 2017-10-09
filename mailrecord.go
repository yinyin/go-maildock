package maildock

import (
	"time"
	"encoding/json"
)

type MailRecord struct {
	FromAddress string	`json:"from_address"`
	FromName string	`json:"from_name"`
	MailBody string	`json:"mail_body"`
	CreateAt time.Time `json:"create_at"`
}

func (r *MailRecord) MarshalJSON() ([]byte, error) {
	var d = struct {
		FromAddress string	`json:"from_address"`
		FromName *string	`json:"from_name"`
		MailBody string	`json:"mail_body"`
		CreateAt int64`json:"create_at"`
	} {
		FromAddress: r.FromAddress,
		MailBody: r.MailBody,
		CreateAt: r.CreateAt.Unix(),
	}
	if "" != r.FromName {
		d.FromName = &r.FromName
	}
	return json.Marshal(d)
}

func NewMailRecord(fromAddress, fromName, mailBody string, createAtEpoch int64) (r *MailRecord) {
	c := time.Unix(createAtEpoch, 0)
	return &MailRecord{
		FromAddress: fromAddress,
		FromName:fromName,
		MailBody:mailBody,
		CreateAt: c,
	}
}