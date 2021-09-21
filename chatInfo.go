package main

type ChatInfo struct {
	ChatID       uint `gorm:"primarykey"`
	CreatedAt    string
	Username     string
	RoomName     string
	Conversation string
}
