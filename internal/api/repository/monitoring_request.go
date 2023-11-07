package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"threat-monitoring/internal/models"
	"time"
)

// вывод списка всех заявок без услуг включенных в них + фильтрация по статусу и дате формирования
func (r *Repository) GetMonitoringRequests(status string, startDate, endDate time.Time) ([]models.MonitoringRequest, error) {
	var monitoringRequests []models.MonitoringRequest

	if status != "" {
		if startDate.IsZero() {
			if endDate.IsZero() {
				// фильтрация только по статусу
				res := r.db.Where("status = ?", status).Find(&monitoringRequests)
				return monitoringRequests, res.Error
			}

			// фильтрация по статусу и endDate
			res := r.db.Where("status = ?", status).Where("formation_date < ?", endDate).
				Find(&monitoringRequests)
			return monitoringRequests, res.Error
		}

		// фильтрация по статусу и startDate
		if endDate.IsZero() {
			res := r.db.Where("status = ?", status).Where("formation_date > ?", startDate).
				Find(&monitoringRequests)
			return monitoringRequests, res.Error
		}

		// фильтрация по статусу, startDate и endDate
		res := r.db.Where("status = ?", status).Where("formation_date BETWEEN ? AND ?", startDate, endDate).
			Find(&monitoringRequests)
		return monitoringRequests, res.Error
	}

	if startDate.IsZero() {
		if endDate.IsZero() {
			// без фильтрации
			res := r.db.Find(&monitoringRequests)
			return monitoringRequests, res.Error
		}

		// фильтрация по endDate
		res := r.db.Where("formation_date < ?", endDate).
			Find(&monitoringRequests)
		return monitoringRequests, res.Error
	}

	if endDate.IsZero() {
		// фильтрация по startDate
		res := r.db.Where("formation_date > ?", startDate).
			Find(&monitoringRequests)
		return monitoringRequests, res.Error
	}

	//фильтрация по startDate и endDate
	res := r.db.Where("formation_date BETWEEN ? AND ?", startDate, endDate).
		Find(&monitoringRequests)
	return monitoringRequests, res.Error
}

// вывод одной заявки со списком её услуг
func (r *Repository) GetMonitoringRequestById(requestId int) (models.MonitoringRequest, []models.Threat, error) {
	var monitoringRequest models.MonitoringRequest
	var threats []models.Threat

	//ищем такую заявку
	result := r.db.First(&monitoringRequest, "request_id =?", requestId)
	if result.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return models.MonitoringRequest{}, nil, result.Error
	}
	//ищем услуги в заявке
	res := r.db.
		Table("monitoring_requests_threats").
		Select("threats.*").
		Joins("JOIN threats ON monitoring_requests_threats.threat_id = threats.threat_id").
		Where("monitoring_requests_threats.request_id = ?", requestId).
		Find(&threats)
	if res.Error != nil {
		r.logger.Error("error while getting threats for monitoring requests")
		return models.MonitoringRequest{}, nil, res.Error
	}

	return monitoringRequest, threats, nil
}

// добавление угрозы для мониторинга в заявку и создание заявки если ее не было
func (r *Repository) AddThreatToRequest(request models.MonitoringRequestCreateMessage) error {
	var monitoringRequest models.MonitoringRequest
	r.db.Where("creator_id = ?", request.CreatorId).Where("status = ?", "created").First(&monitoringRequest)
	fmt.Println(monitoringRequest)

	if monitoringRequest.RequestId == 0 {
		newMonitoringRequest := models.MonitoringRequest{
			CreatorId:    request.CreatorId,
			AdminId:      models.GetAdminId(),
			Status:       "created",
			CreationDate: time.Now(),
		}
		res := r.db.Create(&newMonitoringRequest)
		if res.Error != nil {
			return res.Error
		}
		monitoringRequest = newMonitoringRequest
	}

	monitoringRequestsThreats := models.MonitoringRequestsThreats{
		RequestId: monitoringRequest.RequestId,
		ThreatId:  request.ThreatId,
	}

	res := r.db.Create(&monitoringRequestsThreats)
	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"monitoring_requests_threats_request_id_threat_id_key\" (SQLSTATE 23505)" {
		return errors.New("данная услуга уже добавлена в заявку")

	}

	return res.Error
}

// присваивает заявке статус удалено
func (r *Repository) DeleteMonitoringRequest(id int) error {
	var request models.MonitoringRequest
	res := r.db.First(&request, "creator_id =? and status = 'created'", id)
	if res.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return res.Error
	}

	request.Status = "deleted"
	result := r.db.Save(request)
	return result.Error
}

// изменение статуса клиента
func (r *Repository) UpdateMonitoringRequestClient(id int, status string) error {
	var monitoringRequest models.MonitoringRequest
	err := r.db.First(&monitoringRequest, "creator_id = ? and status = 'created'", id)
	if err.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return err.Error
	}

	monitoringRequest.Status = status
	if status == "formated" {
		monitoringRequest.FormationDate = time.Now()
	}
	res := r.db.Save(&monitoringRequest)

	return res.Error
}

// получение черновика заявки
func (r *Repository) GetMonitoringRequestDraft(userId int) (int, error) {
	var monitoringRequest models.MonitoringRequest
	err := r.db.First(&monitoringRequest, "creator_id = ? and status = 'created'", userId)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		r.logger.Error("error while getting monitoring request draft", err)
		return 0, err.Error
	}

	return monitoringRequest.RequestId, nil
}

// изменение статуса модератора
func (r *Repository) UpdateMonitoringRequestAdmin(id int, status string) error {
	var monitoringRequest models.MonitoringRequest

	err := r.db.First(&monitoringRequest, "request_id = ? and status = 'formated'", id)
	if err.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return err.Error
	}

	monitoringRequest.Status = status
	res := r.db.Save(&monitoringRequest)

	return res.Error
}

// удаление услуги из заявки
func (r *Repository) DeleteThreatFromRequest(userId, threatId int) (models.MonitoringRequest, []models.Threat, error) {
	var request models.MonitoringRequest
	r.db.Where("creator_id = ? and status = 'created'", userId).First(&request)

	if request.RequestId == 0 {
		return models.MonitoringRequest{}, nil, errors.New("no such request")
	}

	var monitoringRequestThreats models.MonitoringRequestsThreats
	err := r.db.Where("request_id = ? AND threat_id = ?", request.RequestId, threatId).First(&monitoringRequestThreats).Error
	if err != nil {
		return models.MonitoringRequest{}, nil, errors.New("такой угрозы нет в данной заявке")
	}

	err = r.db.Where("request_id = ? AND threat_id = ?", request.RequestId, threatId).
		Delete(models.MonitoringRequestsThreats{}).Error

	if err != nil {
		return models.MonitoringRequest{}, nil, err
	}

	return r.GetMonitoringRequestById(request.RequestId)
}
