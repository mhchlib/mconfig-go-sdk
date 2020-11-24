package client

type Config interface {
	String(key string, defaultV ...string) string
	Int64(key string, defaultV ...int64) int64
	Bool(key string, defaultV ...bool) bool
	Map(key string, defaultV ...map[string]interface{}) map[string]interface{}
	List(key string, defaultV ...[]interface{}) []interface{}
}
