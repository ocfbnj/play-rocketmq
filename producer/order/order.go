package main

import (
	"context"
	"log"

	"github.com/ocfbnj/play-rocketmq/consts"
	"github.com/ocfbnj/play-rocketmq/utils"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

func main() {
	rlog.SetLogLevel("error")

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{consts.Endpoint})),
		producer.WithQueueSelector(producer.NewHashQueueSelector()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer p.Shutdown()

	if err := p.Start(); err != nil {
		log.Fatalln(err)
	}

	msg := primitive.NewMessage(consts.Topic, []byte("Hello World"))
	msg.WithShardingKey("key5")

	res, err := p.SendSync(context.TODO(), msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(utils.MarshalJson(res))
}
