// Package utils package utils
package utils

import "testing"

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

func TestGetNewClose(t *testing.T) {
	type args struct {
		close float64
		unit  int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no close",
			args: args{
				close: 0,
				unit:  1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		close := tt.args.close
		unit := tt.args.unit
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNewClose(close, unit); got != want {
				t.Errorf("GetNewClose() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetMaxByOpen(t *testing.T) {
	type args struct {
		open float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no open",
			args: args{
				open: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		open := tt.args.open
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxByOpen(open); got != want {
				t.Errorf("GetMaxByOpen() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetMinByOpen(t *testing.T) {
	type args struct {
		open float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no open",
			args: args{
				open: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		open := tt.args.open
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMinByOpen(open); got != want {
				t.Errorf("GetMinByOpen() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetDiff(t *testing.T) {
	type args struct {
		close float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no close",
			args: args{
				close: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		close := tt.args.close
		want := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDiff(close); got != want {
				t.Errorf("GetDiff() = %v, want %v", got, want)
			}
		})
	}
}
