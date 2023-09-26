package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetConnectionString() (string, error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("error while getting connection string:", err)
		return "", err
	}
	return viper.GetString("db.connection_string"), nil
}
