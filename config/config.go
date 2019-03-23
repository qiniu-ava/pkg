package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/pkg/errors"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "config", "./config", `config file path, default to "./config"`)
}

// LoadConfigFile loads a json config file located by commandline flag '-config'
func LoadConfigFile(v interface{}) error {
	buf, e := ioutil.ReadFile(filePath)
	if e != nil {
		return errors.Wrap(e, "read config file failed, did you set the right path?")
	}
	if e := json.Unmarshal(buf, v); e != nil {
		return errors.Wrap(e, "unmarshal config file failed, please check the file path and content")
	}
	return nil
}
