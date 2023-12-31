package liteDB

import (
	"log"
)

func (db *LiteDB) Lock(name string, ignorelock bool) {
	if !ignorelock {
		db.Mu.Lock()
		log.Println("Got ", name, " lock")
	} else {
		log.Println("Ignoring ", name, " lock")
	}
}

func (db *LiteDB) Unlock(name string, ignorelock bool) {
	if !ignorelock {
		db.Mu.Unlock()
		log.Println("Rem ", name, " lock")
	}
}
