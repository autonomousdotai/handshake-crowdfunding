package api

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type EtherHandler struct {
	Client *pubsub.Client
	Topic  *pubsub.Topic
}

func NewEthHandler(pubsubClient *pubsub.Client, subName string) (*EtherHandler, error) {
	handler := EtherHandler{}

	handler.Client = pubsubClient

	topic := pubsubClient.Topic(subName)
	if topic == nil || topic.ID() != subName {
		topic, err := pubsubClient.CreateTopic(context.Background(), subName)
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		} else {
			handler.Topic = topic
		}
	} else {
		handler.Topic = topic
	}

	sub, err := pubsubClient.CreateSubscription(context.Background(), subName, pubsub.SubscriptionConfig{Topic: topic})

	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}

	return &handler, nil
}

func (etherHandler *EtherHandler) Process() (error) {

	return nil
}
