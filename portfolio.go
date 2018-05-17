package main


const (
	ASSET_CRYPTO = iota
	ASSET_STOCK = iota
)


type Portfolio struct {
	cryptocurrency map[string][]Asset
	stock map[string][]Asset
}


type Asset struct {
	price uint
	count uint
	class uint
}
