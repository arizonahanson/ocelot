/*
Copyright © 2022 Arizona Hanson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/starlight/ocelot/pkg/ocelot"
)

var (
	version string
	Ocelot  *ocelot.Ocelot
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ocelot",
	Short:   "Brokerage CLI",
	Long:    `A command-line trading platform written in Go`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Ocelot = ocelot.GetOcelot()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	// Search config in home directory with filename ".ocelot.toml".
	viper.AddConfigPath(home)
	viper.SetConfigName(".ocelot")
	viper.SetConfigType("toml")
	// read in environment variables that match.
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	viper.ReadInConfig()
}

func Quit() {
	fmt.Fprintln(os.Stderr, "exiting")
	os.Exit(0)
}
