package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No stock ticker provided. Exiting program.")
		os.Exit(1)
	}

	vanguardUrl := "investor.vanguard.com"
	stockTicker := os.Args[1]
	vanguardTickerPath := fmt.Sprintf("/vmf/api/%s/delayedPrice", stockTicker)
	url := fmt.Sprintf("https://%s/%s", vanguardUrl, vanguardTickerPath)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Could not get stock ticker data.")
	}

	var apiResponse Response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Can not unmarshal JSON.")
	}

	var result Result
	for _, quote := range apiResponse.Quotes {
		result.Ticker = stockTicker
		result.Price = quote.Equity.Pricing.AskPrice
		result.Status = "200"
	}

	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
	//fmt.Printf("Getting price for stock ticker %s.\n", stockTicker)

}

type Result struct {
	Ticker string `json:"ticker"`
	Price  string `json:"price"`
	Status string `json:"status"`
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
