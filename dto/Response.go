package dto

type Response struct {
	Message          string             `bson:"message" json:"message,omitempty"`
	TradeEvaluation  []*TradeEvaluation `bson:"tradeEvaluation" json:"tradeEvaluation,omitempty"`
	ATradeEvaluation []*TradeEvaluation `bson:"atradeEvaluation" json:"atradeEvaluation,omitempty"`
	ChartinkResponse *ChartinkResponse  `bson:"chartinkResponse" json:"chartinkResponse,omitempty"`
}
