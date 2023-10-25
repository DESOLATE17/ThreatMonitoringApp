package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
	"threat-monitoring/internal/models"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignedRole models.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		token, err := jwt.ParseWithClaims(jwtStr, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}

		myClaims := token.Claims.(*models.JwtClaims)

		if !myClaims.IsAdmin && assignedRole == 1 || myClaims.IsAdmin && assignedRole == 0 {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Printf("user %v is not admin", myClaims.UserId)

			return
		}
	}
}
