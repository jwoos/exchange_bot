package assets


import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)


type IEXRequest struct {
	Information []string
	Symbols []string
}


type IEXMarketBatch struct {
	Batch map[string]*IEXMarketBatchTypes
}

type IEXMarketBatchTypes struct {
	Quote *IEXQuote  `json:"quote,omitempty"`
	Price *float64 `json:"price,omitempty"`
}


type IEXQuote struct {
	Symbol           string  `json:"symbol"`
	CompanyName      string  `json:"companyName"`
	PrimaryExchange  string  `json:"primaryExchange"`
	Sector           string  `json:"sector"`
	CalculationPrice string  `json:"calculationPrice"`
	Open             float64 `json:"open"`
	OpenTime         uint    `json:"openTime"`
	Close            float64 `json:"close"`
	CloseTime        uint    `json:"closeTime"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	LatestPrice      float64 `json:"latestPrice"`
	LatestSource     string  `json:"latestSource"`
	LatestTime       string  `json:"latestTime"`
	LatestUpdate     uint    `json:"latestUpdate"`
	LatestVolume     uint    `json:"latestVolume"`
	IEXRealtimePrice float64 `json:"iexRealtimePrice"`
	IEXRealtimeSize  uint    `json:"iexRealtimeSize"`
	IEXLastUpdated   uint    `json:"iexLastUpdated"`
	DelayedPrice     float64 `json:"delayedPrice"`
	DelayedPriceTime uint    `json:"delayedPriceTime"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`
	IEXMarketPercent float64 `json:"iexMarketPercent"`
	IEXVolume        uint    `json:"iexVolume"`
	AvgTotalVolume   uint    `json:"avgTotalVolume"`
	IEXBidPrice      float64 `json:"iexBidPrice"`
	IEXBidSize       uint    `json:"iexBidSize"`
	IEXAskPrice      float64 `json:"iexAskPrice"`
	IEXAskSize       uint    `json:"iexAskSize"`
	MarketCap        uint    `json:"marketCap"`
	PeRatio          float64 `json:"peRatio"`
	Week52High       float64 `json:"week52High"`
	Week52Low        float64 `json:"week52Low"`
	YtdChange        float64 `json:"ytdChange"`
}

func (mb *IEXMarketBatch) Fetch(config IEXRequest) error {
	req, err := http.NewRequest("GET", IEX_API_BASE, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("symbols", strings.Join(config.Symbols, ","))
	q.Add("types", strings.Join(config.Information, ","))

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &mb.Batch)
	if err != nil {
		return err
	}

	return nil
}
