package data

type GlobalQuote struct {
	StockData Stock `json:"Global Quote"`
}

// Stock is the struct equivalent of the json body returned by
// Alpha Vantage's Quote Endpoint API method
type Stock struct {
	Symbol        string `json:"01. symbol"`
	Open          string `json:"02. open"`
	High          string `json:"03. high"`
	Low           string `json:"04. low"`
	Price         string `json:"05. price"`
	Volume        string `json:"06. volume"`
	LTD           string `json:"07. latest trading day"`
	PrevClose     string `json:"08. previous close"`
	Change        string `json:"09. change"`
	PercentChange string `json:"10. change percent"`
	Destination   string `json:"Destination"`
}

// MoreStock contains important financial metrics returned by AV's
// Company Overview endpoint
type MoreStock struct {
	Ticker            string `json:"Symbol"`
	Name              string `json:"Name"`
	Exchange          string `json:"Exchange"`
	Sector            string `json:"Sector"`
	MarketCap         string `json:"MarketCapitalization"`
	PERatio           string `json:"PERatio"`
	PEGRatio          string `json:"PEGRatio"`
	DivPerShare       string `json:"DividendPerShare"`
	EPS               string `json:"EPS"`
	RevPerShare       string `json:"RevenuePerShareTTM"`
	ProfitMargin      string `json:"ProfitMargin"`
	YearHigh          string `json:"52WeekHigh"`
	YearLow           string `json:"52WeekLow"`
	SharesOutstanding string `json:"SharesOutstanding"`
	PriceToBookRatio  string `json:"PriceToBookRatio"`
	Beta              string `json:"Beta"`
}
