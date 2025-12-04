package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"log"
	"lootor_achievements/gen/go/achievements"
	"lootor_achievements/internal/config"
	"lootor_achievements/internal/grpcserver"
	"lootor_achievements/internal/models"
	"lootor_achievements/internal/repositories"
	"lootor_achievements/internal/services"
	"lootor_achievements/pkg/database"

	"net"
	"time"
)

type App struct {
	Echo *echo.Echo
}

func NewEchoApp(cfg *config.Config) (*App, error) {
	_ = godotenv.Load()

	e := echo.New()

	e.GET("/swagger/achievements*", echoSwagger.WrapHandler)
	e.Server.MaxHeaderBytes = 1 << 20
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://dev.lootor.me",
			"https://lootor.me",
			"https://www.lootor.me",
			"https://www.dev.lootor.me",
			"http://localhost:3000",
			"http://localhost:4173",
			"*",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// db init
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		if err = database.Reconnect(&cfg.Database); err != nil {
			log.Printf("Reconnection failed: %v", err)
		}
	}

	achievementsRepo := repositories.NewAchievementsRepository(db)

	err = db.AutoMigrate(&models.Achievements{}, &models.AchievementsLinks{})
	if err != nil {
		return nil, err
	}

	achievementsService := services.NewAchievementsService(achievementsRepo)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	grpcServer := grpcserver.NewAchievementsService(achievementsService)

	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	s := grpc.NewServer()
	achievements.RegisterAchievementsServiceServer(s, grpcServer)

	log.Println("Starting gRPC server on port 50055")
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}

	return &App{Echo: e}, nil
}
