package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/miracl/conflate"
)

// AppConfig ...
type AppConfig struct {
	data map[string]json.RawMessage
}

// NewAppConfig ...
func NewAppConfig() *AppConfig {
	merge := conflate.New()
	_ = merge.AddFiles("config.json")

	for _, env := range os.Environ() {
		_ = merge.AddGo(env)
	}

	merged, _ := merge.MarshalJSON()
	var data map[string]json.RawMessage
	_ = json.Unmarshal(merged, &data)

	return &AppConfig{data}
}

// Value ...
func (cfg *AppConfig) Value(conf interface{}) error {
	if reflect.TypeOf(conf).Kind() != reflect.Ptr {
		return errors.New("config is not a pointer type")
	}
	configName := reflect.TypeOf(conf).Elem().Name()

	raw, ok := cfg.data[configName]
	if !ok {
		return nil
	}

	return json.Unmarshal(raw, conf)
}
