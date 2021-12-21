package config

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
}

var Conf = &config{}

// Раньше здесь была вся логика по обработке переменных среды и флагам,
// но переехала в main, т.к. при запуске тестов через go test ./... возникала ошибка
