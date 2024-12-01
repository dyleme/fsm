package basic

import (
	"fmt"
	"slices"
)

const (
	Crash  State = "crash"
	Moving State = "moving"
	Stay   State = "stay"
)

var allStates = map[State]struct{}{
	Crash:  {},
	Moving: {},
	Stay:   {},
}

type event struct {
	Name string
	src  []State
	dst  State
}

var (
	toCrashEvent = &event{
		Name: "ToCrash",
		src: []State{
			Moving,
		},
		dst: Crash,
	}

	toMovingEvent = &event{
		Name: "ToMoving",
		src: []State{
			Stay,
			Moving,
		},
		dst: Moving,
	}

	toStayEvent = &event{
		Name: "ToStay",
		src: []State{
			Moving,
		},
		dst: Stay,
	}
)

var possibleEvents = map[State][]*event{
	Crash: {},
	Moving: {
		toCrashEvent,
		toMovingEvent,
		toStayEvent,
	},
	Stay: {
		toMovingEvent,
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

// CanToMoving returns value if state can be change to the ToMoving.
func CanToMoving(state State) bool {
	return slices.Contains(possibleEvents[state], toMovingEvent)
}

// ToMoving returns new state if it can be reached or err otherise.
func ToMoving(state State) (State, error) {
	events, ok := possibleEvents[state]
	if !ok {
		return state, UnknownStateError{string(state)}
	}

	if slices.Contains(events, toMovingEvent) {
		return toMovingEvent.dst, nil
	}

	return state, ProhibittedEventError{string(state), toMovingEvent.Name}
}

// CanToStay returns value if state can be change to the ToStay.
func CanToStay(state State) bool {
	return slices.Contains(possibleEvents[state], toStayEvent)
}

// ToStay returns new state if it can be reached or err otherise.
func ToStay(state State) (State, error) {
	events, ok := possibleEvents[state]
	if !ok {
		return state, UnknownStateError{string(state)}
	}

	if slices.Contains(events, toStayEvent) {
		return toStayEvent.dst, nil
	}

	return state, ProhibittedEventError{string(state), toStayEvent.Name}
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
