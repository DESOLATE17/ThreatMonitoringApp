package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"threat-monitoring/internal/models"
	"time"
)

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Creates a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  models.UserSignUp  true  "User information"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      409  {object}  error
// @Failure      500  {object}  error
// @Router       /signUp [post]
func (h *Handler) SignUp(c *gin.Context) {
	var newClient models.UserSignUp
	var err error

	if err = c.BindJSON(&newClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных о новом пользователе"})
		return
	}

	if newClient.Password, err = h.hasher.Hash(newClient.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	if err = h.repo.SignUp(c.Request.Context(), models.User{
		Login:    newClient.Login,
		Name:     newClient.Name,
		Password: newClient.Password,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "нельзя создать пользователя с таким логином"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "пользователь успешно создан"})
}

// SignIn godoc
// @Summary      User sign-in
// @Description  Authenticates a user and generates an access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  models.UserLogin  true  "User information"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      401  {object}  error
// @Failure      500  {object}  error
// @Router       /signIn [post]
func (h *Handler) SignIn(c *gin.Context) {
	var clientInfo models.UserLogin
	var err error

	if err = c.BindJSON(&clientInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "неверный формат данных")
		return
	}

	if clientInfo.Password, err = h.hasher.Hash(clientInfo.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	user, err := h.repo.GetByCredentials(c.Request.Context(), models.User{Password: clientInfo.Password, Login: clientInfo.Login})
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка авторизации"})
		return
	}

	token, err := h.tokenManager.NewJWT(user.UserId, user.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при формировании токена"})
		return
	}

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "127.0.0.1:8080", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "клиент успешно авторизован"})
}

// Logout godoc
// @Summary      Logout
// @Description  Logs out the user by blacklisting the access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	jwtStr, err := c.Cookie("AccessToken")
	if !strings.HasPrefix(jwtStr, jwtPrefix) || err != nil { // если нет префикса то нас дурят!
		c.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа
		return
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	_, _, err = h.tokenManager.Parse(jwtStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	// сохраняем в блеклист редиса
	err = h.redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, time.Hour)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

const (
	ServerToken = "abahjsvbdwekvnva"
	ServiceUrl  = "http://127.0.0.1:8081/addPayment"
)

func (h *Handler) UserPayment(c *gin.Context) {
	// принимает заявку и отправляет её в сервис
	client := &http.Client{}
	req, err := http.NewRequest("PUT", ServiceUrl, c.Request.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Server-Token", ServerToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode == 200 {
		c.JSON(http.StatusOK, gin.H{"message": "заявка принята в обработку"})
		return
	}
	c.AbortWithError(http.StatusInternalServerError, errors.New("заявка не принята в обработку"))
}

// ручка вызывается сервисом на python
func (h *Handler) FinishUserPayment(c *gin.Context) {
	token := c.GetHeader("Server-Token")
	if token != ServerToken {
		c.AbortWithError(http.StatusForbidden, errors.New("неверный токен"))
		return
	}

	var request models.MonitoringRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	// сохраняем в базу
	err := h.repo.SavePayment(request)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "данные сохранены"})
}
