package db

import (
	"StockMarket/dto"
)

type StockRepo struct {
	MongoRepo
}

var DB StockRepo

func Init() {
	var t *dto.TradeEvaluation
	DB.CreateIndex(&t, "symbol", true)
}
