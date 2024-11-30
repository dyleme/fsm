package fsm

// A --> B
// B --> C
// A --> C
// C --> A: last

// still --> moving: move
// moving --> moving: move
// moving --> still: c
// moving --> crash
type State string
