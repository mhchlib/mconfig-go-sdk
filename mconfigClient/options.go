package mconfigClient

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
	namespace         string
	registryUrl       []string
	registryType      Register_Type
	metadata          map[string]string
	appKey            string
	configKeys        []string
	retryIntervalTime time.Duration
	enableRetry       bool
	enableNameSpace   bool
	enableRegistry    bool
	directLinkAddress string
	logger            log.Logger
}

func NewOptions() *Options {
	o := &Options{}
	return o
}

func Registry(registerType Register_Type, registerUrl []string) Option {
	return func(options *Options) {
		options.registryType = registerType
		options.registryUrl = registerUrl
		options.enableRegistry = true
	}
}

func DirectLinkAddress(address string) Option {
	return func(options *Options) {
		options.directLinkAddress = address
	}
}

func NameSpace(namespace string) Option {
	return func(options *Options) {
		options.namespace = namespace
		options.enableNameSpace = true
	}
}

func Metadata(key string, value string) Option {
	return func(options *Options) {
		abfilters := options.metadata
		if abfilters == nil {
			abfilters = map[string]string{}
		}
		abfilters[key] = value
		options.metadata = abfilters
	}
}

func AppKey(appKey string) Option {
	return func(options *Options) {
		options.appKey = appKey
	}
}

func RetryIntervalTime(t time.Duration) Option {
	return func(options *Options) {
		options.retryIntervalTime = t
		options.enableRetry = true
	}
}

func ConfigKey(keys ...string) Option {
	return func(options *Options) {
		configKeys := options.configKeys
		if configKeys == nil {
			configKeys = []string{}
		}
		configKeys = append(configKeys, keys...)
		options.configKeys = configKeys
	}
}

func Logger(log log.Logger) Option {
	return func(options *Options) {
		options.logger = log
	}
}
