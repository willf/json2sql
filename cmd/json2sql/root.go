package json2sql

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willf/json2sql/pkg/json2sql"
)

var table string
var create bool

var rootCmd = &cobra.Command{
	Use:   "json2sql",
	Short: "json2sql is a tool to convert json to sql",
	Long:  `json2sql is a tool to convert json to sql insert statements`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("json2sql is a tool to convert json to sql")
		json2sql.MainLoop(table, create)
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&table, "table", "t", "", "table name to use")
	rootCmd.Flags().BoolVarP(&create, "create", "c", false, "Include CREATE TABLE statement")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
