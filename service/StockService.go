package service

import (
	"StockMarket/db"
	"StockMarket/dto"
	"context"
	"errors"
	"log"
	"math"
	"net/url"
	"time"
)

type StockService struct{}

// ListTradeEvaluation
func (a *StockService) ListTradeEvaluation(ctx context.Context, request *dto.Request) (*dto.Response, error) {
	var tradeEvaluation []*dto.TradeEvaluation
	query := db.M{
		"recommendation": db.M{
			"$in": []interface{}{"Sell", "Buy", "Strong Buy", "Strong Sell"},
		},
	}
	err := db.DB.FindAllPagination(&tradeEvaluation, query, db.M{"profit": -1}, 0, 0)
	if err != nil {
		return nil, errors.New("tradeEvaluation not found")
	}

	var atradeEvaluation []*dto.TradeEvaluation
	aquery := db.M{
		"recommendation": db.M{
			"$in": []interface{}{"Avoid"},
		},
	}
	err = db.DB.FindAllPagination(&atradeEvaluation, aquery, nil, 0, 0)
	if err != nil {
		return nil, errors.New("tradeEvaluation not found")
	}
	sanitizeEvaluations(tradeEvaluation)
	sanitizeEvaluations(atradeEvaluation)
	return &dto.Response{TradeEvaluation: tradeEvaluation, ATradeEvaluation: atradeEvaluation}, nil
}

func sanitizeFloat(f float64) float64 {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return 0 // Or a fallback like -1 or math.SmallestNonzeroFloat64
	}
	return f
}

func sanitizeTradeEvaluation(te *dto.TradeEvaluation) {
	te.CurrentPrice = sanitizeFloat(te.CurrentPrice)
	te.High5Min = sanitizeFloat(te.High5Min)
	te.AverageVolume = sanitizeFloat(te.AverageVolume)
	te.CurrentVolume = sanitizeFloat(te.CurrentVolume)
	te.MovingAvg50 = sanitizeFloat(te.MovingAvg50)
	te.RSI = sanitizeFloat(te.RSI)
	te.ResistanceLevel = sanitizeFloat(te.ResistanceLevel)
	te.SupportLevel = sanitizeFloat(te.SupportLevel)
	te.ProfitTarget = sanitizeFloat(te.ProfitTarget)
	te.StopLoss = sanitizeFloat(te.StopLoss)
	te.Quantity = sanitizeFloat(te.Quantity)
	te.Charges = sanitizeFloat(te.Charges)
	te.Profit = sanitizeFloat(te.Profit)
}

func sanitizeEvaluations(list []*dto.TradeEvaluation) {
	for _, te := range list {
		sanitizeTradeEvaluation(te)
	}
}

func (a *StockService) Start(ctx context.Context, request *dto.Request) (*dto.Response, error) {
	if request.Name != "" {
		nseService := &NseService{}
		if _, err := nseService.EvaluateStock(ctx, &dto.Request{Name: request.Name, Amount: request.Amount, Cookie: request.Cookie}); err != nil {
			log.Printf("Error evaluating stock %s: %v\n", request.Name, err)
			return nil, err
		}
		resp, err := a.ListTradeEvaluation(ctx, &dto.Request{})
		if err == nil && resp != nil {
			broadcastToClients(resp)
		} else {
			log.Println("Failed to list trade evaluations:", err)
		}
		return &dto.Response{Message: "started"}, nil
	}
	go func() {
		ticker := time.NewTicker(20 * time.Second) // run every 15s (you can adjust)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("StockService Start: context done, stopping loop")
				return

			case <-ticker.C:
				log.Println("Starting evaluation cycle...")

				var tradeEvaluation []*dto.TradeEvaluation
				err := db.DB.FindAllPagination(&tradeEvaluation, nil, db.M{"profit": -1}, 0, 0)
				if err != nil {
					log.Println("tradeEvaluation not found start")
				}

				m := make(map[string]string)

				for _, d := range tradeEvaluation {
					m[d.Symbol] = "DB"
				}

				form := url.Values{}
				form.Set("scan_clause", `( {cash} ( [0] 5 minute close > [-1] 5 minute max( 20 , [0] 5 minute close ) and [0] 5 minute volume > [0] 5 minute sma( volume,20 ) and latest volume > 1000000 ) )`)
				form.Set("debug_clause", `groupcount( 1 where [0] 5 minute close > [-1] 5 minute max( 20 , [0] 5 minute close )),groupcount( 1 where [0] 5 minute volume > [0] 5 minute sma( volume,20 )),groupcount( 1 where daily volume > 1000000)`)
				fetchChartData(ctx, m, form, "CHART_INK_HIGH_VOLUME")

				form2 := url.Values{}
				form2.Set("scan_clause", `( {57960} ( [0] 15 minute close > [-1] 15 minute max( 20 , [0] 15 minute close ) and [0] 15 minute volume > [0] 15 minute sma( volume,20 ) ) )`)
				form2.Set("debug_clause", `groupcount( 1 where [0] 15 minute close > [-1] 15 minute max( 20 , [0] 15 minute close )),groupcount( 1 where [0] 15 minute volume > [0] 15 minute sma( volume,20 ))`)
				fetchChartData(ctx, m, form2, "CHART_INK_15_min")

				form3 := url.Values{}
				form3.Set("scan_clause", `( {57960} ( [0] 15 minute close > [0] 15 minute open and [0] 15 minute high <= [0] 15 minute close * 1.0005 and [0] 15 minute low >= [0] 15 minute open * 0.9995 and [0] 15 minute close > latest open * 1.05 and ( [0] 15 minute high - [0] 15 minute open ) * .65 <= ( [0] 15 minute close - [0] 15 minute open ) and [0] 15 minute close > [0] 15 minute open and [0] 15 minute volume >= 1000 ) )`)
				form3.Set("debug_clause", `groupcount( 1 where [0] 15 minute close > [0] 15 minute open),groupcount( 1 where [0] 15 minute high <= [0] 15 minute close * 1.0005),groupcount( 1 where [0] 15 minute low >= [0] 15 minute open * 0.9995),groupcount( 1 where [0] 15 minute close > daily open * 1.05),groupcount( 1 where ( [0] 15 minute high - [0] 15 minute open ) * .65 <= ( [0] 15 minute close - [0] 15 minute open )),groupcount( 1 where [0] 15 minute close > [0] 15 minute open),groupcount( 1 where [0] 15 minute volume >= 1000)`)
				fetchChartData(ctx, m, form3, "CHART_INK_15_BULLISH")

				form4 := url.Values{}
				form4.Set("scan_clause", `( {57960} ( [0] 15 minute close < [0] 15 minute open * 0.99 and [0] 15 minute high <= [0] 15 minute open * 1.0005 and [0] 15 minute low >= [0] 15 minute close * 0.9995 and [0] 15 minute volume >= 1000 ) )`)
				form4.Set("debug_clause", `groupcount( 1 where [0] 15 minute close < [0] 15 minute open * 0.99),groupcount( 1 where [0] 15 minute high <= [0] 15 minute open * 1.0005),groupcount( 1 where [0] 15 minute low >= [0] 15 minute close * 0.9995),groupcount( 1 where [0] 15 minute volume >= 1000)`)
				fetchChartData(ctx, m, form4, "CHART_INK_15_BEARISH")

				form5 := url.Values{}
				form5.Set("scan_clause", `( {57960} ( [0] 5 minute close > [-1] 5 minute max( 20 , [0] 5 minute close ) and [0] 5 minute volume > [0] 5 minute sma( volume,20 ) ) )`)
				form5.Set("debug_clause", `groupcount( 1 where [0] 5 minute close > [-1] 5 minute max( 20 , [0] 5 minute close )),groupcount( 1 where [0] 5 minute volume > [0] 5 minute sma( volume,20 ))`)
				fetchChartData(ctx, m, form5, "CHART_INK_5_min")

				form6 := url.Values{}
				form6.Set("scan_clause", `( {57960} ( [0] 5 minute close > [0] 5 minute open and [0] 5 minute high <= [0] 5 minute close * 1.0005 and [0] 5 minute low >= [0] 5 minute open * 0.9995 and [0] 5 minute close > latest open * 1.05 and ( [0] 5 minute high - [0] 5 minute open ) * .65 <= ( [0] 5 minute close - [0] 5 minute open ) and [0] 5 minute close > [0] 5 minute open and [0] 5 minute volume >= 1000 ) )`)
				form6.Set("debug_clause", `groupcount( 1 where [0] 5 minute close > [0] 5 minute open),groupcount( 1 where [0] 5 minute high <= [0] 5 minute close * 1.0005),groupcount( 1 where [0] 5 minute low >= [0] 5 minute open * 0.9995),groupcount( 1 where [0] 5 minute close > daily open * 1.05),groupcount( 1 where ( [0] 5 minute high - [0] 5 minute open ) * .65 <= ( [0] 5 minute close - [0] 5 minute open )),groupcount( 1 where [0] 5 minute close > [0] 5 minute open),groupcount( 1 where [0] 5 minute volume >= 1000)`)
				fetchChartData(ctx, m, form6, "CHART_INK_BULLISH_5_min")

				form7 := url.Values{}
				form7.Set("scan_clause", `( {57960} ( [0] 5 minute close < [0] 5 minute open * 0.99 and [0] 5 minute high <= [0] 5 minute open * 1.0005 and [0] 5 minute low >= [0] 5 minute close * 0.9995 and [0] 5 minute volume >= 1000 ) )`)
				form7.Set("debug_clause", `groupcount( 1 where [0] 5 minute close < [0] 5 minute open * 0.99),groupcount( 1 where [0] 5 minute high <= [0] 5 minute open * 1.0005),groupcount( 1 where [0] 5 minute low >= [0] 5 minute close * 0.9995),groupcount( 1 where [0] 5 minute volume >= 1000)`)
				fetchChartData(ctx, m, form7, "CHART_INK_BEARISH_5_min")

				for d, o := range m {
					nseService := &NseService{}
					if _, err := nseService.EvaluateStock(ctx, &dto.Request{Name: d, Amount: request.Amount, Cookie: request.Cookie, Origin: o}); err != nil {
						log.Printf("Error evaluating stock %s: %v\n", d, err)
					}
				}

				resp, err := a.ListTradeEvaluation(ctx, &dto.Request{})
				if err == nil && resp != nil {
					broadcastToClients(resp)
				} else {
					log.Println("Failed to list trade evaluations:", err)
				}
			}
		}
	}()

	return &dto.Response{Message: "started"}, nil
}

func fetchChartData(ctx context.Context, m map[string]string, form url.Values, source string) {
	chartinkService := &ChartinkService{}
	res, err := chartinkService.CreateStockData(ctx, &dto.Request{FormData: form})
	if err != nil || res == nil || res.ChartinkResponse == nil {
		log.Println("ChartinkService failed:", err)
	}

	for _, data := range res.ChartinkResponse.Data {
		if existingSource, ok := m[data.NSECode]; ok {
			m[data.NSECode] = existingSource + "," + source
		} else {
			m[data.NSECode] = source
		}
	}
}

func broadcastToClients(data *dto.Response) {
	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Println("WebSocket write error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func (a *StockService) Reset(ctx context.Context, request *dto.Request) (*dto.Response, error) {
	var tradeEvaluation []*dto.TradeEvaluation
	err := db.DB.DeleteMany(&tradeEvaluation, nil)
	if err != nil {
		return nil, errors.New("tradeEvaluation not found")
	}
	return &dto.Response{Message: "Deleted"}, nil
}
