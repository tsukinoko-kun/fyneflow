package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tsukinoko-kun/fyneflow"
)

const (
	FlowKeyA = iota
	FlowKeyB
)

func main() {
	a := app.New()
	w := a.NewWindow("Flow Test")

	// wrap the fyne.Window using a flow
	// the generic type is the type for the flow keys
	flow := fyneflow.NewFlow(w)
	defer flow.Close()

	// add an item to the flow with the key A
	flow.Set(FlowKeyA, func() fyne.CanvasObject {
		// use flow shared data of key "txt"
		return container.NewVBox(
			widget.NewLabel("I am A"),
			widget.NewButton("Go to B", func() {
				// swich to flow item with key B
				if err := flow.GoTo(FlowKeyB); err != nil {
					panic(err)
				}
			}),
		)
	})

	// add an item to the flow with the key B
	flow.Set(FlowKeyB, func() fyne.CanvasObject {
		// use flow shared data of key "txt"
		return container.NewVBox(
			widget.NewLabel("I am B"),
			widget.NewButton("Go to A", func() {
				// swich to flow item with key A
				if err := flow.GoTo(FlowKeyA); err != nil {
					panic(err)
				}
			}),
		)
	})

	// display the window
	w.ShowAndRun()
}
