package main

const (
	ASSET_CRYPTO = iota
	ASSET_STOCK  = iota
)

type Asset struct {
	price float64
	count uint
	class uint
}

func newAsset(price float64, count uint, class uint) *Asset {
	asset := &Asset{
		price: price,
		count: count,
		class: class,
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
