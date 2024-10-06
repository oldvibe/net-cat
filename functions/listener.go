package ncat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	messages []string
	clients  map[string]*Client
	mutex    sync.Mutex
	Content  []byte
}

type Client struct {
	conn net.Conn
	name string
}





func CreateNewServer() *Server {
	return &Server{
		messages: make([]string, 0),
		clients:  make(map[string]*Client),
	}
}

func (s *Server) Listen(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.StartConnection(connection)
	}
}

func (s *Server) RemoveConnection(client *Client, msg string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, c := range s.clients {
		if c != client {
			c.conn.Write([]byte(msg))
			s.WaitingForInput(c)
		}
	}
	delete(s.clients, client.name)
}

func (s *Server) StartConnection(cc net.Conn) {
	defer cc.Close()
	pngn := strings.TrimRight(string(s.Content), "\n") + "\n"
	cc.Write([]byte("Welcome to TCP-Chat!\n"))
	cc.Write([]byte(  "\033[1;94m" + pngn +  "\033[0m"))
	cc.Write([]byte("[ENTER YOUR NAME]: "))

	reader := bufio.NewReader(cc)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}
	name = strings.TrimSpace(name)

	client := &Client{conn: cc, name: name}

	s.mutex.Lock()
	if len(s.clients) >= 10 {
		cc.Write([]byte("Chat is already full!\n"))
		s.mutex.Unlock()
		return
	}
	s.clients[name] = client
	s.mutex.Unlock()

	joinMessage := fmt.Sprintf("\n%s has joined our chat...\n", name)
	s.BroadcastMessage(joinMessage, client)

	// Send previous messages to the new client
	for _, msg := range s.messages {
		cc.Write([]byte(msg))
	}

	s.WaitingForInput(client)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		if s.IsValidMsg(msg) {
			formattedMsg := fmt.Sprintf("\n[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, msg)
			s.messages = append(s.messages, formattedMsg)
			s.BroadcastMessage(formattedMsg, client)
			s.WaitingForInput(client)
		} else {
			s.WaitingForInput(client)
		}

	}
	leaveMessage := fmt.Sprintf("\n%s has left our chat...\n", name)
	s.RemoveConnection(client, leaveMessage)
}

func (s *Server) BroadcastMessage(msg string, sender *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, client := range s.clients {
		if client != sender {
			client.conn.Write([]byte(msg))
			s.WaitingForInput(client)
		}
	}
}

func (s *Server) WaitingForInput(client *Client) {
	prompt := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.name)
	client.conn.Write([]byte(prompt))
}

func (s *Server) IsValidMsg(message string) bool {
	message = strings.TrimSpace(message)
	return len(message) > 0 && s.IsPrintable(message)
}

func (s *Server) IsPrintable(message string) bool {
	for _, c := range message {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}
