package main

import (
	"fmt"
)

func fillSlice(numbers []int) []int {
	// Fill the slice with some values (replace with your logic)
	for i := 0; i < 5; i++ {
		numbers = append(numbers, i*2)
	}
	fmt.Println("Filled slice:", numbers)
	return numbers
}

func maina() {
	// Create an empty slice
	// numbers := []int{}

	// Call the function to fill the slice
	numbers := fillSlice([]int{})

	fmt.Println("Filled slice:", numbers)
}
