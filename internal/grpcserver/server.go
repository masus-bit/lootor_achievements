package grpcserver

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"lootor_achievements/gen/go/achievements"
	"lootor_achievements/internal/services"
)

type AchievementsService struct {
	achievements.UnimplementedAchievementsServiceServer
	service *services.AchievementsService
}

func NewAchievementsService(svc *services.AchievementsService) *AchievementsService {
	return &AchievementsService{service: svc}
}

func (s *AchievementsService) AddOrUpdateAchievement(
	ctx context.Context,
	req *achievements.AddOrUpdateAchievementRequest,
) (*achievements.AddOrUpdateAchievementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.AddOrUpdateAchievement(req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *AchievementsService) GetUserAchievements(
	ctx context.Context,
	req *achievements.GetUserAchievementsRequest,
) (*achievements.GetUserAchievementsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.GetAllUserAchievements(req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *AchievementsService) GetAchievedUserAchievements(
	ctx context.Context,
	req *achievements.GetUserAchievementsRequest,
) (*achievements.GetUserAchievementsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.GetAchievedUserAchievements(req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *AchievementsService) GetOneAchievement(
	ctx context.Context,
	req *achievements.GetOneAchievementRequest,
) (*achievements.GetOneAchievementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.GetOneAchievement(req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *AchievementsService) GetAllAchievementsItems(
	ctx context.Context,
	req *achievements.GetAllAchievementsItemsRequest,
) (*achievements.GetAllAchievementsItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.GetAllAchievementsItems()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (s *AchievementsService) AddZeroAchievements(
	ctx context.Context,
	req *achievements.GetUserAchievementsRequest,
) (*achievements.GetUserAchievementsRequest, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	resp, err := s.service.AddZeroAchievements(req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}
