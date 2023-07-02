package tpl

func MainTemplate() []byte {
	return []byte(`package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}

func CmdRootTemplate() []byte {
	return []byte(`package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "v0.0.0"
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "{{ .PkgName }}",
	Long:  ` + "`{{ .PkgName }}`" + `,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "{{.AppName}}.yaml", "config file.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	// Set the key replacer for env variables.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
`)
}

func CmdServeTemplate() []byte {
	return []byte(`package cmd

import (
	"github.com/spf13/cobra"
	"{{ .PkgName }}/server"
)

var startingPort string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		startAddr := "0.0.0.0:"
		server.StartServer(startAddr + startingPort)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&startingPort, "port", "p", "8080", "Port to start the server on")

	rootCmd.AddCommand(serveCmd)
}
`)
}
func CmdVersionTemplate() []byte {
	return []byte(`package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version of {{ .AppName }}",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("{{ .AppName }} version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`)
}
