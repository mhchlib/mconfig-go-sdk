package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-go-sdk/mconfigClient"
	"github.com/spf13/viper"
	"time"
)

func main() {
	viper.AddConfigPath("conf/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("Fatal error config file: %s \n", err)
	}
	mconfig := mconfigClient.NewClient(
		mconfigClient.NameSpace("local_test"),
		//mconfigClient.Registry(mconfigClient.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		mconfigClient.DirectLinkAddress("127.0.0.1:8081"),
		//mconfig.Metadata("port", "8080"), //meta data
		mconfigClient.Metadata("ip", "192.0.0.1"),
		mconfigClient.AppKey("app_tPss5k5H_DS"),
		mconfigClient.ConfigKey("config_tPssCRQrGxh"),
		mconfigClient.RetryIntervalTime(15*time.Second),
	)
	mconfig.AdapterMconfigMergeToViper()

	//for {
	//	log.Info(viper.Get("name.first"))
	//	<-time.After(5 * time.Second)
	//}
	mconfig.OnWatchConfigChange("name.first", func(key string, value interface{}) {
		log.Info("change", key, value)
	})
	select {}
}
