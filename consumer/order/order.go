package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ocfbnj/play-rocketmq/consts"
	"github.com/ocfbnj/play-rocketmq/utils"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

func main() {
	rlog.SetLogLevel("error")

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{consts.Endpoint})),
		consumer.WithGroupName(consts.ConsumerGroup),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumerOrder(true),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Shutdown()

	err = c.Subscribe(consts.Topic, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)

			log.Printf("orderly context: %s\n", utils.MarshalJson(orderlyCtx))
			log.Printf("msgs: %s\n", utils.MarshalJson(msgs))
			time.Sleep(1 * time.Minute)

			return consumer.ConsumeSuccess, nil
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err := c.Start(); err != nil {
		log.Fatalln(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	fmt.Println((<-sig).String())
}
