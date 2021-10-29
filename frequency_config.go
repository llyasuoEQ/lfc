package lfc

import (
	"sort"
)

type Fconfig struct {
	AppName string      `toml:"app_name" yaml:"app_name" json:"app_name"`
	Rules   []rules     `toml:"rules" yaml:"rules" json:"rules"`
	Period  string      `toml:"period" yaml:"period" json:"period"`
	Redis   redisConfig `toml:"redis" yaml:"redis" json:"redis"`
}

type rules struct {
	Name      string   `toml:"name" yaml:"name" json:"name"`
	Period    string   `toml:"period" yaml:"period" json:"period"`
	Threshold int64    `toml:"threshold" yaml:"threshold" json:"threshold"`
	Code      int64    `toml:"code" yaml:"code" json:"code"`
	Fields    []string `toml:"fields" yaml:"fields" json:"fields"`
}

func (r *rules) getSortFields() (fields []string) {
	if r == nil {
		return
	}
	fields = r.Fields
	sort.Strings(fields)
	return
}
