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
	events := make(map[string]event)
	states := make(map[string]state)

	for _, s := range fd.td.Events {
		states[s.Src] = state{Name: strings.Title(s.Src), Value: s.Src}
		states[s.Dst] = state{Name: strings.Title(s.Dst), Value: s.Dst}
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

	stateSlice := make([]state, 0, len(states))
	for _, s := range states {
		stateSlice = append(stateSlice, s)
	}

	sort.Slice(eventSlice, func(i, j int) bool {
		return eventSlice[i].Name < eventSlice[j].Name
	})

	sort.Slice(stateSlice, func(i, j int) bool {
		return stateSlice[i].Name < stateSlice[j].Name
	})

	d := data{
		PkgName:        fd.packageName,
		States:         stateSlice,
		Events:         eventSlice,
		PossibleEvents: possibleEvents(stateSlice, eventSlice),
		GenDynamic:     false,
		GenType:        false,
	}

	return d, nil
}
