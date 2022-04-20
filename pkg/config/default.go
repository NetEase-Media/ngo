package config

import (
	"time"
)

var defaultConfiguration = Empty()

func SetConfiguration(cfg *Configuration) {
	defaultConfiguration = cfg
}

func GetConfiguration() *Configuration {
	return defaultConfiguration
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return defaultConfiguration.Get(key) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return defaultConfiguration.GetString(key) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return defaultConfiguration.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return defaultConfiguration.GetInt(key) }

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 { return defaultConfiguration.GetInt32(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return defaultConfiguration.GetInt64(key) }

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint { return defaultConfiguration.GetUint(key) }

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 { return defaultConfiguration.GetUint32(key) }

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 { return defaultConfiguration.GetUint64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return defaultConfiguration.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return defaultConfiguration.GetTime(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return defaultConfiguration.GetDuration(key) }

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int { return defaultConfiguration.GetIntSlice(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return defaultConfiguration.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return defaultConfiguration.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string {
	return defaultConfiguration.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return defaultConfiguration.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint { return defaultConfiguration.GetSizeInBytes(key) }

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface{}, opt ...DecoderConfigOption) error {
	return defaultConfiguration.UnmarshalKey(key, rawVal, opt...)
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func Unmarshal(rawVal interface{}, opts ...DecoderConfigOption) error {
	return defaultConfiguration.Unmarshal(rawVal, opts...)
}

func Exists(key string) bool {
	return defaultConfiguration.Get(key) != nil
}

func Sub(key string) *Configuration {
	return defaultConfiguration.Sub(key)
}

func SubSlice(key string) []*Configuration {
	return defaultConfiguration.SubSlice(key)
}

func OnChange(fn func(*Configuration)) {
	defaultConfiguration.OnChange(fn)
}
