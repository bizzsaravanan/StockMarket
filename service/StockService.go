package service

import (
	"StockMarket/db"
	"StockMarket/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	cookies string = `_ga=GA1.1.1449965907.1742914394; _ga_E0LYHCLJY3=GS1.1.1743048259.3.0.1743048262.0.0.0; _ga_WM2NSQKJEK=GS1.1.1743048263.5.1.1743049265.0.0.0; nseQuoteSymbols=[{"symbol":"FUSI-RE","identifier":"","type":"equity"},{"symbol":"BHARTIARTL","identifier":"","type":"equity"},{"symbol":"DRSCARGO","identifier":"","type":"equity"},{"symbol":"BAJAJFINSV","identifier":"","type":"equity"}]; AKA_A2=A; nsit=q9fXTCitiiSI_F-v_PBiYU3m; RT="z=1&dm=nseindia.com&si=f09c6216-44fe-4695-836b-3f3edec05350&ss=m9lcjp15&sl=2&se=8c&tt=23c&bcn=%2F%2F684d0d4c.akstat.io%2F&ld=ixy&nu=31bqm9a&cl=pez&ul=1lpz&hd=1ltg"; nseappid=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkubnNlIiwiYXVkIjoiYXBpLm5zZSIsImlhdCI6MTc0NDg5MzY0MSwiZXhwIjoxNzQ0OTAwODQxfQ.GC3_CDUNX8D-sg_4ih9UV875PzgBB8i5lYDtTGgY7Rk; _ga_87M7PJ3R97=GS1.1.1744893249.13.1.1744893640.10.0.0; _abck=279A6C95DA46B78A5A53B5953CC08322~0~YAAQzELHF53/x0GWAQAAamzIQw35USMpAniv/8cElG8y5j51fpOnoJxRu2Va2jdK4iiDzIBdd1eyfEPCC9uPZtQ4mSCaVHJDLXLVnAEH/gL8qG4QyenhFQn8X6CvfOFk8CKEzT+S3SqVWUay3UXRRcczCDftYXh/0TGihUI9wE/XfSBq1Et3X3sJ53wAl8zaPu6iqDVc/ux1t/zoF19Qqb+PzRDqMVRHkXyk5oczOmPCZuKV7lwiJHFWympPwojgiVoGeySUMLXReSS7xfBlLfETGxD5FPv9w5/vbgeOaFRyCRgXLWnY8AxH9tB63e8byNcwizdaVi/Kdt7gMkrXYxzLyMA/YhBt7Ou6Rt9N/9r3tNTtki7WmrE/CjhI8W/nBl6V1BhXECYCRhWLvVUlP+6lW2SNDUXOfaPJaey0lyCOqgPcvbJxnbWCSV/2Xkn0Jzc/mwnsoc9YuL9qtfVOxG57/48TyjRDcst4BLOetAdMoEAdFFvoTSSEsRYbi6euaKdrMQH7lne4LIEHmDf54KbXkXl4qhSsVKUNJ0pZH/MglfQKhcAjfdoM6oXi3BXzThdTzWA=~-1~-1~-1; ak_bmsc=6B5711B82876BA5BDA5FF0CFF90C30D7~000000000000000000000000000000~YAAQvELHF9kBxkGWAQAArpnYQxsI04i40UP0/Au3RdDXPd+2zFXFNH5SKulpTQed2RfEdGRshPWGoXo/TGPt+bo93DZ1PTBZCLqXdA3Kb4lRoFnCmOS8CxfuUKvt37ACuTori+4SFY/CslpL63dc35LUtrwc69V0lwRSAbanzrA6xErXOu7L8QR8TR2KmzybcinKlmKlqTJKPIr1V+kWsdFwiUUasdWRAJN+y3by/DGn/V5TZLIXSAjEOajvbqGRyZjUyAL5doVEcwJuLF6raT8mF5+AYUvNZZ24iNZrg1zM3K9U74VNCjncGNs2R1a+S9uN6nSxKpcTOPVFYrkvMBJDcYNr9eNJ43c7CjrSo0MO4TxrGigV/MgSCeOd6FNRvhk5MKLM6ZfBoAE=; bm_sv=44EB4034B73F790C9962C2FD0C8B7566~YAAQl7xWaEhXpj2WAQAADI3rQxv4CMZgwyoEgAFbu9EXHpc9GpO1UgziSs05eUMVoWHBgWYCGeZc4gL5uq62eDVbK85Okik190ycbwox/RgFLK44aopkEVtmXELHKpfBpOnowrIj/2VJO0ult39Vtp8M/6NYNVtRUqYAWjPAK1qujoEsOLeKC/j9Qelm2jUyZHWx7mCdDOewt3zeUW1Tgk6MUPY23NoYBgKJJg554tOaHJHigCZMF404SfKrC6/ptouE~1; bm_sz=EDA141C767EAC0A70996C3B9140A5310~YAAQl7xWaElXpj2WAQAADY3rQxs4S2wwyBATE5pjvlMar3N4+p0SPDaSqNOJS5ycE3IyDlYvmDcfsz4maU7cCHR3U+HCy8hsfnc4pK3BPqKpT+iZ2Bmg4cf+87+M1MGPzWTpQE1LBmCQPSOQrQKxlyH8j8m4aw9xFk7i1Gm19/Pi7/WMq8GXDmibTSONCgRqeNF7aVPULF/BJNZ2Ofdq/MT3LuLh0OTnRdnx6Z3K48plFYdlMfbIJPwInm14SXaF24ZxEgCiThs0nvt88/4sodtN8txHOcIDRzwtg/RMdYLErpF4Vj8BFgOu89JcLgDgvXubtnziqMTR1Go0mH6BiVWgv3L9KLdpMA/FcpXHwt0tMxp/j5v6pAO6Qcv/hH6yYO7/Ug0dynLm69e3+6UZJPK7Q9DrEz9bNUbiddxWKS3UFzxhxlA//qxHpHjCsEfH7UnWlwC3vUmz+zwRoD4=~3683906~4276548`
)

type StockService struct{}

func (s *StockService) CreateStockData(ctx context.Context, request *dto.Request) {
	symbol := request.Name
	data, err := fetchNSEData(symbol)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	currentPrice := data.PriceInfo.LastPrice
	high5min := data.PriceInfo.IntraDayHighLow.Min   // Assuming this represents the 5-minute high
	avgVolume := data.SecurityInfo.TotalTradedVolume // Placeholder for average volume
	currentVolume := data.SecurityInfo.TotalTradedVolume
	movingAvg50 := data.PriceInfo.VWAP // Placeholder for 50-day moving average
	rsi := calculateRSI(symbol, 14)
	resistanceLevel := data.PriceInfo.WeekHighLow.Max
	supportLevel := data.PriceInfo.WeekHighLow.Min
	profitTarget := resistanceLevel
	stopLoss := supportLevel

	log.Println(currentPrice, high5min, avgVolume, currentVolume, movingAvg50, rsi, resistanceLevel, supportLevel, profitTarget, stopLoss)
	buyScore := calculateBuyScore(currentPrice, high5min, avgVolume, currentVolume, movingAvg50, rsi, resistanceLevel, supportLevel)
	sellScore := calculateSellScore(currentPrice, high5min, avgVolume, currentVolume, movingAvg50, rsi, profitTarget, stopLoss)

	recommendation := getTradeRecommendation(buyScore, sellScore)
	fmt.Printf("Buy Score: %d, Sell Score: %d\n", buyScore, sellScore)
	fmt.Printf("Trade Recommendation: %s\n", recommendation)

	evaluation := &dto.TradeEvaluation{
		Symbol:          symbol,
		CurrentPrice:    currentPrice,
		High5Min:        high5min,
		AverageVolume:   avgVolume,
		CurrentVolume:   currentVolume,
		MovingAvg50:     movingAvg50,
		RSI:             rsi,
		ResistanceLevel: resistanceLevel,
		SupportLevel:    supportLevel,
		ProfitTarget:    profitTarget,
		StopLoss:        stopLoss,
		BuyScore:        buyScore,
		SellScore:       sellScore,
		Recommendation:  recommendation,
		EvaluatedAt:     time.Now().In(time.FixedZone("IST", 5*60*60+30*60)),
	}
	err = db.DB.FindAndUpdate(
		evaluation,
		db.M{"symbol": symbol},
		db.M{"$set": evaluation},
		options.FindOneAndUpdate().SetUpsert(true),
	)
	if err != nil {
		log.Println("DB error", err)
	}
}

func fetchNSEData(symbol string) (*dto.NSEData, error) {
	url := fmt.Sprintf("https://www.nseindia.com/api/quote-equity?symbol=%s", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers like the ones used by Chrome
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/json")

	// Manually add cookies from Chrome (example)
	req.Header.Set("Cookie", cookies)

	// Create the HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var data dto.NSEData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func fetchNSEChartData(symbol string) (*dto.NSEChartResponse, error) {
	url := "https://charting.nseindia.com//Charts/ChartData/"

	// TODO: need fix date based on holiday
	// Current time in IST
	ist := time.FixedZone("IST", 5*60*60+30*60)
	nowIST := time.Now().In(ist)
	toDate := nowIST.Unix()       // in seconds
	fromDate := toDate - 24*60*60 // 24 hours ago in seconds

	// Convert to strings if needed
	fromDateStr := strconv.FormatInt(fromDate, 10)
	toDateStr := strconv.FormatInt(toDate, 10)

	log.Println("fromDateStr", fromDateStr, toDateStr, nowIST)

	payload := []byte(`{
		"exch": "N",
		"tradingSymbol": "` + symbol + `-EQ",
		"fromDate": ` + fromDateStr + `,
		"toDate": ` + toDateStr + `,
		"timeInterval": 5,
		"chartPeriod": "I",
		"chartStart": 0
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Request error:", err)
	}

	// Set headers like the ones used by Chrome
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/json")

	// Manually add cookies from Chrome (example)
	req.Header.Set("Cookie", cookies)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP error:", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data dto.NSEChartResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
	}
	return &data, nil
}

func GetLast14DescendingCloses(closePrices []float64) []float64 {
	if len(closePrices) < 20 {
		return closePrices // return as-is if not enough data
	}

	start := len(closePrices) - 20
	if start < 0 {
		start = 0
	}
	last20 := closePrices[start:]
	return last20
}

func calculateRSI(symbol string, period int) float64 {
	c, _ := fetchNSEChartData(symbol)
	closes := GetLast14DescendingCloses(c.C)
	if len(closes) < period+1 {
		return 0.0 // Not enough data
	}

	var gains, losses float64
	for i := 1; i <= period; i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

// Trade Recommendation Function
func getTradeRecommendation(buyScore, sellScore int) string {
	if buyScore >= 8 && sellScore <= 3 {
		return "Strong Buy"
	} else if buyScore >= 5 && buyScore > sellScore {
		return "Buy"
	} else if sellScore >= 8 && buyScore <= 3 {
		return "Strong Sell"
	} else if sellScore >= 5 && sellScore > buyScore {
		return "Sell"
	} else {
		return "Avoid"
	}
}

// Function to calculate buy score
func calculateBuyScore(currentPrice, high5min, avgVolume, currentVolume, movingAvg50, rsi, resistanceLevel, supportLevel float64) int {
	score := 0

	// Price breakout logic
	if currentPrice > high5min {
		score += 3
	} else if currentPrice >= high5min-0.01 {
		score += 1
	}

	// Volume spike check
	volumeFactor := currentVolume / avgVolume
	if volumeFactor >= 1.5 {
		score += 3
	} else if volumeFactor >= 1.0 {
		score += 2
	}

	// Moving average trend confirmation
	if currentPrice > movingAvg50 {
		score += 2
	} else {
		score -= 2
	}

	// RSI strength
	if rsi < 30 {
		score += 2
	} else if rsi > 70 {
		score -= 1
	} else {
		score += 1
	}

	// Resistance and support check
	if currentPrice < resistanceLevel+0.05 {
		score -= 2
	} else if currentPrice > supportLevel-0.05 {
		score += 1
	}

	return score
}

// Function to calculate sell score
func calculateSellScore(currentPrice, high5min, avgVolume, currentVolume, movingAvg50, rsi, profitTarget, stopLoss float64) int {
	score := 0

	// Below breakout high
	if currentPrice < high5min {
		score += 3
	}

	// Volume drop
	volumeFactor := currentVolume / avgVolume
	if volumeFactor < 0.5 {
		score += 2
	} else if volumeFactor < 1.0 {
		score += 1
	}

	// Below 50 MA
	if currentPrice < movingAvg50 {
		score += 2
	}

	// RSI overbought/sold
	if rsi > 70 {
		score += 2
	} else if rsi < 30 {
		score -= 2
	}

	// Hitting targets
	if currentPrice >= profitTarget {
		score += 3
	} else if currentPrice <= stopLoss {
		score += 3
	}

	return score
}
