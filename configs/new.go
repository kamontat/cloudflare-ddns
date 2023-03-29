package configs

import "github.com/kc-workspace/go-lib/mapper"

func New(m mapper.Mapper) (result *Config, err error) {
	result = new(Config)
	err = mapper.ToStruct(m, result)

	return
}
