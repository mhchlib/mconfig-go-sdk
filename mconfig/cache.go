package mconfig

import (
	"sync"
)

type FieldType int

type ConfigCache struct {
	Cache map[string]*FieldInterface
	sync.RWMutex
}

type OriginConfigCache struct {
	m map[string]string
	sync.RWMutex
}

func (originConfigCache OriginConfigCache) Put(key string, value string) {
	originConfigCache.Lock()
	defer originConfigCache.Unlock()
	originConfigCache.m[key] = value
}
