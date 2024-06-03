package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/deepsquads/clickhouse-utility/cmd/rmv"
)

var (
	rootCmd = &cobra.Command{
		Use:   "main",
		Short: "ClickHouse utility",
	}
)

func main() {
	// Workaround for removing the limits from k8s
	runtime.GOMAXPROCS(1)

	rootCmd.AddCommand(rmv.GetCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s\n", err)
	}
}
