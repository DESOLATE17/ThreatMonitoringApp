package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"threat-monitoring/internal/api"
	"threat-monitoring/internal/models"
)

type Handler struct {
	repo api.Repo
}

func NewHandler(repo api.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")

	r.GET("/home", h.GetCardsList)
	r.GET("/card/:id", h.GetCardById)
	r.POST("/card/:id", h.DeleteCard)

	r.Static("/image", "./resources")
	return r
}

func (h *Handler) DeleteCard(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = h.repo.DeleteThreatByID(id)
	if err != nil { // если не получилось
		log.Printf("cant get product by id %v", err)
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/home")
}

func (h *Handler) GetCardsList(c *gin.Context) {
	query := c.Query("query")
	resultThreats := make([]models.Threat, 0)

	// получаем данные по товару
	threats, err := h.repo.GetThreats()
	if err != nil { // если не получилось
		log.Printf("cant get product by id %v", err)
		c.Error(err)
		return
	}
	fmt.Println(threats)

	if query != "" {
		for _, v := range threats {
			if strings.Contains(v.Name, query) {
				resultThreats = append(resultThreats, v)
			}
		}
	} else {
		resultThreats = threats
	}
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title": "Main website",
		"items": resultThreats,
	})
}

func (h *Handler) GetCardById(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	threat, err := h.repo.GetThreatByID(id)
	if err != nil { // если не получилось
		log.Printf("cant get product by id %v", err)
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "card.gohtml", gin.H{
		"title": "Main website",
		"card":  threat,
	})
}
