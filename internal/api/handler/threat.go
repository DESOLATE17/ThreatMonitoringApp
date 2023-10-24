package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"threat-monitoring/internal/models"
	"threat-monitoring/internal/utils"
)

func (h *Handler) DeleteThreat(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = h.repo.DeleteThreatByID(id)
	if err != nil {
		log.Printf("cant get threat by id %v", err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "threat deleted successfully")
}

func (h *Handler) GetThreatsList(c *gin.Context) {
	query := c.Query("query")

	threats, err := h.repo.GetThreatsList(query)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	requestId, err := h.repo.GetMonitoringRequestDraft(models.GetClientId())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if requestId == 0 {
		c.JSON(http.StatusOK, gin.H{"threats": threats})
		return
	}
	c.JSON(http.StatusOK, gin.H{"threats": threats, "draftId": requestId})
}

func (h *Handler) GetThreatById(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	threat, err := h.repo.GetThreatByID(id)
	if err != nil { // если не получилось
		h.logger.Printf("cant get threat by id %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, threat)
}

func (h *Handler) AddThreat(c *gin.Context) {
	var newThreat models.Threat
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newThreat.Name = c.Request.FormValue("name")
	if newThreat.Name == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("имя угрозы не может быть пустым"))
		return
	}

	newThreat.Description = c.Request.FormValue("description")

	count := c.Request.FormValue("count")
	newThreat.Count, err = strconv.Atoi(count)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	price := c.Request.FormValue("price")
	newThreat.Price, err = strconv.Atoi(price)
	if err != nil || newThreat.Price == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("цена указана неверно"))
		return
	}

	if newThreat.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err = h.repo.AddThreat(newThreat); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "новая услуга успешно добавлена")
}

// изменяет данные про угрозу
func (h *Handler) UpdateThreat(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")

	var updateThreat models.Threat
	threatId := c.Param("id")
	updateThreat.ThreatId, err = strconv.Atoi(threatId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	updateThreat.Name = c.Request.FormValue("name")
	updateThreat.Description = c.Request.FormValue("description")

	count := c.Request.FormValue("count")
	if count != "" {
		updateThreat.Count, err = strconv.Atoi(count)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	price := c.Request.FormValue("price")
	if price != "" {
		updateThreat.Price, err = strconv.Atoi(price)
		if err != nil || updateThreat.Price == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("цена указана неверно"))
			return
		}
	}
	if header != nil && header.Size != 0 {
		if updateThreat.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// delete old image from db

		url := h.repo.DeleteThreatImage(updateThreat.ThreatId)

		// delete image from minio
		h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
	}

	if err = h.repo.UpdateThreat(updateThreat); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "услуга успешно изменена"})
}
