package assets


import (
	"encoding/json"
	"testing"
)


func TestCCMulti(t *testing.T) {
	raw := `{"BTC":{"USD":8413.95},"ETH":{"USD":696.3}}`

	var multi CCMulti

	t.Run(
		"Unmarshal Multi",
		func(innerT *testing.T) {
			err := json.Unmarshal([]byte(raw), &multi.Batch)
			if err != nil {
				innerT.Errorf("Error thrown: %v", err)
			}

			if len(multi.Batch) != 2 {
				innerT.Errorf("Wrong number of keys in map: %d", len(multi.Batch))
			}

			for _, sym := range multi.Batch {
				_, okay := sym["USD"]
				if !okay {
					innerT.Errorf("Key not found")
				}
			}
		},
	)

	t.Run(
		"Marshal Multi",
		func(innerT *testing.T) {
			marshalled, err := json.Marshal(multi.Batch)
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

func TestCCMultiFetch(t *testing.T) {
	var multi CCMulti
	var err error

	t.Run(
		"Fetch Multi",
		func(innerT *testing.T) {
			err = multi.Fetch(CCRequest{
				FromSymbols: []string{"BTC", "ETH"},
				ToSymbols: []string{"USD"},
			})
			if err != nil {
				innerT.Errorf("Error fetching data: %v", err)
			}

			if len(multi.Batch) != 2 {
				innerT.Errorf("Error fetching data: %v", err)
			}

			for _, v := range multi.Batch {
				_, okay := v["USD"]
				if !okay {
					innerT.Errorf("Missing key")
				}
			}
		},
	)
}
