package main

import (
	// "bufio"
	"encoding/csv"
	"fmt"
	"github.com/eiannone/keyboard"
	"golang.org/x/term"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

func timer(duration int, done chan<- bool) {
	time.Sleep(time.Duration(duration) * time.Second)
	done <- true
}

func typetest(text string, time_sec int) {
	done := make(chan bool)

	if time_sec > 0 {
		go timer(time_sec, done)
	}
	width, height, termerr := term.GetSize(0)

	if height > 0 {

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
	breakFlag := false

	fmt.Printf(text)
	cursorToBeginning()
	cursorUp(maxLen / width) //2

	for {
		select {
		case <-done:
			fmt.Println("DONEEEE")
			breakFlag = true
			break
		default:
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

			// TODO when i press space go to the next word instead of putting a space where i am now
			if string(char) != string(text[cursor_pos]) {
				hist += "_"
			} else {
				hist += string(char)
			}

			cursorToBeginning()
			cursorUp((cursor_pos) / width)

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
				cursorDown((cursor_pos) / width)
				break
			}

			output := hist + text[cursor_pos:maxLen]

			// TODO: for time type mode, make it "infinite" scrolling
			// if (len(output) > (3 * width)){
			// 	// if the cursor is on the third line or lower
			// 	if (len(hist) > (2 * width)){
			// 		lines_typed := len(hist)%width + 1
			// 		fmt.Printf(output[lines_typed:])
			// 	} else{
			// 		fmt.Printf(output[:(3*width)])
			// 	}
			// }
			fmt.Printf(output)
			cursorToBeginning()
			cursorUp(maxLen/width - ((cursor_pos) / width))
			cursorRight(len(hist) % width)
		}
		if breakFlag {
			break
		}
	}

	time_taken := time.Since(start).Seconds()
	mins_taken := time.Since(start).Minutes()
	num_words := len(strings.Split(text, " "))
	num_chars := len(strings.ReplaceAll(text, " ", ""))

	fmt.Printf("\n\nTime taken: %.2f seconds", time_taken)
	fmt.Printf("\nWords typed: %d words", num_words)
	fmt.Printf("\nCharacters typed: %d characters", num_chars)
	fmt.Printf("\nCPM : %f CPM", float64(calcNumCorrectChars(hist, text))/mins_taken)
	fmt.Printf("\nWPM: %f WPM", float64(calcNumCorrectWords(hist, text))/mins_taken)
	fmt.Printf("\nAccuracy: %.2f%", (float64(calcNumCorrectChars(hist, text))/float64(num_chars)) * 100)
}

func calcNumCorrectWords(typed string, original_text string) int {
	typed_arr := strings.Split(typed, " ")
	original_arr := strings.Split(original_text, " ")
	var shorter_arr_len int
	result := 0

	if len(typed_arr) < len(original_arr) {
		shorter_arr_len = len(typed_arr)
	} else {
		shorter_arr_len = len(original_arr)
	}
	for i := 0; i < shorter_arr_len; i++ {
		if typed_arr[i] == original_arr[i] {
			result++
		}
	}

	return result
}

func calcNumCorrectChars(typed string, original_text string) int {
	result := 0
	for i := 0; i < len(original_text); i++ {
		if typed[i] == original_text[i] {
			if typed[i] != ' ' {
				result++
			}
		}
	}
	return result
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
	var cursor_line int = cursor_pos / chars_per_line
	return cursor_line
}

func cursorDown(num int) {
	if num > 0 {
		fmt.Printf("\033[%dB", num)
	}
}

func cursorUp(num int) {
	if num > 0 {
		fmt.Printf("\033[%dA", num)
	}
}

func cursorToBeginning() {
	fmt.Printf("\r")
}

func cursorRight(num int) {
	if num > 0 {
		fmt.Printf("\033[%dC", num)
	}
}

func main() {
	gracefulShutdown()
	numwordlist := 200
	wordlist := "wordlists/200.csv"
	var words []string
	var numwords int

	for {
		var i string
		printIntro()
		fmt.Scanln(&i)
		if i == "1" {
			fmt.Println("Please enter the number of time you wish to type for in seconds: ")
			var time_sec int
			fmt.Scanln(&time_sec)
			numwords = 10 * time_sec

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
			typetest(strings.Join(selectedWords, " "), time_sec)
		} else if i == "2" {
			fmt.Println("Please enter the number of words you wish to type for: ")
			fmt.Scanln(&numwords)

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
			typetest(strings.Join(selectedWords, " "), -1)
		} else if i == ":q" {
			fmt.Println("Thank you for using CLI Type!")
			os.Exit(0)
		} else {
			fmt.Println("Unknown command.")
		}
	}
}
