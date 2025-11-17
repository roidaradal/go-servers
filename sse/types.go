package sse

import "fmt"

type ClientChan chan string

type StreamClient struct {
	Topic
	Channel ClientChan
}

type Topic struct {
	Group string
	Focus string
}

type Message struct {
	Topic
	Content string
}

type MockData struct {
	Title   string
	Message string
}

func (t Topic) String() string {
	return fmt.Sprintf("<%s/%s>", t.Group, t.Focus)
}
