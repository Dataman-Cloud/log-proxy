package store

import (
	"golang.org/x/net/context"
)

const key = "store"

// Setter defines a context that enables setting values
type Setter interface {
	Set(string, interface{})
}

// returns the Store associated with this context
func FromContext(c context.Context) Store {
	return c.Value(key).(Store)
}

// adds the Store to this context if it supports this setter interface
func ToContext(c Setter, store Store) {
	c.Set(key, store)
}
