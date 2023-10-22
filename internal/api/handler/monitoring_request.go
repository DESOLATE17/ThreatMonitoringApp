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
	}

	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "end_date cant be before start_date")
		return
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
	request, threats, err := h.repo.GetMonitoringRequestById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
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
	request.CreatorId = models.GetClientId()

	if request.ThreatId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "услуга не может быть пустой")
		return
	}

	err = h.repo.AddThreatToRequest(request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "услуга добавлена в заявку")
	return
}

// логически удаляет заявку
func (h *Handler) DeleteMonitoringRequest(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}
	err2 := h.repo.DeleteMonitoringRequest(id)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, "заявка успешно удалена")
}

// изменяет статус клиента в заявке
func (h *Handler) UpdateMonitoringRequestClient(c *gin.Context) {
	threatIdStr := c.Param("id")
	threatId, err := strconv.Atoi(threatIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}

	var newRequestStatus models.MonitoringRequest
	err = c.BindJSON(&newRequestStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if newRequestStatus.Status != "formated" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Поменять статус можно только на 'сформирован'")
		return
	}

	err = h.repo.UpdateMonitoringRequestClient(threatId, newRequestStatus.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Статус изменен")
}

// изменяет статус администратором в заявке
func (h *Handler) UpdateMonitoringRequestAdmin(c *gin.Context) {
	threatIdStr := c.Param("id")
	threatId, err := strconv.Atoi(threatIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}

	var newRequestStatus models.MonitoringRequest
	err = c.BindJSON(&newRequestStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if newRequestStatus.Status != "accepted" && newRequestStatus.Status != "canceled" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Поменять статус можно только на 'принят' и 'отменен'")
		return
	}
	err = h.repo.UpdateMonitoringRequestAdmin(threatId, newRequestStatus.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Статус изменен")
	return
}

func (h *Handler) DeleteThreatFromRequest(c *gin.Context) {
	threatIdStr := c.Param("threatId")
	threatId, err := strconv.Atoi(threatIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	requestIdStr := c.Param("requestId")
	requestId, err := strconv.Atoi(requestIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.repo.DeleteThreatFromRequest(requestId, threatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Угроза удалена из заявки")
}
