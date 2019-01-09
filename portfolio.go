package main

const (
	ASSET_CRYPTO = iota
	ASSET_STOCK  = iota
)

type Asset struct {
	price float64
	symbol string
	count float64
}

func newAsset(symbol string, price float64, count float64) *Asset {
	asset := &Asset{
		symbol: symbol,
		price: price,
		count: count,
	}

	return asset
}

type Portfolio struct {
	cryptocurrency map[string][]*Asset
	stock          map[string][]*Asset
}

func newPortfolio() *Portfolio {
	portfolio := &Portfolio{
		cryptocurrency: make(map[string][]*Asset),
		stock:          make(map[string][]*Asset),
	}

	return portfolio
}

func (p *Portfolio) appendStock(assets ...*Asset) {
	for _, asset := range assets {
		p.stock[asset.symbol] = append(p.stock[asset.symbol], asset)
	}
}

func (p *Portfolio) appendCrypto(assets ...*Asset) {
	for _, asset := range assets {
		p.cryptocurrency[asset.symbol] = append(p.cryptocurrency[asset.symbol], asset)
	}
}
