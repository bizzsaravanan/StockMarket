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
	"math"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type NseService struct{}

const (
	nseBaseURL         = "https://www.nseindia.com"
	nsechartBaseURL    = "https://charting.nseindia.com"
	nseQuoteURL        = "https://www.nseindia.com/api/quote-equity?symbol=%s"
	nseChartURL        = "https://charting.nseindia.com//Charts/ChartData/"
	userAgent          = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
	maxRetryAttempts   = 3
	defaultRSIPeriod   = 14
	averageVolumeRange = 10
)

var nseCookie string

func init() {
	nseCookie = refreshNSECookies()
}

func refreshNSECookies() string {
	return `_ga=GA1.1.1449965907.1742914394; _ga_E0LYHCLJY3=GS1.1.1743048259.3.0.1743048262.0.0.0; _ga_WM2NSQKJEK=GS1.1.1743048263.5.1.1743049265.0.0.0; nseQuoteSymbols=[{"symbol":"ETERNAL","identifier":"","type":"equity"},{"symbol":"WIPRO","identifier":"","type":"equity"},{"symbol":"BAJAJ-AUTO","identifier":"","type":"equity"},{"symbol":"ICIBK0204","identifier":"","type":"equity"},{"symbol":"ICICIBANK","identifier":"","type":"equity"}]; RT="z=1&dm=nseindia.com&si=f09c6216-44fe-4695-836b-3f3edec05350&ss=m9lr2o7x&sl=0&se=8c&tt=0&bcn=%2F%2F684d0d4c.akstat.io%2F&ul=3f8uw&hd=3f8zp"; AKA_A2=A; _abck=279A6C95DA46B78A5A53B5953CC08322~0~YAAQDCozaqsi/RaWAQAAHJdrSA13qX8VYt10pS0Ajf1fnvafnwHDO4919CNt/6p1X/ApBmiSgjD+yNKgw9V0/ZnE1oyDJvpvNMgRZFPS8Yo4n7ahwSy//DEto3UrdhFxDXk5h01blDSC9gP6qwSlGaaOlAxMozzUc6aCwZT2oIC2WqqZrAiop7o9XcQwdgHzEric8jA5nMfw+tT9CL3GvJ06YvUBOtgoHG7rJTVyQd1IwUPR71do3wmsN17wxXCTPN5EIj9IOMefVdKIpiUyTSgOfD+XVsAuYd2vKFGXWN4uSFqbL4TT1VSn5qrmMFXx50hBBAD1nLrS3tfpJRxKB9Y3FodEi3iv/kppBnSpXfKe/i/kgeb9/9pfj9XKurIFvxayam4gpVq2wWGTD4aByaMQ6W6ldK8B1N1dWAqtJyWuHbdZhAz2McAenjsLQ7S8mtJnLvvT1Yo1puZkcMf4pNA5rVsEr+z7OadkJk4crHgOhnyvcr9wgfvhCiZsLrZ2HXW4d6qATiAu/XNaDA7wvimdAvKdd/mNRgV6+pPc~-1~-1~-1; nsit=E85tYYa9Bj1yojgSdAhBWavv; nseappid=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkubnNlIiwiYXVkIjoiYXBpLm5zZSIsImlhdCI6MTc0NDk3MTczNCwiZXhwIjoxNzQ0OTc4OTM0fQ.HJs5YJgKDQfLMsbYhfq62fALRvoIW6OB8wgDcJStwA4; bm_mi=E4342691FC7CD976551FE46D36ED29CB~YAAQDCozal0j/RaWAQAAHaBrSBsx+vLONCwKiowUARAEwIlRFyJ0Q/Pf2zH7UruZmdPJVHhZxMw2OGW22iXEodYYCyfn6kaOYIGFHqEbcZckCHZRN9cmeQWozcw/btx90hfkhNFJm0Eswtcj5evjXcJ4FYdQ1wv1wNMcDXDMIIYuRvZbkqKeg5Cabx+qo00JgUWNspelDRqpAYaj0vClA4802eo6D9+dr16LE82lCujZRD66F4G/gxPoNMgA+uIzdzec361I+Z9ArPSsw3GKMyL6DgA8pCx3qbLOura90jSAIgJAtEyAjWHt4avRIKdNWZ1XQNAEGNip9v7uwO28jco=~1; bm_sz=D95C79D716A19A9F45B3C47D26CF4A98~YAAQDCozal8j/RaWAQAAHaBrSBsOYA8D+EuAsQqdtf8xNL7bK4xiARsGtPmEYsdZtdqMxcDXEH+K9n3Ho8R7aijRBqqQzYMDUUSUu3nTbo+3wGrmtlKheemVeVNQl3BotetIPw1U1C/e8p3TagUqUr6kUHty9YptjBRGQtyKewWbXpKNB5QlzkxQeRVpY/aieR6+U7Pz0jzhbwaSNA9UZ2WBKB597HqR7EtC/q567Sg62NBOsSOME6mjxyOmis/1iZYQc7bXl3dLiAf1XuiItIepk2BCC7FNjD9ciuR5I1IY7PjcErY1Y7U5CKdn8SUIjcflNBSZrGoIc0TbQK4fw9syOTZTf13k7hZbtCA1rNiTxlmvz1EPd6mUvGRIDjq6F9zxAKv0dT8Ns7GQtD75QXI=~4536372~3359302; _ga_87M7PJ3R97=GS1.1.1744971734.17.1.1744971734.60.0.0; ak_bmsc=35ACCD830F67B0BD8C444EBBA5A2CD0D~000000000000000000000000000000~YAAQDCozap0j/RaWAQAAeKNrSBvWu79EofNyNZ6NQl+OQEGEOyHQy+1GtMGsr5gUcGob/yQnIwgqsn/w0+0DzTGJx2ojbqCwPWOtCfKxDbVKM20gEYUITTzFHhce6o2/8C8b1RTZ/GbVm6uN4ZMEaHRhfDovVziVmBL495yhhtDqRL3/I3io0jrE4u0YHv8XuZGeXl9ZE/qTiJ8k5oRa15rUIjxcIerdBSzOMxEKnOxzAISniaW/wAzrDf0XiyIv+6o8CZrQcDkX++f7kGHnpnB+b0fLSSoBWiDcsOFi5FbQpN48Tz56866B175CueO3qrFWJ4Sdmj1QnHYSnPBK8PlksZr47jBTqPrwlHsUbIok29RhIlfjkxiD8CQuoID3smDkQTNNYygwhFprovvuwTEchtr4Ur2IZVIJS+cB9Skgf/S28yb3IGnaTBTxSBSJUoS9PeT5FQgXKJhPJObex+/QtrZtubKYs+kklmVrchjb73oYqjDF3W0UKZtaf9KsEXAOJ/FuJw==; bm_sv=D868F151403D47842AC49403C481961A~YAAQDCozaosn/RaWAQAAmuJrSBspPkv6KnJkINpC2avgZej2BO/uC/Ti+6KBJXzYUF11+WOiDCz5XAHC59c2plZG5fMeF/NgTCiXHM02R0B39pZnb7PEbdOop60yMt/aH7FY7eAA4f3sdzK0qHq1M1lTbP37jGNF9lvVxsqt0PqWNO8vBZD3UDrCTAFmtovNScoFs89fxeN8vTTbOM8bN/lzvE5h+SQu3h5RBDwwiZ19mc6x7fTIprAMA8H/LiS+NFQ=~1`
}

func (s *NseService) EvaluateStock(ctx context.Context, request *dto.Request) (*dto.Response, error) {
	if request.Cookie != "" {
		nseCookie = request.Cookie
	}
	stockName := request.Name
	data, err := fetchNSEQuote(stockName)
	if err != nil {
		return nil, err
	}
	closes, volumes, _, _, currVolume, err := fetchNSEChart(stockName, "5", -6, false)
	if err != nil {
		return nil, err
	}

	avgPrice := calculateAverage(closes)

	currentPrice := data.PriceInfo.LastPrice
	high5min := data.PriceInfo.IntraDayHighLow.Max // Assuming this represents the 5-minute high
	avgVolume := calculateAverage(volumes)         // Placeholder for average volume
	// currVolume := data.SecurityInfo.TotalTradedVolume
	movingAvg50 := data.PriceInfo.VWAP // Placeholder for 50-day moving average
	rsi := calculateRsi(closes, 14)
	resistanceLevel := data.PriceInfo.WeekHighLow.Max
	supportLevel := data.PriceInfo.WeekHighLow.Min
	profitTarget := resistanceLevel
	stopLoss := supportLevel
	buyScore, sellScore := calculateScores(rsi, currentPrice, avgPrice, currVolume, avgVolume)
	recommendation := getRecommendation(buyScore, sellScore)
	if currentPrice < 100 {
		recommendation = "Avoid"
	}

	fmt.Printf("StockName: %s, Buy Score: %d, Sell Score: %d, Trade Recommendation: %s\n", stockName, buyScore, sellScore, recommendation)

	currentStatus := Evaluate1MinContinuation(stockName, recommendation)

	capital := request.Amount
	marginMultiplier := 4.9
	effectiveCapital := capital * marginMultiplier
	quantity := math.Floor(effectiveCapital / currentPrice)

	profitTarget, stopLoss = calculateTargetStopLossRR(currentPrice, recommendation, 0.005)
	charges := calculateGrowwIntradayCharges(currentPrice, profitTarget, quantity, recommendation)

	var grossProfit float64
	if recommendation == "Sell" || recommendation == "Strong Sell" {
		grossProfit = (currentPrice - profitTarget) * quantity
	} else {
		grossProfit = (profitTarget - currentPrice) * quantity
	}

	// Net profit after deducting charges
	netProfit := math.Round((grossProfit-charges)*100) / 100

	var tradeEvaluation *dto.TradeEvaluation
	db.DB.FindOne(&tradeEvaluation, db.M{"symbol": stockName})
	isNew := true
	if tradeEvaluation != nil {
		isNew = false
	}

	evaluation := &dto.TradeEvaluation{
		Name:            data.Info.CompanyName,
		Symbol:          stockName,
		CurrentPrice:    currentPrice,
		High5Min:        high5min,
		AverageVolume:   avgVolume,
		CurrentVolume:   currVolume,
		MovingAvg50:     movingAvg50,
		RSI:             rsi,
		ResistanceLevel: resistanceLevel,
		SupportLevel:    supportLevel,
		ProfitTarget:    profitTarget,
		StopLoss:        stopLoss,
		BuyScore:        buyScore,
		SellScore:       sellScore,
		Recommendation:  recommendation,
		CurrentStatus:   currentStatus,
		Quantity:        quantity,
		Charges:         charges,
		Profit:          netProfit,
		EvaluatedAt:     time.Now().In(time.FixedZone("IST", 5*60*60+30*60)),
		IsNew:           isNew,
		Origin:          request.Origin,
	}
	err = db.DB.FindAndUpdate(
		&dto.TradeEvaluation{},
		db.M{"symbol": stockName},
		db.M{
			"$set": evaluation,
		},
		options.FindOneAndUpdate().SetUpsert(true),
	)
	if err != nil {
		log.Println("DB error", err)
	}
	return nil, nil
}

func calculateGrowwIntradayCharges(bprice, sprice, quantity float64, recommendation string) float64 {
	buyPrice := bprice
	sellPrice := sprice

	if recommendation == "Sell" || recommendation == "Strong Sell" {
		buyPrice = sprice
		sellPrice = bprice
	}

	buyTurnover := buyPrice * quantity
	sellTurnover := sellPrice * quantity
	turnover := buyTurnover + sellTurnover

	// Brokerage: 0.1% per side or ₹20 max, ₹2 min
	calcBrokerage := func(turnover float64) float64 {
		charge := 0.001 * turnover
		if charge < 2 {
			return 2
		}
		if charge > 20 {
			return 20
		}
		return charge
	}
	brokerage := calcBrokerage(buyTurnover) + calcBrokerage(sellTurnover)

	// STT on sell side only
	stt := 0.00025 * sellTurnover

	// Stamp duty on buy side only
	stampDuty := 0.00003 * buyTurnover

	// Exchange txn charges (NSE)
	exchangeTxn := 0.0000297 * turnover

	// SEBI & IPFT charges
	sebi := 0.000001 * turnover
	ipft := 0.000001 * turnover

	// GST on (brokerage + exchangeTxn + sebi + ipft)
	gst := 0.18 * (brokerage + exchangeTxn + sebi + ipft)

	totalCharges := brokerage + stt + stampDuty + exchangeTxn + sebi + ipft + gst

	return math.Round(totalCharges*100) / 100
}

func calculateTargetStopLossRR(currentPrice float64, recommendation string, riskPercent float64) (targetPrice, stopLossPrice float64) {
	if recommendation == "Buy" || recommendation == "Strong Buy" {
		stopLossPrice = currentPrice * (1 - riskPercent)
		risk := currentPrice - stopLossPrice
		targetPrice = currentPrice + 2*risk
	} else if recommendation == "Sell" || recommendation == "Strong Sell" {
		stopLossPrice = currentPrice * (1 + riskPercent)
		risk := stopLossPrice - currentPrice
		targetPrice = currentPrice - 2*risk
	}
	return
}

func retryRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	for i := 0; i < maxRetryAttempts; i++ {
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			return resp, nil
		}
		log.Println(resp)
		log.Println(err)
		time.Sleep(2 * time.Second)
		nseCookie = refreshNSECookies()
	}
	return nil, fmt.Errorf("request failed after %d attempts", maxRetryAttempts)
}

func fetchNSEQuote(stockName string) (*dto.NSEData, error) {
	url := fmt.Sprintf(nseQuoteURL, stockName)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", nseBaseURL)
	req.Header.Set("Cookie", nseCookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/json")

	resp, err := retryRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result dto.NSEData
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Quote unmarshal error:", err)
		log.Println("Raw response:", string(body))
		return nil, err
	}
	return &result, nil
}

func fetchNSEChart(stockName, minute string, days int, isHigh bool) ([]float64, []float64, []float64, []float64, float64, error) {
	// TODO: need fix date based on holiday
	// Current time in IST
	// Define IST timezone
	nowIST := time.Now().UTC()

	toDate := nowIST.AddDate(0, 0, 1).Unix()      // current time in seconds
	fromDate := nowIST.AddDate(0, 0, days).Unix() // days offset

	// Build JSON payload
	payload := map[string]interface{}{
		"exch":          "N",
		"tradingSymbol": stockName + "-EQ",
		"fromDate":      fromDate,
		"toDate":        toDate,
		"timeInterval":  minute,
		"chartPeriod":   "I",
		"chartStart":    0,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Payload marshal error:", err)
		return nil, nil, nil, nil, 0.0, err
	}

	req, err := http.NewRequest("POST", nseChartURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Request error:", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", nsechartBaseURL)
	req.Header.Set("Cookie", nseCookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", nsechartBaseURL)
	client := &http.Client{}
	resp, err := retryRequest(client, req)
	if err != nil {
		return nil, nil, nil, nil, 0.0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result dto.NSEChartResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Chart unmarshal error:", err)
		log.Println("Raw response:", string(body))
		return nil, nil, nil, nil, 0.0, err
	}
	closes, _ := GetLast14DescendingCloses(result.C)
	volume, currentVolume := GetLast14DescendingCloses(result.V)
	if isHigh {
		h, _ := GetLast14DescendingCloses(result.H)
		l, _ := GetLast14DescendingCloses(result.L)
		return closes, volume, h, l, currentVolume, nil
	}
	return closes, volume, nil, nil, currentVolume, nil
}

func CalculateRSI(closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 0
	}
	var gains, losses float64
	for i := 1; i <= period; i++ {
		change := closes[len(closes)-i] - closes[len(closes)-i-1]
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
	return 100 - (100 / (1 + rs))
}

func calculateAverage(slice []float64) float64 {
	if len(slice) == 0 {
		return 0
	}
	var sum float64
	for _, val := range slice {
		sum += val
	}
	return sum / float64(len(slice))
}

func calculateScores(rsi, price, avgPrice, currVolume, avgVolume float64) (int, int) {
	buyScore, sellScore := 0, 0
	if rsi < 30 {
		buyScore += 2
	} else if rsi < 50 {
		buyScore++
	} else if rsi > 70 {
		sellScore += 2
	} else if rsi > 50 {
		sellScore++
	}

	if price > avgPrice*1.1 {
		buyScore += 2
	} else if price > avgPrice*1.03 {
		buyScore++
	} else if price < avgPrice*0.9 {
		sellScore += 2
	} else if price < avgPrice*0.97 {
		sellScore++
	}

	if currVolume > avgVolume*1.5 {
		buyScore += 2
	} else if currVolume > avgVolume*1.2 {
		buyScore++
	} else if currVolume < avgVolume*0.5 {
		sellScore += 2
	} else if currVolume < avgVolume*0.8 {
		sellScore++
	}

	return buyScore, sellScore
}

func getRecommendation(buyScore, sellScore int) string {
	if buyScore >= 5 && sellScore == 0 {
		return "Strong Buy"
	} else if buyScore >= 3 && sellScore <= 1 {
		return "Buy"
	} else if sellScore >= 5 && buyScore == 0 {
		return "Strong Sell"
	} else if sellScore >= 3 && buyScore <= 1 {
		return "Sell"
	}
	return "Avoid"
}

func GetLast14DescendingCloses(closePrices []float64) ([]float64, float64) {
	if len(closePrices) < 20 {
		if len(closePrices) == 0 {
			return closePrices, 0.0
		}
		return closePrices, closePrices[0] // return as-is if not enough data
	}

	start := len(closePrices) - 20
	if start < 0 {
		start = 0
	}
	last20 := closePrices[start:]
	for i, j := 0, len(last20)-1; i < j; i, j = i+1, j-1 {
		last20[i], last20[j] = last20[j], last20[i]
	}
	return last20, last20[0]
}

func calculateRsi(closes []float64, period int) float64 {
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

func Evaluate1MinContinuation(stockName, direction string) string {
	// Fetch last 5 minutes of 1-min chart data
	closes, _, highs, lows, _, err := fetchNSEChart(stockName, "1", -5, true)

	if err != nil || len(closes) < 2 || len(highs) < 2 || len(lows) < 2 {
		log.Println("Insufficient data or error fetching chart:", err)
		return "HOLD"
	}

	// Ensure data consistency
	latestC := closes[len(closes)-1]
	prevH := highs[len(highs)-2]
	prevL := lows[len(lows)-2]

	// Small buffer to reduce false HOLDs due to minor fluctuations
	buffer := 0.001 // 0.1%

	switch direction {
	case "Buy", "Strong Buy":
		if latestC > prevH*(1+buffer) {
			return "BUY" // breakout above previous high
		} else if latestC < prevL*(1-buffer) {
			return "SELL" // strong reversal signal
		}
	case "Sell", "Strong Sell":
		if latestC < prevL*(1-buffer) {
			return "SELL" // breakdown below previous low
		} else if latestC > prevH*(1+buffer) {
			return "BUY" // reversal
		}
	}

	return "HOLD"
}

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}
