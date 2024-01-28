package model

type User struct {
	Id       string `gorm:"PRIMARY_KEY"` //`gorm:"column:id;PRIMARY_KEY;AUTO_INCREMENT"`
	NickName string //`gorm:"unique_index:uk_NickName"`
	HeadImg  string
}
