package redis

import (
	"github.com/go-redis/redis"
	"github.com/jakoblorz/singapoor/stream"
)

type Subscriber struct {
	channel *redis.PubSub
	manager *stream.SubscriberHost
}

func (s Subscriber) GetMessageChannel() (<-chan interface{}, error) {

	channelIncoming := s.channel.Channel()
	channelOutgoing := make(chan interface{})

	go func() {

		var message *redis.Message

		for {
			message = <-channelIncoming
			channelOutgoing <- message.Payload
		}
	}()

	return channelOutgoing, nil
}

func (s Subscriber) AddSubscriber(fn func(interface{}) error) chan error {
	return s.manager.AddSubscriber(fn)
}

func (s Subscriber) NotifyOnMessageRecieve(msg interface{}) error {
	return s.manager.NotifyOnMessageRecieve(msg)
}

func (s Subscriber) NotifyOnStreamClose() error {
	return s.manager.NotifyOnStreamClose()
}