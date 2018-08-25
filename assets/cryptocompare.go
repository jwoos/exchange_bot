package assets


import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)


type CCRequest struct {
	FromSymbols []string
	ToSymbols []string
}


type CCMulti struct {
	Batch map[string]map[string]float32
}


func (multi *CCMulti) Fetch(config CCRequest) error {
	req, err := http.NewRequest("GET", CRYPTOCOMPARE_API_BASE, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("fsyms", strings.Join(config.FromSymbols, ","))
	q.Add("tsyms", strings.Join(config.ToSymbols, ","))

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

	err = json.Unmarshal(body, &multi.Batch)
	if err != nil {
		return err
	}

	return nil
}
