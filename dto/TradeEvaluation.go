package dto

import "time"

type TradeEvaluation struct {
	Name            string
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
	CurrentStatus   string
	EvaluatedAt     time.Time
	Quantity        float64
	Charges         float64
	Profit          float64
	IsNew           bool
	Origin          string
}
