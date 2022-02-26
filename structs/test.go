package structs

import (
	"fmt"
	"reflect"
)

type A struct {
	Foo string
	Bar int
}

func test() {
	// fields names
	fmt.Println(Names(&A{}))

	// struct names
	fmt.Println(Name(&A{}))

	// fiels types
	x := Fields(&A{})
	kinds := make([]reflect.Kind, 0)
	for _, v := range x {
		kinds = append(kinds, v.Kind())
	}
	fmt.Println(kinds)
}

//[Foo Bar]
//A
//[string int]
