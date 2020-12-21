package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/keitaro1020/lambda-golang-slf-practice/service/domain"
	"github.com/keitaro1020/lambda-golang-slf-practice/service/infra/db/models"
)

type Config struct {
	User     string
	Pass     string
	Endpoint string
	Name     string
}

type txConn struct {
	*sql.Tx
}

func connectDB(config *Config) (*sql.DB, error) {
	connectStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s",
		config.User,
		config.Pass,
		config.Endpoint,
		"3306",
		config.Name,
		"utf8",
	)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return nil, err
	}

	// set debug
	models.SetDebug(log.StandardLogger().Out)

	return db, nil
}

func NewTransaction(config *Config) func(ctx context.Context, txFunc func(ctx context.Context, tx domain.Tx) error) (err error) {
	return func(ctx context.Context, txFunc func(ctx context.Context, tx domain.Tx) error) (err error) {
		db, err := connectDB(config)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
			if err != nil {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()

		err = txFunc(ctx, &txConn{tx})
		return err
	}
}

func (tx *txConn) Executor() interface{} {
	return tx.Tx
}

func sqlTx(tx domain.Tx) (*sql.Tx, error) {
	txi := tx.Executor()

	sqlTx, ok := txi.(*sql.Tx)
	if !ok {
		return nil, errors.New("invalid connection")
	}
	return sqlTx, nil
}
