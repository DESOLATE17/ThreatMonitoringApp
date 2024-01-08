package repository

import (
	"strconv"
	"threat-monitoring/internal/models"
)

// изменение информации об угрозе
func (r *Repository) UpdateThreat(updateThreat models.Threat) error {
	var threat models.Threat
	res := r.db.First(&threat, "threat_id =?", updateThreat.ThreatId)
	if res.Error != nil {
		return res.Error
	}

	if updateThreat.Name != "" {
		threat.Name = updateThreat.Name
	}

	if updateThreat.Description != "" {
		threat.Description = updateThreat.Description
	}

	if updateThreat.Image != "" {
		threat.Image = updateThreat.Image
	}

	if updateThreat.Count != 0 {
		threat.Count = updateThreat.Count
	}

	if updateThreat.Price != 0 {
		threat.Price = updateThreat.Price
	}

	result := r.db.Save(threat)
	return result.Error
}

func (r *Repository) GetThreatsList(query, lowPrice, highPrice string, isAdmin bool) ([]models.Threat, error) {
	threats := make([]models.Threat, 0)
	if lowPrice == "" {
		lowPrice = "0"
	}

	if highPrice == "" {
		highPrice = "1000000"
	}

	if query != "" {
		res := r.db.Where("is_deleted = ?", "false").Where("name LIKE ? AND price BETWEEN ? AND ?", "%"+query+"%", lowPrice, highPrice).Find(&threats)
		return threats, res.Error
	}

	res := r.db.Where("is_deleted = ? AND price BETWEEN ? AND ?", "false", lowPrice, highPrice).Find(&threats)
	return threats, res.Error
}

func (r *Repository) GetThreatByID(threatId int) (models.Threat, error) {
	threat := models.Threat{}

	err := r.db.First(&threat, "threat_id = ?", strconv.Itoa(threatId)).Error
	if err != nil {
		r.logger.Error(err)
		return threat, err
	}

	return threat, nil
}

func (r *Repository) DeleteThreatByID(threatId int) error {
	err := r.db.Exec("UPDATE threats SET is_deleted=true WHERE threat_id = ?", threatId).Error
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}

func (r *Repository) AddThreat(newThreat models.Threat) error {
	result := r.db.Create(&newThreat)
	return result.Error
}

func (r *Repository) DeleteThreatImage(threatId int) string {
	threat := models.Threat{}

	r.db.First(&threat, "threat_id = ?", strconv.Itoa(threatId))
	return threat.Image
}
