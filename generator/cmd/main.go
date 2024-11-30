package main

import (
	"log"

	"github.com/dyleme/fsm/generator"
)

func main() {

	comParser, err := generator.NewCommentsParser("State", "mermaid")
	if err != nil {
		log.Fatal(err)
	}

	rParser := &generator.RealParser{
		CommentParser: comParser,
	}

	gen := generator.NewFSMGenerator(rParser)

	err = gen.Do()
	if err != nil {
		log.Fatal(err)
	}
}
