package configs

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Названия переменных окружения
const (
	// Сервер.
	SERVER_PORT = "SERVER_PORT"
	SERVER_HOST = "SERVER_HOST"

	// База данных.
	DB_HOST = "DB_HOST"
	DB_PORT = "DB_PORT"
	DB_NAME = "DB_NAME"
	DB_USER = "DB_USER"
	DB_PASS = "DB_PASS"
	DB_SSL  = "DB_SSL"

	// Дебаг.
	LOG_LEVEL = "LOG_LEVEL"
)

// Значения по умолчанию.
const (
	// Сервер.
	DEFAULT_SERVER_PORT = "8080"
	DEFAULT_SERVER_HOST = "127.0.0.1"

	// База данных.
	DEFAULT_DB_HOST = "database"
	DEFAULT_DB_PORT = "5432"
	DEFAULT_DB_NAME = "postgres"
	DEFAULT_DB_USER = "postgres"
	DEFAULT_DB_PASS = "postgres"
	DEFAULT_DB_SSL  = "disable"

	// Дебаг.
	DEFAULT_LOG_LEVEL = "debug"
)

// Структура, описывающая конфигурацию сервера.
type ServerConfig struct {
	Host string
	Port string
}

// Структура, описывающая конфигурацию базы данных.
type DataBaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLmode  string
}

// Стрктура, описывающая конфигурацию приложения.
type Config struct {
	server   ServerConfig
	dataBase DataBaseConfig
}

// Конструктор
func New() *Config {
	cfg := &Config{}
	cfg.setLogConfiguration()
	cfg.parseServerEnvs()
	cfg.parseDatabaseEnvs()

	return cfg
}

func (cfg *Config) GetServerCfg() ServerConfig {
	return cfg.server
}

func (cfg *Config) GetDatabaseCfg() DataBaseConfig {
	return cfg.dataBase
}

func (cfg *Config) setLogConfiguration() {
	var logLevel string

	cfg.parseEnv(LOG_LEVEL, DEFAULT_LOG_LEVEL, &logLevel)

	// Парсинг уровня логирования.
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.DebugLevel
	}

	// Установка уровня логирования глобально.
	logrus.SetLevel(ll)
	logrus.Debugf("Уровень логирования: %s", logLevel)
}

// Приватный метод, парсящий переменные окружения для сервера.
func (cfg *Config) parseServerEnvs() {
	cfg.parseEnv(SERVER_PORT, DEFAULT_SERVER_PORT, &cfg.server.Port)
	cfg.parseEnv(SERVER_HOST, DEFAULT_SERVER_HOST, &cfg.server.Host)

	logrus.Debugf("Конфигурация сервера: %#v", cfg.server)
}

// Приватный метод, парсящий переменные окружения для базы данных.
func (cfg *Config) parseDatabaseEnvs() {
	cfg.parseEnv(DB_HOST, DEFAULT_DB_HOST, &cfg.dataBase.Host)
	cfg.parseEnv(DB_PORT, DEFAULT_DB_PORT, &cfg.dataBase.Port)
	cfg.parseEnv(DB_NAME, DEFAULT_DB_NAME, &cfg.dataBase.Name)
	cfg.parseEnv(DB_USER, DEFAULT_DB_USER, &cfg.dataBase.Username)
	cfg.parseEnv(DB_PASS, DEFAULT_DB_PASS, &cfg.dataBase.Password)
	cfg.parseEnv(DB_SSL, DEFAULT_DB_SSL, &cfg.dataBase.SSLmode)

	logrus.Debugf("Конфигурация базы данных: %#v", cfg.dataBase)
}

// Парсинг переменной окружения
func (cfg *Config) parseEnv(env string, defaultValue string, writeTo *string) {
	env, ok := os.LookupEnv(env)
	if !ok {
		env = defaultValue
	}
	*writeTo = env
}
