package main

import (
	"fmt"
	"strings"
)

// A thing that can yell
type Yeller interface {
	yell(message string)
}

func yell(message string) {
	fmt.Println(strings.ToUpper(message))
}

type Person struct {
	Yeller // Option 2: embed the Yeller interface
}

// Option 1: Person implements Yeller interface by adding
// Explicit, but also stuck with this implementation
// func (p Person) yell(message string) {
// 	fmt.Println(strings.ToUpper(message))
// }

func main() {
	p := Person{}
	// p.yell("hello")
	p.yell("goodbye")
}