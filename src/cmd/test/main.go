package main

import (
	cobra "github.com/surdeus/cobra/src/api"
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

	cfg := cobra.DefaultConfig()
	cfg.Root = "testdb"
	db := cobra.New(cfg)
	db.Run()
	defer db.Stop()

	err = db.Set(db.Path("key1"), Q{1, 2})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Get(db.Path("key1"), &q)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Set(db.Path("key1", "key2"), Q{3, 4})
	if err != nil {
		log.Fatal(err)
	}

	db.Set(db.Path("key1", "key4"), Q{5, 6})
	db.Set(db.Path("key1", "key5"), Q{7, 8})
	err = db.Get(db.SPath([]string{"key1", "key2"}), &q1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db.Has(db.Path("key1")))
	fmt.Println(db.Has(db.Path("key1", "key2")))
	fmt.Println(db.Has(db.Path("key3")))
	fmt.Println(db.HasNot(db.Path("key3")))
	fmt.Println(db.HasNot(db.Path("key1", "key3")))
	fmt.Printf("%d %d\n", q.X, q.Y)
	fmt.Printf("%d %d\n", q1.X, q1.Y)
	keys, err := db.List(db.Path("key1"))
	if err != nil {
		log.Fatal(err)
	}
	for v := range keys {
		fmt.Println(v)
	}
}
