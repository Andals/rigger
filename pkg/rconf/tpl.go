package rconf

import (
	"andals/gobox/misc"
	"encoding/json"
	//     "fmt"
	"io/ioutil"
)

type tplItem struct {
	Tpl  string
	Dst  string
	Ln   string
	Sudo bool
}

var tplConf map[string]*tplItem

func init() {
	tplConf = make(map[string]*tplItem)
}

func parseTplConf(path string) {
	if !misc.FileExist(path) {
		panic("Tpl conf not exists in " + path)
	}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, &tplConf)
	if nil != err {
		panic("Parse tpl conf error: " + err.Error())
	}

	for key, item := range tplConf {
		pkeyPrefix := "tpl " + key + " item "
		item.Tpl = ParseValueByDefinedWithPanic(pkeyPrefix+" tpl", item.Tpl)
		item.Dst = ParseValueByDefinedWithPanic(pkeyPrefix+" dst", item.Dst)
		item.Ln = ParseValueByDefinedWithPanic(pkeyPrefix+" ln", item.Ln)

		tplConf[key] = item
	}
}

func GetTplConf() map[string]*tplItem {
	return tplConf
}
