package main

import (
	"fmt"
)

const (
	ASSET_CRYPTO = iota
	ASSET_STOCK  = iota
)

type Asset struct {
	symbol string
	count float64
}

func newAsset(symbol string, count float64) *Asset {
	asset := &Asset{
		symbol: symbol,
		count: count,
	}

	return asset
}

type Portfolio struct {
	cryptocurrency map[string]*Asset
	stock          map[string]*Asset
}

func newPortfolio() *Portfolio {
	portfolio := &Portfolio{
		cryptocurrency: make(map[string]*Asset),
		stock:          make(map[string]*Asset),
	}

	return portfolio
}

func (p *Portfolio) addStock(symbol string, count float64) (float64, error) {
	asset, ok := p.stock[symbol]
	if !ok {
		asset = newAsset(symbol, count)
		p.stock[symbol] = asset
	} else {
		asset.count += count
	}

	return asset.count, nil
}

func (p *Portfolio) removeStock(symbol string, count float64) (float64, error) {
	asset, ok := p.stock[symbol]
	if !ok {
		return 0, fmt.Errorf("You do not have any %s stocks", symbol)
	}

	if count > asset.count {
		return 0, fmt.Errorf("You have %.0f of %s stocks and are trying to sell %0.f", asset.count, symbol, count)
	}

	asset.count -= count

	return asset.count, nil
}

func (p *Portfolio) addCrypto(symbol string, count float64) (float64, error) {
	asset, ok := p.cryptocurrency[symbol]
	if !ok {
		asset = newAsset(symbol, count)
		p.cryptocurrency[symbol] = asset
	} else {
		asset.count += count
	}

	return asset.count, nil
}

func (p *Portfolio) removeCrypto(symbol string, count float64) (float64, error) {
	asset, ok := p.cryptocurrency[symbol]
	if !ok {
		return 0, fmt.Errorf("You do not have any %s stocks", symbol)
	}

	if count > asset.count {
		return 0, fmt.Errorf("You have %.2f of %s cryptos and are trying to sell %2.f", asset.count, symbol, count)
	}

	asset.count -= count

	return asset.count, nil
}
