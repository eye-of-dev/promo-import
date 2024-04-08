package models

type PromoCodes struct {
	PromoCode string `gorm:"column:promocode"`
}

func (PromoCodes) TableName() string {
	return "promocodes"
}
