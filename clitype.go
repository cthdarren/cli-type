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
Welcome to CLI Type. To change the mode of typing, type :q

[1] Time
[2] Words

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
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc {
			fmt.Printf("Escaped")
			break
		}
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
