package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/bmj2728/catfetch/pkg/shared/api"
	"github.com/bmj2728/catfetch/pkg/shared/catdb"
	"github.com/bmj2728/catfetch/pkg/shared/ui"
	"go.etcd.io/bbolt"
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

	go func() {
		err := catDB.DB().View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("cats"))
			err := b.ForEachBucket(func(cb []byte) error {
				c := b.Bucket(cb)
				err := c.ForEachBucket(func(vb []byte) error {
					v := c.Bucket(vb)
					err := v.ForEach(func(k, v []byte) error {
						if string(k) == "metadata" {
							fmt.Printf("Key: %s - Value: %s\n", string(k), string(v))
						} else {
							fmt.Printf("Key: %s - Value: %vKB\n", string(k), len(v)/1024)
						}
						return nil
					})
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Default().Println(err)
		}
	}()

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
