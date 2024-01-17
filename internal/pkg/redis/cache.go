package redis

import (
	"context"
	"encoding/json"
	"threat-monitoring/internal/models"
	"time"
)

func (c *RedisClient) GetThreats(ctx context.Context) ([]models.Threat, error) {
	result, err := c.client.Get(ctx, "threats").Result()
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	var threats []models.Threat
	err = json.Unmarshal([]byte(result), &threats)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	return threats, nil
}

func (c *RedisClient) SetThreats(ctx context.Context, threats []models.Threat) error {
	jsonBytes, err := json.Marshal(threats)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	return c.client.Set(ctx, "threats", jsonBytes, time.Hour).Err()
}
