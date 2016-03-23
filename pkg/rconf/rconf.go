package rconf

import (
//     "fmt"
)

var extArgs map[string]string

func init() {
	extArgs = make(map[string]string)
}

func Init(rconfDir string, args map[string]string) {
	extArgs = args

	parseVarConf(rconfDir + "/var.json")
}
