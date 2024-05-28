package main

import "fmt"

type Animal struct {
	Name string
	Age  int
}

func (animal *Animal) changeName(name string) {

	animal.Name = name
}

func main() {

	var animal = Animal{
		Name: "cat",
		Age:  2,
	}
	animal.changeName("dog")
	animal.changeName("Pig")
	fmt.Printf("%+v", animal)

}
