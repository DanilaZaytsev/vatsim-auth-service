package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	yc "github.com/ydb-platform/ydb-go-yc"
)

var (
	DB    *ydb.Driver
	Table table.Client
)

func InitYDB(ctx context.Context) error {
	dsn := os.Getenv("YDB_DSN")
	saKeyFile := os.Getenv("YDB_SA_KEY")

	driver, err := ydb.Open(ctx,
		dsn,
		yc.WithServiceAccountKeyFileCredentials(saKeyFile),
		yc.WithInternalCA(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to YDB: %w", err)
	}

	DB = driver
	Table = driver.Table()
	log.Println("✅ Connected to YDB")
	return nil
}
