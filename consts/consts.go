package consts

import (
	"fmt"
	"log"
	"net"
)

const (
	Topic         = "TopicTest"
	ConsumerGroup = "TopicTestGroup"
	NameSrvHost   = "rmq-namesrv"
	NameSrvPort   = "9876"
)

var NameSrvEndpoint string

func init() {
	addrs, err := net.LookupHost(NameSrvHost)
	if err != nil {
		log.Fatalln(err)
	}

	NameSrvEndpoint = fmt.Sprintf("%s:%s", addrs[0], NameSrvPort)
	log.Printf("name server address is %s", NameSrvEndpoint)
}
