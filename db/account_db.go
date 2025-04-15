package db

type AccountData struct {
	Account string `gorm:"primarykey; comment:account name"`
	Uid     int64  `gorm:"unique; comment:user uuid"`
}
