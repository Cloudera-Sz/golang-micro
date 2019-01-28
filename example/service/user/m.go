package user

import "time"

type User struct {
	Id         int64     `remark:"Id" gorm:"column:Id;type:bigint(20);primary_key"`
	UserName   string    `remark:"UserName" gorm:"column:UserName;type:varchar(50)"`
	CreateTime time.Time `remark:"CreateTime" gorm:"column:CreateTime;type:datetime(6)"`
	UpdateTime time.Time `remark:"UpdateTime" gorm:"column:UpdateTime;type:datetime(6)"`
}
