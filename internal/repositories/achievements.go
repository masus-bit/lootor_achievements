package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor_achievements/internal/models"
	"time"
)

type AchievementsRepository struct {
	db *gorm.DB
}

func NewAchievementsRepository(db *gorm.DB) *AchievementsRepository {
	return &AchievementsRepository{db: db}
}

func (r *AchievementsRepository) AddOrUpdateAchievement(achievement *models.AchievementsLinks) (
	*models.AchievementsLinks,
	error,
) {
	var existing models.AchievementsLinks
	err := r.db.
		Where(
			"user_login = ? AND achievement_code = ?",
			achievement.UserLogin, achievement.AchievementCode,
		).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {

		var ach models.Achievements
		if err := r.db.Where("code = ?", achievement.AchievementCode).First(&ach).Error; err != nil {
			return nil, err
		}
		achievement.AchievementID = ach.ID
		achievement.Achievement = ach

		if achievement.ID == uuid.Nil {
			achievement.ID = uuid.New()
		}
		if achievement.CreatedAt.IsZero() {
			achievement.CreatedAt = time.Now()
		}

		if err := r.db.Create(achievement).Error; err != nil {
			return nil, err
		}

		return achievement, nil

	} else {
		newValue := achievement.CurrentValueInt

		existing.CurrentValueInt = newValue
		existing.Exp = achievement.Exp

		if achievement.Level > existing.Level {
			existing.Level = achievement.Level
		}

		if achievement.CurrentValueBool {
			existing.CurrentValueBool = true
		}

		if err := r.db.Save(&existing).Error; err != nil {
			return nil, err
		}

		r.db.Preload("Achievement").First(&existing, existing.ID)

		return &existing, nil
	}
}

func (r *AchievementsRepository) GetAchievementsByUserLogin(userLogin string) ([]models.AchievementsLinks, error) {
	var achievements []models.AchievementsLinks

	query := r.db.Model(models.AchievementsLinks{}).
		Where("user_login = ?", userLogin).
		Preload("Achievement").
		Order("COALESCE(level, 0) DESC, created_at ASC").
		Find(&achievements)
	if query.Error != nil {
		return nil, query.Error
	}

	return achievements, nil
}

func (r *AchievementsRepository) GetOneAchievement(userLogin, id string) (*models.AchievementsLinks, error) {
	var achievement *models.AchievementsLinks

	query := r.db.Model(models.AchievementsLinks{}).Where(
		"user_login = ? AND achievement_code = ?",
		userLogin,
		id,
	).Preload("Achievement").First(&achievement)
	if query.Error != nil {
		return nil, query.Error
	}

	return achievement, nil
}

func (r *AchievementsRepository) GetAchievementItem(code string) (*models.Achievements, error) {
	var achievement *models.Achievements
	query := r.db.Model(models.Achievements{}).Where("code = ?", code)
	query.First(&achievement)
	if query.Error != nil {
		return nil, query.Error
	}
	return achievement, nil
}

func (r *AchievementsRepository) GetAllAchievementsItems() ([]models.Achievements, error) {
	var achievements []models.Achievements
	query := r.db.Model(models.Achievements{}).Find(&achievements)
	if query.Error != nil {
		return nil, query.Error
	}
	return achievements, nil
}

func (r *AchievementsRepository) GetAllAchievementsItemsMap() (map[string]models.Achievements, error) {
	var achievements []models.Achievements
	var achievementsMap map[string]models.Achievements
	query := r.db.Model(models.Achievements{}).Find(&achievements)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, a := range achievements {
		achievementsMap[a.ID.String()] = a
	}
	return achievementsMap, nil
}

func (r *AchievementsRepository) GetTotalAchievedAchievements(code string, level int) (int64, error) {
	var count int64
	err := r.db.Model(&models.AchievementsLinks{}).
		Where("achievement_code = ?", code).
		Where("level = ?", level).
		Count(&count).Error

	return count, err
}

func (r *AchievementsRepository) CreateMultipleRecords(achievements []models.AchievementsLinks) (
	[]models.AchievementsLinks,
	error,
) {
	if len(achievements) == 0 {
		return nil, errors.New("empty ach list")
	}

	err := r.db.Create(&achievements)
	if err.Error != nil {
		return nil, err.Error
	}

	return achievements, nil
}
