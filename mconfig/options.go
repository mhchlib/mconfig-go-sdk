package mconfig

import (
	log "github.com/mhchlib/logger"
	"time"
)

// Register_Type ...
type Register_Type string

var (
	// RegisterType_Etcd ...
	RegisterType_Etcd Register_Type = "etcd"
)

const (
	// Default_Retry_Time ...
	Default_Retry_Time = 5 * time.Second
	// Default_NameSpace ...
	Default_NameSpace = "com.github.mhchlib"
)

// Options ...
type Options struct {
	namespace         string
	registryAddress   string
	metadata          map[string]string
	appKey            string
	configKeys        []string
	envKey            string
	retryIntervalTime time.Duration
	enableRetry       bool
	enableNameSpace   bool
	directLinkAddress string
	logger            log.Logger
}

// NewOptions ...
func NewOptions() *Options {
	o := &Options{}
	return o
}

// Registry ...
func Registry(registerUrl string) Option {
	return func(options *Options) {
		options.registryAddress = registerUrl
	}
}

// DirectLinkAddress ...
func DirectLinkAddress(address string) Option {
	return func(options *Options) {
		options.directLinkAddress = address
	}
}

// NameSpace ...
func NameSpace(namespace string) Option {
	return func(options *Options) {
		options.namespace = namespace
		options.enableNameSpace = true
	}
}

// Metadata ...
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

// AppKey ...
func AppKey(appKey string) Option {
	return func(options *Options) {
		options.appKey = appKey
	}
}

// EnvKey ...
func EnvKey(envKey string) Option {
	return func(options *Options) {
		options.envKey = envKey
	}
}

// RetryIntervalTime ...
func RetryIntervalTime(t time.Duration) Option {
	return func(options *Options) {
		options.retryIntervalTime = t
		options.enableRetry = true
	}
}

// ConfigKey ...
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

// Logger ...
func Logger(log log.Logger) Option {
	return func(options *Options) {
		options.logger = log
	}
}
