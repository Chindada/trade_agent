// Package sinopacapi package sinopacapi
package sinopacapi

// OrderBody OrderBody
type OrderBody struct {
	Stock    string  `json:"stock,omitempty" yaml:"stock"`
	Price    float64 `json:"price,omitempty" yaml:"price"`
	Quantity int64   `json:"quantity,omitempty" yaml:"quantity"`
}

// OrderCancelBody OrderCancelBody
type OrderCancelBody struct {
	OrderID string `json:"order_id,omitempty" yaml:"order_id"`
}

// FetchLastCloseBody FetchLastCloseBody
type FetchLastCloseBody struct {
	StockNumArr []string `json:"stock_num_arr,omitempty" yaml:"stock_num_arr"`
	DateArr     []string `json:"date_arr,omitempty" yaml:"date_arr"`
}

// FetchLastCountBody FetchLastCountBody
type FetchLastCountBody struct {
	StockNumArr []string `json:"stock_num_arr,omitempty" yaml:"stock_num_arr"`
}

// FetchKbarBody FetchKbarBody
type FetchKbarBody struct {
	StockNum  string `json:"stock_num,omitempty" yaml:"stock_num"`
	StartDate string `json:"start_date,omitempty" yaml:"start_date"`
	EndDate   string `json:"end_date,omitempty" yaml:"end_date"`
}

// FetchBody FetchBody
type FetchBody struct {
	StockNum string `json:"stock_num,omitempty" yaml:"stock_num"`
	Date     string `json:"date,omitempty" yaml:"date"`
}

// SubscribeBody SubscribeBody
type SubscribeBody struct {
	StockNumArr []string `json:"stock_num_arr,omitempty" yaml:"stock_num_arr"`
}
