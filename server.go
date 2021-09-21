package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	roomMembersPair := make(map[string]*room)
	allRooms, _ := findAllRooms()

	for _, r := range allRooms {
		roomMembersPair[r.RoomName] = &room{
			name:    r.RoomName,
			members: make(map[net.Addr]*client),
		}
	}
	return &server{
		rooms:    roomMembersPair,
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msgToCurrentRoom(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msgToClient("NAME is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msgToClient(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msgToClient("ROOM is required. usage: /join Room")
		return
	}
	roomName := args[1]

	r := s.rooms[roomName]

	if r == nil {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
		insertRoomInfo(&RoomInfo{
			RoomName:  roomName,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msgToClient(fmt.Sprintf("welcome to %s", r.name))

	chatContents, _ := loadRoomInfos(roomName)
	for _, chat := range chatContents {
		t := chat.CreatedAt
		c.msgToClientForReloadMsg(t + " > " + chat.Username + ": " + chat.Conversation)
	}
}

func insertRoomInfo(roomInfo *RoomInfo) error {
	if result := db.Create(roomInfo); result.Error != nil {
		return result.Error
	}
	return nil
}

func loadRoomInfos(roomName string) ([]ChatInfo, error) {
	var allInfos []ChatInfo
	if result := db.Where("room_name = ?", roomName).Find(&allInfos); result.Error != nil {
		return nil, result.Error
	}
	return allInfos, nil
}

func (s *server) listRooms(c *client, args []string) {
	var rooms, _ = findAllRooms()
	var allRooms []string
	for _, room := range rooms {
		allRooms = append(allRooms, room.RoomName)
	}

	c.msgToClient(fmt.Sprintf("available rooms are: %s", strings.Join(allRooms, ", ")))
}

func findAllRooms() ([]RoomInfo, error) {
	var allRooms []RoomInfo
	if result := db.Find(&allRooms); result.Error != nil {
		return nil, result.Error
	}
	return allRooms, nil
}

func (s *server) msgToCurrentRoom(c *client, args []string) {
	if c.room == nil {
		c.err(fmt.Errorf("you must join the room first"))
		return
	}

	if len(args) < 2 {
		c.msgToClient("MSG is required, usage: /msg MSG")
		return
	}

	chatInfo := &ChatInfo{
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		Username:     c.nick,
		RoomName:     c.room.name,
		Conversation: strings.Join(args[1:], " "),
	}
	insertChatInfo(chatInfo)

	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
}

func insertChatInfo(chatInfo *ChatInfo) error {
	if result := db.Create(chatInfo); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr())

	s.quitCurrentRoom(c)

	c.msgToClient("bye bye~")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
