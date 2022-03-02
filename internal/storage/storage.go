package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zavyalov-den/url-shortener/internal/config"
)

type DB struct {
	db    map[string]string
	debug bool
}

func (db *DB) Save(key, value string) {
	db.db[key] = value
	if !db.debug {
		db.saveToFile()
	}
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

	defer file.Close()

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(config.Conf.FileStoragePath, flag, 0777)
	if err != nil {
		if _, createErr := os.Create(config.Conf.FileStoragePath); createErr != nil {
			panic("can't read or create storage file.")
		}
	}

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

	data, err := os.ReadFile(config.Conf.FileStoragePath)
	if err != nil {
		if _, createErr := os.Create(config.Conf.FileStoragePath); createErr != nil {
			panic("can't read or create storage file.")
		}

		panic(err)
	}

	err = json.Unmarshal(data, &storage)
	if err != nil {
		fmt.Println(err.Error())
	}

	db.db = storage
}

func NewStorage(debug bool) *DB {
	storage := &DB{
		db: make(map[string]string),
	}
	if !debug {
		storage.readFromFile()
	}

	storage.debug = debug

	return storage
}
