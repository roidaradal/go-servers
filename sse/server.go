package sse

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/str"
)

func RunServer(host string, port int) {
	corsCfg := cors.DefaultConfig()
	corsCfg.MaxAge = 12 * time.Hour
	corsCfg.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"Accept",
		"User-Agent",
		"Cache-Control",
	}
	corsCfg.ExposeHeaders = []string{
		"Content-Length",
	}
	corsCfg.AllowMethods = []string{
		"GET",
		"POST",
	}
	corsCfg.AllowAllOrigins = true

	server := gin.Default()
	server.Use(cors.New(corsCfg))
	address := fmt.Sprintf("%s:%d", host, port)

	stream := startStreamServer()
	s := server.Group("/stream")
	s.Use(stream.HTTPMiddleware())
	s.GET("/items", handleStream)
	go stream.Run()

	server.Run(address)
}

type StreamServer struct {
	MessageChan       chan Message
	NewClientsChan    chan StreamClient
	ClosedClientsChan chan ClientChan
	AllClients        map[ClientChan]Topic
	TopicClients      map[Topic]map[ClientChan]bool
}

func startStreamServer() *StreamServer {
	stream := &StreamServer{
		MessageChan:       make(chan Message),
		NewClientsChan:    make(chan StreamClient),
		ClosedClientsChan: make(chan ClientChan),
		AllClients:        make(map[ClientChan]Topic),
		TopicClients:      make(map[Topic]map[ClientChan]bool),
	}
	go stream.Listen()
	return stream
}

func (s *StreamServer) Listen() {
	for {
		select {
		case newClient := <-s.NewClientsChan:
			clientChan := newClient.Channel
			topic := newClient.Topic
			s.AllClients[clientChan] = topic
			if dict.NoKey(s.TopicClients, topic) {
				s.TopicClients[topic] = make(map[ClientChan]bool)
			}
			s.TopicClients[topic][clientChan] = true
			fmt.Printf("Add client: %s, Total: %d\n", topic, len(s.AllClients))
		case channel := <-s.ClosedClientsChan:
			topic, ok := s.AllClients[channel]
			if ok && dict.HasKey(s.TopicClients, topic) {
				delete(s.TopicClients[topic], channel)
				if len(s.TopicClients[topic]) == 0 {
					delete(s.TopicClients, topic)
				}
			}
			delete(s.AllClients, channel)
			close(channel)
			fmt.Printf("Closed client, Total: %d\n", len(s.AllClients))
		case message := <-s.MessageChan:
			for clientChan := range s.TopicClients[message.Topic] {
				clientChan <- message.Content
			}
		}
	}
}

func (s *StreamServer) Run() {
	interval := 1 * time.Second
	for {
		start := clock.TimeNow()
		if len(s.AllClients) == 0 {
			clock.Sleep(interval, start)
			continue
		}
		for topic := range s.TopicClients {
			if len(s.TopicClients[topic]) == 0 {
				continue
			}
			data := mockData(topic)
			message, err := str.JSON(data)
			if err != nil {
				continue
			}
			s.MessageChan <- Message{
				Topic:   topic,
				Content: message,
			}
		}
		clock.Sleep(interval, start)
	}
}

func (s *StreamServer) HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		clientChan := make(ClientChan)
		s.NewClientsChan <- StreamClient{
			Channel: clientChan,
			Topic:   newTopic(c),
		}
		defer func() {
			s.ClosedClientsChan <- clientChan
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
}

func newTopic(c *gin.Context) Topic {
	return Topic{
		Group: c.Query("group"),
		Focus: c.Query("focus"),
	}
}

func handleStream(c *gin.Context) {
	v, ok := c.Get("clientChan")
	if !ok {
		return
	}
	clientChan, ok := v.(ClientChan)
	if !ok {
		return
	}
	c.Stream(func(w io.Writer) bool {
		if message, ok := <-clientChan; ok {
			c.SSEvent("message", message)
			return true
		}
		return false
	})
}

func mockData(topic Topic) MockData {
	return MockData{
		Title:   fmt.Sprintf("Mock data for %s", topic),
		Message: fmt.Sprintf("Created on %s", clock.DateTimeNow()),
	}
}
