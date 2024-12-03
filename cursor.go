package main

import "fmt"

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
