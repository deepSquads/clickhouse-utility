package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Config struct {
	Host     string `env:"CLICKHOUSE_HOST,required"`
	Database string `env:"CLICKHOUSE_DB,required"`
	User     string `env:"CLICKHOUSE_USER,required"`
	Password string `env:"CLICKHOUSE_PASSWORD,required"`
}

func NewClickHouseClient(ctx context.Context, cfg Config) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Host},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("cannot connect to clickhouse [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func tableExists(ctx context.Context, client clickhouse.Conn, table string) (bool, error) {
	res, err := client.Query(ctx, fmt.Sprintf("exists table %s", table))
	if err != nil {
		return false, fmt.Errorf("client.Query: %w", err)
	}
	res.Next()
	var exists uint8
	err = res.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("res.Scan: %w", err)
	}
	return exists == 1, nil
}
