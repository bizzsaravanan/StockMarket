package dto

type NSEChartResponse struct {
	S string    `json:"s"` // status
	T []int64   `json:"t"` // timestamps
	O []float64 `json:"o"` // open prices
	H []float64 `json:"h"` // high prices
	L []float64 `json:"l"` // low prices
	C []float64 `json:"c"` // close prices
}
