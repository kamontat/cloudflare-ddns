package configs

type Config struct {
	Entities    []Entity    `json:"entities" yaml:"entities"`
	Settings    Setting     `json:"settings" yaml:"settings"`
	Secrets     Secret      `json:"secrets" yaml:"secrets"`
	Development Development `json:"development" yaml:"development"`
}
