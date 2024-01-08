package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"threat-monitoring/internal/models"
	"threat-monitoring/internal/utils"
)

// DeleteThreat godoc
// @Summary      Delete threat by ID
// @Description  Deletes a threat with the given ID
// @Tags         Threats
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Threat ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/threats/{id} [delete]
func (h *Handler) DeleteThreat(c *gin.Context) {
	threatId := c.Param("id")
	id, err := strconv.Atoi(threatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	err = h.repo.DeleteThreatByID(id)
	if err != nil {
		h.logger.Error("cant get threat by id %v", err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "угроза успешно удалена"})
}

// GetThreatsList godoc
// @Summary      Get threats list
// @Description  Retrieves a list of threats based on the provided query.
// @Tags         Threats
// @Accept       json
// @Produce      json
// @Param        query   query    string  false  "Query string to filter threats"
// @Param        lowPrice   query    number  false  "LowPrice to filter threats"
// @Param        highPrice   query    number  false  "HighPrice string to filter threats"
// @Success      200  {object}  map[string]any
// @Failure      500  {object}  error
// @Router       /api/threats [get]
func (h *Handler) GetThreatsList(c *gin.Context) {
	query := c.Query("query")
	lowPriceStr := c.Query("lowPrice")
	highPriceStr := c.Query("highPrice")

	lowPrice, err := strconv.Atoi(lowPriceStr)
	if err != nil {
		lowPrice = 0
	}

	highPrice, err := strconv.Atoi(highPriceStr)
	if err != nil {
		highPrice = 1000000
	}

	if lowPrice > highPrice {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("неверно заданы параметры фильтрации по цене"))
		return
	}

	threats, err := h.repo.GetThreatsList(query, lowPriceStr, highPriceStr, c.GetBool(adminCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	requestId, err := h.repo.GetMonitoringRequestDraft(c.GetInt(userCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if requestId == 0 {
		c.JSON(http.StatusOK, gin.H{"threats": threats})
		return
	}
	c.JSON(http.StatusOK, gin.H{"threats": threats, "draftId": requestId})
}

// GetThreatById godoc
// @Summary      Get threat by ID
// @Description  Retrieves a threat by its ID
// @Tags         Threats
// @Produce      json
// @Param        id   path    int     true        "Threat ID"
// @Success      200  {object}  models.Threat
// @Failure      400  {object}  error
// @Router       /api/threats/{id} [get]
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

// AddThreat godoc
// @Summary      Add new threat
// @Description  Add a new threat with image, name, description, summary, count, and price
// @Tags         Threats
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file false "Threat image"
// @Param        name formData string true "Threat name"
// @Param        description formData string false "Threat description"
// @Param        summary formData string false "Threat summary"
// @Param        count formData string true "Threat count"
// @Param        price formData string true "Threat price"
// @Success      201  {string}  map[string]an
// @Failure      400  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /api/threats [post]
func (h *Handler) AddThreat(c *gin.Context) {
	var newThreat models.Threat
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		if newThreat.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
			return
		}
	}

	newThreat.Name = c.Request.FormValue("name")
	if newThreat.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "имя угрозы не может быть пустым"})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "цена указана неверно"})
		return
	}

	if err = h.repo.AddThreat(newThreat); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "новая услуга успешно добавлена"})
}

// UpdateThreat godoc
// @Summary      Update threat by ID
// @Description  Updates a threat with the given ID
// @Tags         Threats
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path        int     true        "ID"
// @Param        name        formData    string  false       "name"
// @Param        description formData    string  false       "description"
// @Param        count       formData    string  false       "count"
// @Param        price       formData    string  false       "price"
// @Param        image       formData    file    false       "image"
// @Success      200         {object}    map[string]any
// @Failure      400         {object}    error
// @Router       /api/threats/{id} [put]
func (h *Handler) UpdateThreat(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")

	var updateThreat models.Threat
	threatId := c.Param("id")
	updateThreat.ThreatId, err = strconv.Atoi(threatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	updateThreat.Name = c.Request.FormValue("name")
	updateThreat.Description = c.Request.FormValue("description")

	count := c.Request.FormValue("count")
	if count != "" {
		updateThreat.Count, err = strconv.Atoi(count)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	}

	price := c.Request.FormValue("price")
	if price != "" {
		updateThreat.Price, err = strconv.Atoi(price)
		if err != nil || updateThreat.Price == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "цена указана неверно"})
			return
		}
	}
	if header != nil && header.Size != 0 {
		if updateThreat.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		// delete old image from db

		url := h.repo.DeleteThreatImage(updateThreat.ThreatId)

		// delete image from minio
		h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
	}

	if err = h.repo.UpdateThreat(updateThreat); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "услуга успешно изменена"})
}
