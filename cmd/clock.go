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
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/starlight/ocelot/pkg/ocelot"
)

var Ocelot *ocelot.Ocelot

// clockCmd represents the clock command
var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Show the current Market Time, according to the API",
	Long: `Show the current Market Time, according to the API.
Includes round-trip time, one-way delay and response lag.`,
	Run: func(cmd *cobra.Command, args []string) {
		clock, err := Ocelot.GetClock()
		if err != nil {
			log.Printf("error: %s", err)
		}
		fmt.Printf("Market Time: %s\n", clock.Market.Timestamp.Round(time.Second))
		if clock.Market.IsOpen {
			fmt.Printf("OPEN until %s\n", clock.Market.NextClose)
		} else {
			fmt.Printf("CLOSED until %s\n", clock.Market.NextOpen)
		}
		fmt.Printf("RTT: %s\n", clock.RTT)
		fmt.Printf("OWD: %s\n", clock.OWD)
		fmt.Printf("LAG: %s\n", clock.LAG)
	},
}

func init() {
	Ocelot = ocelot.GetOcelot()
	Ocelot.GetClock() // exercise api
	showCmd.AddCommand(clockCmd)
}
