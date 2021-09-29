package gorpUtil

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zkcrescent/chaos/utils"
	gorp "gopkg.in/gorp.v2"
)

func MigrateQuery(fs ...TableCheck) *utils.CommandRegister {
	var (
		dsn        string
		skipExists bool
	)

	var migrateQueryCmd = &cobra.Command{
		Use:   "migrate-query",
		Short: "cli command that export sql queries for creating tables",
		Long:  "cli command that export sql queries for creating tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			if skipExists && dsn == "" {
				return utils.Error("skip-exists is true while dsn is empty.")
			}
			db, err := sql.Open("mysql", dsn)
			if err != nil {
				return err
			}

			dm := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
			queries, err := Tables.CreateTableQueries(dm, skipExists, fs...)
			if err != nil {
				return err
			}
			fmt.Println("Queries: \n", strings.Join(queries, "\n"))
			return nil
		},
	}

	return &utils.CommandRegister{
		Command: migrateQueryCmd,
		ParseFlag: func() error {
			migrateQueryCmd.Flags().StringVarP(&dsn, "dsn", "D", "", "Your mysql connection string.")
			migrateQueryCmd.Flags().BoolVarP(&skipExists, "skip-exists", "S", false, "Skip queries if tables already exist, require dsn.")
			return migrateQueryCmd.Flags().Parse(os.Args)
		},
	}
}
