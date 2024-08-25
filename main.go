package fyneflow

import (
	"errors"
	"math"

	"fyne.io/fyne/v2"
)

type (
	FlowKey int

	FlowItemGenerator func(flow *Flow) fyne.CanvasObject

	FlowItem struct {
		generator FlowItemGenerator
	}

	Flow struct {
		window  fyne.Window
		current FlowKey
		next    chan FlowKey
		items   map[FlowKey]*FlowItem

		close    bool
	}
)

const (
	FlowKeyNone FlowKey = math.MinInt
)

// NewFlow creates a new Flow.
// The generic type is the type for the flow keys.
// The Flow is associated with the given fyne.Window.
// Use the returned Flow to add FlowItems and to switch between them.
// Use the Close method to close the Flow. This will not close the associated fyne.Window!
func NewFlow(w fyne.Window) *Flow {
	f := new(Flow)
	f.current = FlowKeyNone
	f.window = w
	f.items = make(map[FlowKey]*FlowItem)
	f.next = make(chan FlowKey, 1)
	go f.loop()
	return f
}

// Close closes the Flow.
// This will not close the associated fyne.Window!
func (f *Flow) Close() {
	f.next <- FlowKeyNone

	f.close = true
	for k := range f.items {
		delete(f.items, k)
	}
}

func (f *Flow) loop() {
	for !f.close {
		next := <-f.next
		if next != f.current {
			f.apply(next)
		}
	}
}

// Set adds a FlowItem to the Flow.
// If the key already exists, the FlowItem will be overwritten.
func (f *Flow) Set(key FlowKey, generator FlowItemGenerator) *FlowItem {
	fi := new(FlowItem)
	fi.generator = generator
	apply := len(f.items) == 0
	f.items[key] = fi
	if apply {
		f.next <- key
	}
	return fi
}

// GoTo sets the content of the window to the content of the FlowItem associated with the given key.
// If the key is not found, an error is returned.
func (f *Flow) GoTo(key FlowKey) error {
	if f.current == key {
		return nil
	}

	if _, ok := f.items[key]; ok {
		f.next <- key
		return nil
	} else {
		return errors.New("flow: key not found")
	}
}

func (f *Flow) apply(next FlowKey) {
	fi, ok := f.items[next]
	if !ok || fi == nil {
		return
	}

	obj := fi.generator(f)

	f.window.SetContent(obj)
	f.current = next
}

// Current returns the key of the FlowItem that is currently displayed.
func (f *Flow) Current() FlowKey {
	return f.current
}
