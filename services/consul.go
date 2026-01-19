package services

import (
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type ConsulConfig struct {
	Address string
	Folder  string
}

var ConsulConnectionConfig *ConsulConfig

// Загрузка конфигурации Consul
func LoadConsulServiceConfig() {
	// Установка значений по умолчанию
	viper.SetDefault("CONSUL_ADDRESS", "localhost:8500")
	viper.SetDefault("CONSUL_FOLDER", "config")

	ConsulConnectionConfig = &ConsulConfig{
		Address: viper.GetString("CONSUL_ADDRESS"),
		Folder:  viper.GetString("CONSUL_FOLDER"),
	}

	log.Printf("Consul configuration: Address=%s, Folder=%s", ConsulConnectionConfig.Address, ConsulConnectionConfig.Folder)
}

// GetConsulSecret получает секрет из Consul
func GetConsulSecret(key string) (string, error) {
	// Подключение к Consul
	client, err := api.NewClient(&api.Config{
		Address:   ConsulConnectionConfig.Address,
		Namespace: "",
	})
	if err != nil {
		return "", err
	}

	kv, _, err := client.KV().Get(ConsulConnectionConfig.Folder+"/"+key, nil)
	if err != nil {
		return "", err
	}
	if kv == nil {
		return "", nil // Если ключ не найден
	}
	return string(kv.Value), nil
}

func MustGetConsulSecret(key string) string {

	value, err := GetConsulSecret(key)
	if err != nil || value == "" {
		log.Fatalf("Failed to get secret from Consul: %v", err)
	}
	return value
}
