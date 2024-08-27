package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

var coinNames []string = []string{"bitcoin", "ethereum"}
var CurrentCoinsValues []Coin
var LastUpdateTime time.Time

type Coin struct {
	ID                           string      `json:"id"`
	Symbol                       string      `json:"symbol"`
	Name                         string      `json:"name"`
	Image                        string      `json:"image"`
	CurrentPrice                 float64     `json:"current_price"`
	MarketCap                    int64       `json:"market_cap"`
	MarketCapRank                int64       `json:"market_cap_rank"`
	FullyDilutedValuation        int64       `json:"fully_diluted_valuation"`
	TotalVolume                  float64     `json:"total_volume"`
	High24h                      float64     `json:"high_24h"`
	Low24h                       float64     `json:"low_24h"`
	PriceChange24h               float64     `json:"price_change_24h"`
	PriceChangePercentage24h     float64     `json:"price_change_percentage_24h"`
	MarketCapChange24h           float64     `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64     `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64     `json:"circulating_supply"`
	TotalSupply                  float64     `json:"total_supply"`
	MaxSupply                    float64     `json:"max_supply"`
	ATH                          float64     `json:"ath"`
	ATHChangePercentage          float64     `json:"ath_change_percentage"`
	ATHDate                      string      `json:"ath_date"`
	ATL                          float64     `json:"atl"`
	ATLChangePercentage          float64     `json:"atl_change_percentage"`
	ATLDate                      string      `json:"atl_date"`
	ROI                          interface{} `json:"roi"`
	LastUpdated                  string      `json:"last_updated"`
}

func SetCoinsRequest() error { // Return coins info from coingecko
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=bitcoin%2Cethereum&category=layer-1&x_cg_demo_api_key="

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	fmt.Println(len(body))

	if err != nil {
		return err
	}

	var coins []Coin
	err = json.Unmarshal(body, &coins)
	fmt.Println(len(coins))

	if err != nil {
		return err
	}

	var result []Coin

	for _, value := range coins {
		if len(result) == len(coinNames) {
			break
		}

		if slices.Contains(coinNames, value.ID) {
			result = append(result, value)
		}
	}

	CurrentCoinsValues = result

	LastUpdateTime = time.Now()
	return nil
}
