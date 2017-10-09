package mysqldb

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"log"
	"net/mail"
	"time"
	"context"
	"strings"
	maildock "github.com/yinyin/go-maildock"
)

// Convert string into Null-able string object. Empty string will be treated as NULL.
func toNullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: (v != "")}
}

// Enables parse time and set connection time zone to UTC
func enhanceDSN(dsnValue string) string {
	return dsnValue + "?parseTime=true&loc=UTC&time_zone=%22%2B00%3A00%22"
}

type mysqlConnection struct {
	conn *sql.DB
}

func newMySQLConnectionWithDSN(dsnValue string) (conn * mysqlConnection, err error) {
	dbconn, err := sql.Open("mysql", enhanceDSN(dsnValue))
	if nil != err {
		return nil, err
	}
	conn = &mysqlConnection{
		conn: dbconn,
	}
	return conn, nil
}

func newMySQLConnection(dsnValues []string) (conn * mysqlConnection, err error) {
	for _, dsnValue := range dsnValues {
		if conn, err = newMySQLConnectionWithDSN(dsnValue); nil == err {
			return conn, nil
		}
	}
	if nil == err {
		err = errors.New("None of DSN are able to connect database")
	}
	return nil, err
}

func (c * mysqlConnection) AppendMail(ctx context.Context, fromAddress *mail.Address, toAddresses []*mail.Address, mailBody string) (err error) {
	var tx *sql.Tx
	if tx, err = c.conn.Begin(); nil != err {
		return err
	}
	defer tx.Rollback()
	currentEpoch := time.Now().Unix()
	res, err := tx.ExecContext(ctx, "INSERT INTO `MailContent`(`from_address`, `from_name`, `mail_body`, `create_at`) VALUES(?, ?, ?, ?)",
			strings.ToLower(fromAddress.Address), toNullString(fromAddress.Name), mailBody, currentEpoch)
	if err != nil {
		log.Printf("failed on insert mail record into database: %v", err)
		return err
	}
	mailSn, err := res.LastInsertId()
	if nil != err {
		log.Printf("failed on getting mail-sn from last insert id: %v", err)
		return err
	}
	for _, toAddr := range  toAddresses {
		if _, err = tx.ExecContext(ctx, "INSERT INTO `MailRecipient`(`mail_sn`, `to_address`, `to_name`, `create_at`) VALUES(?, ?, ?, ?)",
			mailSn, strings.ToLower(toAddr.Address), toNullString(toAddr.Name), currentEpoch); nil != err {
			log.Printf("failed on saving recipient record: %v", err)
		}
	}
	tx.Commit()
	return nil
}

func (c * mysqlConnection) PurgeMail(ctx context.Context, d time.Duration) (err error) {
	boundEpoch := time.Now().Add(-d).Unix()
	_, err = c.conn.ExecContext(ctx, "DELETE FROM `MailContent` WHERE (`create_at` < ?)", boundEpoch)
	return err
}

func (c * mysqlConnection) SearchForRecipient(ctx context.Context, recipientAddress string) (mailRecords []*maildock.MailRecord, err error) {
	rows, err := c.conn.QueryContext(ctx, "SELECT C.`from_address`, C.`from_name`, C.`mail_body`, C.`create_at` FROM `MailContent` AS C LEFT JOIN `MailRecipient` AS R ON C.`mail_sn` = R.`mail_sn` WHERE (R.`to_address` = ?) ORDER BY C.`create_at` DESC LIMIT 10",
		strings.ToLower(recipientAddress))
	if nil != err {
		return nil, err
	}
	defer rows.Close()
	mailRecords = make([]*maildock.MailRecord, 0)
	for rows.Next() {
		var fromAddress, mailBody string
		var fromName sql.NullString
		var createAtEpoch int64
		if err = rows.Scan(&fromAddress, &fromName, &mailBody, &createAtEpoch); nil != err {
			log.Printf("failed on fetching mail record: %v", err)
			return nil, err
		}
		mailRecords = append(mailRecords, maildock.NewMailRecord(fromAddress, fromName.String, mailBody, createAtEpoch))
	}
	return mailRecords, nil
}

func (c * mysqlConnection) Close() (err error) {
	return c.conn.Close()
}
