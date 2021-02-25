package mconfig

import "sync"

type ConfigWatch struct {
	oldV     interface{}
	callback ConfigChangeCallBack
}

type ConfigChangeCallBack func(key string, value interface{})

type ConfigWatchMap struct {
	sync.RWMutex
	m map[string]*ConfigWatch
}

func (c *ConfigWatchMap) AddConfigChangeCallBack(key string, val interface{}, callback ConfigChangeCallBack) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = &ConfigWatch{
		oldV:     val,
		callback: callback,
	}
}
