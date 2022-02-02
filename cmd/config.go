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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config key [value]",
	Short: "Configure ocelot behavior",
	Long:  `Configure ocelot behavior.`,
	Args:  cobra.RangeArgs(1, 2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// key
		if len(args) == 0 {
			ret := []string{}
			for _, key := range viper.AllKeys() {
				if strings.Contains(key, toComplete) {
					ret = append(ret, key)
				}
			}
			return ret, cobra.ShellCompDirectiveNoFileComp
		}
		// value
		switch viper.Get(args[0]).(type) {
		default:
			return nil, cobra.ShellCompDirectiveNoFileComp
		case bool:
			// only complete bools
			return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if len(args) == 1 {
			value := viper.Get(key)
			if value == nil {
				// not a key with a value set
				return fmt.Errorf("Unknown config key \"%s\"", key)
			}
			fmt.Println(value)
		} else if len(args) == 2 {
			value := args[1]
			// don't change the type from the default
			switch viper.Get(key).(type) {
			default:
				// trying to set a group value
				return fmt.Errorf("Invalid config key \"%s\"", key)
			case nil:
				// not a key with a value set
				return fmt.Errorf("Unknown config key \"%s\"", key)
			case string:
				viper.Set(key, value)
				break
			case bool:
				boolValue, err := cast.ToBoolE(value)
				if err != nil {
					return fmt.Errorf("Invalid boolean value \"%s\" for key \"%s\"", value, key)
				}
				viper.Set(key, boolValue)
				break
			case int64:
				intValue, err := cast.ToInt64E(value)
				if err != nil {
					return fmt.Errorf("Invalid integer value \"%s\" for key \"%s\"", value, key)
				}
				viper.Set(key, intValue)
				break
			case float64:
				floatValue, err := cast.ToFloat64E(value)
				if err != nil {
					return fmt.Errorf("Invalid float value \"%s\" for key \"%s\"", value, key)
				}
				viper.Set(key, floatValue)
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
