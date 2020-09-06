package isucongo

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soichisumi/go-util/logger"
	"go.uber.org/zap"
	"time"
)

func Connect(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectWithTimeout(driverName, dsn string, timeout time.Duration) (*sql.DB, error) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	tTimeOut := time.After(timeout)
	for {
		select {
		case <-tTimeOut:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-t.C:
			db, err := Connect(driverName, dsn)
			if err == nil {
				return db, nil
			}
			logger.Info("failed to connect database. retrying...", zap.Error(errors.Wrapf(err, "failed to connect to db %s", dsn)))
		}
	}
}