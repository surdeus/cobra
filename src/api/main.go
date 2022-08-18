package db

import (
	"encoding/gob"
	//"reflect"
	"os"
	"errors"
	"strings"
	"github.com/surdeus/goblin/src/tool/ftest"
)

type DB struct {
	Config Config
	Running bool
}

type Config struct {
	EntrySuffix, SubSuffix string
	Root string
	Sep string
	EmptyKey Key
}

type Key = string

var (
	ErrEmptyKey = errors.New("empty key provided")
	ErrNotExist = errors.New("the entry does not exist")
	ErrNoAccess = errors.New("got no access")
	ErrNoParent = errors.New("parent for the sub does not exist")
)

func DefaultConfig() Config {
	return Config {
		EntrySuffix : ".gob",
		SubSuffix : ".sub",
		Sep : string(os.PathSeparator),
		EmptyKey : "",
	}
}

func (db *DB)formPath(s []Key, k Key) string {
	var (
		ret string
	)

	ret = db.Config.Root + db.Config.Sep

	if len(s) > 0 {
		ret += strings.Join([]string(s),
			db.Config.SubSuffix + db.Config.Sep)
		ret += db.Config.SubSuffix
		if !db.KeyIsEmpty(k) {
			ret += db.Config.Sep
		}
	}

	if !db.KeyIsEmpty(k) {
		ret += string(k) + db.Config.EntrySuffix
	}


	return ret
}

func New(cfg Config) *DB {
	var db *DB

	db = &DB{
		Config : cfg,
	}

	os.Mkdir(db.Config.Root, 0700)

	return db
}

func (db *DB) KeyIsEmpty(k Key) bool {
	return k == db.Config.EmptyKey
}

func (db *DB) Run() error {
	db.Running = true
	go func() {
		for db.Running {
		}
	}()
	return nil
}

func (db *DB) setFile(s []Key, v any) error {
	f, err := os.OpenFile(
		db.formPath(s[:len(s)-1], s[len(s)-1]),
		os.O_WRONLY | os.O_CREATE,
		0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	err = enc.Encode(v)
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

func (db *DB) Set(s []Key, v any) error {
	if len(s) == 0 || db.KeyIsEmpty(s[len(s)-1]) {
		return ErrEmptyKey
	}

	return db.setFile(s, v)
}

func (db *DB) Get(s []Key, v any) error {
	if len(s) == 0 || db.KeyIsEmpty(s[len(s)-1]) {
		return ErrEmptyKey
	}

	return db.getFile(s, v)
}

func (db *DB) Sub(s []Key) (int, error) {
	i, err := db.subDir(s)
	if i >= 0 {
			return i, err
	}

	return -1, nil
}

func (db *DB) CheckSub(s []Key) int {
	for i, _ := range s {
		if !ftest.IsDir(db.formPath(s[:i+1], "")) {
			return i
		}
	}

	return -1
}

func (db *DB) subDir(s []Key) (int, error) {
	if len(s) == 0 {
		return -1, nil
	}
	i := db.CheckSub(s[1:])
	if i >= 0 {
		return i, ErrNoParent
	}
	err := os.Mkdir(
		db.formPath(s, db.Config.EmptyKey),
		0700)
	if err != nil {
		return -1, err
	}

	return -1, nil
}

func (db *DB) Stop() error {
	db.Running = false
	return nil
}

