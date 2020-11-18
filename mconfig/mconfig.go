package mconfig

import "log"

type Mconfig struct {
	opts *Options
}

type Option func(*Options)

func NewMconfig(opts ...Option) Config {
	//获取默认options
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	if options.EnableRetry == false {
		options.RetryNum = Default_Retry_Num
	}

	if options.EnableNameSpace == false {
		options.NameSpace = Default_NameSpace
	}

	if options.EnableRegistry == false {
		log.Fatal("[mconfig] you should set an Registry address to link mconfig server")
	}

	m := &Mconfig{}
	m.opts = options
	m.initMconfigLink()
	return m
}

func (m *Mconfig) String(key string, defaultVs ...string) string {
	var defaultV string
	if len(defaultVs) >= 1 {
		defaultV = defaultVs[0]
	} else {
		defaultV = ""
	}
	cache, err := m.getValueFromCache(key, FieldType_String)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_String)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_String)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_String)
		if err != nil {
			return defaultV
		}
		data, _ = (cache).(*FieldInterface_String)
	}
	return data.Value

}

func (m *Mconfig) Int64(key string, defaultVs ...int64) int64 {
	var defaultV int64
	if len(defaultVs) >= 1 {
		defaultV = defaultVs[0]
	} else {
		defaultV = 0
	}
	cache, err := m.getValueFromCache(key, FieldType_Int)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_Int)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_Int)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_Int)
		if err != nil {
			return defaultV
		}
		data, _ = cache.(*FieldInterface_Int)
	}
	return data.Value
}

func (m *Mconfig) Bool(key string, defaultVs ...bool) bool {
	var defaultV bool
	if len(defaultVs) >= 1 {
		defaultV = defaultVs[0]
	} else {
		defaultV = false
	}
	cache, err := m.getValueFromCache(key, FieldType_Bool)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_Bool)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_Bool)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_Bool)
		if err != nil {
			return defaultV
		}
		data, _ = cache.(*FieldInterface_Bool)
	}
	return data.Value
}

func (m *Mconfig) Map(key string, defaultVs ...map[string]interface{}) map[string]interface{} {
	var defaultV map[string]interface{}
	if len(defaultVs) >= 1 {
		defaultV = defaultVs[0]
	} else {
		defaultV = nil
	}
	cache, err := m.getValueFromCache(key, FieldType_Map)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_Map)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_Map)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_Map)
		if err != nil {
			return defaultV
		}
		data, _ = cache.(*FieldInterface_Map)
	}
	return data.Value
}

func (m *Mconfig) List(key string, defaultVs ...[]interface{}) []interface{} {
	var defaultV []interface{}
	if len(defaultVs) >= 1 {
		defaultV = defaultVs[0]
	} else {
		defaultV = nil
	}

	cache, err := m.getValueFromCache(key, FieldType_List)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_List)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_List)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_List)
		if err != nil {
			return defaultV
		}
		data, _ = cache.(*FieldInterface_List)
	}
	return data.Value
}
