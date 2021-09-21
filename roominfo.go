package main

type RoomInfo struct {
	RoomID    uint `gorm:"primarykey"`
	CreatedAt string
	RoomName  string
}
