//    TODOs
//
//    Infinite scroll on time mode
//    .rc file for keeping track of settings
//    start test with most recent settings when the program is run

package main

import (
	// "bufio"
	"embed"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
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
	escaped := false


	if (time_sec > 0 && len(text) > (3 * width)){
		fmt.Printf(text[:(3*width)])
		cursorToBeginning()
		cursorUp(2)
	}else{
		fmt.Printf(text)
		cursorToBeginning()
		cursorUp(maxLen / width) 
	}

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
				cursorToBeginning()
				if (len(text) < (3*width) || time_sec <= 0) {
					cursorUp((cursor_pos)/width)
				}

				if key == keyboard.KeyEsc {
					breakFlag = true
					escaped = true
					break
				}

				if key == keyboard.KeySpace {
					if text[cursor_pos] != ' '{
						for {
							if text[cursor_pos] == ' ' || cursor_pos >= maxLen - 1{
								break
							}
							hist += "_"
							cursor_pos += 1
						}
					} 
					char = ' '
				}



				if key == keyboard.KeyCtrlH{
					if len(hist) > 2 && hist[len(hist)-2] == ' '{
							cursor_pos -= 1
							hist = hist[:len(hist)-1]
					}
					for {
						if (len(hist) < 1 || cursor_pos <= 0){
							break
						}
						hist = hist[:len(hist)-1]
						cursor_pos -= 1
						
						if len(hist) == 0 || hist[len(hist)-1] == ' '{
							break
						}
					}
				
				} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
					if cursor_pos < 1 || len(hist) < 1{
						break
					}
					cursor_pos -= 1
					hist = hist[:len(hist)-1]
				} else {
					if string(char) != string(text[cursor_pos]) {
						hist += "_"
					} else {
						hist += string(char)
					}
					cursor_pos += 1
				}

				output := hist + text[cursor_pos:maxLen]
				lines_typed := len(hist)/width

				if cursor_pos == maxLen {
					cursorDown(lines_typed)
					fmt.Printf(output)
					breakFlag = true
					break
				}


				// If less than 3 widths worth of content
				if (len(output) < (3*width) || time_sec <= 0){
					//Go to the top of the output of 3 or less lines
					// Cursor is now on the bottom
					fmt.Printf(output)
					cursorToBeginning()
					// // place cursor on the correct line 
					cursorUp(maxLen/width - ((cursor_pos) / width))
					// // place cursor on the correct column
					cursorRight(len(hist) % width)
				} else {
					// TODO breaks if you spam spacebar and reach the end
					if (cursor_pos + 1 > width){
						widths_typed := lines_typed*width
						ending_index := widths_typed+(3*width)
						if ending_index >= maxLen{
							ending_index = maxLen-1
						}
						output = output[widths_typed:ending_index]
					} else{
						output = output[0:3*width]
					}
					fmt.Printf(output)
					cursorToBeginning()
					if (cursor_pos + 1 > width){
						cursorUp(2)
					} else {
						cursorUp((lines_typed)+2)
					}
					cursorRight(cursor_pos % width)
				}
				break
			default:
				break
		}
		if breakFlag{
			break
		}
	}
	
	if (!escaped){
		time_taken := time.Since(start).Seconds()
		mins_taken := time.Since(start).Minutes()
		printResults(hist, text, time_taken, mins_taken)
		return
	}
	fmt.Printf(("\n\n\n"))
}

func printResults(hist string, text string, time_taken float64, mins_taken float64){
	// remove all underscores, don't include underscores in score calc
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
			numwords = 30 * time_sec

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
