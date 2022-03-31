package data

// BasicStock is the struct equivalent of the json body returned by
// Alpha Vantage's Quote Endpoint API method
type BasicStock struct {
	Symbol        string `json:"symbol"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Price         string `json:"price"`
	Volume        string `json:"volume"`
	PrevClose     string `json:"previous close"`
	Change        string `json:"change"`
	PercentChange string `json:"change percent"`
}
