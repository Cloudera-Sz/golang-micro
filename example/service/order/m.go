package order

import "time"

type Order struct {
	Id         int64     `remark:"Id" gorm:"column:Id;type:bigint(20);primary_key"`
	UserId     int64     `remark:"UserId" gorm:"column:UserId;type:bigint(50)"`
	CreateTime time.Time `remark:"CreateTime" gorm:"column:CreateTime;type:datetime(6)"`
	UpdateTime time.Time `remark:"UpdateTime" gorm:"column:UpdateTime;type:datetime(6)"`
}
