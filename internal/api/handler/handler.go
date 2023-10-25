package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "threat-monitoring/docs"
	"threat-monitoring/internal/api"
	"threat-monitoring/internal/api/repository"
	"threat-monitoring/internal/models"
	minio "threat-monitoring/internal/pkg/minio"
	redis "threat-monitoring/internal/pkg/redis"
)

type Handler struct {
	repo   api.Repo
	minio  minio.Client
	redis  redis.Client
	logger *logrus.Entry
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("./config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}

func NewHandler(logger *logrus.Logger) *Handler {
	vp := viper.New()
	if err := initConfig(vp); err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	repo, err := repository.NewRepository(logger, vp)
	if err != nil {
		logger.Error(err)
	}

	minioConfig := minio.InitConfig(vp)

	minioClient, err := minio.NewMinioClient(context.Background(), minioConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	redisConfig := redis.InitRedisConfig(vp, logger)

	redisClient, err := redis.NewRedisClient(context.Background(), redisConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	return &Handler{
		repo:   repo,
		minio:  minioClient,
		logger: logger.WithField("component", "handler"),
		redis:  redisClient,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/image", "./resources")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// услуги - угрозы
	r.GET("/threats", h.GetThreatsList)
	r.GET("/threats/:id", h.GetThreatById)
	r.Use(h.WithAuthCheck(models.Admin)).DELETE("/threats/:id", h.DeleteThreat)
	r.Use(h.WithAuthCheck(models.Admin)).POST("/threats", h.AddThreat)
	r.Use(h.WithAuthCheck(models.Admin)).PUT("/threats/:id", h.UpdateThreat)
	r.Use(h.WithAuthCheck(models.Client)).POST("/threats/request", h.AddThreatToRequest)

	// заявки - мониторинг угроз
	r.Use(h.WithAuthCheck(models.Admin)).GET("/monitoring-requests", h.GetMonitoringRequestsList)
	// разный доступ, у админа к любой, у юзера только к своей ???
	r.GET("/monitoring-requests/:id", h.GetMonitoringRequestById)
	r.Use(h.WithAuthCheck(models.Client)).DELETE("/monitoring-requests", h.DeleteMonitoringRequest)
	r.Use(h.WithAuthCheck(models.Client)).PUT("/monitoring-requests/client", h.UpdateMonitoringRequestClient)
	// важно id admin
	r.Use(h.WithAuthCheck(models.Admin)).PUT("/monitoring-requests/admin", h.UpdateMonitoringRequestAdmin)

	// м-м

	r.Use(h.WithAuthCheck(models.Client)).DELETE("/monitoring-request-threats/threats/:threatId", h.DeleteThreatFromRequest)

	// авторизация и регистрация
	r.POST("/login", h.Login)
	r.POST("/sign_up", h.Register)
	return r
}
