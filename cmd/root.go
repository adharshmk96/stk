/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var SemVer = "v0.0.0"

func GetSemverInfo() string {
	if SemVer != "v0.0.0" {
		return SemVer
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		return version.Main.Version
	}
	return SemVer
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stk",
	Short: "STK is a CLI tool and library for building and managing server applications",
	Long: `STK provides various tools for building and managing server applications.
You can use the gsk (go server kit) package for building rest api servers. 
stk CLI has tools to do project and module generation, sql migration management.
Refer the documentation to see the usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			fmt.Println(GetSemverInfo())
			return
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".stk.yaml", "config file.")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "display stk version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	// Set the key replacer
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
