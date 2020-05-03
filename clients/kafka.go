package clients

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
)

var topicsToListen = [...]string{
	"trade-created",
}

var listeners = map[string]*kafka.Reader{}

type TradeCreatedMessage struct {
	ID       uuid.UUID `json:"id"`
	UserId   uint      `json:"user_id"`
	WalletId int       `json:"wallet_id"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Amount   uint      `json:"amount"`
}

func InitKafka() {
	for _, topic := range topicsToListen {
		listeners[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			GroupID: "wallet",
			Topic:   topic,
		})
		go listenForTopic(topic)
	}
}

// TODO: Not working
func listenForTopic(topic string) {
	for {
		fmt.Println("Listening for topic", topic)
		m, err := listeners[topic].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
