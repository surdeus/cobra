package main

import (
	garbo "github.com/surdeus/garbo/src/api"
	"log"
	"fmt"
)

type Q struct {
	X, Y int
}

func main() {
	var (
		err error
		q Q
	)

	cfg := garbo.DefaultConfig()
	cfg.Root = "testdb"
	db := garbo.New(cfg)
	db.Run()
	defer db.Stop()

	err = db.Set("key1", Q{1, 2})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Get("key1", &q)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d %d\n", q.X, q.Y)
}
