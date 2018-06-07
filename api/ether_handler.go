package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/ninjadotorg/handshake-crowdfunding/utils"
)

type EtherHandler struct {
	BubsubClient       *pubsub.Client
	PubsubSubscription *pubsub.Subscription
}

func NewEthHandler(pubsubClient *pubsub.Client, topicName, subscriptionName string) (*EtherHandler, error) {
	handler := EtherHandler{}

	handler.BubsubClient = pubsubClient

	topic := pubsubClient.Topic(topicName)
	existed, err := topic.Exists(context.Background())
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}
	if topic == nil || !existed {
		var err error
		topic, err = pubsubClient.CreateTopic(context.Background(), topicName)
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}

	sub := pubsubClient.Subscription(subscriptionName)
	existed, err = sub.Exists(context.Background())
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}
	if sub == nil || !existed {
		var err error
		sub, err = pubsubClient.CreateSubscription(context.Background(), subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}

	handler.PubsubSubscription = sub

	return &handler, nil
}

func (etherHandler *EtherHandler) Receive() error {
	err := etherHandler.PubsubSubscription.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message : %s", m.Data)
		m.Ack()
		etherHandler.Process(m.Data)
	})
	if err != nil {
		log.Println("NewEthHandler", err)
		return err
	}
	return nil
}

func (etherHandler *EtherHandler) Process(bytes []byte) error {
	logData := map[string]interface{}{}
	err := json.Unmarshal(bytes, &logData)
	if err != nil {
		log.Println("NewEthHandler.Process()", err)
		return err
	}
	event := logData["event"].(string)
	fromAddress := logData["from_address"].(string)
	data, ok := logData["data"].(map[string]interface{})
	if !ok {
		return errors.New("data is missed")
	}
	switch event {
	case "__init":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_CROWD {
				offchainId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventInit(hid, offchainId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__shake":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			val, ok = data["state"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			state := int(val)
			balance, ok := data["balance"].(float64)
			if !ok {
				return errors.New("balance is invalid")
			}
			if offchainType == utils.OFFCHAIN_CROWD_SHAKE {
				crowdFundingShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventShake(hid, state, balance, crowdFundingShakeId, fromAddress)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__unshake":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			val, ok = data["state"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			state := int(val)
			balance, ok := data["balance"].(float64)
			if !ok {
				return errors.New("balance is invalid")
			}
			if offchainType == utils.OFFCHAIN_USER {
				userId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventUnShake(hid, state, balance, userId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__cancel":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			val, ok = data["state"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			state := int(val)
			if offchainType == utils.OFFCHAIN_USER {
				userId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventCancel(hid, state, userId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__refund":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			val, ok = data["state"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			state := int(val)
			if offchainType == utils.OFFCHAIN_USER {
				userId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventRefund(hid, state, userId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__stop":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			val, ok = data["state"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			state := int(val)
			if offchainType == utils.OFFCHAIN_CROWD {
				crowdFundingId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventStop(hid, state, crowdFundingId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	case "__withdraw":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			amount, ok := data["amount"].(float64)
			if !ok {
				return errors.New("state is invalid")
			}
			if offchainType == utils.OFFCHAIN_CROWD {
				crowdFundingId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = crowdService.ProcessEventWithdraw(hid, amount, crowdFundingId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}
			return nil
		}
		break
	}
	return nil
}
