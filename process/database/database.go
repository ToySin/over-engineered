package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config is a struct that contains database configuration
type Config struct {
	DatabaseType string `envconfig:"DATABASE_TYPE" envDefault:"mysql"`
	Host         string `envconfig:"DATABASE_HOST" envDefault:"localhost"`
	Port         string `envconfig:"DATABASE_PORT" envDefault:"3306"`
	User         string `envconfig:"DATABASE_USER" envDefault:"root"`
	Password     string `envconfig:"DATABASE_PASSWORD" envDefault:"password"`
	DBName       string `envconfig:"DATABASE_NAME" envDefault:"test"`
}

func createDBConnection(cfg *Config) (*gorm.DB, error) {
	var dsn string
	switch cfg.DatabaseType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}
}

// SQLClient is a struct that provides query functions
type SQLClient struct {
	db *gorm.DB
}

// NewSQLClient is a function to create a new SQLClient with the given configuration
func NewSQLClient(cfg *Config) (*SQLClient, error) {
	db, err := createDBConnection(cfg)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&TMessages{})
	return &SQLClient{db: db}, nil
}

// Close is a function to close the database connection
func (s *SQLClient) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// SaveMessage is a function to save a message to the database
func (s *SQLClient) SaveMessage(message string) error {
	// Get the last sequence number
	var lastMessage TMessages
	if err := s.db.Last(&lastMessage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			lastMessage.Sequence = 0
		} else {
			return err
		}
	}

	// Increment the sequence number
	sequence := lastMessage.Sequence + 1
	return s.db.Create(&TMessages{Sequence: sequence, Message: message}).Error
}

// GetMessages is a function to get all messages from the database
func (s *SQLClient) GetMessages() ([]TMessages, error) {
	var messages []TMessages
	if err := s.db.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
