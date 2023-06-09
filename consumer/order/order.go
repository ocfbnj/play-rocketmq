package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ocfbnj/play-rocketmq/consts"
	"github.com/ocfbnj/play-rocketmq/utils"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

var wg = sync.WaitGroup{}

func interceptor(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
	wg.Add(1)
	defer wg.Done()

	return next(ctx, req, reply)
}

func main() {
	rlog.SetLogLevel("error")

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{consts.NameSrvEndpoint})),
		consumer.WithGroupName(consts.ConsumerGroup),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumerOrder(true),
		consumer.WithInterceptor(interceptor),
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
			time.Sleep(10 * time.Second)
			log.Println("done")

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

	c.Shutdown()
	wg.Wait()
}
