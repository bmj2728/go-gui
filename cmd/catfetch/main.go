package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/bmj2728/catfetch/pkg/shared/api"
	"github.com/bmj2728/catfetch/pkg/shared/catdb"
	"github.com/bmj2728/catfetch/pkg/shared/ui"
)

func main() {

	catDB, err := catdb.OpenDB("cats.db")
	if err != nil {
		log.Default().Println(err)
	}
	defer func(catDB *catdb.CatDB) {
		err := catDB.Close()
		if err != nil {
			log.Default().Println(err)
		}
	}(catDB)

	// Fetch available tags
	go func() {
		api.FetchCAASTags(30 * time.Second)
	}()

	// Make a window and run the loop
	go func() {
		// Create window
		w := new(app.Window)
		w.Option(app.Title("CatFetch"), app.Size(unit.Dp(400), unit.Dp(500)))

		if err := ui.Run(w, catDB); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()

}
