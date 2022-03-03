package storage

import (
	"encoding/json"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"os"
)

type DB struct {
	db map[string]string
}

func (db *DB) Save(key, value string) {
	db.db[key] = value
	db.saveToFile()
}

func (db *DB) Get(key string) (string, error) {
	longURL, ok := db.db[key]
	if !ok {
		return "", fmt.Errorf("failed to get an URL")
	}

	return longURL, nil
}

func (db *DB) saveToFile() {
	var file *os.File

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(config.C.FileStoragePath, flag, 0755)
	if err != nil {
		panic("failed to open storage file")
	}

	defer file.Close()

	data, err := json.Marshal(db.db)
	if err != nil {
		panic("failed to read from DB")
	}

	_, err = file.Write(data)
	if err != nil {
		panic("failed to write to DB")
	}
}

func (db *DB) readFromFile() {
	storage := make(map[string]string)

	data, err := os.ReadFile(config.C.FileStoragePath)
	if err != nil {
		if _, createErr := os.Create(config.C.FileStoragePath); createErr != nil {
			panic("can't read or create storage file.")
		}
	}

	err = json.Unmarshal(data, &storage)
	if err != nil {
		// who cares
		fmt.Println(err)
	}

	db.db = storage
}

func NewStorage(fileStorage bool) *DB {
	storage := &DB{
		db: make(map[string]string),
	}
	if fileStorage {
		storage.readFromFile()
	}

	return storage
}
