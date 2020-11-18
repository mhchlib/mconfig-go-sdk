package mconfig

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"log"
	"sync"
)

type FieldType int

type ConfigCache struct {
	Cache map[string]*FieldInterface
	sync.RWMutex
}

type OriginConfigCache struct {
	Data map[string]string
	sync.RWMutex
}

func (m *Mconfig) getValueFromCache(key string, fieldType FieldType) (FieldInterface, error) {
	cache := m.opts.Cache
	cache.RLock()
	data, ok := cache.Cache[key]
	cache.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	return *data, nil

}

func (m *Mconfig) reloadConfigData(key string, fieldType FieldType) (ret FieldInterface, err error) {
	log.Println("[mconfig] reload config from mconfig server")
	defer func() {
		if err == nil {
			//load the data to cache
			m.opts.Cache.Lock()
			m.opts.Cache.Cache[key] = &ret
			m.opts.Cache.Unlock()
		}
	}()
	configs := m.opts.ConfigsData
	var value gjson.Result
	flag := false
	configs.RLock()
	for _, v := range configs.Data {
		value = gjson.Get(v, key)
		if value.Exists() {
			flag = true
			break
		}
	}
	configs.RUnlock()
	if flag {
		switch fieldType {
		case FieldType_String:
			return &FieldInterface_String{
				Value: value.String(),
			}, nil

		case FieldType_Int:
			return &FieldInterface_Int{
				Value: value.Int(),
			}, nil
		case FieldType_Bool:
			return &FieldInterface_Bool{
				Value: value.Bool(),
			}, nil
		case FieldType_List:
			if value.IsArray() {
				raw := value.Raw
				list := []interface{}{}
				err := json.Unmarshal([]byte(raw), &list)
				if err != nil {
					return &FieldInterface_List{
						Value: nil,
					}, nil
				} else {
					return &FieldInterface_List{
						Value: list,
					}, nil
				}
			} else {
				return &FieldInterface_List{
					Value: nil,
				}, nil
			}
		case FieldType_Map:
			if value.IsObject() {
				raw := value.Raw
				maps := make(map[string]interface{})
				err := json.Unmarshal([]byte(raw), &maps)
				if err != nil {
					return &FieldInterface_Map{
						Value: nil,
					}, nil
				} else {
					return &FieldInterface_Map{
						Value: maps,
					}, nil
				}
			} else {
				return &FieldInterface_Map{
					Value: nil,
				}, nil
			}
		}

	}
	return nil, errors.New("not found from origin config")
}
