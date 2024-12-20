package main

import (
	"os"
	"time"

	"github.com/InTeam-Russia/go-backend-template/internal/applogger"
	"github.com/InTeam-Russia/go-backend-template/internal/auth"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/InTeam-Russia/go-backend-template/internal/cors"
	"github.com/InTeam-Russia/go-backend-template/internal/db"
	"github.com/InTeam-Russia/go-backend-template/internal/ml"
	"github.com/InTeam-Russia/go-backend-template/internal/polls"
	"github.com/InTeam-Russia/go-backend-template/internal/recommendations"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := config.LoadFromEnv()
	if err != nil {
		panic(err)
	}

	logger := applogger.Create(config.LogLevel)

	pgPool, err := db.CreatePool(config.PostgresUrl, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pgPool.Close()

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	r := gin.New()

	cors.Setup(r, config)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	cookieConfig := auth.DefaultCookieConfig()
	cookieConfig.Secure = config.SessionCookieSecure
	cookieConfig.Domain = config.SessionCookieDomain

	var mlService ml.Service
	if config.MockML {
		mlService = ml.NewMockService(logger)
	} else {
		mlService = ml.NewAPIService(logger, config.MLBaseURL)
	}

	userRepo := user.NewPGRepo(pgPool, logger)
	pollsRepo := polls.NewPGRepo(pgPool, logger)
	sessionRepo := session.NewRedisRepo(redisClient, logger)

	auth.SetupRoutes(r, userRepo, sessionRepo, mlService, logger, cookieConfig)
	polls.SetupRoutes(r, pollsRepo, sessionRepo, mlService, logger)
	recommendations.SetupRoutes(r, sessionRepo, mlService, userRepo, logger)

	err = r.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
