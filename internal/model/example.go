package model

import "time"

// Example is the sample GORM model used to demonstrate the request flow.
type Example struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updated_at"`
}

// TableName returns the database table name for Example.
func (Example) TableName() string {
	return "examples"
}
