package model

import (
	"time"
)

type Task struct {
	Id          int64 // `gorm:"primary_key AUTO_INCREMENT"`
	Title       *string
	Description *string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
