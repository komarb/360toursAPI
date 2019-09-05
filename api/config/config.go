package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host 		string		`json:"host"`
	Port 		string		`json:"port"`
	DBName 		string		`json:"dbname"`
	CollectionName 	string 	`json:"collectionName"`
}

func GetConfig() *Config {
	var config Config
	data, err := ioutil.ReadFile("api/config.json")
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(data, &config.DB)
	if err != nil {
		log.Panic(err)
	}
	return &config
}