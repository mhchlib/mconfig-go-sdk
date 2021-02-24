package mconfigClient

import (
	"github.com/spf13/viper"
	"sync"
)

type Adapter interface {
	AdapterMconfigMergeToViper(viperArr ...*viper.Viper) error
}

type AdapterCallBack func(key string, value interface{})

type AdapterCallBacks struct {
	sync.RWMutex
	m []AdapterCallBack
}

func (c *AdapterCallBacks) Add(callback AdapterCallBack) {
	c.Lock()
	defer c.Unlock()
	c.m = append(c.m, callback)
}
