package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	_ "threat-monitoring/docs"
	"threat-monitoring/internal/api"
	"threat-monitoring/internal/api/repository"
	"threat-monitoring/internal/models"
	"threat-monitoring/internal/pkg/auth"
	"threat-monitoring/internal/pkg/hash"
	minio "threat-monitoring/internal/pkg/minio"
	redis "threat-monitoring/internal/pkg/redis"
)

type Handler struct {
	logger *logrus.Entry

	minio minio.Client
	redis redis.Client
	repo  api.Repo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("/home/dasha/GolandProjects/WebAppDevelopment/ThreatMonitoringApp/config")
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

	tokenManager, err := auth.NewManager(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		logger.Fatalln(err)
	}

	return &Handler{
		repo:         repo,
		minio:        minioClient,
		logger:       logger.WithField("component", "handler"),
		redis:        redisClient,
		hasher:       hash.NewSHA256Hasher(os.Getenv("SALT")),
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/image", "./resources")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// услуги - угрозы
	r.GET("/api/threats", h.GetThreatsList)
	r.GET("/api/threats/:id", h.GetThreatById)
	r.DELETE("/threats/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.DeleteThreat)
	r.POST("/threats", h.WithAuthCheck([]models.Role{models.Admin}), h.AddThreat)
	r.PUT("/threats/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateThreat)
	r.POST("/threats/request/:threatId", h.WithAuthCheck([]models.Role{models.Client}), h.AddThreatToRequest)

	// заявки - мониторинг угроз
	r.GET("/monitoring-requests", h.WithAuthCheck([]models.Role{models.Admin, models.Client}), h.GetMonitoringRequestsList)
	// разный доступ, у админа к любой, у юзера только к своей
	r.GET("/monitoring-requests/:id", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.GetMonitoringRequestById)
	r.DELETE("/monitoring-requests", h.WithAuthCheck([]models.Role{models.Client}), h.DeleteMonitoringRequest)
	r.PUT("/monitoring-requests/client", h.WithAuthCheck([]models.Role{models.Client}), h.UpdateMonitoringRequestClient)
	r.PUT("/monitoring-requests/admin/:requestId", h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateMonitoringRequestAdmin)

	// м-м

	r.DELETE("/monitoring-request-threats/threats/:threatId", h.WithAuthCheck([]models.Role{models.Client}), h.DeleteThreatFromRequest)

	// авторизация и регистрация
	r.POST("/signIn", h.SignIn)
	r.POST("/signUp", h.SignUp)
	r.POST("/logout", h.Logout)

	// асинхронный сервис
	r.PUT("/monitoring-requests/user-payment-start", h.WithAuthCheck([]models.Role{models.Client}), h.UserPayment) // обращение к асинхронному сервису
	r.PUT("/monitoring-requests/user-payment-finish", h.FinishUserPayment)                                         // обращение к асинхронному сервису

	return r
}
