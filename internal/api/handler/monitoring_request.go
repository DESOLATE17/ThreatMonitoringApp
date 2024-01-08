package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"threat-monitoring/internal/models"
	"time"
)

// GetMonitoringRequestsList godoc
// @Summary      Get list of monitoring requests
// @Description  Retrieves a list of monitoring requests based on the provided parameters
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Param        status      query  string    false  "Monitoring request status"
// @Param        start_date  query  string    false  "Start date in the format '2006-01-02T15:04:05Z'"
// @Param        end_date    query  string    false  "End date in the format '2006-01-02T15:04:05Z'"
// @Success      200  {object}  []models.MonitoringRequest
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/monitoring-requests [get]
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

	monitoringRequests, err := h.repo.GetMonitoringRequests(status, startDate, endDate, c.GetInt(userCtx), c.GetBool(adminCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, monitoringRequests)
}

// GetMonitoringRequestById godoc
// @Summary      Get monitoring request by ID
// @Description  Retrieves a monitoring request with the given ID
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Monitoring Request ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/monitoring-requests/{id} [get]
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

// AddThreatToRequest godoc
// @Summary      Add threat to request
// @Description  Adds a threat to a monitoring request
// @Tags         Threats
// @Accept       json
// @Produce      json
// @Param        threatId  path  int  true  "Threat ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/threats/request/{threatId} [post]
func (h *Handler) AddThreatToRequest(c *gin.Context) {
	var request models.MonitoringRequestCreateMessage
	var err error

	request.CreatorId = c.GetInt(userCtx)
	idStr := c.Param("threatId")
	request.ThreatId, err = strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

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

// DeleteMonitoringRequest godoc
// @Summary      Delete monitoring request by user ID
// @Description  Deletes a monitoring request for the given user ID
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/monitoring-requests [delete]
func (h *Handler) DeleteMonitoringRequest(c *gin.Context) {
	userId := c.GetInt(userCtx)
	err := h.repo.DeleteMonitoringRequest(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "заявка успешно удалена"})
}

// UpdateMonitoringRequestClient godoc
// @Summary      Update monitoring request status by client
// @Description  Updates the status of a monitoring request by client on formated
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Param        newStatus    body    models.NewStatus  true    "New status of the monitoring request"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  error
// @Router       /api/monitoring-requests/client [put]
func (h *Handler) UpdateMonitoringRequestClient(c *gin.Context) {
	var newRequestStatus models.NewStatus
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

// UpdateMonitoringRequestAdmin godoc
// @Summary      Update monitoring request status by ID
// @Description  Updates the status of a monitoring request with the given ID on "accepted"/"closed"/"canceled"
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Param        requestId  path  int  true  "Request ID"
// @Param        newRequestStatus  body  models.NewStatus  true  "New request status"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /monitoring-requests/admin/{requestId} [put]
func (h *Handler) UpdateMonitoringRequestAdmin(c *gin.Context) {
	var newRequestStatus models.NewStatus
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

	if newRequestStatus.Status != "accepted" && newRequestStatus.Status != "canceled" {
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

// DeleteThreatFromRequest godoc
// @Summary      Delete threat from request
// @Description  Deletes a threat from a request based on the user ID and threat ID
// @Tags         MonitoringRequests
// @Accept       json
// @Produce      json
// @Param        threatId  path  int  true  "Threat ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Router       /api/monitoring-request-threats/threats/{threatId} [delete]
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
