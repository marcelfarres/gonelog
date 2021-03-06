package log

import (
	"github.com/One-com/gonelog/syslog"
)

// Once events has been created, a Handler can ensure it's shipped to the log system.
// Formatters are a special kind og Handlers which ends the handler pipeline and
// convert the *Event to []byte (and does something with the bytes)

// Handler is the interface needed to be a part of the Handler chain.
// Formatters implement this to reveive events.
// The Handlers in this file pass the Event to further handlers
type Handler interface {
	Log(e Event) error
}

//---
type handlerFunc func(e Event) error

// HandlerFunc generates a Handler from a function, by calling it when Log is called.
func HandlerFunc(fn func(e Event) error) Handler {
	return handlerFunc(fn)
}
func (h handlerFunc) Log(e Event) error {
	return h(e)
}

//---

// FilterHandler lets a function evaluate whether to discard the Event or pass it on
// to a next Handler
func FilterHandler(fn func(e Event) bool, h Handler) Handler {
	return HandlerFunc(func(e Event) error {
		if fn(e) {
			return h.Log(e)
		}
		return nil
	})
}

// LvlFilterHandler discards events with a level above maxLvl
func LvlFilterHandler(maxLvl syslog.Priority, h Handler) Handler {
	return FilterHandler(func(e Event) (pass bool) {
		return e.Lvl <= maxLvl
	}, h)
}

// MultiHandler distributes the event to several Handlers
// if an error happen the last error is returned.
func MultiHandler(hs ...Handler) Handler {
	return HandlerFunc(func(e Event) error {
		var maybe_err error
		for _, h := range hs {
			err := h.Log(e)
			if err != nil {
				maybe_err = err
			}
		}
		return maybe_err
	})
}
