package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Achievements struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code string    `gorm:"type:varchar(50);unique" json:"code"`
}
