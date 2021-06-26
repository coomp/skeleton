package skeleton

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type (
	Skeleton struct {
		handles  []HandlerFunc
		consumer rocketmq.PushConsumer
	}
	HandlerFunc func(Context) error
)

func (s *Skeleton) Use(handle ...HandlerFunc) {
	s.handles = append(s.handles, handle...)
}

func New() (s *Skeleton) {
	// daemon hold request from switch
	// when message received, process all callback
	s = &Skeleton{}
	s.handles = append(s.handles, NewDefaultHandle())
	return
}

func (s *Skeleton) StartWatch() {

	mqEndpointStr := GetOrDefault(MQ_ENDPOINT, "192.168.147.129:9876")
	mqEndpoints := strings.SplitN(mqEndpointStr, ",", -1)
	mqTopic := GetOrDefault(MQ_TOPIC, "coomp")

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(mqEndpoints)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerOrder(true),
	)

	err := c.Subscribe(mqTopic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
			// msg to context
			req := DefaultRequest{
				Data: msgs[i].String(),
			}

			for _, handle := range s.handles {
				handle(req)
			}
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func (s *Skeleton) StopWatch() {
	if s.consumer != nil {
		err := s.consumer.Shutdown()
		if err != nil {
			fmt.Printf("shutdown Consumer error: %s", err.Error())
		}
	}
}
