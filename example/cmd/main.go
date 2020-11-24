package main

import (
	"github.com/mhchlib/mconfig-go-sdk/client"
	"log"
	"time"
)

func main() {
	config := client.NewMconfig(
		client.Registry(client.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		client.ABFilters("port", "8080"), //meta data
		client.ABFilters("ip", "192.0.0.1"),
		client.AppKey("1002"),
		client.ConfigKey("1000"),
		client.RetryTime(15*time.Second),
	)
	//`{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	for {
		b := config.String("name")
		log.Print(b)
		time.Sleep(time.Second * 3)
	}
}
