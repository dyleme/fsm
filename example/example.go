package main

import (
	"log"

	fsm "github.com/dyleme/fsm/generated"
)

func main() {
	state := fsm.Still

	var err error
	state, err = fsm.Stop(state)
	if err != nil {
		log.Fatal(err)
	}

	state = fsm.Moving
	state, err = fsm.Move(state)
	if err != nil {
		log.Fatal(err)
	}

	state, err = fsm.ToCrash(state)
	if err != nil {
		log.Fatal(err)
	}

}
