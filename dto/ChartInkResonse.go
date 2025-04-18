package dto

type ChartinkResponse struct {
	Draw            int         `json:"draw"`
	RecordsTotal    int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Data            []StockData `json:"data"`
	Link            string      `json:"link"`
}

type StockData struct {
	Sr      int     `json:"sr"`
	NSECode string  `json:"nsecode"`
	Name    string  `json:"name"`
	BSECode string  `json:"bsecode"`
	PerChg  float64 `json:"per_chg"`
	Close   float64 `json:"close"`
	Volume  int64   `json:"volume"`
}
