package services

import (
	"fmt"
	"github.com/google/uuid"
	"lootor_achievements/gen/go/achievements"
	"lootor_achievements/internal/models"
	"lootor_achievements/internal/repositories"
	"time"
)

type AchievementsService struct {
	achievementsRepo *repositories.AchievementsRepository
}

func NewAchievementsService(achievementsRepo *repositories.AchievementsRepository) *AchievementsService {
	return &AchievementsService{
		achievementsRepo: achievementsRepo,
	}
}

func (s *AchievementsService) AddOrUpdateAchievement(dto *achievements.AddOrUpdateAchievementRequest) (
	*achievements.AddOrUpdateAchievementResponse,
	error,
) {
	achievementItem, err := s.achievementsRepo.GetAchievementItem(dto.GetCode())
	if err != nil {
		return nil, fmt.Errorf("oшибка получения ориг ачивы: %v", err)
	}
	curBool := false

	if dto.GetLevel() > 0 {
		curBool = true
	}

	achievement := models.AchievementsLinks{
		AchievementID:    achievementItem.ID,
		AchievementCode:  dto.GetCode(),
		UserLogin:        dto.GetUserLogin(),
		Achievement:      *achievementItem,
		Exp:              int(dto.GetXp()),
		Level:            int(dto.GetLevel()),
		CurrentValueInt:  int(dto.GetValue()),
		CurrentValueBool: curBool,
	}

	resultAchievement, err := s.achievementsRepo.AddOrUpdateAchievement(&achievement)
	if err != nil {
		return nil, fmt.Errorf("ошибка добавления/обновления ачивы: %v", err)
	}

	convertedAchievement := s.convertAchievementToProto(resultAchievement, &resultAchievement.Achievement)

	result := &achievements.AddOrUpdateAchievementResponse{
		AchievedAchievement: convertedAchievement,
	}

	return result, nil
}

func (s *AchievementsService) GetAchievedUserAchievements(dto *achievements.GetUserAchievementsRequest) (
	*achievements.GetUserAchievementsResponse,
	error,
) {
	var allUserAchievements []*achievements.AchievedAchievement
	userAchievements, err := s.achievementsRepo.GetAchievementsByUserLogin(dto.GetUserLogin())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ачив пользователя: %v", err)
	}

	for _, a := range userAchievements {
		temp := s.convertAchievementToProto(&a, &a.Achievement)
		allUserAchievements = append(allUserAchievements, temp)
	}

	return &achievements.GetUserAchievementsResponse{AchievedAchievements: allUserAchievements}, nil

}

func (s *AchievementsService) GetAllUserAchievements(dto *achievements.GetUserAchievementsRequest) (
	*achievements.GetUserAchievementsResponse,
	error,
) {
	var allUserAchievements []*achievements.AchievedAchievement
	userAchievements, err := s.achievementsRepo.GetAchievementsByUserLogin(dto.GetUserLogin())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ачив пользователя: %v", err)
	}

	allAchievements, err := s.achievementsRepo.GetAllAchievementsItems()
	if err != nil {
		return nil, err
	}

	achievementsUserMap := make(map[string]models.AchievementsLinks)

	for _, a := range userAchievements {
		achievementsUserMap[a.AchievementID.String()] = a
	}

	for _, a := range allAchievements {
		achievementID := a.ID.String()
		var temp *achievements.AchievedAchievement
		if userAchievement, exists := achievementsUserMap[achievementID]; exists {
			temp = s.convertAchievementToProto(&userAchievement, &a)
		} else {
			temp = s.convertAchievementToProto(
				&models.AchievementsLinks{
					CreatedAt:        time.Time{},
					UpdatedAt:        time.Time{},
					ID:               uuid.UUID{},
					AchievementID:    uuid.UUID{},
					AchievementCode:  "",
					UserLogin:        "",
					Exp:              0,
					Level:            0,
					CurrentValueBool: false,
					CurrentValueInt:  0,
				}, &a,
			)
		}
		allUserAchievements = append(allUserAchievements, temp)
	}

	return &achievements.GetUserAchievementsResponse{AchievedAchievements: allUserAchievements}, nil

}

func (s *AchievementsService) GetOneAchievement(dto *achievements.GetOneAchievementRequest) (
	*achievements.GetOneAchievementResponse,
	error,
) {
	achievement, err := s.achievementsRepo.GetOneAchievement(dto.GetUserLogin(), dto.GetId())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ачивы: %v", err)
	}

	result := s.convertAchievementToProto(achievement, &achievement.Achievement)

	return &achievements.GetOneAchievementResponse{AchievedAchievement: result}, nil
}

func (s *AchievementsService) GetAllAchievementsItems() (*achievements.GetAllAchievementsItemsResponse, error) {
	achievementsItems, err := s.achievementsRepo.GetAllAchievementsItems()
	var resultAchievements []*achievements.AchievementItem
	if err != nil {
		return nil, fmt.Errorf("ошибка олучения всех ачив: %v", err)
	}

	for _, a := range achievementsItems {
		resultAchievements = append(
			resultAchievements, &achievements.AchievementItem{
				Id:   a.ID.String(),
				Code: a.Code,
			},
		)
	}
	return &achievements.GetAllAchievementsItemsResponse{Achievements: resultAchievements}, nil
}

func (s *AchievementsService) AddZeroAchievements(dto *achievements.GetUserAchievementsRequest) (
	*achievements.GetUserAchievementsRequest,
	error,
) {
	achievementsItems, err := s.achievementsRepo.GetAllAchievementsItems()
	if err != nil {
		return nil, err
	}
	var achievementsSlice []models.AchievementsLinks
	for _, a := range achievementsItems {
		achievementsSlice = append(
			achievementsSlice, models.AchievementsLinks{
				AchievementID:    a.ID,
				AchievementCode:  a.Code,
				UserLogin:        dto.UserLogin,
				Achievement:      a,
				Exp:              0,
				Level:            0,
				CurrentValueInt:  0,
				CurrentValueBool: false,
			},
		)
	}
	_, err = s.achievementsRepo.CreateMultipleRecords(achievementsSlice)
	if err != nil {
		return nil, err
	}
	return &achievements.GetUserAchievementsRequest{UserLogin: dto.UserLogin}, nil
}

func (s *AchievementsService) convertAchievementToProto(
	link *models.AchievementsLinks,
	achievement *models.Achievements,
) *achievements.AchievedAchievement {
	var totalAchieved int64
	var linkLevel int

	if link.Level == 0 {
		linkLevel = 1
	}
	totalAchieved, err := s.achievementsRepo.GetTotalAchievedAchievements(link.AchievementCode, linkLevel)
	if err != nil {
		totalAchieved = 0
	}

	var id string
	if link.ID == uuid.Nil {
		id = ""
	} else {
		id = link.ID.String()
	}
	result := &achievements.AchievedAchievement{
		Id: id,
		Achievement: &achievements.AchievementItem{
			Id:   achievement.ID.String(),
			Code: achievement.Code,
		},
		UserLogin:        link.UserLogin,
		Xp:               int64(link.Exp),
		Level:            int64(link.Level),
		TotalAchieved:    totalAchieved,
		CreatedAt:        link.CreatedAt.String(),
		UpdatedAt:        link.UpdatedAt.String(),
		CurrentValueBool: link.CurrentValueBool,
		CurrentValueInt:  int64(link.CurrentValueInt),
	}
	return result
}
