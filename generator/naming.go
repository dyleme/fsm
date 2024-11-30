package generator

import (
	"sort"
	"strings"
)

type ConflictingNamesError struct {
	firstName  string
	secondName string
}

func (error ConflictingNamesError) Error() string {
	return "conflicting names: " + error.firstName + " and " + error.secondName
}

func BetterNaming(fd fileData) (data, error) {
	states := make([]state, 0, len(fd.td.States))
	events := make(map[string]event)

	for _, s := range fd.td.States {
		states = append(states, state{Name: strings.Title(s), Value: s})
	}

	for _, le := range fd.td.Events {
		e := events[le.Dst]

		var src state
		for _, s := range states {
			if s.Value == le.Src {
				src = s
				break
			}
		}

		value := strings.Title(le.Label)
		if value == "" {
			value = "To" + strings.Title(le.Dst)
		}
		name := strings.ToLower(value[:1]) + value[1:] + "Event"

		if e.Name != "" && e.Name != name {
			return data{}, ConflictingNamesError{e.Name, name}
		}

		if e.Value != "" && e.Value != value {
			return data{}, ConflictingNamesError{e.Value, value}
		}

		var dst state
		for _, s := range states {
			if s.Value == le.Dst {
				dst = s
				break
			}
		}

		e.Name = name
		e.Value = value
		e.Src = append(e.Src, src)
		e.Dst = dst

		events[le.Dst] = e
	}

	eventSlice := make([]event, 0, len(events))
	for _, e := range events {
		eventSlice = append(eventSlice, e)
	}

	sort.Slice(eventSlice, func(i, j int) bool {
		return eventSlice[i].Name < eventSlice[j].Name
	})

	sort.Slice(states, func(i, j int) bool {
		return states[i].Name < states[j].Name
	})

	d := data{
		PkgName:        fd.packageName,
		States:         states,
		Events:         eventSlice,
		PossibleEvents: possibleEvents(states, eventSlice),
		GenDynamic:     false,
		GenType:        false,
	}

	return d, nil
}
