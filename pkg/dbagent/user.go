package dbagent

import (
	"gorm.io/gorm"
)

// User User
type User struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	UserName string `json:"user_name" yaml:"user_name" gorm:"column:user_name"`

	FirstName string `json:"first_name" yaml:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" yaml:"last_name" gorm:"column:last_name"`

	Login Login `json:"login" yaml:"login" gorm:"-"`
}

// Login Login
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// TableName TableName
func (User) TableName() string {
	return "basic_user"
}

// InsertUser InsertUser
func (c *DBAgent) InsertUser(record *User) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllUser DeleteAllUser
func (c *DBAgent) DeleteAllUser() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&User{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetAllUser GetAllUser
func (c *DBAgent) GetAllUser() ([]User, error) {
	var tmp []User
	err := c.DB.Model(&User{}).Find(&tmp).Error
	return tmp, err
}

// GetUserByAccount GetUserByAccount
func (c *DBAgent) GetUserByAccount(account string) (User, error) {
	var tmp User
	err := c.DB.Model(&User{}).Where("account = ?", account).Find(&tmp).Error
	return tmp, err
}
