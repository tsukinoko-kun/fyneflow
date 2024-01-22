package fyneflow

import (
	"errors"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type FlowItemGenerator[K comparable] func() fyne.CanvasObject

type FlowItem[K comparable] struct {
	generator FlowItemGenerator[K]
}

type Flow[K comparable] struct {
	window  fyne.Window
	flowMut sync.Mutex
	current K
	next    K
	items   map[K]*FlowItem[K]

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
func NewFlow[K comparable](w fyne.Window) *Flow[K] {
	f := new(Flow[K])
	f.window = w
	f.flowMut = sync.Mutex{}
	f.items = make(map[K]*FlowItem[K])
	go f.loop()
	return f
}

// Close closes the Flow.
// This will not close the associated fyne.Window!
func (f *Flow[K]) Close() {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	f.close = true
	for k := range f.items {
		delete(f.items, k)
	}
}

func (f *Flow[K]) loop() {
	for !f.close {
		if f.next != f.current {
			f.apply()
		}
	}
}

// Set adds a FlowItem to the Flow.
// If the key already exists, the FlowItem will be overwritten.
func (f *Flow[K]) Set(key K, generator FlowItemGenerator[K]) *FlowItem[K] {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	fi := new(FlowItem[K])
	fi.generator = generator
	apply := len(f.items) == 0
	f.items[key] = fi
	if apply {
		f.next = key
	}
	return fi
}

// GoTo sets the content of the window to the content of the FlowItem associated with the given key.
// If the key is not found, an error is returned.
func (f *Flow[K]) GoTo(key K) error {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	if f.current == key {
		return nil
	}

	if _, ok := f.items[key]; ok {
		f.next = key
		return nil
	} else {
		return errors.New("flow: key not found")
	}
}

func (f *Flow[K]) apply() {
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
func (f *Flow[K]) Current() K {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	return f.current
}

// Next returns the key of the FlowItem that will be displayed next.
func (f *Flow[K]) Next() K {
	f.flowMut.Lock()
	defer f.flowMut.Unlock()

	return f.next
}

// UseStateStr returns a binding.String that is associated with the given key.
// If the key does not exist, a new binding.String is created with the given default value.
// If the key does exist, the default value is ignored.
// The returned binding.String is shared between all FlowItems of the Flow.
func (f *Flow[K]) UseStateStr(key string, def string) binding.String {
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
func (f *Flow[K]) UseStateInt(key string, def int) binding.Int {
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
