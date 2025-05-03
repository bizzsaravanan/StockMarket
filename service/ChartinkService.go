package service

import (
	"StockMarket/dto"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ChartinkService struct{}

func (s *ChartinkService) CreateStockData(ctx context.Context, request *dto.Request) (*dto.Response, error) {
	// Step 1: Fetch the CSRF token
	client := &http.Client{}
	getReq, _ := http.NewRequest("GET", "https://chartink.com/screener/process", nil)
	getResp, err := client.Do(getReq)
	if err != nil {
		fmt.Println("Error fetching CSRF token:", err)
		return nil, err
	}
	defer getResp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(getResp.Body)
	bodyString := string(bodyBytes)

	// Extract the CSRF token from the HTML
	re := regexp.MustCompile(`name="csrf-token" content="(.*?)"`)
	matches := re.FindStringSubmatch(bodyString)
	if len(matches) < 2 {
		fmt.Println("CSRF token not found")
		return nil, err
	}
	csrfToken := matches[1]

	postReq, err := http.NewRequest("POST", "https://chartink.com/screener/process", strings.NewReader(request.FormData.Encode()))
	if err != nil {
		fmt.Println("CSRF token not found")
		return nil, err
	}

	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	postReq.Header.Set("X-CSRF-TOKEN", csrfToken)
	postReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	postReq.Header.Set("Referer", "https://chartink.com/screener/15-min-stock-breakouts-with-high-volume")
	postReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	// Copy cookies from GET response (important!)
	for _, cookie := range getResp.Cookies() {
		postReq.AddCookie(cookie)
	}

	// Send the POST request
	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("POST request failed:", err)
		return nil, err
	}
	defer postResp.Body.Close()

	postBody, _ := ioutil.ReadAll(postResp.Body)

	var result *dto.ChartinkResponse
	err = json.Unmarshal(postBody, &result)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}
	return &dto.Response{ChartinkResponse: result}, nil
}
