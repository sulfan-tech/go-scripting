package mysql

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func LoadDatabaseConfig() *DatabaseConfig {
	env := "../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file on mysql config: ", err)
	}

	return &DatabaseConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// SetupDatabase initializes and returns a database connection using the configuration.
func SetupDatabase() (*gorm.DB, error) {
	cfg := LoadDatabaseConfig()
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.DBName + "?parseTime=True"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
