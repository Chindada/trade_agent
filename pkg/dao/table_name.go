// Package dao package dao
package dao

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Stock) TableName() string {
	return "basic_stock"
}
