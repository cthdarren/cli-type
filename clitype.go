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
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	var maxLen int = len(text)
	var cursor_pos int = 0
	var hist string = ""
	var start_timer bool = true
	var start = time.Now()

	fmt.Printf(text)
	fmt.Printf("\033[%dD", maxLen)

	for {
		char, key, err := keyboard.GetKey()
		if start_timer{
			start = time.Now()
			start_timer = false
		}

		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		if key == keyboard.KeySpace{
			hist += " " 
		} else {
			if string(char) != string(text[cursor_pos]){
				hist += "_"
			} else{
				hist += string(char)
			}
		}

		if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if cursor_pos > 0{
				cursor_pos -= 1
			}
			if len(hist) > 1 {
				hist = hist[:len(hist)-2]
			}
		} else {
			cursor_pos += 1
		}

		if cursor_pos == maxLen{
			break
		}


		fmt.Printf("\r" + hist + text[cursor_pos:maxLen])
		fmt.Printf("\033[%dD", maxLen - cursor_pos)

	}
	time_taken := time.Since(start).Seconds()
	mins_taken := time.Since(start).Minutes()
	num_words := len(strings.Split(text, " "))
	num_chars := len(text)

	fmt.Printf("\n\n Time taken: %.2f seconds", time_taken)
	fmt.Printf("\n Words typed: %d words",num_words)
	fmt.Printf("\n Characters typed: %d characters",num_chars)
	fmt.Printf("\n CPM : %f CPM",float64(num_chars)/mins_taken)
	fmt.Printf("\n WPM: %f WPM",float64(num_words)/mins_taken)
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
			time := 3 
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
