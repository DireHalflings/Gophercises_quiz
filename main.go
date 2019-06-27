package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	//load flags
	timerFlag := flag.Int("timer", 30, "an int")
	fileFlag := flag.String("file", "problems.csv", "a string")
	shuffle := flag.Bool("shuffle", false, "a bool")

	flag.Parse()

	timer := *timerFlag
	fileName := *fileFlag

	//load csv file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	var questions []question

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var quizEntry question
		quizEntry.question = record[0]
		quizEntry.answer = cleanup(record[1])

		questions = append(questions, quizEntry)
	}

	//if shuffle flag, do the shuffle
	if *shuffle {
		questions = shuffleArray(questions)
	}

	//Greeting Message...
	fmt.Println("Welcome to Quiz Taker!")
	fmt.Println("You will have " + strconv.Itoa(timer) + " seconds for each question. There are " + strconv.Itoa(len(questions)) + " questions.")
	fmt.Print("Press 'Enter' to start your quiz now...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	clear()

	correct := 0

	//Start asking the questions...
	for i := 0; i < len(questions); i++ {

		fmt.Println(questions[i].question)

		c1 := make(chan string, 1)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			text = cleanup(text)
			c1 <- text
		}()

		select {
		case res := <-c1:
			if res == questions[i].answer {
				correct++
			}
		case <-time.After(time.Duration(timer) * time.Second):
			fmt.Println("Out of time!")
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		clear()
	}

	//Here is your grade, goodbye!
	fmt.Println("You got " + strconv.Itoa(correct) + " correct out of " + strconv.Itoa(len(questions)) + "!")
	fmt.Print("Press 'Enter' to exit now...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	clear()
}

type question struct {
	question string
	answer   string
}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func cleanup(answer string) string {
	answer = strings.Trim(answer, " ")
	answer = strings.ToUpper(answer)
	return answer
}

func shuffleArray(originalArray []question) []question {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]question, len(originalArray))
	n := len(originalArray)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(originalArray))
		ret[i] = originalArray[randIndex]
		originalArray = append(originalArray[:randIndex], originalArray[randIndex+1:]...)
	}
	return ret
}
