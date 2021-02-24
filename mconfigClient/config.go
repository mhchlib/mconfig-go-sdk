package mconfigClient

type Config interface {
	//base
	String(key string, defaultV ...string) string
	Int64(key string, defaultV ...int64) int64
	Bool(key string, defaultV ...bool) bool
	Map(key string, defaultV ...map[string]interface{}) map[string]interface{}
	SliceList(key string, defaultV ...[]interface{}) []interface{}
	Interface(key string, defaultV interface{}) interface{}
}

type WatchChange interface {
	OnWatchConfigChange(key string, f ConfigChangeCallBack)
}
