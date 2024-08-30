package main

import "fmt"

func printIntro(){
	fmt.Print(`
Welcome to CLI Type. To change the mode of typing, type :q

[1] Time
[2] Words

Select a mode: `)
}

func timetype(){
	var typed string
	fmt.Println("testing words for you to type this is how is done if fan cup food lunch hungry rice snowboard goggles\r")	
	fmt.Scan(&typed)
	fmt.Println("type test complete")
}

func wordstype(){

}

func main(){
	for {
		var i string
		printIntro()
		fmt.Scan(&i)
		if i == "1" {
			fmt.Println("Chose Time")
			timetype()
		} else if i == "2" {
			fmt.Println("Chose Words")
		} else if i == ":q" {
			fmt.Println("exiting...")
		} else {
			fmt.Println("Unknown command.")
		}
	}
}
