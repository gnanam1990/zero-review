package core

import "time"

// ChatMessage is one chat turn.
type ChatMessage struct {
	Role      string // "ai" | "user"
	Text      string
	Timestamp time.Time
}

// Toast is a transient status message.
type Toast struct {
	Message string
	Kind    string // success | error | info
	Until   time.Time
}

// LoadingStep is one item in the progress timeline.
type LoadingStep struct {
	Label  string
	Done   bool
	Active bool
}
