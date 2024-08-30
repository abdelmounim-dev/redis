package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/abdelmounim-dev/redis/pkg/handlers"
	"github.com/abdelmounim-dev/redis/pkg/parser"
)

type Message struct {
	content    string
	senderAddr net.Addr
}

type Server struct {
	connections map[net.Addr]*net.Conn
	maxConns    int
	listener    net.Listener
	msgChan     chan Message
	errChan     chan error
	connChan    chan net.Conn
	mu          sync.Mutex
}

func NewServer(port int, maxConns int) (*Server, error) {
	p := strconv.Itoa(port)
	listener, err := net.Listen("tcp4", ":"+p)
	if err != nil {
		return nil, err
	}

	return &Server{
		connections: make(map[net.Addr]*net.Conn, maxConns),
		maxConns:    maxConns,
		listener:    listener,
		msgChan:     make(chan Message),
		errChan:     make(chan error),
		connChan:    make(chan net.Conn),
		mu:          sync.Mutex{},
	}, nil
}

func (s *Server) Run() error {
	go s.handleErrors()
	go s.establishConnections()

	for len(s.connections) <= s.maxConns {
		conn := <-s.connChan
		s.addConnection(conn)
		go s.handleConn(conn)
	}

	return nil
}

func (s *Server) handleErrors() {
	for {
		err := <-s.errChan
		log.Print("ERROR: ", err)
	}
}

func (s *Server) establishConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.errChan <- fmt.Errorf("CONN: %v", err)
			continue
		}
		log.Println("Connection Established: ", conn.RemoteAddr())

		// TODO: handle reaching maximum connections

		s.connChan <- conn
	}
}

func (s *Server) addConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[conn.RemoteAddr()] = &conn
}

func (s *Server) removeConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn.RemoteAddr() != nil {
		conn.Close()
	}
	delete(s.connections, conn.RemoteAddr())
}

func (s *Server) handleConn(c net.Conn) {
	defer func() {
		c.Close()
		s.removeConnection(c)
		log.Print("Connection Closed")
	}()

	reader := bufio.NewReader(c)
	for {
		// TODO: handle buffer

		p := parser.NewParser(reader)
		t, err := p.NextToken()
		if err != nil {
			s.errChan <- fmt.Errorf("READ: %v", err)
		}
		log.Println(*t)
		rt, err := handlers.HandleCommand(t)
		if err != nil {
			s.errChan <- fmt.Errorf("READ: %v", err)
		}

		res, err := rt.Serialize()
		if err != nil {
			s.errChan <- fmt.Errorf("READ: %v", err)
		}

		_, err = c.Write(res)
		if err != nil {
			s.errChan <- fmt.Errorf("WRITE: %v", err)
		}

	}
}
