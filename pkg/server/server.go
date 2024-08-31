package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abdelmounim-dev/redis/pkg/handlers"
	"github.com/abdelmounim-dev/redis/pkg/parser"
	"github.com/abdelmounim-dev/redis/pkg/storage"
)

type Message struct {
	content    string
	senderAddr net.Addr
}

type Server struct {
	mu          sync.Mutex
	connections map[net.Addr]*net.Conn
	connNum     atomic.Int32
	maxConns    int32
	timeout     time.Duration
	listener    net.Listener
	msgChan     chan Message
	errChan     chan error
	log         *log.Logger
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc

	// attributes specific to redis
	store    storage.Store
	handlers *handlers.Handlers
}

func NewServer(address string, maxConns int32, timeout time.Duration) (*Server, error) {
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	store := storage.NewKeyValueStore()
	handlers := handlers.NewHandlers(store)
	return &Server{
		connections: make(map[net.Addr]*net.Conn, maxConns),
		maxConns:    maxConns,
		listener:    listener,
		msgChan:     make(chan Message),
		errChan:     make(chan error),
		mu:          sync.Mutex{},
		log:         log.New(os.Stderr, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile),
		ctx:         ctx,
		cancel:      cancel,
		timeout:     timeout,

		store:    store,
		handlers: handlers,
	}, nil
}

func (s *Server) Run() {
	s.wg.Add(2)
	go s.handleErrors()
	go s.establishConnections()

	for {
	}
}

func (s *Server) Kill() error {
	s.log.Println("Initiating graceful shutdown...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	// Signal all goroutines to stop
	s.cancel()

	// Stop accepting new connections
	s.listener.Close()

	// Signal all connection handlers to stop
	s.mu.Lock()
	for _, conn := range s.connections {
		s.msgChan <- Message{content: "shutdown", senderAddr: (*conn).RemoteAddr()}
	}
	s.mu.Unlock()

	// Wait for all goroutines to finish or timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.log.Println("Graceful shutdown completed")
	case <-ctx.Done():
		s.log.Println("Graceful shutdown timed out, forcing exit")
		// Force close any remaining connections
		s.mu.Lock()
		for _, conn := range s.connections {
			(*conn).Close()
		}
		s.mu.Unlock()
	}

	// Drain message channel
	close(s.msgChan)
	for range s.msgChan {
		// Drain remaining messages
	}

	// Close error channel
	close(s.errChan)

	return nil
}

func (s *Server) handleErrors() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.errChan:
			s.log.Print("ERROR: ", err)
		}
	}
}

func (s *Server) canAcceptConnection() bool {
	return s.connNum.Load() < s.maxConns
}

func (s *Server) establishConnections() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				s.errChan <- fmt.Errorf("CONN: %v", err)
				continue
			}

			if !s.canAcceptConnection() {
				s.log.Printf("Max connections reached. Rejecting connection from %s", conn.RemoteAddr())
				conn.Close()
				continue
			}
			s.log.Println("Connection Established: ", conn.RemoteAddr())

			s.addConnection(conn)
			s.wg.Add(1)
			go s.handleConn(conn)
		}
	}
}

func (s *Server) addConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[conn.RemoteAddr()] = &conn
	s.connNum.Add(1)
}

func (s *Server) removeConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn.RemoteAddr() != nil {
		conn.Close()
	}
	delete(s.connections, conn.RemoteAddr())
	s.connNum.Add(-1)
}

func (s *Server) handleConn(c net.Conn) {
	defer func() {
		c.Close()
		s.removeConnection(c)
		s.log.Print("Connection Closed")
		s.wg.Done()
	}()

	reader := bufio.NewReader(c)
	for {
		select {
		case <-s.ctx.Done():
			s.log.Printf("Closing connection to %s\n", c.RemoteAddr())
			return
		default:
			err, res := s.handleMessage(reader)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// This is a timeout, just continue the loop
					continue
				}
				// For other errors, send to error channel and return
				s.errChan <- fmt.Errorf("HANDLE MESSAGE: %v", err)
				return
			}

			// Use a timeout for writing to the connection
			err = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				s.errChan <- fmt.Errorf("SET DEADLINE: %v", err)
				return
			}

			_, err = c.Write(res)
			if err != nil {
				s.errChan <- fmt.Errorf("WRITE: %v", err)
				return
			}
		}
	}
}

func (s *Server) handleMessage(reader *bufio.Reader) (error, []byte) {
	p := parser.NewParser(reader)
	t, err := p.NextToken()
	if err != nil {
		return fmt.Errorf("READ: %v", err), nil
	}
	s.log.Println(*t)
	rt, err := s.handlers.HandleCommand(t)
	if err != nil {
		return fmt.Errorf("HANDLE COMMAND: %v", err), nil
	}

	res, err := rt.Serialize()
	if err != nil {
		return fmt.Errorf("SERIALIZE: %v", err), nil
	}
	return nil, res
}
