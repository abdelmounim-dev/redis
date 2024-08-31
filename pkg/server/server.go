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
	mu           sync.Mutex
	connections  map[net.Addr]*net.Conn
	maxConns     int
	listener     net.Listener
	msgChan      chan Message
	errChan      chan error
	shutdownChan chan struct{}
	log          *log.Logger
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
		mu:          sync.Mutex{},
	}, nil
}

func (s *Server) Run() {
	go s.handleErrors()
	go s.establishConnections()
}

func (s *Server) Kill() error {
	close(s.shutdownChan)

	return nil
}

func (s *Server) handleErrors() {
	for {
		select {
		case <-s.shutdownChan:
			return
		case err := <-s.errChan:
			s.log.Print("ERROR: ", err)
		}
	}
}

func (s *Server) establishConnections() {
	for {
		select {
		case <-s.shutdownChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				s.errChan <- fmt.Errorf("CONN: %v", err)
				continue
			}
			s.log.Println("Connection Established: ", conn.RemoteAddr())

			// TODO: handle reaching maximum connections

			s.addConnection(conn)
			go s.handleConn(conn)
		}
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
		s.log.Print("Connection Closed")
	}()

	reader := bufio.NewReader(c)
	for {
		select {
		case <-s.shutdownChan:
			return
		default:
			// TODO: handle buffer

			err, res := s.handleMessage(reader)

			_, err = c.Write(res)
			if err != nil {
				s.errChan <- fmt.Errorf("WRITE: %v", err)
			}
		}

	}
}

func (s *Server) handleMessage(reader *bufio.Reader) (error, []byte) {
	p := parser.NewParser(reader)
	t, err := p.NextToken()
	if err != nil {
		s.errChan <- fmt.Errorf("READ: %v", err)
	}
	s.log.Println(*t)
	rt, err := handlers.HandleCommand(t)
	if err != nil {
		s.errChan <- fmt.Errorf("READ: %v", err)
	}

	res, err := rt.Serialize()
	if err != nil {
		s.errChan <- fmt.Errorf("READ: %v", err)
	}
	return err, res
}
