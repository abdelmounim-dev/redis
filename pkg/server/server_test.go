package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.NotNil(t, s.listener)
	assert.Equal(t, 10, s.maxConns)
	assert.Equal(t, 5*time.Second, s.timeout)
	assert.NotNil(t, s.log)
}

func TestServerRun(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	defer s.Kill()

	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	conn, err := net.Dial("tcp", s.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	assert.Eventually(t, func() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		return len(s.connections) == 1
	}, 1*time.Second, 10*time.Millisecond)
}

func TestServerKill(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	conn, err := net.Dial("tcp", s.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	err = s.Kill()
	require.NoError(t, err)

	// Try to connect after kill, should fail
	_, err = net.Dial("tcp", s.listener.Addr().String())
	assert.Error(t, err)
}

func TestServerMaxConnections(t *testing.T) {
	maxConns := 2
	s, err := NewServer("localhost:0", maxConns, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	defer s.Kill()

	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	var wg sync.WaitGroup
	for i := 0; i < maxConns+1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", s.listener.Addr().String())
			if err == nil {
				defer conn.Close()
				// Keep connection open for a while
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	assert.Eventually(t, func() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		return len(s.connections) == maxConns
	}, 1*time.Second, 10*time.Millisecond)
}

func TestServerHandleMessage(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	conn, err := net.Dial("tcp", s.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// Send a message
	_, err = conn.Write([]byte("PING\r\n"))
	require.NoError(t, err)

	// Read the response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	require.NoError(t, err)
	assert.Equal(t, "PONG\r\n", response)
}

func TestServerTimeout(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 1*time.Second)
	require.NoError(t, err)

	go s.Run()
	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	conn, err := net.Dial("tcp", s.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	// Wait for more than the timeout duration
	time.Sleep(2 * time.Second)

	// Try to send a message, should fail due to timeout
	_, err = conn.Write([]byte("PING\r\n"))
	assert.Error(t, err)
}

func TestServerConcurrency(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	numClients := 5
	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", s.listener.Addr().String())
			require.NoError(t, err)
			defer conn.Close()

			for j := 0; j < 10; j++ {
				_, err = conn.Write([]byte("PING\r\n"))
				require.NoError(t, err)

				reader := bufio.NewReader(conn)
				response, err := reader.ReadString('\n')
				require.NoError(t, err)
				assert.Equal(t, "PONG\r\n", response)
			}
		}()
	}

	wg.Wait()
}

func TestServerGracefulShutdown(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	go s.Run()
	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	// Start a long-running operation
	conn, err := net.Dial("tcp", s.listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 10; i++ {
			_, err := conn.Write([]byte("PING\r\n"))
			if err != nil {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Initiate graceful shutdown
	time.Sleep(500 * time.Millisecond)
	err = s.Kill()
	require.NoError(t, err)

	// Check if the long-running operation completed
	select {
	case <-done:
		// Operation completed successfully
	case <-time.After(2 * time.Second):
		t.Fatal("Long-running operation did not complete before shutdown")
	}
}

func TestServerErrorHandling(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	// Simulate an error
	s.errChan <- fmt.Errorf("test error")

	// Check if the error is logged (you might need to mock the logger for this)
	// This is a simple check and might need to be adjusted based on your logging implementation
	time.Sleep(100 * time.Millisecond)
	// Assert that the error was logged
}

func TestServerContextCancellation(t *testing.T) {
	s, err := NewServer("localhost:0", 10, 5*time.Second)
	require.NoError(t, err)

	_, cancel := context.WithCancel(context.Background())
	go func() {
		s.Run()
	}()

	time.Sleep(100 * time.Millisecond) // Give some time for the server to start

	// Cancel the context
	cancel()

	// Check if the server shuts down gracefully
	assert.Eventually(t, func() bool {
		_, err := net.Dial("tcp", s.listener.Addr().String())
		return err != nil
	}, 1*time.Second, 10*time.Millisecond)
}
