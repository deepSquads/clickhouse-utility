package rmv

import (
	"context"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/deepsquads/clickhouse-utility/internal"
)

var requiredFlags = []string{"database", "table", "definition", "select"}

func GetCmd() *cobra.Command {
	rmv := internal.Rmv{}

	cmd := &cobra.Command{
		Use:   "rmv",
		Short: "Create a refreshable materialized view",
		Run: func(cmd *cobra.Command, args []string) {

			var cfg internal.Config
			err := env.Parse(&cfg)
			if err != nil {
				log.Fatal().Err(err).Msg("env.Parse()")
				return
			}

			ch, err := internal.NewClickHouseClient(context.TODO(), cfg)
			if err != nil {
				log.Fatal().Err(err).Msg("clickhouse.NewClickHouseClient()")
				return
			}

			err = internal.CreateRefreshableMaterializedView(context.TODO(), ch, rmv)
			if err != nil {
				log.Fatal().Err(err).Msg("CreateRefreshableMaterializedView()")
				return
			}
		},
	}

	cmd.Flags().StringVarP(&rmv.Database, "database", "", "", "Database name")
	cmd.Flags().StringVarP(&rmv.TableName, "table", "t", "", "Table name of the view")
	cmd.Flags().StringVarP(&rmv.TableDefinition, "definition", "d", "", "Schema and engine of the table. This part will be concatenated into a create table statement")
	cmd.Flags().StringVarP(&rmv.SelectQuery, "select", "s", "", "Select query to populate the table")

	for _, flag := range requiredFlags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			log.Fatal().Err(err).Msg("cmd.MarkFlagRequired()")
		}

	}

	return cmd
}
