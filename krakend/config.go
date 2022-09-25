package krakend

import (
	"encoding/json"
	"errors"

	ipfilter "github.com/NEOMorphey/krakend-ipfilter"
	"github.com/luraproject/lura/config"
)

// Namespace is ipfilter's config key in extra config
const Namespace = "github_com/NEOMorphey/krakend-ipfilter"

// ErrNoConfig is returned when there is no config defined for the module
var ErrNoConfig = errors.New("no config defined for the module")


// ParseConfig build ip filter's Config
func ParseConfig(e config.ExtraConfig) (ipfilter.Config, error) {
	res := ipfilter.Config{}
	v, ok := e[Namespace].(map[string]interface{})
	if !ok {
		return nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

