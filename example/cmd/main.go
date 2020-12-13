package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-go-sdk/client"
	"time"
)

func main() {
	config := client.NewMconfig(
		//client.NameSpace("local_test"),
		//client.Registry(client.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		client.DirectLinkAddress("10.98.174.23:8080"),
		client.ABFilters("port", "8080"), //meta data
		client.ABFilters("ip", "192.0.0.1"),
		client.AppKey("test"),
		client.ConfigKey("1000"),
		client.RetryTime(15*time.Second),
	)
	old := ""
	for {
		b := config.String("name")
		if b != old {
			log.Info(b)
			old = b
		}
		time.Sleep(time.Second * 3)
	}
}
