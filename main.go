package fyneflow

import (
	"errors"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type FlowItemGenerator func() fyne.CanvasObject

type FlowItem struct {
	generator FlowItemGenerator
}

type Flow struct {
	window  fyne.Window
	flowMut sync.Mutex
	current string
	next    string
	signal  chan struct{}
	items   map[string]*FlowItem

	stateMut sync.Mutex
	strState map[string]binding.String
	intState map[string]binding.Int
	close    bool
}

// NewFlow creates a new Flow.
// The generic type is the type for the flow keys.
// The Flow is associated with the given fyne.Window.
// Use the returned Flow to add FlowItems and to switch between them.
// Use the Close method to close the Flow. This will not close the associated fyne.Window!
func NewFlow(w fyne.Window) *Flow {
	f := new(Flow)
	f.window = w
	f.flowMut = sync.Mutex{}
	f.items = make(map[string]*FlowItem)
	f.signal = make(chan struct{}, 1)
	go f.loop()
	return f
}

// Close closes the Flow.
// This will not close the associated fyne.Window!
func (f *Flow) Close() {
	f.signal <- struct{}{}
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	f.close = true
	for k := range f.items {
		delete(f.items, k)
	}
}

func (f *Flow) loop() {
	for !f.close {
		<-f.signal
		if f.next != f.current {
			f.apply()
		}
	}
}

// Set adds a FlowItem to the Flow.
// If the key already exists, the FlowItem will be overwritten.
func (f *Flow) Set(key string, generator FlowItemGenerator) *FlowItem {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	fi := new(FlowItem)
	fi.generator = generator
	apply := len(f.items) == 0
	f.items[key] = fi
	if apply {
		f.next = key
	}
	f.signal <- struct{}{}
	return fi
}

// GoTo sets the content of the window to the content of the FlowItem associated with the given key.
// If the key is not found, an error is returned.
func (f *Flow) GoTo(key string) error {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	if f.current == key {
		return nil
	}

	if _, ok := f.items[key]; ok {
		f.next = key
		f.signal <- struct{}{}
		return nil
	} else {
		return errors.New("flow: key not found")
	}
}

func (f *Flow) apply() {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	fi, ok := f.items[f.next]
	if !ok || fi == nil {
		return
	}

	obj := fi.generator()

	f.window.SetContent(obj)
	f.current = f.next
}

// Current returns the key of the FlowItem that is currently displayed.
func (f *Flow) Current() string {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	return f.current
}

// Next returns the key of the FlowItem that will be displayed next.
func (f *Flow) Next() string {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	return f.next
}

// UseStateStr returns a binding.String that is associated with the given key.
// If the key does not exist, a new binding.String is created with the given default value.
// If the key does exist, the default value is ignored.
// The returned binding.String is shared between all FlowItems of the Flow.
func (f *Flow) UseStateStr(key string, def string) binding.String {
	f.stateMut.Lock()
	defer f.stateMut.Unlock()

	if f.strState == nil {
		f.strState = make(map[string]binding.String)
	}

	if _, ok := f.strState[key]; !ok {
		f.strState[key] = binding.BindString(&def)
		_ = f.strState[key].Set(def)
	}

	return f.strState[key]
}

// UseStateInt returns a binding.Int that is associated with the given key.
// If the key does not exist, a new binding.Int is created with the given default value.
// If the key does exist, the default value is ignored.
// The returned binding.Int is shared between all FlowItems of the Flow.
func (f *Flow) UseStateInt(key string, def int) binding.Int {
	f.stateMut.Lock()
	defer f.stateMut.Unlock()

	if f.intState == nil {
		f.intState = make(map[string]binding.Int)
	}

	if _, ok := f.intState[key]; !ok {
		f.intState[key] = binding.BindInt(&def)
		_ = f.intState[key].Set(def)
	}

	return f.intState[key]
}
