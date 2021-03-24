package mconfig

import (
	"github.com/spf13/viper"
	"sync"
)

// Adapter ...
type Adapter interface {
	AdapterMconfigMergeToViper(viperArr ...*viper.Viper) error
}

// AdapterCallBack ...
type AdapterCallBack func(key string, value interface{})

// AdapterCallBacks ...
type AdapterCallBacks struct {
	sync.RWMutex
	m []AdapterCallBack
}

// Add ...
func (c *AdapterCallBacks) Add(callback AdapterCallBack) {
	c.Lock()
	defer c.Unlock()
	c.m = append(c.m, callback)
}
