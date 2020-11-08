//
// Copyright (c) 2016-2020 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var collectorURI string

// RootCmd is the base command for this application.
var RootCmd = &cobra.Command{
	Use:   "stress",
	Short: "Perform stress tests with various tracker implementations",
	Long: `This app allows you to stress test the Snowplow Golang Tracker
in various ways by trying different combinations of emitter types
and settings.  This helps to explore and benchmark the most efficient
way to track Snowplow events with this tracker for your implementation.

Check the 'help' command to see what tests are available!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

// Execute performs the base command execution before any sub-commands are invoked.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&collectorURI, "collector", "c", "localhost:50510", "Collector URI to emit events too")
}
