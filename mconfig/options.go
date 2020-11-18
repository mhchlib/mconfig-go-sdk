package mconfig

type Register_Type string

var (
	RegisterType_Etcd Register_Type = "etcd"
	//RegisterType_Consoul RegisterType = "consoul"
)

const (
	Default_Retry_Num = 5
	Default_NameSpace = "com.github.mhchlib.mconfig"
)

type Options struct {
	NameSpace       string
	RegistryUrl     string
	RegistryType    Register_Type
	ABFilters       map[string]string
	AppKey          string
	ConfigKeys      []string
	ConfigsData     *OriginConfigCache
	Cache           *ConfigCache
	RetryNum        int
	EnableRetry     bool
	EnableNameSpace bool
	EnableRegistry  bool
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

func Registry(registerType Register_Type, registerUrl string) Option {
	return func(options *Options) {
		options.RegistryType = registerType
		options.RegistryUrl = registerUrl
		options.EnableRegistry = true
	}
}

func NameSpace(namespace string) Option {
	return func(options *Options) {
		options.NameSpace = namespace
		options.EnableNameSpace = true
	}
}

func ABFilters(key string, value string) Option {
	return func(options *Options) {
		abfilters := options.ABFilters
		if abfilters == nil {
			abfilters = map[string]string{}
		}
		abfilters[key] = value
		options.ABFilters = abfilters
	}
}

func AppKey(appKey string) Option {
	return func(options *Options) {
		options.AppKey = appKey
	}
}

func Retry(num int) Option {
	return func(options *Options) {
		options.RetryNum = num
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
