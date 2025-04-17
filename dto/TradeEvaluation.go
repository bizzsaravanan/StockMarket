package dto

import "time"

type TradeEvaluation struct {
	Symbol          string
	CurrentPrice    float64
	High5Min        float64
	AverageVolume   float64
	CurrentVolume   float64
	MovingAvg50     float64
	RSI             float64
	ResistanceLevel float64
	SupportLevel    float64
	ProfitTarget    float64
	StopLoss        float64
	BuyScore        int
	SellScore       int
	Recommendation  string
	EvaluatedAt     time.Time
}
