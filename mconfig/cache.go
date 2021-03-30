package mconfig

import (
	"sync"
)

// FieldType ...
type FieldType int

// ConfigCache ...
type ConfigCache struct {
	Cache map[string]*FieldInterface
	sync.RWMutex
}

// OriginConfigCache ...
type OriginConfigCache struct {
	m map[string]string
	sync.RWMutex
}

// Put ...
func (originConfigCache OriginConfigCache) Put(key string, value string) {
	originConfigCache.Lock()
	defer originConfigCache.Unlock()
	originConfigCache.m[key] = value
}
