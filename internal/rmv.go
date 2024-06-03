package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Rmv struct {
	Database        string
	TableName       string
	TableDefinition string
	SelectQuery     string
}

func CreateRefreshableMaterializedView(ctx context.Context, client clickhouse.Conn, rmv Rmv) error {
	tempTable := fmt.Sprintf("%s.%s_%d", rmv.Database, rmv.TableName, time.Now().Unix())
	err := client.Exec(ctx, fmt.Sprintf("create or replace table %s %s", tempTable, rmv.TableDefinition))
	if err != nil {
		return fmt.Errorf("create new table: %w", err)
	}
	err = client.Exec(ctx, fmt.Sprintf("insert into %s %s", tempTable, rmv.SelectQuery))
	if err != nil {
		return fmt.Errorf("insert from select: %w", err)
	}
	targetTable := fmt.Sprintf("%s.%s", rmv.Database, rmv.TableName)
	exists, err := tableExists(ctx, client, targetTable)
	if err != nil {
		return fmt.Errorf("table exists: %w", err)
	}
	if exists {
		err = client.Exec(ctx, fmt.Sprintf("exchange tables %s and %s", targetTable, tempTable))
		if err != nil {
			return fmt.Errorf("exchange tables: %w", err)
		}
		err = client.Exec(ctx, fmt.Sprintf("drop table %s", tempTable))
		if err != nil {
			return fmt.Errorf("drop temp table: %w", err)
		}
	} else {
		err = client.Exec(ctx, fmt.Sprintf("rename table %s to %s", tempTable, targetTable))
		if err != nil {
			return fmt.Errorf("rename table: %w", err)
		}
	}
	return nil
}
