package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*SQLClient, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&TMessages{})
	return &SQLClient{db: db}, nil
}

func TestSaveMessage(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err, "Failed to set up test database")

	message := "Hello, Terraform!"
	err = client.SaveMessage(message)
	assert.NoError(t, err, "Failed to save message")

	var result TMessages
	err = client.db.First(&result).Error
	assert.NoError(t, err, "Failed to retrieve saved message")
	assert.Equal(t, message, result.Message, "Saved message does not match")
}

func TestGetMessages(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err, "Failed to set up test database")

	messages := []string{"Message 1", "Message 2", "Message 3"}
	for _, msg := range messages {
		err := client.SaveMessage(msg)
		assert.NoError(t, err, "Failed to save message")
	}

	results, err := client.GetMessages()
	assert.NoError(t, err, "Failed to retrieve messages")

	assert.Equal(t, len(messages), len(results), "Message count mismatch")

	for i, result := range results {
		assert.Equal(t, i+1, result.Sequence, "Sequence number mismatch")
		assert.Equal(t, messages[i], result.Message, "Message content mismatch")
	}
}

func TestClose(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err, "Failed to set up test database")

	err = client.Close()
	assert.NoError(t, err, "Failed to close database connection")
}
