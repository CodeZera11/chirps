package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

func NewDB(path string) (*DB, error){
	db := &DB{
		path: path,
		mux: &sync.RWMutex{},
	}

	err := db.ensureDB()

	return db, err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist){
		return db.createDB()
	}

	return nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}

	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)

	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}

	dat, err := os.ReadFile(db.path)

	if err != nil {
		return dbStructure, err
	}

	err = json.Unmarshal(dat, &dbStructure)

	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}