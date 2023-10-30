/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mkconfigCmd represents the mkconfig command
var mkconfigCmd = &cobra.Command{
	Use:   "mkconfig",
	Short: "Generate a config file for stk",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := cmd.Flag("config").Value.String()

		viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

		// Set your settings
		viper.Set("migrator.workdir", "./stk-migrations")
		viper.Set("migrator.database", "sqlite3")
		viper.Set("migrator.sqlite.filepath", "stk.db")

		// Create the config file
		err := viper.WriteConfigAs(configFile)
		if err != nil {
			fmt.Printf("Error while writing config file: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mkconfigCmd)
}
