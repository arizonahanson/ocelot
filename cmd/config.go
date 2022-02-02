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

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [key [value]]",
	Short: "Configure ocelot behavior.",
	Long:  `Configure ocelot behavior.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if !viper.IsSet(key) {
			// key not bound by viper
			return fmt.Errorf("Unknown config key")
		}
		if len(args) == 1 {
			value := viper.Get(key)
			fmt.Println(value)
		} else if len(args) == 2 {
			value := args[1]
			// don't change the type from the default
			switch viper.Get(key).(type) {
			default:
				// trying to set a group value
				return fmt.Errorf("Unknown config type")
			case string:
				viper.Set(key, value)
				break
			case bool:
				viper.Set(key, cast.ToBool(value))
				break
			}
			viper.WriteConfig()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
