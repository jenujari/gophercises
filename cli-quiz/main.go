package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filePath := flag.String("csv", "questions.csv", "Provide  csv file path that contains quiz")
	timeLimit := flag.Int("limit", 30, "Provide time limit for the quiz")
	flag.Parse()

	file, err := os.OpenFile(*filePath, os.O_RDONLY, os.ModeTemporary)
	handleError(err)

	csvReader := csv.NewReader(file)
	lines, e := csvReader.ReadAll()
	handleError(e)

	questionList := parseLines(lines)
	counter := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

quizLoop:
	for i, e := range questionList {
		fmt.Printf("Problem #%d: %s = ", i+1, e.q)
		ansChan := make(chan string)

		go func() {
			var ans string
			fmt.Scanf("%s\n", &ans)
			ansChan <- strings.TrimSpace(ans)
		}()

		select {
		case <-timer.C:
			{
				fmt.Printf("\nYou run out of time %d secons \nProgram is now exiting.\n", *timeLimit)
				break quizLoop
			}
		case ans := <-ansChan:
			{
				if ans == e.a {
					counter++
				}
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", counter, len(questionList))
}

func parseLines(lns [][]string) []Question {
	questions := make([]Question, len(lns))
	for i, r := range lns {
		questions[i] = Question{
			q: r[0],
			a: strings.TrimSpace(r[1]),
		}
	}
	return questions
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

type Question struct {
	q string
	a string
}
