package main

import (
	"flag"
	"fmt"
	"os"
	"password_gen/internal/generator"
)

func main() {
	var length int
	flag.IntVar(&length, "l", 16, "password length (6-128)")
	flag.IntVar(&length, "length", 16, "password length (6-128)")
	flag.Parse()

	if length < 6 || length > 128 {
		fmt.Printf("length must be between 6 and 128\n")
		os.Exit(1)
	}

	g := generator.New()
	password, err := g.Generate(length)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
	}
	fmt.Println(password)
}
