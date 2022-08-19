package api

import (
	"strings"
)

func (db *DB)formPath(s []Key, k Key) string {
	var (
		ret string
	)
	
	cfg := &(db.config)

	ret = cfg.Root + cfg.Sep
	if len(s) > 0 {
		ret += strings.Join([]string(s),
			cfg.SubSuffix + cfg.Sep)
		ret += cfg.SubSuffix
		if !db.KeyIsEmpty(k) {
			ret += cfg.Sep
		}
	}

	if !db.KeyIsEmpty(k) {
		ret += string(k) + cfg.EntrySuffix
	}


	return ret
}