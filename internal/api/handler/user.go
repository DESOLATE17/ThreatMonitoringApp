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

func (h *Handler) SignUp(c *gin.Context) {
	var newClient models.User
	var err error

	if err = c.BindJSON(&newClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных о новом пользователе"})
		return
	}

	if newClient.Password, err = h.hasher.Hash(newClient.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	if err = h.repo.SignUp(c.Request.Context(), newClient); err != nil {
		if errors.Is(err, models.ErrClientAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) SignIn(c *gin.Context) {
	var clientInfo models.User
	var err error

	if err = c.BindJSON(&clientInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "неверный формат данных")
		return
	}

	if clientInfo.Password, err = h.hasher.Hash(clientInfo.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	user, err := h.repo.GetByCredentials(c.Request.Context(), clientInfo)
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
