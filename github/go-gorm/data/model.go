package data

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint64         `gorm:"primarykey"`
	CreatedAt time.Time      `type:"datetime"`
	UpdatedAt time.Time      `type:"datetime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
