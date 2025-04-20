package service

import (
	"StockMarket/db"
	"StockMarket/dto"
	"context"
	"errors"
	"log"
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
	return &dto.Response{TradeEvaluation: tradeEvaluation, ATradeEvaluation: atradeEvaluation}, nil
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
		ticker := time.NewTicker(30 * time.Second) // run every 15s (you can adjust)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("StockService Start: context done, stopping loop")
				return

			case <-ticker.C:
				log.Println("Starting evaluation cycle...")

				chartinkService := &ChartinkService{}
				res, err := chartinkService.CreateStockData(ctx, &dto.Request{})
				if err != nil || res == nil || res.ChartinkResponse == nil {
					log.Println("ChartinkService failed:", err)
					return
				}

				for _, data := range res.ChartinkResponse.Data {
					nseService := &NseService{}
					if _, err := nseService.EvaluateStock(ctx, &dto.Request{Name: data.NSECode, Amount: request.Amount, Cookie: request.Cookie}); err != nil {
						log.Printf("Error evaluating stock %s: %v\n", data.NSECode, err)
						return
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
