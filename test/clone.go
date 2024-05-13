package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Field1 int
	Field2 string
}

func cloneStruct(src interface{}) interface{} {
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Struct {
		panic("Input value is not a struct")
	}

	dstValue := reflect.New(srcValue.Type()).Elem()
	dstValue.Set(srcValue)

	return dstValue.Interface()
}

func main() {
	original := MyStruct{Field1: 42, Field2: "hello"}
	clone := cloneStruct(original).(MyStruct)

	fmt.Println("Original struct:", original)
	fmt.Println("Cloned struct:", clone)

	// Modify the cloned struct
	clone.Field1 = 100
	clone.Field2 = "world"

	fmt.Println("\nAfter modification:")
	fmt.Println("Original struct:", original)
	fmt.Println("Cloned struct:", clone)
}
