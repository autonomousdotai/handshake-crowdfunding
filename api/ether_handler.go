package api

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type EtherHandler struct {
	BubsubClient       *pubsub.Client
	BubsubSubscription *pubsub.Subscription
}

func NewEthHandler(pubsubClient *pubsub.Client, topicName, subscriptionName string) (*EtherHandler, error) {
	handler := EtherHandler{}

	handler.BubsubClient = pubsubClient

	topic := pubsubClient.Topic(topicName)
	if topic == nil || topic.ID() != topicName {
		var err error
		topic, err = pubsubClient.CreateTopic(context.Background(), topicName)
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}

	sub := pubsubClient.Subscription(subscriptionName)
	existed, err := sub.Exists(context.Background())
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}
	if !existed {
		var err error
		sub, err = pubsubClient.CreateSubscription(context.Background(), subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message : %s", m.Data)
		m.Ack()
		handler.Process(string(m.Data))
	})
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}

	return &handler, nil
}

func (etherHandler *EtherHandler) Process(message string) (error) {

	return nil
}
