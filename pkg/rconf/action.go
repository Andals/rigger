package rconf

import (
	"andals/gobox/misc"
	"encoding/json"
	//     "fmt"
	"io/ioutil"
	"strconv"
)

type mkdirItem struct {
	Dir  string
	Mask string
	Sudo bool
}

type actionConf struct {
	Mkdir []*mkdirItem
	Exec  []string
}

var aconf *actionConf

func init() {
	aconf = new(actionConf)
}

func parseActionConf(path string) {
	if !misc.FileExist(path) {
		panic("Action conf not exists in " + path)
	}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, &aconf)
	if nil != err {
		panic("Parse action conf error: " + err.Error())
	}

	for i, item := range aconf.Mkdir {
		pkeyPrefix := "action mkdir " + strconv.Itoa(i) + " item "
		item.Dir = ParseValueByDefinedWithPanic(pkeyPrefix+" dir", item.Dir)
		item.Mask = ParseValueByDefinedWithPanic(pkeyPrefix+" mask", item.Mask)

		aconf.Mkdir[i] = item
	}
	for i, cmd := range aconf.Exec {
		pkey := "action exec " + strconv.Itoa(i) + " cmd "
		cmd = ParseValueByDefinedWithPanic(pkey, cmd)

		aconf.Exec[i] = cmd
	}
}

func GetActionConf() *actionConf {
	return aconf
}
