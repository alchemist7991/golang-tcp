package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
	msgCh      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		msgCh:      make(chan Message),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	s.acceptLoop()
	<-s.quitCh
	close(s.msgCh)
	return nil
}

// To accept connections made by a client
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Error in accepting connection")
			continue
		}
		fmt.Println("New Connection from :", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

// read data from the connection that was made
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error in reading buffer message")
			continue
		}
		msg := Message{
			from: conn.RemoteAddr().String(),
			payload: buff[:n],
		}
		s.msgCh <- msg
		conn.Write([]byte("Received message\n"))
	}
}

func (s *Server) readMsg() {
	for {
		msg := <-s.msgCh
		fmt.Printf("New Message from (%s): %s", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":3000")
	go server.readMsg()
	log.Fatal(server.Start())
}
