package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var (
	ticker string
	qty    float64
)

func main() {
	flag.StringVar(&ticker, "ticker", "", "stock ticker")
	flag.Float64Var(&qty, "qty", 1, "stock quantity")
	flag.Parse()

	var result Result

	if len(ticker) == 0 {
		printError(&result, "No stock ticker provided. Exiting program.")
		os.Exit(1)
	}

	vanguardUrl := "investor.vanguard.com"
	vanguardTickerPath := fmt.Sprintf("/vmf/api/%s/delayedPrice", ticker)
	url := fmt.Sprintf("https://%s/%s", vanguardUrl, vanguardTickerPath)
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	client := &http.Client{
		Transport: &http.Transport{},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		printError(&result, "Could not read stock ticker data.")
		os.Exit(0)
	}

	var apiResponse Response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		printError(&result, "Could not unmarshal stock ticker data JSON.")
		os.Exit(0)
	}

	if len(result.Error) == 0 {
		quote := apiResponse.Quotes[0]
		price, _ := strconv.ParseFloat(quote.Equity.Pricing.AskPrice, 64)

		result.Success = true
		result.Ticker = ticker
		result.Price = price
		result.Quantity = qty
		result.Balance = qty * price
	}

	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func printError(result *Result, error string) {
	result.Success = false
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}

type Result struct {
	Ticker   string  `json:"ticker,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Success  bool    `json:"success,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
	Error    string  `json:"error,omitempty"`
}

type Response struct {
	Quotes []struct {
		Cusip       string `json:"cusip"`
		DisplayName string `json:"displayName"`
		Equity      struct {
			DividendPayDateTime   string `json:"dividendPayDateTime"`
			DividendPerShare      string `json:"dividendPerShare"`
			EarningsDateIndicator string `json:"earningsDateIndicator"`
			ExDividendDateTime    string `json:"exDividendDateTime"`
			InvestmentSubType     string `json:"investmentSubType"`
			IsCaveatEmptor        bool   `json:"isCaveatEmptor"`
			PreferredSecurityData struct {
			} `json:"preferredSecurityData"`
			Pricing struct {
				AdjustedClosePrice      string `json:"adjustedClosePrice"`
				AskPrice                string `json:"askPrice"`
				BidPrice                string `json:"bidPrice"`
				CloseAskPrice           string `json:"closeAskPrice"`
				CloseBidPrice           string `json:"closeBidPrice"`
				CurrencyCode            string `json:"currencyCode"`
				DayHighPrice            string `json:"dayHighPrice"`
				DayLowPrice             string `json:"dayLowPrice"`
				LastTradePrice          string `json:"lastTradePrice"`
				NetChange               string `json:"netChange"`
				OfficialClosePrice      string `json:"officialClosePrice"`
				OpenPrice               string `json:"openPrice"`
				PreviousClosePrice      string `json:"previousClosePrice"`
				PriceChangePercentToday string `json:"priceChangePercentToday"`
				PriceDirection          string `json:"priceDirection"`
				YearHighPrice           string `json:"yearHighPrice"`
				YearLowPrice            string `json:"yearLowPrice"`
			} `json:"pricing"`
			TradingDates struct {
				LastTradeDateTime     string `json:"lastTradeDateTime"`
				NextEarningsDateTime  string `json:"nextEarningsDateTime"`
				PreviousCloseDateTime string `json:"previousCloseDateTime"`
				YearHighDateTime      string `json:"yearHighDateTime"`
				YearLowDateTime       string `json:"yearLowDateTime"`
			} `json:"tradingDates"`
			TradingExchanges struct {
				AskPriceExchange  string `json:"askPriceExchange"`
				BidPriceExchange  string `json:"bidPriceExchange"`
				LastTradeExchange string `json:"lastTradeExchange"`
			} `json:"tradingExchanges"`
			TradingHalt                 bool   `json:"tradingHalt"`
			TradingHaltLimitUpLimitDown bool   `json:"tradingHaltLimitUpLimitDown"`
			UpcIndicator                string `json:"upcIndicator"`
			Volume                      struct {
				AskSize                  string `json:"askSize"`
				BidSize                  string `json:"bidSize"`
				LastTradeVolume          string `json:"lastTradeVolume"`
				TodaysVolume             string `json:"todaysVolume"`
				VolumemovingAverage10Day string `json:"volumemovingAverage10Day"`
				VolumemovingAverage25Day string `json:"volumemovingAverage25Day"`
				VolumemovingAverage50Day string `json:"volumemovingAverage50Day"`
			} `json:"volume"`
			Yield string `json:"yield"`
		} `json:"equity"`
		ProductType string `json:"productType"`
		Ticker      string `json:"ticker"`
		Vendor      string `json:"vendor"`
	} `json:"quotes"`
	Errors []any `json:"errors"`
}
