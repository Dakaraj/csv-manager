// Copyright © 2018 Anton Kramarev <kramarev.anton@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	amount        uint
	delimiter     string
	runeDelimiter rune
	header        bool
	fileFolder    string
	fileName      string
	backup        bool
)

// VERSION represents a current version of application
const VERSION = "0.0.2"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "csv-manager",
	Version: VERSION,
	Short:   "Manage your csv files in different ways",
	Long: `Divide CSV files into equal parts, randomize your file,
perform other operations with CSV files`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
