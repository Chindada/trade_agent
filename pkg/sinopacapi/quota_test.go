// Package sinopacapi package sinopacapi
package sinopacapi

import (
	"testing"
)

func TestGetStockBuyCost(t *testing.T) {
	type args struct {
		price float64
		qty   int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "no price",
			args: args{
				price: 0,
				qty:   1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		price := tt.args.price
		qty := tt.args.qty
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStockBuyCost(price, qty); got != want {
				t.Errorf("GetStockBuyCost() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetStockSellCost(t *testing.T) {
	type args struct {
		price float64
		qty   int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "no price",
			args: args{
				price: 0,
				qty:   1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		price := tt.args.price
		qty := tt.args.qty
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStockSellCost(price, qty); got != want {
				t.Errorf("GetStockSellCost() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetStockTradeFeeDiscount(t *testing.T) {
	type args struct {
		price float64
		qty   int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "no price",
			args: args{
				price: 0,
				qty:   1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		price := tt.args.price
		qty := tt.args.qty
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStockTradeFeeDiscount(price, qty); got != want {
				t.Errorf("GetStockSellCost() = %v, want %v", got, want)
			}
		})
	}
}
