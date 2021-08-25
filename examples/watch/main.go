package main

import (
	"github.com/mhchlib/mconfig-go-sdk/mconfig"
	"github.com/prometheus/common/log"
	"time"
)

func main() {

	mconfigClient := mconfig.NewClient(
		mconfig.NameSpace("local_test"),
		//mconfigClient.Registry(mconfigClient.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		mconfig.DirectLinkAddress("127.0.0.1:8081"),
		//mconfig.Metadata("port", "8080"), //meta data
		mconfig.Metadata("ip", "192.0.0.1"),
		mconfig.AppKey("app_tPss5k5H_DS"),
		mconfig.ConfigKey("config_tPssCRQrGxh"),
		mconfig.RetryIntervalTime(15*time.Second),
	)
	mconfigClient.OnWatchConfigChange("name", func(key string, value interface{}) {
		log.Info(key, value)
	})

	select {}
}
