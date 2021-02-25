package mconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/register"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type MconfigClientV1 struct {
	opts            *Options
	apdterCallBacks AdapterCallBacks
	configsData     *OriginConfigCache
	cache           *ConfigCache
	configWatchMap  *ConfigWatchMap
}

func NewClient(opts ...Option) MconfigClient {
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	if options.enableRetry == false {
		options.retryIntervalTime = Default_Retry_Time
	}

	if options.enableNameSpace == false {
		options.namespace = Default_NameSpace
	}

	if options.logger == nil {
		options.logger = log.NewLogger(
			log.EnableCodeData(false),
			log.MetaData("provider", "mconfig"),
		)
	}
	m := &MconfigClientV1{}
	m.cache = &ConfigCache{
		Cache: make(map[string]*FieldInterface),
	}
	m.configsData = &OriginConfigCache{
		m: make(map[string]string),
	}
	m.configWatchMap = &ConfigWatchMap{
		m: make(map[string]*ConfigWatch),
	}
	m.opts = options
	m.initMconfigEngine()
	return m
}

func (m *MconfigClientV1) getValueFromCache(key string, fieldType FieldType) (FieldInterface, error) {
	cache := m.cache
	cache.RLock()
	data, ok := cache.Cache[key]
	cache.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	return *data, nil

}

func (m *MconfigClientV1) reloadConfigData(key string, fieldType FieldType) (ret FieldInterface, err error) {
	log := m.opts.logger
	log.Info("reload config from mconfig server")
	defer func() {
		if err == nil {
			//load the data to cache
			m.cache.Lock()
			m.cache.Cache[key] = &ret
			m.cache.Unlock()
		}
	}()
	configs := m.configsData
	var value gjson.Result
	flag := false
	configs.RLock()
	for _, v := range configs.m {
		v = strings.ReplaceAll(v, "'", "\"")
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
		case FieldType_Interface:
			return &FieldInterface_Interface{
				Value: value,
			}, nil
		}

	}
	return nil, errors.New("not found from origin config")
}

func (m *MconfigClientV1) AdapterMconfigMergeToViper(viperArr ...*viper.Viper) error {
	if len(viperArr) == 0 {
		viperArr = append(viperArr, viper.GetViper())
	}
	var tmpCallBackArr []AdapterCallBack
	for _, v := range viperArr {
		callback := func(key string, value interface{}) {
			in := strings.NewReader(value.(string))
			err := v.MergeConfig(in)
			if err != nil {
				log.Error("mconfig merge config to viper error:", err)
			}
		}
		tmpCallBackArr = append(tmpCallBackArr, callback)
		m.apdterCallBacks.Add(callback)
	}

	//init
	m.configsData.RLock()
	defer m.configsData.RUnlock()
	for key, value := range m.configsData.m {
		for _, callback := range tmpCallBackArr {
			callback(key, value)
		}
	}
	return nil
}

func (m *MconfigClientV1) initAddressProvider() func(serviceName string) (*register.ServiceVal, error) {
	log := m.opts.logger
	if m.opts.enableRegistry {
		regClient, err := register.InitRegister(
			register.SelectEtcdRegister(),
			register.ResgisterAddress(m.opts.registryUrl),
			register.Namespace(m.opts.namespace),
		)
		if err != nil {
			log.Fatal("register fail")
		}
		return func(serviceName string) (*register.ServiceVal, error) {
			val, err := regClient.GetService(serviceName)
			if err != nil {
				return nil, err
			}
			log.Info("get service address:", fmt.Sprintf("%+v", val))
			return val, nil
		}
	}
	if m.opts.directLinkAddress == "" {
		log.Fatal("you should provider a direct link address or an register center address...")
	}
	return func(serviceName string) (*register.ServiceVal, error) {
		return &register.ServiceVal{
			Address: m.opts.directLinkAddress,
		}, nil
	}
}

func (m *MconfigClientV1) initMconfigEngine() {
	log := m.opts.logger
	addressProvider := m.initAddressProvider()
	request := &server.WatchConfigStreamRequest{
		AppKey:     m.opts.appKey,
		ConfigKeys: m.opts.configKeys,
		Metadata:   m.opts.metadata,
	}
	//添加连接断开重试机制
	retryTime := m.opts.retryIntervalTime
	once := true
	enableRetry := false
	started := make(chan interface{})
	go func(m *MconfigClientV1, started chan interface{}) {
		for {
			if enableRetry {
				log.Info("mconfig request fail, will retry request after ", retryTime)
				<-time.After(retryTime)
			}
			enableRetry = true
			service, err := addressProvider("mconfig")
			if err != nil {
				log.Info(err)
				continue
			}
			withTimeout, _ := context.WithTimeout(context.Background(), time.Second*3)
			dial, err := grpc.DialContext(withTimeout, service.Address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Info(err, " addr:", fmt.Sprintf("%+v", service))
				continue
			}
			mConfigService := server.NewMConfigClient(dial)
			stream, err := mConfigService.WatchConfigStream(context.Background())
			if err != nil {
				log.Info(err)
				continue
			}
			err = stream.SendMsg(request)
			if err != nil {
				log.Info(err)
				continue
			}
			for {
				recv, err := stream.Recv()
				log.Info(recv)
				if err != nil {
					log.Info(err)
					break
				}

				configs := recv.Configs
				originConfigData := m.configsData
				for _, config := range configs {
					originConfigData.Put(config.ConfigKey, config.Val)
					m.executeConfigAdapter(config.ConfigKey, config.Val)
				}

				go m.executeConfigWatchNotify()

				if once {
					started <- &struct{}{}
					once = false
				}
				//refer the cache
				//todo: 这里可以优化为主动去更新cache中内容，但是这个需要加大cache大小，带上类型
				// support soon....
				m.cache.Lock()
				m.cache.Cache = map[string]*FieldInterface{}
				m.cache.Unlock()
				log.Info("refresh mconfig cache success")
			}
		}
	}(m, started)
	<-started
	close(started)
}

func (m *MconfigClientV1) String(key string, defaultVs ...string) string {
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

func (m *MconfigClientV1) Int64(key string, defaultVs ...int64) int64 {
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

func (m *MconfigClientV1) Bool(key string, defaultVs ...bool) bool {
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

func (m *MconfigClientV1) Map(key string, defaultVs ...map[string]interface{}) map[string]interface{} {
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

func (m *MconfigClientV1) SliceList(key string, defaultVs ...[]interface{}) []interface{} {
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

func (m *MconfigClientV1) Interface(key string, defaultV interface{}) interface{} {
	cache, err := m.getValueFromCache(key, FieldType_Interface)
	if err != nil {
		cache, err = m.reloadConfigData(key, FieldType_Interface)
		if err != nil {
			return defaultV
		}
	}
	data, ok := cache.(*FieldInterface_Interface)
	if !ok {
		cache, err = m.reloadConfigData(key, FieldType_Interface)
		if err != nil {
			return defaultV
		}
		data, _ = cache.(*FieldInterface_Interface)
	}
	return data.Value
}

func (m *MconfigClientV1) OnWatchConfigChange(key string, f ConfigChangeCallBack) {
	m.configWatchMap.AddConfigChangeCallBack(key, m.Interface(key, struct{}{}), f)
}

func (m *MconfigClientV1) executeConfigAdapter(key string, value string) {
	m.apdterCallBacks.RLock()
	defer m.apdterCallBacks.RUnlock()
	for _, callback := range m.apdterCallBacks.m {
		callback(key, value)
	}
}

func (m *MconfigClientV1) executeConfigWatchNotify() {
	m.configWatchMap.Lock()
	defer m.configWatchMap.Unlock()
	for key, configWatch := range m.configWatchMap.m {
		newVal := m.Interface(key, struct{}{})
		//log.Info(key, configWatch.oldV, newVal)
		if newVal != configWatch.oldV {
			configWatch.callback(key, newVal)
			configWatch.oldV = newVal
		}
	}
}
