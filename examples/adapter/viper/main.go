package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-go-sdk/mconfig"
	"github.com/spf13/viper"
	"time"
)

func main() {
	viper.AddConfigPath("conf/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("Fatal error config file: %s \n", err)
	}
	mconfigClient := mconfig.NewClient(
		//mconfig.NameSpace("com.github.mhchlib"),
		mconfig.NameSpace("local_test"),
		mconfig.Registry("etcd://etcd.u.hcyang.top:31770"),
		//mconfig.DirectLinkAddress("127.0.0.1:8081"),
		//mconfig.DirectLinkAddress("mconfig.u.hcyang.top:8080"),
		//mconfig.Metadata("port", "8080"), //meta data
		mconfig.Metadata("ip", "192.0.0.1"),
		mconfig.AppKey("app_tPss5k5H_DS"),
		mconfig.ConfigKey("config_tPssCRQrGxh"),
		mconfig.RetryIntervalTime(15*time.Second),
	)
	mconfigClient.AdapterMconfigMergeToViper()

	for {
		log.Info(viper.Get("name.first"))
		<-time.After(5 * time.Second)
	}
	select {}
}
