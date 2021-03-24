package mconfig

import "sync"

// ConfigWatch ...
type ConfigWatch struct {
	oldV     interface{}
	callback ConfigChangeCallBack
}

// ConfigChangeCallBack ...
type ConfigChangeCallBack func(key string, value interface{})

// ConfigWatchMap ...
type ConfigWatchMap struct {
	sync.RWMutex
	m map[string]*ConfigWatch
}

// AddConfigChangeCallBack ...
func (c *ConfigWatchMap) AddConfigChangeCallBack(key string, val interface{}, callback ConfigChangeCallBack) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = &ConfigWatch{
		oldV:     val,
		callback: callback,
	}
}
