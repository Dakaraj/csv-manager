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
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// writeCSVGoroutine function takes a slice of source file and writes
// a new one with index supplied
func writeCSVGoroutine(index int, lines [][]string, headerLine []string, c chan uint) {
	newFilePath := fmt.Sprintf("%s/%03d.%s", fileFolder, index, fileName)
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

	csvWriter.WriteAll(lines) // writing all provided lines

	c <- 1 // this message means the task is complete
}

func divide(cmd *cobra.Command, args []string) {
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
	if fileLength < int(amount) {
		fmt.Printf("File lines amount (%d) is smaller than requested parts (%d)\n", amount, fileLength)
		os.Exit(0)
	}

	// based on the total amount of lines in file dividing this
	// equally between all destination files
	// remainder of equation is spread between files as well
	quotient := fileLength / int(amount)
	remainder := fileLength % int(amount)
	var linesPerFile = make([]int, amount)
	for i := 0; i < int(amount); i++ {
		if remainder > 0 {
			linesPerFile[i] = quotient + 1
			remainder--
		} else {
			linesPerFile[i] = quotient
		}
	}

	// writers are launched as goroutines
	// each call receives a unique slice of initial records
	curIndex := 0
	c := make(chan uint)
	for i, val := range linesPerFile {
		go writeCSVGoroutine(i+1, records[curIndex:curIndex+val], headerLine, c)
		curIndex = curIndex + val
	}

	// this loop awaits all goroutines to complete successfuly
	var total uint
	for total != amount {
		total += <-c
	}

	fmt.Println("New files created successfully")
}

// divideCmd represents the divide command
var divideCmd = &cobra.Command{
	Use:   "divide [path to csv]",
	Short: "Divide CSV file into equal parts",
	Long: `This method divides a CSV file into the number of equal parts.
Parts amount shoud be provided with arguments.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// validate amount of arguments
		if len(args) != 1 {
			return errors.New("A path should be provided as a single string argument")
		}

		// validate path
		if _, err := os.Stat(args[0]); err != nil {
			return errors.New("Provided path is invalid")
		}

		// validate amount of destination files
		if amount < 2 || amount > 999 {
			return errors.New("Invalid value for parts amount. Should be between 2 and 999")
		}

		// validate length of delimiter
		if len(delimiter) != 1 {
			return errors.New("Delimiter should only be one character long")
		}

		return nil
	},
	Run: divide,
}

func init() {
	rootCmd.AddCommand(divideCmd)

	// defining flags for application
	divideCmd.Flags().UintVarP(&amount, "amount", "a", 2, "Amount of equal parts to divide a file. Should be between 2 and 999")
	divideCmd.Flags().BoolVarP(&header, "field-names", "f", false, "Use if file contains a header line with field names")
	divideCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Single character to be used as delimiter")
}
