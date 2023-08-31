package models

// Stock is the struct equivalent to the body returned by the gRPC API
// swagger:model
type Stock struct {
	Symbol        string
	Open          float64
	High          float64
	Low           float64
	Price         float64
	Volume        float64
	LTD           string
	PrevClose     float64
	Change        float64
	PercentChange string
	Destination   string
}

// MoreStock contains important financial metrics
// swagger:model
type MoreStock struct {
	Ticker            string
	Name              string
	Exchange          string
	Sector            string
	MarketCap         float64
	PERatio           float64
	PEGRatio          float64
	DivPerShare       float64
	EPS               float64
	RevPerShare       float64
	ProfitMargin      float64
	YearHigh          float64
	YearLow           float64
	SharesOutstanding float64
	PriceToBookRatio  float64
	Beta              float64
}
