package adapter

import (
	"sync"
)

type CallBack func(key string, value interface{})

type CallBacks struct {
	sync.RWMutex
	m []CallBack
}

func (c *CallBacks) Add(callback CallBack) {
	c.Lock()
	defer c.Unlock()
	c.m = append(c.m, callback)
}

var callbacks *CallBacks

func init() {
	callbacks = &CallBacks{
		m: make([]CallBack, 0),
	}
}

func ExecuteAdapterCallBack(key string, value interface{}) {
	callbacks.RLock()
	defer callbacks.RUnlock()
	for _, callback := range callbacks.m {
		callback(key, value)
	}
}
