package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	name string
	age  int
}

func (p *person) SetName(name string) *person {
	p.name = name
	return p
}

func (p *person) SetAge(age int) *person {
	p.age = age
	return p
}

func main() {
	p := &person{}

	p.SetAge(10).SetName("test")

	p.SetName("test")
	fmt.Printf("%#v", p)
}

func main1() {

	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	sp := &s
	fmt.Println(sp.age)
	sp.age = 51
	fmt.Println(sp.age, s.age)
}
