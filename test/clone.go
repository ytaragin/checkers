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

func aTest() {
	m1 := MyStruct{3, "hello"}
	m2 := MyStruct{4, "there"}

	ms := []MyStruct{m1, m2}

	fmt.Printf("%+v\n", ms)

	for _, i := range ms {
		fmt.Printf("%+v\n", i)
		i.Field1++
		fmt.Printf("%+v\n", i)
	}

	fmt.Printf("%+v\n", ms)

	for i := 0; i < len(ms); i++ {
		fmt.Printf("%+v\n", ms[i])
		ms[i].Field1++
		fmt.Printf("%+v\n", ms[i])
	}

	fmt.Printf("%+v\n", ms)

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

	aTest()
}
