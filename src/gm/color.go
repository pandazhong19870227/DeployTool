package main

import (
	"fmt"
)

func yellowBegin() {
	fmt.Printf("\033[1;33m")
}

func greenBegin() {
	fmt.Printf("\033[1;32m")
}

func redBegin() {
	fmt.Printf("\033[1;31m")
}

func colorEnd() {
	fmt.Printf("\033[0m")
}
