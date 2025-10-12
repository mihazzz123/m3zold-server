package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DBConfig
	Logger   LoggerConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Port  int
	Debug bool
}

type HTTPConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	Url             string
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	CheckConnection time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
	Output string
}

type AuthConfig struct {
	JWTSecret string
	TokenTTL  time.Duration
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	port, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	debug := getEnv("APP_DEBUG", "false") == "true"

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "10"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))

	checkConnect, _ := time.ParseDuration(getEnv("DB_CHECK_CONNECTION", "5m"))
	readTimeout, _ := time.ParseDuration(getEnv("HTTP_READ_TIMEOUT", "5s"))
	writeTimeout, _ := time.ParseDuration(getEnv("HTTP_WRITE_TIMEOUT", "10s"))
	idleTimeout, _ := time.ParseDuration(getEnv("HTTP_IDLE_TIMEOUT", "120s"))
	tokenTTL, _ := time.ParseDuration(getEnv("TOKEN_TTL", "15m"))

	// Формируем DB_URL если он не указан
	dbURL := getEnv("DB_URL", "")
	if dbURL == "" {
		dbURL = buildDBURL(
			getEnv("DB_HOST", "localhost"),
			strconv.Itoa(dbPort),
			getEnv("DB_USER", ""),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_NAME", ""),
			getEnv("DB_SSLMODE", "disable"),
		)
	}

	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "clean-arch-go"),
			Env:   getEnv("APP_ENV", "development"),
			Port:  port,
			Debug: debug,
		},
		HTTP: HTTPConfig{
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Database: DBConfig{
			Url:             dbURL,
			Driver:          getEnv("DB_DRIVER", "postgres"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            dbPort,
			User:            getEnv("DB_USER", ""),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", ""),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    maxOpen,
			MaxIdleConns:    maxIdle,
			CheckConnection: checkConnect,
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", ""),
			TokenTTL:  tokenTTL,
		},
	}
}

// buildDBURL строит DSN строку для PostgreSQL
func buildDBURL(host, port, user, password, dbname, sslmode string) string {
	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
