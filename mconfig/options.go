package mconfig

import (
	log "github.com/mhchlib/logger"
	"time"
)

type Register_Type string

var (
	RegisterType_Etcd Register_Type = "etcd"
	//RegisterType_Consoul RegisterType = "consoul"
)

const (
	Default_Retry_Time = 5 * time.Second
	Default_NameSpace  = "com.github.mhchlib"
)

type Options struct {
	Namespace         string
	RegistryUrl       []string
	RegistryType      Register_Type
	Metadata          map[string]string
	AppKey            string
	ConfigKeys        []string
	ConfigsData       *OriginConfigCache
	Cache             *ConfigCache
	RetryTime         time.Duration
	EnableRetry       bool
	EnableNameSpace   bool
	EnableRegistry    bool
	DirectLinkAddress string
	Logger            log.Logger
}

func NewOptions() *Options {
	o := &Options{}
	o.Cache = &ConfigCache{
		Cache: make(map[string]*FieldInterface),
	}
	o.ConfigsData = &OriginConfigCache{
		Data: make(map[string]string),
	}
	return o
}

func Registry(registerType Register_Type, registerUrl []string) Option {
	return func(options *Options) {
		options.RegistryType = registerType
		options.RegistryUrl = registerUrl
		options.EnableRegistry = true
	}
}

func DirectLinkAddress(address string) Option {
	return func(options *Options) {
		options.DirectLinkAddress = address
	}
}

func NameSpace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
		options.EnableNameSpace = true
	}
}

func Metadata(key string, value string) Option {
	return func(options *Options) {
		abfilters := options.Metadata
		if abfilters == nil {
			abfilters = map[string]string{}
		}
		abfilters[key] = value
		options.Metadata = abfilters
	}
}

func AppKey(appKey string) Option {
	return func(options *Options) {
		options.AppKey = appKey
	}
}

func RetryTime(t time.Duration) Option {
	return func(options *Options) {
		options.RetryTime = t
		options.EnableRetry = true
	}
}

func ConfigKey(keys ...string) Option {
	return func(options *Options) {
		configKeys := options.ConfigKeys
		if configKeys == nil {
			configKeys = []string{}
		}
		configKeys = append(configKeys, keys...)
		options.ConfigKeys = configKeys
	}
}

func Logger(log log.Logger) Option {
	return func(options *Options) {
		options.Logger = log
	}
}
