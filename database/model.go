package database

import "gorm.io/gorm"

// TMessages represents a database table that contains sequence and message
type TMessages struct {
	gorm.Model
	Sequence int
	Message  string `gorm:"type:varchar(255)"`
}

// TableName is a function to return the table name
func (m *TMessages) TableName() string {
	return "messages"
}
