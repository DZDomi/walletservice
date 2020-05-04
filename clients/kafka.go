package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DZDomi/walletservice/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var topicsToListen = [...]string{
	"trade-created",
}

var topicsToWrite = [...]string{
	"wallets-updated",
}

var listeners = map[string]*kafka.Reader{}
var writers = map[string]*kafka.Writer{}

type TradeCreatedMessage struct {
	ID         uint   `json:"id"`
	User       uint   `json:"user_id"`
	FromWallet uint   `json:"from_wallet_id"`
	ToWallet   uint   `json:"to_wallet_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     uint   `json:"amount"`
}

type Wallet struct {
	User   uint   `json:"user_id"`
	Wallet uint   `json:"wallet_id"`
	Amount uint   `json:"amount"`
	Action string `json:"action"`
}

type TriggerEvent struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

type WalletsUpdatedEvent struct {
	TriggeredBy *TriggerEvent `json:"triggered_by"`
	Wallets     []Wallet      `json:"wallets"`
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
	for _, topic := range topicsToWrite {
		writers[topic] = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
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
		switch topic {
		case "trade-created":
			go deductWalletBalances(m)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func sendToTopic(topic string, message string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	// TODO: Think about the error handling here
	go writers[topic].WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(id.String()),
		Value: []byte(message),
	})
	//if err != nil {
	//	fmt.Println("Error while sending wallet updated event:", err.Error())
	//}
	return nil
}

func deductWalletBalances(message kafka.Message) {

	id := string(message.Key)
	fmt.Println("Got message:", id)

	trade := &TradeCreatedMessage{}
	if err := json.Unmarshal(message.Value, trade); err != nil {
		fmt.Println("Unable to decode message:", message.Value)
		return
	}

	fromLock, err := GetLock(string(trade.FromWallet), time.Minute)
	if err != nil {
		fmt.Println("Unable to obtain lock for:", trade.FromWallet)
	}
	defer ReleaseLock(fromLock)
	toLock, err := GetLock(string(trade.ToWallet), time.Minute)
	if err != nil {
		fmt.Println("Unable to obtain lock for:", trade.ToWallet)
	}
	defer ReleaseLock(toLock)

	tx := models.DB.Begin()

	fromWallet := &models.Wallet{}
	toWallet := &models.Wallet{}
	models.DB.Where("id = ?", trade.FromWallet).First(&fromWallet)
	models.DB.Where("id = ?", trade.ToWallet).First(&toWallet)

	if fromWallet.ID == 0 || toWallet.ID == 0 {
		fmt.Println("Unable to find wallets:", trade.FromWallet, ":", trade.ToWallet)
		return
	}

	// TODO: Make this correct
	fromWallet.Balance -= trade.Amount
	toWallet.Balance += trade.Amount

	models.DB.Save(fromWallet)
	models.DB.Save(toWallet)

	tx.Commit()

	wallets := []Wallet{
		// TODO: Fix amount
		{
			User:   fromWallet.User,
			Wallet: fromWallet.ID,
			Amount: trade.Amount,
			Action: "subtract",
		},
		{
			User:   toWallet.User,
			Wallet: toWallet.ID,
			Amount: trade.Amount,
			Action: "add",
		},
	}
	walletsUpdatedEvent, err := json.Marshal(&WalletsUpdatedEvent{
		TriggeredBy: &TriggerEvent{
			ID:   trade.ID,
			Type: "trade",
		},
		Wallets: wallets,
	})
	if err != nil {
		fmt.Println("Error while marshalling:,", err.Error())
		return
	}

	if err := sendToTopic("wallets-updated", string(walletsUpdatedEvent)); err != nil {
		fmt.Println("Error while sending")
	}
}
