package models

import "github.com/kc-workspace/go-lib/mapper"

func NewConfig(m mapper.Mapper) (result *Config, err error) {
	result = new(Config)
	err = mapper.ToStruct(m, result)
	return
}
