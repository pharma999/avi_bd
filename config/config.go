package config

import "fmt"

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	GameTickMs       int
	CountdownSeconds int
	CrashMin         float64
	CrashMax         float64
	InitialBalance   float64
}

func LoadConfig() *Config {
	if err := LoadEnv(); err != nil {
		panic(err)
	}

	cfg := &Config{
		AppPort:          GetEnv("APP_PORT", "8080"),
		DBHost:           GetEnv("DB_HOST", "localhost"),
		DBPort:           GetEnv("DB_PORT", "5432"),
		DBUser:           GetEnv("DB_USER", "postgres"),
		DBPassword:       GetEnv("DB_PASSWORD", "postgres"),
		DBName:           GetEnv("DB_NAME", "aviator_db"),
		JWTSecret:        GetEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		GameTickMs:       GetEnvInt("GAME_TICK_MS", 100),
		CountdownSeconds: GetEnvInt("COUNTDOWN_SECONDS", 5),
		CrashMin:         1.5,
		CrashMax:         100.0,
		InitialBalance:   1000.0,
	}

	fmt.Println("DBHost:", cfg.DBHost)
	fmt.Println("DBPort:", cfg.DBPort)
	fmt.Println("DBUser:", cfg.DBUser)
	fmt.Println("DBName:", cfg.DBName)
	fmt.Println("DBPassword loaded:", cfg.DBPassword)

	return cfg
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}
