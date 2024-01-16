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
// @Router       /api/signUp [post]
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
		Password: newClient.Password,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "нельзя создать пользователя с таким логином"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "пользователь успешно создан"})
}

// CheckAuth godoc
// @Summary      Check user authentication
// @Description  Retrieves user information based on the provided user context
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  string
// @Router       /api/check-auth [get]
func (h *Handler) CheckAuth(c *gin.Context) {
	var userInfo = models.User{UserId: c.GetInt(userCtx)}

	userInfo, err := h.repo.GetUserInfo(c.Request.Context(), userInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "неверный формат данных")
		return
	}

	c.JSON(http.StatusOK, userInfo)
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
// @Router       /api/signIn [post]
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

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "клиент успешно авторизован", "isAdmin": user.IsAdmin, "login": user.Login, "userId": user.UserId})
}

// Logout godoc
// @Summary      Logout
// @Description  Logs out the user by blacklisting the access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /api/logout [post]
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
