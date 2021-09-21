# TCP-Chat

> This project is based on https://github.com/plutov/packagemain/tree/master/20-tcp-chat

## functionality

```
1. provide every chatinfo with send time
2. store registered rooms
3. reload room's history chatinfos when user joins registered rooms
```
## build command

```
1. go build .
2. ./chat
```

## usage command

```
1. /nick <name>, call youself as <name>.
2. /join <room>, join <room> if <room> exists, otherwise create <room> and join.
3. /rooms, list all existing rooms.
4. /msg <meg>, communicate <msg> to other users in the current room.
5. /quit, quit chat connection.
```
