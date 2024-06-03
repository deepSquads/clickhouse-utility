package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var ch driver.Conn
var rmv = Rmv{
	Database:        "db",
	TableName:       "rmv",
	TableDefinition: "(timestamp DateTime, value UInt32) engine=Memory",
	SelectQuery:     "select max(timestamp) timestamp, sum(value) value from raw",
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("env.Parse()")
		return
	}

	ch, err = NewClickHouseClient(ctx, cfg)
	if err != nil {
		log.Fatal().Msgf("NewClickHouseClient: %s", err)
	}

	err = ch.Exec(ctx, `
create or replace table raw
(
	timestamp DateTime,
	value UInt32
)
engine=Memory;
	`)
	if err != nil {
		log.Fatal().Msgf("ch.Exec: %s", err)
	}

	code := m.Run()
	os.Exit(code)
}

func validateTables(t *testing.T) {
	rows, err := ch.Query(context.Background(), fmt.Sprintf("show tables from %s", rmv.Database))
	assert.NoError(t, err)
	prefix := fmt.Sprintf("%s_", rmv.TableName)
	var table string
	for rows.Next() {
		assert.NoError(t, rows.Scan(&table))
		assert.Equal(t, false, strings.HasPrefix(table, prefix))
	}
}

func testRmvValues(t *testing.T, expectTimestamp time.Time, expectedValue uint32) {
	rows, err := ch.Query(context.Background(), fmt.Sprintf("select timestamp, value from %s.%s", rmv.Database, rmv.TableName))
	assert.NoError(t, err)
	var timestamp time.Time
	var value uint32
	for rows.Next() {
		assert.NoError(t, rows.Scan(&timestamp, &value))
	}
	assert.Equal(t, expectTimestamp, timestamp)
	assert.Equal(t, expectedValue, value)

}

func TestCreateRefreshableMaterializedView(t *testing.T) {
	assert.NoError(t, ch.Exec(context.Background(), `
	insert into raw values
		('2024-05-22 00:00:00', 5),
		('2024-05-22 00:01:00', 2),
		('2024-05-22 00:03:00', 8),
	;`))
	assert.NoError(t, CreateRefreshableMaterializedView(context.Background(), ch, rmv))
	validateTables(t)
	testRmvValues(t, time.Time(time.Date(2024, time.May, 22, 0, 3, 0, 0, time.UTC)), uint32(15))
	assert.NoError(t, ch.Exec(context.Background(), `
	insert into raw values
		('2024-05-22 00:05:00', 8),
		('2024-05-22 00:08:00', 9),
	;`))
	assert.NoError(t, CreateRefreshableMaterializedView(context.Background(), ch, rmv))
	validateTables(t)
	testRmvValues(t, time.Time(time.Date(2024, time.May, 22, 0, 8, 0, 0, time.UTC)), uint32(32))
}
