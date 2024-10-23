//    TODOs
//
//    ctrl + backspace
//    Infinite scroll on time mode
//    .rc file for keeping track of settings
//    start test with most recent settings when the program is run

package main

import (
	// "bufio"
	"embed"
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

//go:embed wordlists/200.csv
var wordlist200 embed.FS

//go:embed wordlists/1000.csv
var wordlist1000 embed.FS

//go:embed wordlists/2000.csv
var wordlist2000 embed.FS

//go:embed wordlists/5000.csv
var wordlist5000 embed.FS

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

type KeyPress struct {
	Char rune
	Key keyboard.Key
	Err error
}

func waitForKey(keypressed chan <- KeyPress) {
	for {
		char, key, err := keyboard.GetKey()
		keypressed <- KeyPress{Char: char, Key: key, Err: err}
	}
}

func typetest(text string, time_sec int) {
	done := make(chan bool)
	keypressed := make(chan KeyPress)

	if time_sec > 0 {
		go timer(time_sec, done)
	}

	go waitForKey(keypressed) 

	width, height, termerr := term.GetSize(0)

	if height > 0 {}

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
		if start_timer {
			start = time.Now()
			start_timer = false
		}

		select {
			case <- done:
				breakFlag = true
				cursorDown((3*width - cursor_pos) / width)
				break
			case pressed := <- keypressed:
				char := pressed.Char
				key := pressed.Key
				err := pressed.Err

				if start_timer {
					start = time.Now()
					start_timer = false
				}

				if err != nil {
					panic(err)
				}

				if key == keyboard.KeyEsc {
					breakFlag = true
					break
				}

				if key == keyboard.KeySpace {
					if text[cursor_pos] != ' '{
						for {
							if text[cursor_pos] == ' '{
								break
							}
							hist += "_"
							cursor_pos += 1
						}
					} 
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
					breakFlag = true
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
				break
			default:
				break
		}
		if breakFlag{
			break
		}
	}

	time_taken := time.Since(start).Seconds()
	mins_taken := time.Since(start).Minutes()
	typed_no_underscores := strings.ReplaceAll(hist, "_", "")
	num_words_typed := len(strings.Fields(typed_no_underscores))
	num_chars_typed := len(strings.ReplaceAll(typed_no_underscores, " ", ""))
	// this includes num chars skipped with spacebar
	num_chars_in_hist := len(strings.ReplaceAll(hist, " ", ""))

	fmt.Printf("\n\nTime taken: %.2f seconds", time_taken)
	fmt.Printf("\nWords typed: %d words", num_words_typed)
	fmt.Printf("\nCharacters typed: %d characters", num_chars_typed)
	fmt.Printf("\nCPM : %.2f CPM", float64(calcNumCorrectChars(hist, text))/mins_taken)
	fmt.Printf("\nWPM: %.2f WPM", float64(calcNumCorrectWords(hist, text))/mins_taken)
	if num_chars_typed == 0{
		fmt.Printf("\nAccuracy: -")
	} else {
		fmt.Printf("\nAccuracy: %.2f%%", (float64(calcNumCorrectChars(hist, text)*100)/float64(num_chars_in_hist)))
	}
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
	var shorter_len int  
	if len(typed) < len(original_text){
		shorter_len = len(typed)
	} else{
		shorter_len = len(original_text)
	}
	for i := 0; i < shorter_len; i++ {
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

			data, err := wordlist200.ReadFile("wordlists/200.csv")
			if err != nil {
				fmt.Println("Error while reading the file", err)
				return
			}

			reader := csv.NewReader(strings.NewReader(string(data)))

			records, err := reader.ReadAll()
			// records := strings.Join(string(data), ",")
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

			data, err := wordlist200.ReadFile("wordlists/200.csv")
			if err != nil {
				fmt.Println("Error while reading the file", err)
				return
			}

			reader := csv.NewReader(strings.NewReader(string(data)))
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
