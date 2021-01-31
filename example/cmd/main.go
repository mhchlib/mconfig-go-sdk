package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-go-sdk/client"
	"time"
)

func main() {
	config := client.NewMconfig(
		client.NameSpace("local_test"),
		client.Registry(client.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		//client.DirectLinkAddress("127.0.0.1:8081"),
		//client.Metadata("port", "8080"), //meta data
		//client.Metadata("ip", "192.0.0.1"),
		client.AppKey("appKey"),
		client.ConfigKey("configKey"),
		client.RetryTime(15*time.Second),
	)
	old := ""
	for {
		b := config.String("name")
		//log.Info(b)
		if b != old {
			log.Info(b)
			old = b
		}
		time.Sleep(time.Second * 3)
	}
}
