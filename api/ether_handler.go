package api

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"encoding/json"
	"errors"
	"github.com/ninjadotorg/handshake-crowdfunding/utils"
	"strconv"
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
		handler.Process(m.Data)
	})
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}

	return &handler, nil
}

func (etherHandler *EtherHandler) Process(bytes []byte) (error) {
	logData := map[string]interface{}{}
	err := json.Unmarshal(bytes, &logData)
	if err != nil {
		log.Println("NewEthHandler.Process()", err)
		return err
	}
	event := logData["event"]
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
				offchainId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}


				err = crowdService.ProcessEventShake(hid, state, balance, offchainId)
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
