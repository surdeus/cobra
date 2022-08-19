package api

import (
	"os"
	"errors"
	"encoding/gob"
	"github.com/surdeus/goblin/src/tool/ftest"
)

func (db *DB) checkSubDir(s []Key) int {
	for i, _ := range s {
		if !ftest.IsDir(db.formPath(s[:i+1], "")) {
			return i
		}
	}

	return -1
}

func (db *DB) mkSubDir(s []Key) (int, error) {
	// The root always exists.
	if len(s) == 0 {
		return -1, nil
	}

	i := db.checkSubDir(s[:len(s)-1])
	if i >= 0 {
		return i, os.ErrNotExist
	}

	pth := db.formPath(s, db.config.EmptyKey)
	err := os.Mkdir(pth, db.config.SubPerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return -1, err
	}

	return -1, nil
}

func (db *DB) setFile(s []Key, v any) error {
	pth := db.formPath(s[:len(s)-1], s[len(s)-1])
	f, err := os.OpenFile(
		pth,
		os.O_WRONLY | os.O_CREATE,
		db.config.EntryPerm)
	if err != nil && err != os.ErrExist {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	err = enc.Encode(v)
	if err != nil {
		return err
	}

	_, err = db.mkSubDir(s)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getFile(s []Key, v any) error {

	f, err := os.Open(
		db.formPath(s[:len(s)-1], s[len(s)-1]) )
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(v)

	return nil
}