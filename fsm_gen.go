package fsm

import (
	"fmt"
	"slices"
)

const (
	Crash  State = "crash"
	Moving State = "moving"
	Still  State = "still"
)

var allStates = map[State]struct{}{
	Crash:  {},
	Moving: {},
	Still:  {},
}

type event struct {
	Name string
	src  []State
	dst  State
}

var (
	cEvent = &event{
		Name: "C",
		src: []State{
			Moving,
		},
		dst: Still,
	}

	moveEvent = &event{
		Name: "Move",
		src: []State{
			Still,
			Moving,
		},
		dst: Moving,
	}

	toCrashEvent = &event{
		Name: "ToCrash",
		src: []State{
			Moving,
		},
		dst: Crash,
	}
)

var possibleEvents = map[State][]*event{
	Crash: {},
	Moving: {
		cEvent,
		moveEvent,
		toCrashEvent,
	},
	Still: {
		moveEvent,
	},
}

// UnknownEventError is errror that provided state not exists.
type UnknownStateError struct {
	State string
}

func (e UnknownStateError) Error() string {
	return fmt.Sprintf("Unknown state: '%s'", e.State)
}

// ProhibittedEventError is error that signals that Event can't be done from State.
type ProhibittedEventError struct {
	State string
	Event string
}

func (e ProhibittedEventError) Error() string {
	return fmt.Sprintf("Prohibitted event '%s' in state '%s'", e.Event, e.State)
}

// CanC returns value if state can be change to the C.
func CanC(state State) bool {
	return slices.Contains(possibleEvents[state], cEvent)
}

// C returns new state if it can be reached or err otherise.
func C(state State) (State, error) {
	events, ok := possibleEvents[state]
	if !ok {
		return state, UnknownStateError{string(state)}
	}

	if slices.Contains(events, cEvent) {
		return cEvent.dst, nil
	}

	return state, ProhibittedEventError{string(state), cEvent.Name}
}

// CanMove returns value if state can be change to the Move.
func CanMove(state State) bool {
	return slices.Contains(possibleEvents[state], moveEvent)
}

// Move returns new state if it can be reached or err otherise.
func Move(state State) (State, error) {
	events, ok := possibleEvents[state]
	if !ok {
		return state, UnknownStateError{string(state)}
	}

	if slices.Contains(events, moveEvent) {
		return moveEvent.dst, nil
	}

	return state, ProhibittedEventError{string(state), moveEvent.Name}
}

// CanToCrash returns value if state can be change to the ToCrash.
func CanToCrash(state State) bool {
	return slices.Contains(possibleEvents[state], toCrashEvent)
}

// ToCrash returns new state if it can be reached or err otherise.
func ToCrash(state State) (State, error) {
	events, ok := possibleEvents[state]
	if !ok {
		return state, UnknownStateError{string(state)}
	}

	if slices.Contains(events, toCrashEvent) {
		return toCrashEvent.dst, nil
	}

	return state, ProhibittedEventError{string(state), toCrashEvent.Name}
}

// Parse returns corresponding state to the provided string,
// or error eitherise.
func Parse(state string) (State, error) {
	if _, ok := allStates[State(state)]; !ok {
		return "", UnknownStateError{state}
	}

	return State(state), nil
}

// IsLastState return true if no futher transitions could be done from provided state.
func IsLastState(state State) bool {
	return len(possibleEvents[state]) == 0
}
