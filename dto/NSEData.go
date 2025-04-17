package dto

type NSEData struct {
	PriceInfo struct {
		LastPrice       float64 `json:"lastPrice"`
		VWAP            float64 `json:"vwap"`
		IntraDayHighLow struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"intraDayHighLow"`
		WeekHighLow struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"weekHighLow"`
	} `json:"priceInfo"`
	SecurityInfo struct {
		TotalTradedVolume float64 `json:"totalTradedVolume"`
	} `json:"securityInfo"`
}
