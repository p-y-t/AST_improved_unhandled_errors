package main

import (
	"AST_improved_unhandled_errors/utils"
	"fmt"
	"strconv"
	"strings" // Standard library
)

func test_01() {
	_ = greet()

	sum, _ := utils.Add(5, 10)
	fmt.Printf("Sum from utils: %d\n", sum)

	num, _ := strconv.Atoi("1234")
	fmt.Println("Num is: ", num)

	upper := strings.ToUpper("hello, world!")
	fmt.Printf("Uppercased string: %s\n", upper)

	reversed, _, _ := utils.ReverseString("hello")
	fmt.Printf("Reversed string: %s\n", reversed)
}

// greet prints a greeting message.
func greet() int {
	fmt.Println("Hello from the main file!")
	return 0
}
