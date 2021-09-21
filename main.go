package main

import (
	"log"
	"net"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	_db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db = _db
	if err := db.AutoMigrate(&ChatInfo{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&RoomInfo{}); err != nil {
		panic(err)
	}

	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("started server on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}

/* todo
1. add timeStamp to all msg -- done
2. store chat infos to database(sqlite) -- done
3. store rooms to database(sqlite) -- done
4. load chat infos from db -- done
5. initialize rooms data from db -- done
6. what if join the current room again? -- just rejoin
5. update README
*/
