package coupon

import "time"

type Coupon struct {
	Id         int64     `remark:"Id" gorm:"column:Id;type:bigint(20);primary_key"`
	Price      float64   `remark:"Price" gorm:"column:Price;type:decimal(18,4)"`
	CreateTime time.Time `remark:"CreateTime" gorm:"column:CreateTime;type:datetime(6)"`
	UpdateTime time.Time `remark:"UpdateTime" gorm:"column:UpdateTime;type:datetime(6)"`
}
