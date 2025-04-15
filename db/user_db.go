package db

type UserData struct {
	Uid     int64  `gorm:"primarykey; comment:uuid"`
	Account string `gorm:"unique; comment:account name"`
}
