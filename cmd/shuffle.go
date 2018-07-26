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
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func shuffle(cmd *cobra.Command, args []string) {
	var fileLength int
	var headerLine []string

	// replacing all Windows separators with universal ones
	filePath := strings.Replace(args[0], "\\", "/", -1)
	fileFolder = path.Dir(filePath)
	fileName = path.Base(filePath)
	file, _ := os.Open(filePath)
	defer file.Close()
	csvReader := csv.NewReader(file)

	// setting a new delimiter for Reader
	if len(delimiter) == 1 {
		runeDelimiter = rune(delimiter[0])
		csvReader.Comma = runeDelimiter
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if header {
		headerLine = records[0]
		records = records[1:]
	}

	fileLength = len(records)
	rand.Seed(time.Now().UnixNano())
	randomizedIndexes := rand.Perm(fileLength) // creating an array of randomized indexes

	newFilePath := fmt.Sprintf("%s/tmp.%s", fileFolder, fileName)
	newFile, err := os.Create(newFilePath)
	defer newFile.Close()

	if err != nil {
		fmt.Printf("File creation failed with error: %s\n", err.Error())
		os.Exit(1)
	}

	csvWriter := csv.NewWriter(newFile)

	// setting a new delimiter for Writer
	if len(delimiter) == 1 {
		csvWriter.Comma = runeDelimiter
	}

	// writing a field names line if -f key is provided
	if header {
		csvWriter.Write(headerLine)
	}

	// based on randomized indexes write to file line by line
	for _, val := range randomizedIndexes {
		fmt.Println(records[val])
		csvWriter.Write(records[val])
	}

	csvWriter.Flush()

	// explicitly close files before renaming
	file.Close()
	newFile.Close()

	if backup {
		os.Rename(filePath, filePath+".old")
	} else {
		os.Remove(filePath)
	}

	os.Rename(newFilePath, filePath)
}

// shuffleCmd represents the shuffle command
var shuffleCmd = &cobra.Command{
	Use:   "shuffle",
	Short: "Shuffle lines in file",
	Long: `Shuffle lines in the file provided.
File will be rewritten, old backup is saved by default`,
	Args: func(cmd *cobra.Command, args []string) error {
		// validate amount of arguments
		if len(args) != 1 {
			return errors.New("A path should be provided as a single string argument")
		}

		// validate path
		if _, err := os.Stat(args[0]); err != nil {
			return errors.New("Provided path is invalid")
		}

		// validate length of delimiter
		if len(delimiter) != 1 {
			return errors.New("Delimiter should only be one character long")
		}

		return nil
	},
	Run: shuffle,
}

func init() {
	rootCmd.AddCommand(shuffleCmd)

	shuffleCmd.Flags().BoolVarP(&header, "field-names", "f", false, "Use if file contains a header line with field names")
	shuffleCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Single character to be used as delimiter")
	shuffleCmd.Flags().BoolVarP(&backup, "backup", "b", false, "Use this option if you need to backup an original file")
}
