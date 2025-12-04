package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AchievementsLinks struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	ID               uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AchievementID    uuid.UUID    `gorm:"type:uuid;primaryKey" json:"achievementId"`
	AchievementCode  string       `gorm:"type:varchar(50);primaryKey" json:"achievementCode"`
	UserLogin        string       `gorm:"type:varchar(50);primaryKey" json:"userLogin"`
	Achievement      Achievements `gorm:"foreignKey:AchievementID" json:"achievement"`
	Exp              int          `json:"exp"`
	Level            int          `json:"level"`
	CurrentValueInt  int          `json:"currentValueInt"`
	CurrentValueBool bool         `json:"currentValueBool"`
}
