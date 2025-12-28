package googlesheetreader

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	TG_TOKEN      string
	UNU_TOKEN     string
	UNU_URL       string
	DB_HOST       string
	DB_PASSWORD   string
	DB_DB         int
	SPREADSHEETID string
	SHEETLIST     string
}

var (
	cfg  *Config
	once sync.Once
)

// Load загружает конфигурацию один раз (singleton)
func Load() *Config {
	once.Do(func() {
		cfg = &Config{}

		// Пытаемся загрузить .env файл
		loadEnvFile()

		// Загружаем значения из переменных окружения
		cfg.TG_TOKEN = getEnv("TG_TOKEN", "")
		cfg.UNU_TOKEN = getEnv("UNU_TOKEN", "")
		cfg.UNU_URL = getEnv("UNU_URL", "https://unu.im/api")
		cfg.DB_HOST = getEnv("DB_HOST", "localhost:6379")
		cfg.DB_PASSWORD = getEnv("DB_PASSWORD", "")
		cfg.SPREADSHEETID = getEnv("SPREADSHEETID", "")
		cfg.SHEETLIST = getEnv("SHEETLIST", "")

		// Обрабатываем числовое значение
		dbStr := getEnv("DB_DB", "0")
		if dbNum, err := strconv.Atoi(dbStr); err == nil {
			cfg.DB_DB = dbNum
		} else {
			log.Printf("Invalid DB_DB value: %v, using default 0", err)
			cfg.DB_DB = 0
		}

		// Валидация обязательных полей
		if cfg.TG_TOKEN == "" {
			log.Fatal("TG_TOKEN is required")
		}
		if cfg.UNU_TOKEN == "" {
			log.Fatal("UNU_TOKEN is required")
		}

		log.Printf("Config loaded successfully. DB_HOST: %s", cfg.DB_HOST)
	})

	return cfg
}

func loadEnvFile() {
	// Получаем текущую директорию пакета config
	_, filename, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(filename)

	envPaths := []string{
		filepath.Join(configDir, ".env"),             // рядом с config.go
		filepath.Join(configDir, "..", ".env"),       // на уровень выше (в api/)
		filepath.Join(configDir, "..", "..", ".env"), // в корне go/
		".env",      // текущая рабочая директория
		"./.env",    // текущая директория
		"/app/.env",
		"../../config/.env", // для Docker
	}

	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Printf("Error loading .env file %s: %v", path, err)
			} else {
				log.Printf("Loaded environment from: %s", path)
				return
			}
		}
	}

	log.Println("No .env file found, using only environment variables")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Get() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded. Call config.Load() first")
	}
	return cfg
}

func TestReader(t *testing.T) {
	cfg := Load()
	//cfg := config.Get()
	spreadsheetId := cfg.SPREADSHEETID
	testRowPositive := []string{"5", "6"}
	testRowNegative := []string{"679"}

	for _, value := range testRowPositive {
		resp, err := Reader(spreadsheetId, "BOT", value)
		if assert.NoError(t, err) {
			require.Equal(t, "Реми 2гис", resp.Values[0][0])
		}
	}
	for _, value := range testRowNegative {
		_, err := Reader(spreadsheetId, "BOT", value)
		if !assert.NoError(t, err) {
			require.EqualValues(t, models.ErrorGoogleSheet, err)
		}

	}

}
