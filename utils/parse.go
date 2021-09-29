package utils

import (
	"io/ioutil"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("parse")

// MustParse a yaml file from a path to obj
func MustParse(path string, obj interface{}) {
	bts, err := ioutil.ReadFile(path)
	ExitOnErr(err, "read file")
	ExitOnErr(yaml.Unmarshal(bts, obj), "unmarshal yaml file")

}
