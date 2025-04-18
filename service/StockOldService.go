package service

import (
	"StockMarket/db"
	"StockMarket/dto"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	cookies          string = `_ga=GA1.1.1449965907.1742914394; _ga_E0LYHCLJY3=GS1.1.1743048259.3.0.1743048262.0.0.0; _ga_WM2NSQKJEK=GS1.1.1743048263.5.1.1743049265.0.0.0; AKA_A2=A; _abck=279A6C95DA46B78A5A53B5953CC08322~0~YAAQl7xWaFSdrz2WAQAAS0c2RQ0KTC9gzAycfOckWaLNkRKsCjTcDRcGNVHDfykb10D6xA5E4TyhmM5SZPd0/CpZuMx7puIHGZNHfyxT/DeVJqnzUjFfXR4G5gD303FhW01V3QBK0JavrBOK9dot4tPqGTUmbsy1zGf4ofi9A2sGs0B2xqWKH7uqD2Rr8IBu2nqjzaxODGClCdLTrmQ9KM6e//kffXoYtZvC1/0Ex4i0RB20MEx3PYDJhDTXYrYSPbzsYH/70BCe62mQkJBtmmJv6XW0xGdBlHyyqaTOgAcmRgayUs+Cplu0AGHOQwhGZ3tnw2qBpDo2qH1BSGmMOL4ss2o+leX709P/aH4bITiZhw4QpS/ZhCNtx+ZhDH6bJfWtllpL6uZsilHUtoEIe9HPBPtxdhxta1yE0ICRGWy51yjlINUNYWUAlSK0jj4Xajcr1bl862sIIwGANzU8v6Tiae+Zr/s4Z3ZbNUqkQ9l/Dw8ftvyyOnlo51U110tmAF+E1f++S3k7+0WDmTiyOTIGoUt80pEtBGzkcqww~-1~-1~-1; bm_mi=D773E6F0AC534C29135517E8A179C912~YAAQl7xWaHGdrz2WAQAAj0w2RRvGOJ2G8v+AGIjY8GuMW96kTMfP+9t9smmTFepPIRxzAnCtmChkv3clGl1B2gWI2Ius911fyH0S8TdDAwo/7I1T6GMg9Pu090gY0n28hjH2T9XaZ+/doPgU9TtCFSvFVo8YyACLYynCXe1wExaw9wreEabQjfiK1m8jFMP9dq/cEGmzYWs2Magak7Ar2CSUr40wo7SbvkVGa4uXZ4oFrE8vQ++BCd/cnwsubJaLu/SL3cB2bD6vWIXFnhQ8owb/T/T+uw+6jkmaT9obo3iqE15etJiKD07Adn1fnhIMrt824Y08OlN+q7og362riHkNsdgtkmJE0hEMTx8zvpAHY3Yy9ofWexs=~1; nsit=qSpuPIOC80EOcafTzEXN_X3c; ak_bmsc=9FB46853F07FAE24D756E5A810D2D804~000000000000000000000000000000~YAAQl7xWaKidrz2WAQAAsGg2RRsyWRbdQvgwhHUaT9fNQCngHIHPBtyy+L8TolPkkPlbbG1f1gxcvSG2MPxzwDIb38wxJr6bhtN7BbBY9NbP7/xiWQ+TIUeUXA5OCQTNa4ged7iAw8SMS3uNCDyQFT6AfqcaTGgQEMf1pQ0fZUof4tK58X7V0l87TE7E5wSDUvL+dHWikTSozJsviqLkry257wVII0pTC0oYX+nvwRAYJMb8a8Y1Vbqe4IQEvg1+Sr2K5b1an0nvsmrikIdx2FydKrxMF59PoXpGyUqfHuXraN8h9H9pn6kbVaZxcZUNmcSXScD4yTU7AVlTqDiTz0Vyme1iZ1lWNzZJC/4E2+XDuvdwGJiAQjO2IW6ifw3TrM6FR42GG/wF/Xa0VRRHaSiULL/YquWLD+9MfPiN6CSM1f5vsAyOtxfw0sis2P/ryUxaemFO3R6WzB8V3QyX/BBAcJSYCoTIsTUq9bFJsMGQSmXrQO76p+IhuY5poCU/QFXDpowRyIjNV884Mq81N7vjx6iMAoRpMw==; RT="z=1&dm=nseindia.com&si=f09c6216-44fe-4695-836b-3f3edec05350&ss=m9lr2o7x&sl=1&se=8c&tt=t4&bcn=%2F%2F684d0d41.akstat.io%2F&ld=2av&nu=evuzm1q&cl=9aw"; nseQuoteSymbols=[{"symbol":"FUSI-RE","identifier":"","type":"equity"},{"symbol":"BHARTIARTL","identifier":"","type":"equity"},{"symbol":"DRSCARGO","identifier":"","type":"equity"},{"symbol":"BAJAJFINSV","identifier":"","type":"equity"},{"symbol":"BAJAJ-AUTO","identifier":"","type":"equity"}]; nseappid=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkubnNlIiwiYXVkIjoiYXBpLm5zZSIsImlhdCI6MTc0NDkxNzk0MCwiZXhwIjoxNzQ0OTI1MTQwfQ.22gzPU-_0m_PWNZma1mZ1XuDs5fdSnnwVC0n6poIe_Q; bm_sz=6D23CFCD57E456A09F50B1F8F02EB182~YAAQl7xWaEaerz2WAQAAX8g2RRsAtaUnn/ysl9oCrX5M7E6q+F/eaUCe1V8VewRQp8rn2NaA0J4FTheglANL/HC7kZvKTy7/4utVlhTljK0oM6psGsGUM+9s0GZ8+wkGOTIrYPGHyQe9UEK/qE0DdXK3ubD54uwobMyJ/+TpBBiw3y/vbWvzpsUEDxZfTK8QSoKMyX5Fu0mf2rcqBQWK0hRU83klM+J8HMj3x7rfC90v6L9iAfj8jdK/Fdl0y/h99+a1Ig80PdIq0vw/Vu1iywKxs95WonH4mB0QuCrE3JC07tailLMPtlUpkck8DYn28a9b2MZea+1UabKb2uR5vdwMZHFShkqWfk8QjQxfx65SyyHq3dKWmd7PUrDMVWRqEKJRHmWIy0s5wkZaKLmf+ccUdzOPXKsSFRzKiRxW6Q==~3228983~4601666; _ga_87M7PJ3R97=GS1.1.1744917909.15.1.1744917942.27.0.0; bm_sv=003B8B9EF40D47E78379FA7F6D56A376~YAAQl7xWaE+erz2WAQAAU842RRutbZDXoTAfehKJAwh3ib7FgYFcv3errRyGEMlgOb7/HfG3JxwgIjmGCenH+W52SZ0ozE7eQttEsgrhJP4X0Q6Cgu59G3VCakqhjhJOWmXr1dP6N10VuFfSVm8058bC71Fhdu/oOEtY9pHqHbpibyvwDBlfPTw/5GXJNKtNYF6Y6dPwgXz2OShyTlRniQxhLm/gdSLpReno2Kq7Cvo344j/BwvmmB2x+yphzl6eD15J~1`
	cookieExpiryDate time.Time
)

func (s *NseService) CreateStockData(ctx context.Context, request *dto.Request) {

	/* if time.Now().After(cookieExpiryDate) {
		fmt.Println("Cookie expired, extending 30 minutes...")
		getCookie()
		cookieExpiryDate = time.Now().Add(30 * time.Minute)
	} */
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
	rsi, _ := calculateRSI(symbol, 14)
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

func calculateRSI(symbol string, period int) (float64, float64) {
	c, _ := fetchNSEChartData(symbol)
	closes, currentVolume := GetLast14DescendingCloses(c.C)
	if len(closes) < period+1 {
		return 0.0, currentVolume // Not enough data
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
		return 100, currentVolume
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi, currentVolume
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

// ---- Part 1: _ga format ----
func generateGA() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("GA1.1.%d.%d", rand.Int63n(1e10), rand.Int63n(1e10))
}

func generateSubGA() string {
	return fmt.Sprintf("GS1.1.%d.%d.%d.%d.0.0.0", rand.Int63n(1e10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}

// ---- Part 2: Random token generator ----
func randomToken(size int) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// ---- Part 3: JWT token ----
func generateJWT(secret string) (string, error) {
	claims := jwt.MapClaims{
		"iss": "api.nse",
		"aud": "api.nse",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ---- Part 4: nseQuoteSymbols ----
type Symbol struct {
	Symbol     string `json:"symbol"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

func generateQuoteSymbols() string {
	symbols := []Symbol{
		{"FUSI-RE", "", "equity"},
		{"BHARTIARTL", "", "equity"},
		{"DRSCARGO", "", "equity"},
		{"BAJAJFINSV", "", "equity"},
	}
	data, _ := json.Marshal(symbols)
	return string(data)
}

// ---- Part 5: RT Token ----
func generateRT() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	randStr := func(length int) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = chars[rand.Intn(len(chars))]
		}
		return string(b)
	}
	return fmt.Sprintf("z=1&dm=nseindia.com&si=%s-%s-%s-%s-%s&ss=%s&sl=2&se=8c&tt=23c&bcn=%%2F%%2F%s.akstat.io%%2F&ld=%s&nu=%s&cl=%s&ul=%s&hd=%s",
		randStr(8), randStr(4), randStr(4), randStr(4), randStr(12),
		randStr(8), randStr(8), randStr(8), randStr(3), randStr(3), randStr(4))
}

func getCookie() {
	rand.Seed(time.Now().UnixNano())
	jwtToken, _ := generateJWT("dummy_secret_key")
	cookies = "_ga=" + generateGA() +
		";_ga_E0LYHCLJY3=" + generateSubGA() +
		";_ga_WM2NSQKJEK=" + generateSubGA() +
		";nseQuoteSymbols=" + generateQuoteSymbols() +
		";AKA_A2=A" +
		";nsit=" + randomToken(16) +
		";RT=" + generateRT() +
		";nseappid=" + jwtToken +
		";_ga_87M7PJ3R97=" + generateSubGA() +
		";_abck=" + randomToken(64) +
		";ak_bmsc=" + randomToken(64) +
		";bm_sv=" + randomToken(64) +
		";bm_sz=" + randomToken(64)
}
