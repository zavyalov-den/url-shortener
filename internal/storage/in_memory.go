package storage

import (
	"encoding/json"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"os"
)

type Storage interface {
	GetURL(k string) (string, error)
	GetUserURLs(id int) []UserURL
	SaveURL(id int, url UserURL) error
}

type BasicStorage struct {
	db       map[string]string
	userURLs map[int][]UserURL
}

var _ Storage = (*BasicStorage)(nil)

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (db *BasicStorage) GetURL(key string) (string, error) {
	longURL, ok := db.db[key]
	if !ok {
		return "", fmt.Errorf("failed to get an URL")
	}

	return longURL, nil
}

func (db *BasicStorage) saveToFile() {
	var file *os.File

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(config.Config.FileStoragePath, flag, 0755)
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

func (db *BasicStorage) readFromFile() {
	storage := make(map[string]string)

	data, err := os.ReadFile(config.Config.FileStoragePath)
	if err != nil {
		if _, createErr := os.Create(config.Config.FileStoragePath); createErr != nil {
			panic("can't read or create storage file.")
		}
	}

	err = json.Unmarshal(data, &storage)
	if err != nil {
		fmt.Println(err)
	}

	db.db = storage
}

func (db *BasicStorage) GetUserURLs(id int) []UserURL {
	fmt.Println(id, db.userURLs)
	return db.userURLs[id]
}

func (db *BasicStorage) SaveURL(userID int, url UserURL) error {
	db.db[url.ShortURL] = url.OriginalURL
	db.saveToFile()

	urls := db.userURLs[userID]

	for _, v := range urls {
		if v.ShortURL == url.ShortURL {
			return nil
		}
	}

	urls = append(urls, url)
	db.userURLs[userID] = urls
	return nil
}

func NewStorage() Storage {
	if config.Config.DatabaseDSN != "" {
		fmt.Println("using DB: ", config.Config.DatabaseDSN)
		db := NewDB()
		db.InitDB()
		return db

	} else {
		fmt.Println("in memory storage")
		storage := &BasicStorage{
			db:       make(map[string]string),
			userURLs: make(map[int][]UserURL),
		}
		if config.Config.FileStoragePath != "" {
			storage.readFromFile()
		}

		return storage
	}
}
