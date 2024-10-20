package main

import (
	// "bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"golang.org/x/term"
	"github.com/eiannone/keyboard"
)

func printIntro() {
	fmt.Print(`

=========================================================================
Welcome to CLI Type. To change the mode of typing, type :q

[1] Time
[2] Words

=========================================================================

Select a mode: `)
}

func timetype(text string) {
	width, height, termerr := term.GetSize(0)

	if height > 0{

	}

	if termerr != nil {
        return
    }

	if keeberr := keyboard.Open(); keeberr != nil {
		panic(keeberr)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("\nPress ESC to quit\n")

	var maxLen int = len(text)
	var cursor_pos int = 0
	var hist string = ""
	var start_timer bool = true
	var start = time.Now()

	fmt.Printf(text)
	cursorToBeginning()
	cursorUp(maxLen/width) //2

	for {
		char, key, err := keyboard.GetKey()
		if start_timer {
			start = time.Now()
			start_timer = false
		}

		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		if key == keyboard.KeySpace {
			char = ' '
		}

		if string(char) != string(text[cursor_pos]) {
			hist += "_"
		} else {
			hist += string(char)
		}

		cursorToBeginning()
		cursorUp((cursor_pos)/width)

		if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if cursor_pos > 0 {
				cursor_pos -= 1
			}
			if len(hist) > 1 {
				hist = hist[:len(hist)-2]
			}
		} else {
			cursor_pos += 1
		}

		if cursor_pos == maxLen {
			break
		}

		fmt.Printf(hist + text[cursor_pos:maxLen])
		cursorToBeginning()
		cursorUp(maxLen/width - ((cursor_pos)/width))
		cursorRight(len(hist)%width)
		
		// if (maxLen-cursor_pos > width){
		// 	
		// 	# for number of times maxLen-cursor_pos can be divided by width{
		// 		fmt.Printf("\033[A")
		// 	}
		// }
		// fmt.Printf("\033[%dD", maxLen-cursor_pos)

	}
	time_taken := time.Since(start).Seconds()
	mins_taken := time.Since(start).Minutes()
	num_words := len(strings.Split(text, " "))
	num_chars := len(text)

	fmt.Printf("\n\nTime taken: %.2f seconds", time_taken)
	fmt.Printf("\nWords typed: %d words", num_words)
	fmt.Printf("\nCharacters typed: %d characters", num_chars)
	fmt.Printf("\nCPM : %f CPM", float64(num_chars)/mins_taken)
	fmt.Printf("\nWPM: %f WPM", float64(num_words)/mins_taken)
}

func wordstype(text string) {

}

func gracefulShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)

	go func() {
		<-s
		fmt.Println("Shutting down gracefully due to signal.")
		os.Exit(0)
	}()
}

func numLinesToStartFromCursor(chars_per_line int, text_len int, cursor_pos int) int {
	var cursor_line int = cursor_pos/chars_per_line
	return cursor_line
}

func cursorUp(num int) {
	if num > 0{
		fmt.Printf("\033[%dA", num)
	}
}

func cursorToBeginning() {
	fmt.Printf("\r")
}

func cursorRight(num int){
	if num > 0{
		fmt.Printf("\033[%dC", num)
	}
}

func main() {
	gracefulShutdown()
	for {
		var i string
		printIntro()
		fmt.Scanln(&i)
		if i == "1" {
			numwordlist := 200
			wordlist := "wordlists/200.csv"
			var words []string
			var time int
			fmt.Println("Please enter the amount of time you wish to type for in seconds: ")
			fmt.Scanln(&time)
			numwords := time * 10

			file, err := os.Open(wordlist)

			if err != nil {
				fmt.Println("Error while reading the file", err)
				return
			}

			defer file.Close()

			reader := csv.NewReader(file)

			records, err := reader.ReadAll()

			for _, eachrecord := range records {
				words = append(words, eachrecord...)
			}

			selectedWords := make([]string, numwords)
			for i := 0; i < numwords; i++ {
				index := rand.Intn(numwordlist)
				selectedWords[i] = words[index]
			}
			timetype(strings.Join(selectedWords, " "))
		} else if i == "2" {
			fmt.Println("Chose Words")
		} else if i == ":q" {
			fmt.Println("Thank you for using CLI Type!")
			os.Exit(0)
		} else {
			fmt.Println("Unknown command.")
		}
	}
}
