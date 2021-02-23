package adapter

import (
	log "github.com/mhchlib/logger"
	"github.com/spf13/viper"
	"strings"
)

func EnableMconfigCacheToViper(viperArr ...*viper.Viper) {
	if len(viperArr) == 0 {
		viperArr = append(viperArr, viper.GetViper())
	}
	for _, v := range viperArr {
		callbacks.Add(func(key string, value interface{}) {
			in := strings.NewReader(value.(string))
			err := v.MergeConfig(in)
			if err != nil {
				log.Error("mconfig merge config to viper error:", err)
			}
			log.Info("mconfig merge config to viper success")
		})
	}
}
