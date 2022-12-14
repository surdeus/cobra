package api

import (
	"os"
	"io/fs"
	"fmt"
	"errors"
	"strings"
)

type Key = string

type DB struct {
	config DBConfig
	running bool
}

type DBConfig struct {
	EntrySuffix, SubSuffix string
	Root string
	Sep string
	EmptyKey Key
	EntryPerm, SubPerm fs.FileMode
}

var (
	ErrEmptyKey = errors.New("empty key provided")
	ErrNotExist = errors.New("the entry or sub does not exist")
	ErrNoAccess = errors.New("got no access")
	ErrNoParent = errors.New("parent for the sub does not exist")
)

func DefaultConfig() DBConfig {
	return DBConfig {
		EntrySuffix : ".gob",
		SubSuffix : ".sub",
		Sep : string(os.PathSeparator),
		EmptyKey : "",
		EntryPerm : 0700,
		SubPerm : 0700,
	}
}

func New(cfg DBConfig) *DB {
	var db *DB

	db = &DB{
		config : cfg,
	}

	os.Mkdir(db.config.Root, db.config.SubPerm)

	return db
}

func (db *DB) Run() error {
	db.running = true
	go func() {
		for db.running {
			/* Currently does nothing.
				further will do caching and stuff. */
		}
	}()
	return nil
}

func (db *DB)Path(keys ...string) []Key {
	return []Key(keys)
}

func (db *DB)SPath(keys []string) []Key {
	return []Key(keys)
}

func (db *DB) KeyIsEmpty(k Key) bool {
	return k == db.config.EmptyKey
}

func (db *DB) Set(s []Key, v any) error {
	if len(s) == 0 || db.KeyIsEmpty(s[len(s)-1]) {
		return ErrEmptyKey
	}

	err := db.setFile(s, v)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB)Has(s []Key) bool {
	var (
		i int
	)

	i = db.checkSubDir(s)
	if i >= 0 {
		return false
	}

	return true
}

func (db *DB)HasNot(s []Key) int {
	var (
		i int
	)

	i = db.checkSubDir(s)

	return i
}

func (db *DB)List(s []Key) (chan Key, error) {
	var (
		fname, suf string
		has bool
	)

	has = db.Has(s)
	if !has {
		return nil, ErrNotExist
	}

	suf = db.config.EntrySuffix
	c := make(chan Key)
	go func() {
		files, _ := os.ReadDir(db.formPath(s, ""))
		for _, f := range files {
			fname = f.Name()
			if strings.HasSuffix(fname, suf) {
				c <- strings.TrimSuffix(fname, suf)
			}
		}
		close(c)
	}()

	return c, nil
}

func (db *DB) Get(s []Key, v any) error {
	if len(s) == 0 || db.KeyIsEmpty(s[len(s)-1]) {
		return ErrEmptyKey
	}

	err := db.getFile(s, v)
	if err == os.ErrNotExist {
		err = fmt.Errorf("get %v: %w",
			s, ErrNotExist)
	}

	return nil
}

func (db *DB) Stop() error {
	db.running = false
	return nil
}

