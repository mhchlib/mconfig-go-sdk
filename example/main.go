package main

import (
	"github.com/mhchlib/mconfig-go-sdk/pkg"
	"log"
	"time"
)

func main() {
	config := pkg.NewMconfig(
		pkg.Registry(pkg.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		pkg.ABFilters("port", "8080"), //meta data
		pkg.ABFilters("ip", "192.0.0.1"),
		pkg.AppKey("1002"),
		pkg.ConfigKey("1000"),
		pkg.RetryTime(15*time.Second),
	)
	//`{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	for {
		b := config.String("name")
		log.Print(b)
		time.Sleep(time.Second * 3)
	}
}
