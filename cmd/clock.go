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
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clockCmd represents the clock command
var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Show the current market time, according to the API",
	Long: `Show the current market time, according to the API.
Includes round-trip time, one-way delay and response time-to-live.`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// first request to init api
		maxTTL := viper.GetDuration("clock.max-ttl")
		_, err := Ocelot.GetClock(maxTTL)
		cobra.CheckErr(err)
		// second request measurements used
		clock, err := Ocelot.GetClock(maxTTL)
		cobra.CheckErr(err)
		fmt.Printf("Market Time: %s (%s TTL)\n", clock.Market.Timestamp.Round(time.Second), clock.TTL)
		if clock.Market.IsOpen {
			duration := time.Until(clock.Market.NextClose).Round(time.Second)
			fmt.Printf("Market %s until %s (%s)\n", color.GreenString("OPEN"), clock.Market.NextClose, duration)
		} else {
			duration := time.Until(clock.Market.NextOpen).Round(time.Second)
			fmt.Printf("Market %s until %s (%s)\n", color.HiYellowString("CLOSED"), clock.Market.NextOpen, duration)
		}
	},
}

func init() {
	clockCmd.Flags().DurationP("max-ttl", "m", 0, "maximum market time-to-live delay (default 0s for no maximum)")
	viper.BindPFlag("clock.max-ttl", clockCmd.Flags().Lookup("max-ttl"))
	showCmd.AddCommand(clockCmd)
}
