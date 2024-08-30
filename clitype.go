package main

import "fmt"

func printIntro(){
	fmt.Print(`
Welcome to CLI Type. To change the mode of typing, type :q

[1] Time
[2] Words

Select a mode: `)
}
func main(){
	for {
		var i int
		printIntro()
		fmt.Scan(&i)
		fmt.Println(i)
	}
}
