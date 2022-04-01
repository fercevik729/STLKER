package data

// BasicStock is the struct equivalent of the json body returned by
// Alpha Vantage's Quote Endpoint API method
type Stock struct {
	Symbol        string `json:"01. symbol"`
	Open          string `json:"02. open"`
	High          string `json:"03. high"`
	Low           string `json:"04. low"`
	Price         string `json:"05. price"`
	Volume        string `json:"06. volume"`
	PrevClose     string `json:"08. previous close"`
	Change        string `json:"09. change"`
	PercentChange string `json:"10. change percent"`
}
