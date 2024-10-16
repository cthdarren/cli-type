package main

import (
	// "bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func timetype() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	var text string = "The quick brown fox jumps over the lazy dog"
	var maxLen int = len(text)
	var cursor_pos int = 0
	var hist string = ""

	fmt.Printf(text)
	fmt.Printf("\033[%dD", maxLen)

	for {
		char, key, err := keyboard.GetKey()

		if err != nil {
			panic(err)
		}

		if key == keyboard.KeySpace{
			hist += " " 
		} else{
			hist += string(char)
		}

		if key == keyboard.KeyEsc {
			break
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
}

func wordstype() {

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
			fmt.Println("Chose Time")
			timetype()
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
