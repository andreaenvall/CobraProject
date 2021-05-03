/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A quick little basic quiz",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		quiz()
	},
}

func init() {

	rootCmd.AddCommand(addCmd)
}

type Line struct {
	Question string
	Answer   string
}
type Score struct {
	Name   string
	Points string
}

func quiz() {

	quizfil := flag.String("csv", "quiz.csv", "path to csv file with the quiz")

	flag.Parse()

	csvFile, err := os.Open(*quizfil)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	var lines []Line
	for {
		line, error := csvReader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		lines = append(lines, Line{
			Question: line[0],
			Answer:   line[1],
		})
	}

	count := 0

	for idx, line := range lines {
		fmt.Print(strconv.Itoa(idx+1) + ": " + line.Question + ": ")

		var input string
		fmt.Scan(&input)

		if input == line.Answer {
			count++
		}
	}
	fmt.Println()
	fmt.Println("You scored " + strconv.Itoa(count) + " out of " + strconv.Itoa(len(lines)))
	scores(count)

}

func scores(scores int) {

	csvfile, err := os.Open("scoreboard.csv")

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	var countscore float64 = 0
	var sumOfPlayers float64 = 0

	for {

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stringvar := record[0]

		//här förstår jag inte riktigt varför jag måste ha två variabler fast jag inte vill ??
		intVar, err := strconv.Atoi(stringvar)
		if intVar < scores {
			countscore++
		}
		sumOfPlayers++
	}
	var percentfloat float64 = ((countscore / sumOfPlayers) * 100)
	var y = fmt.Sprint(math.RoundToEven(percentfloat))

	fmt.Println("You were better than " + y + "% of all quizzers")
	insertscore(scores)
}

func insertscore(countscore int) {

	csvfile, err := os.OpenFile("scoreboard.csv", os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	var score [][]string
	score = append(score, []string{strconv.Itoa(int(countscore)), ""})
	writer := csv.NewWriter(csvfile)
	writer.WriteAll(score)
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your score has now been registrated!")
}
