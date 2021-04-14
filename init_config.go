package lfc

import (
	"fmt"
	"sync"

	toml "github.com/pelletier/go-toml"
)

const (
	CONFIG_KEY = "frequency"
)

var (
	globalInitFCOnce sync.Once
	globalFConfig    Fconfig
)

func InitByConfigStr(confStr string) (productRule *ProductRule, err error) {
	var tree *toml.Tree
	tree, err = toml.Load(confStr)
	if err != nil {
		return
	}
	productRule, err = InitConfig(tree)
	return
}

func InitConfig(tree *toml.Tree) (productRule *ProductRule, err error) {
	globalInitFCOnce.Do(func() {
		t, ok := tree.Get(CONFIG_KEY).(*toml.Tree)
		if !ok {
			err = fmt.Errorf("not fount %s key in config toml", CONFIG_KEY)
		}
		err = t.Unmarshal(&globalFConfig)
		if err != nil {
			err = fmt.Errorf("unmarshal from toml tree to struct failed[err=%v]", err)
			return
		}
		// new redis
		globalRedisClient = globalFConfig.Redis.newRedis()

		// formate config
		_, err = globalProductRule.formatByFrequencyConfig(&globalFConfig)
	})
	return
}
