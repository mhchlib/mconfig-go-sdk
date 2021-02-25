package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-go-sdk/mconfig"
	"time"
)

func main() {
	config := mconfig.NewClient(
		mconfig.NameSpace("local_test"),
		//mconfig.Registry(mconfig.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		mconfig.DirectLinkAddress("127.0.0.1:8081"),
		//mconfig.Metadata("port", "8080"), //meta data
		mconfig.Metadata("ip", "192.0.0.1"),
		mconfig.AppKey("app_tPss5k5H_DS"),
		mconfig.ConfigKey("config_tPssCRQrGxh"),
		mconfig.RetryTime(15*time.Second),
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
