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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/image", "./resources")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// услуги - угрозы
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/threats", h.WithAuthCheck([]models.Role{}), h.GetThreatsList)
		apiGroup.GET("/threats/:id", h.GetThreatById)
		apiGroup.DELETE("/threats/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.DeleteThreat)
		apiGroup.POST("/threats", h.WithAuthCheck([]models.Role{models.Admin}), h.AddThreat)
		apiGroup.PUT("/threats/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateThreat)
		apiGroup.POST("/threats/request/:threatId", h.WithAuthCheck([]models.Role{models.Client}), h.AddThreatToRequest)

		// заявки - мониторинг угроз
		apiGroup.GET("/monitoring-requests", h.WithAuthCheck([]models.Role{models.Admin, models.Client}), h.GetMonitoringRequestsList)
		// разный доступ, у админа к любой, у юзера только к своей
		apiGroup.GET("/monitoring-requests/:id", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.GetMonitoringRequestById)
		apiGroup.DELETE("/monitoring-requests", h.WithAuthCheck([]models.Role{models.Client}), h.DeleteMonitoringRequest)
		apiGroup.PUT("/monitoring-requests/client", h.WithAuthCheck([]models.Role{models.Client}), h.UpdateMonitoringRequestClient)
		apiGroup.PUT("/monitoring-requests/admin/:requestId", h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateMonitoringRequestAdmin)

		// м-м

		apiGroup.DELETE("/monitoring-request-threats/threats/:threatId", h.WithAuthCheck([]models.Role{models.Client}), h.DeleteThreatFromRequest)

		// авторизация и регистрация
		apiGroup.POST("/signIn", h.SignIn)
		apiGroup.POST("/signUp", h.SignUp)
		apiGroup.POST("/logout", h.Logout)
		apiGroup.GET("/check-auth", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.CheckAuth)

		// асинхронный сервис
		apiGroup.PUT("/monitoring-requests/user-payment-finish", h.FinishUserPayment) // обращение к асинхронному сервису
	}

	return r
}
