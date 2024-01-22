package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Frank-Mayer/fyneflow"
)

func main() {
	app := app.New()
	window := app.NewWindow("Flow Test")

	// wrap the fyne.Window using a flow
	// the generic type is the type for the flow keys
	flow := fyneflow.NewFlow[string](window)
	defer flow.Close()

	// add an item to the flow with the key "a"
	flow.Set("a", func() fyne.CanvasObject {
		// use flow shared data of key "txt"
		txt := flow.UseStateStr("txt", "test text")
		return container.NewVBox(
			widget.NewLabel("I am A"),
			widget.NewEntryWithData(txt),
			widget.NewButton("Go to B", func() {
				// swich to flow item with key "b"
				if err := flow.GoTo("b"); err != nil {
					panic(err)
				}
			}),
		)
	})

	// add an item to the flow with the key "b"
	flow.Set("b", func() fyne.CanvasObject {
		// use flow shared data of key "txt"
		txt := flow.UseStateStr("txt", "test text")
		return container.NewVBox(
			widget.NewLabel("I am B"),
			widget.NewEntryWithData(txt),
			widget.NewButton("Go to A", func() {
				// swich to flow item with key "a"
				if err := flow.GoTo("a"); err != nil {
					panic(err)
				}
			}),
		)
	})

	// display the window
	window.ShowAndRun()
}
