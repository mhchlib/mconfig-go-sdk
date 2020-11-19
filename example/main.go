package main

import (
	"github.com/mhchlib/mconfig-go-sdk/mconfig"
	"log"
	"time"
)

func main() {
	config := mconfig.NewMconfig(
		mconfig.Registry(mconfig.RegisterType_Etcd, "etcd.u.hcyang.top:31770"),
		mconfig.ABFilters("port", "8080"),
		mconfig.ABFilters("ip", "192.0.0.1"),
		mconfig.AppKey("1002"),
		mconfig.ConfigKey("1000"),
		mconfig.Retry(2),
	)
	//`{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	for {
		b := config.String("name.last")
		log.Print(b)
		time.Sleep(time.Second * 3)
	}
}
