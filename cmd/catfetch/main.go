package main

import (
	"log"
	"os"

	"go-gui/pkg/shared/ui"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {

	go func() {
		// Create window
		w := new(app.Window)
		w.Option(app.Title("Cat Image Viewer"), app.Size(unit.Dp(400), unit.Dp(500)))

		if err := ui.Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
