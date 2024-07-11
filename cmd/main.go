package main

import (
	"log"
	"test/ttracker/pkg/envparser"
	"test/ttracker/pkg/migrator"
)

type A struct {
	name string `strlen:"20"`
	age  int
}

type B struct {
	series string
	number string `strlen:"10"`
}

func main() {
	err := envparser.Load("../config.env")
	if err != nil {
		log.Fatal(err)
	}

	types := []any{
		A{name: "Albert", age: 44},
		B{series: "123 424", number: "123 4242"},
	}

	migrator.Init()
	migrator.Migrate(types)
	migrator.Seed(types)
}
