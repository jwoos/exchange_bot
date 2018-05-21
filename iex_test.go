package assets


import (
	"encoding/json"
	"testing"
)


func TestMarketBatchQuote(t *testing.T) {
	raw := `{"FB":{"quote":{"symbol":"FB","companyName":"Facebook Inc.","primaryExchange":"Nasdaq Global Select","sector":"Technology","calculationPrice":"tops","open":183.49,"openTime":1526650200838,"close":183.76,"closeTime":1526587200340,"high":184.19,"low":183,"latestPrice":182.93,"latestSource":"IEX real time price","latestTime":"10:00:59 AM","latestUpdate":1526652059676,"latestVolume":1823268,"iexRealtimePrice":182.93,"iexRealtimeSize":100,"iexLastUpdated":1526652059676,"delayedPrice":183.244,"delayedPriceTime":1526651173683,"previousClose":183.76,"change":-0.83,"changePercent":-0.00452,"iexMarketPercent":0.0153,"iexVolume":27896,"avgTotalVolume":28037884,"iexBidPrice":182.28,"iexBidSize":100,"iexAskPrice":183.15,"iexAskSize":100,"marketCap":529514253942,"peRatio":26.9,"week52High":195.32,"week52Low":144.51,"ytdChange":0.01289824716128323}},"MSFT":{"quote":{"symbol":"MSFT","companyName":"Microsoft Corporation","primaryExchange":"Nasdaq Global Select","sector":"Technology","calculationPrice":"tops","open":96.3,"openTime":1526650200750,"close":96.18,"closeTime":1526587200443,"high":96.84,"low":96.01,"latestPrice":96.43,"latestSource":"IEX real time price","latestTime":"10:01:09 AM","latestUpdate":1526652069339,"latestVolume":3451062,"iexRealtimePrice":96.43,"iexRealtimeSize":100,"iexLastUpdated":1526652069339,"delayedPrice":96.65,"delayedPriceTime":1526651171334,"previousClose":96.18,"change":0.25,"changePercent":0.0026,"iexMarketPercent":0.01318,"iexVolume":45485,"avgTotalVolume":27386643,"iexBidPrice":96.3,"iexBidSize":200,"iexAskPrice":96.92,"iexAskSize":100,"marketCap":740890735214,"peRatio":28.53,"week52High":98.69,"week52Low":67.14,"ytdChange":0.12915155143873494}}}`

	var mb IEXMarketBatch

	t.Run(
		"Marshal Quote",
		func(innerT *testing.T) {
			err := json.Unmarshal([]byte(raw), &mb.Batch)
			if err != nil {
				innerT.Errorf("Error thrown: %v", err)
			}

			if len(mb.Batch) != 2 {
				innerT.Errorf("Wrong number of keys in map: %d", len(mb.Batch))
			}

			for _, mbt := range mb.Batch {
				if mbt.Price != nil {
					innerT.Errorf("Price should be nil")
				}

				if mbt.Quote == nil {
					innerT.Errorf("Quote should not be nil")
				} else {
					if mbt.Quote.Symbol == "" {
						innerT.Errorf("Quote should be filled in")
					}
				}
			}
		},
	)

	// TODO fix skipped test
	t.Run(
		"Unmarshal Quote",
		func(innerT *testing.T) {
			innerT.Skip("Skipping due to float inconsistency")

			marshalled, err := json.Marshal(mb.Batch)
			if err != nil {
				innerT.Errorf("Error marshalling: %v", err)
			}

			if string(marshalled) != raw {
				innerT.Log(string(marshalled))
				innerT.Log(raw)
				innerT.Errorf("Not equal")
			}
		},
	)
}

func TestMarketBatchPrice(t *testing.T) {
	raw := `{"FB":{"price":182.68},"MSFT":{"price":96.36}}`

	var mb IEXMarketBatch

	t.Run(
		"Marshal Price",
		func(innerT *testing.T) {
			err := json.Unmarshal([]byte(raw), &mb.Batch)
			if err != nil {
				innerT.Errorf("Error thrown: %v", err)
			}

			if len(mb.Batch) != 2 {
				innerT.Errorf("Wrong number of keys in map: %d", len(mb.Batch))
			}

			for _, mbt := range mb.Batch {
				if mbt.Price == nil {
					innerT.Errorf("Price should not be nil")
				} else {
					innerT.Log(*mbt.Price)
				}

				if mbt.Quote != nil {
					innerT.Errorf("Quote should be nil")
				}
			}
		},
	)

	t.Run(
		"Unmarshal Price",
		func(innerT *testing.T) {
			marshalled, err := json.Marshal(mb.Batch)
			if err != nil {
				innerT.Errorf("Error marshalling: %v", err)
			}

			if string(marshalled) != raw {
				innerT.Log(string(marshalled))
				innerT.Log(raw)
				innerT.Errorf("Not equal")
			}
		},
	)
}

func TestMarketBatchFetch(t *testing.T) {
	var mb IEXMarketBatch
	var err error

	t.Run(
		"Fetch Quote",
		func(innerT *testing.T) {
			err = mb.Fetch(IEXRequest{
				Information: []string{"quote"},
				Symbols: []string{"FB", "MSFT"},
			})

			if err != nil {
				innerT.Errorf("Error fetching data: %v", err)
			}

			if len(mb.Batch) != 2 {
				innerT.Errorf("Error fetching data: %v", err)
			}

			for _, v := range mb.Batch {
				if v.Quote == nil {
					innerT.Errorf("Quote should be nil")
				}

				if v.Price != nil {
					innerT.Errorf("Price should be nil")
				}
			}
		},
	)

	t.Run(
		"Fetch Price",
		func(innerT *testing.T) {
			err = mb.Fetch(IEXRequest{
				Information: []string{"price"},
				Symbols: []string{"FB", "MSFT"},
			})

			if err != nil {
				innerT.Errorf("Error fetching data: %v", err)
			}

			if len(mb.Batch) != 2 {
				innerT.Errorf("Error fetching data: %v", err)
			}

			for _, v := range mb.Batch {
				if v.Quote != nil {
					innerT.Errorf("Quote should be nil")
				}

				if v.Price == nil {
					innerT.Errorf("Price should not be nil")
				}
			}
		},
	)

	t.Run(
		"Fetch Quote and Price",
		func(innerT *testing.T) {
			err = mb.Fetch(IEXRequest{
				Information: []string{"price", "quote"},
				Symbols: []string{"FB", "MSFT"},
			})

			if err != nil {
				innerT.Errorf("Error fetching data: %v", err)
			}

			if len(mb.Batch) != 2 {
				innerT.Errorf("Error fetching data: %v", err)
			}

			for _, v := range mb.Batch {
				if v.Quote == nil {
					innerT.Errorf("Quote should not be nil")
				}

				if v.Price == nil {
					innerT.Errorf("Price should not be nil")
				}
			}
		},
	)
}
