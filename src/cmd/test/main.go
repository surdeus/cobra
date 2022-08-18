package main

import (
	garbo "github.com/surdeus/cobra/src/api"
	"log"
	"fmt"
)

type Q struct {
	X, Y int
}

func main() {
	var (
		err error
		q, q1 Q
	)

	cfg := garbo.DefaultConfig()
	cfg.Root = "testdb"
	db := garbo.New(cfg)
	db.Run()
	defer db.Stop()

	err = db.Set([]garbo.Key{"key1"}, Q{1, 2})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Get([]garbo.Key{"key1"}, &q)
	if err != nil {
		log.Fatal(err)
	}

	db.Sub([]garbo.Key{"key1"})
	db.Set([]garbo.Key{"key1", "key2"}, Q{3, 4})
	db.Get([]garbo.Key{"key1", "key2"}, &q1)

	fmt.Printf("%d %d\n", q.X, q.Y)
	fmt.Printf("%d %d\n", q1.X, q1.Y)
}
