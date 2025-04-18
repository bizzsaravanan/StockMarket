package main

import (
	"StockMarket/db"
	"StockMarket/service"
	"log"

	"github.com/spf13/viper"
)

// apiKey := "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"

func main() {
	LoadConfig()
	if err := db.Connect(); err != nil {
		log.Panic("MOngo Connection failed: ", err)
	}
	db.Init()
	service.RegisterService("StockService", &service.StockService{})
	service.Connect()
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Config not found...", err)
	}
}
