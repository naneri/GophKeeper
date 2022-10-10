package record

import (
	"gorm.io/gorm"
	"time"
)

type Record struct {
	ID        uint32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Type      string
	Data      string
	Path      string
	UserId    uint32
}
