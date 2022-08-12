package db

import (
	"encoding/gob"
	//"reflect"
	"os"
)

type DB struct {
	Config Config
}

type Config struct {
	Suffix, Root string
}


func DefaultConfig() Config {
	return Config {
		Suffix : ".gob",
	}
}

func (db *DB)FormPath(p string) string {
	return db.Config.Root +
		string(os.PathSeparator) +
		p + db.Config.Suffix 
}

func New(cfg Config) *DB {
	var db *DB

	db = &DB{
		Config : cfg,
	}

	os.Mkdir(db.Config.Root, 0755)

	return db
}

func (db *DB)Run() error {
	go func() {
		for {
		}
	}()
	return nil
}

func (db *DB)Set(k string, v any) error {

	f, err := os.OpenFile(db.FormPath(k),
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

func (db *DB)Get(k string, v any) error {
	f, err := os.Open(db.FormPath(k))
	if err != nil {
		return nil
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(v)

	return nil
}

func (db *DB)Stop() error {
	return nil
}

