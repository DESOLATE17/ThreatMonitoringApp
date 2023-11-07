package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"threat-monitoring/internal/models"
	"time"
)

// список заявок
func (h *Handler) GetMonitoringRequestsList(c *gin.Context) {
	status := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.DateTime, startDateStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.DateTime, endDateStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		if endDate.Before(startDate) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "end_date не может быть раньше, чем start_date"})
			return
		}
	}

	monitoringRequests, err := h.repo.GetMonitoringRequests(status, startDate, endDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, monitoringRequests)
}

// получение заявки по id с информацией об услуге
func (h *Handler) GetMonitoringRequestById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	request, threats, err := h.repo.GetMonitoringRequestById(id, c.GetInt(userCtx), c.GetBool(adminCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "такая заявка не найдена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"request": request, "threats": threats})
	return
}

// добавляем услугу в заявку (или создаем новую заявку и добавляем в нее услугу)
func (h *Handler) AddThreatToRequest(c *gin.Context) {
	var request models.MonitoringRequestCreateMessage

	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	request.CreatorId = c.GetInt(userCtx)

	if request.ThreatId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "услуга не может быть пустой"})
		return
	}

	err = h.repo.AddThreatToRequest(request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "услуга добавлена в заявку"})
	return
}

// логически удаляет заявку
func (h *Handler) DeleteMonitoringRequest(c *gin.Context) {
	userId := c.GetInt(userCtx)
	err := h.repo.DeleteMonitoringRequest(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "заявка успешно удалена"})
}

// изменяет статус клиента в заявке
func (h *Handler) UpdateMonitoringRequestClient(c *gin.Context) {
	var newRequestStatus models.MonitoringRequest
	err := c.BindJSON(&newRequestStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if newRequestStatus.Status != "formated" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'formated'"})
		return
	}

	err = h.repo.UpdateMonitoringRequestClient(c.GetInt(userCtx), newRequestStatus.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
}

// изменяет статус администратором в заявке
func (h *Handler) UpdateMonitoringRequestAdmin(c *gin.Context) {
	var newRequestStatus models.MonitoringRequest
	err := c.BindJSON(&newRequestStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("requestId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if newRequestStatus.Status != "accepted" && newRequestStatus.Status != "canceled" && newRequestStatus.Status != "closed" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'accepted, 'closed' и 'canceled'"})
		return
	}
	err = h.repo.UpdateMonitoringRequestAdmin(c.GetInt(userCtx), id, newRequestStatus.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
	return
}

func (h *Handler) DeleteThreatFromRequest(c *gin.Context) {
	threatIdStr := c.Param("threatId")
	threatId, err := strconv.Atoi(threatIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := c.GetInt(userCtx)

	request, threats, err := h.repo.DeleteThreatFromRequest(userId, threatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Угроза удалена из заявки", "threats": threats, "monitoring-request": request})
}
