package dap

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/go-delve/delve/service/debugger"
)

func ServeDAP(addr string) {
	s := newServer(addr)
	s.serveTCP()
}

type server struct {
	addr     string
	debugger *debugger.Debugger
}

func newServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) serveTCP() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening for DAP connections on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			log.Println(c)
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
