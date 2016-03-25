package rconf

import (
	"andals/gobox/misc"
	"andals/gobox/shell"
	"encoding/json"
	//     "fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type tplItem struct {
	Tpl  string
	Dst  string
	Ln   string
	Sudo bool
}

type mkdirItem struct {
	Dir  string
	Mask string
	Sudo bool
}

type actionConf struct {
	Mkdir []*mkdirItem
	Exec  []string
}

type RiggerConf struct {
	varConf map[string]string
	tplConf map[string]*tplItem
	aconf   *actionConf

	user           string
	rconfDir       string
	extArgs        map[string]string
	unparsedValues map[string]string
}

func NewRiggerConf(rconfDir string, extArgs map[string]string) *RiggerConf {
	rconf := new(RiggerConf)

	rconf.varConf = make(map[string]string)
	rconf.tplConf = make(map[string]*tplItem)
	rconf.aconf = new(actionConf)

	rconf.user = os.Getenv("USER")
	rconf.rconfDir = rconfDir
	rconf.extArgs = extArgs
	rconf.unparsedValues = make(map[string]string)

	return rconf
}

func (this *RiggerConf) Parse() {
	this.parseVarConf()
	this.parseTplConf()
	this.parseActionConf()
}

/**
* return parsed value and whether delay parsed, if delay parsed, bool is true
 */
func (this *RiggerConf) ParseValueByDefined(key string, value string) (string, bool) {
	re := regexp.MustCompile("\\${([^}]+)}")
	matches := re.FindAllStringSubmatch(value, -1)

	if len(matches) == 0 {
		return value, false
	}

	var rs []string
	for _, item := range matches {
		k := item[1]
		v, ok := this.varConf[k]
		if !ok {
			if k != key { //eg: USER:${USER} get by env
				_, ok = this.unparsedValues[k]
				if ok {
					return value, true
				}
			}
			v = shell.RunCmd("echo $" + k).Output
			v = strings.TrimSpace(v)
			if v == "" {
				panic("Parse " + key + " has undefined field " + k)
			}
		}
		rs = append(rs, item[0])
		rs = append(rs, v)
	}

	return strings.NewReplacer(rs...).Replace(value), false
}

func (this *RiggerConf) GetTplConf() map[string]*tplItem {
	return this.tplConf
}

func (this *RiggerConf) GetActionConf() *actionConf {
	return this.aconf
}

func (this *RiggerConf) parseVarConf() {
	path := this.rconfDir + "/var.json"
	if !misc.FileExist(path) {
		panic("Var conf not exists in " + path)
	}

	var varJsonConf map[string]interface{}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, &varJsonConf)
	if nil != err {
		panic("Parse var conf error: " + err.Error())
	}

	for key, item := range varJsonConf {
		this.unparsedValues[key] = this.parseVarJsonItemtoString(key, item)
	}
	for len(this.unparsedValues) > 0 {
		for key, value := range this.unparsedValues {
			value, delay := this.ParseValueByDefined(key, value)
			if !delay {
				this.varConf[key] = this.parseValueByFunc(key, value)
				delete(this.unparsedValues, key)
			}
		}
	}
}

func (this *RiggerConf) parseVarJsonItemtoString(panicKey string, item interface{}) string {
	var r string

	switch item.(type) {
	case string:
		r = item.(string)
	case map[string]interface{}:
		mv := item.(map[string]interface{})
		v, ok := mv[this.user]
		if !ok {
			v = mv["default"]
		}
		r = v.(string)
	default:
		panic("Key " + panicKey + " item's type not support")
	}

	return strings.TrimSpace(r)
}

func (this *RiggerConf) parseValueByFunc(panicKey string, value string) string {
	re := regexp.MustCompile("__([A-Z]+)__\\(([^)]+)\\)")
	match := re.FindStringSubmatch(value)

	if len(match) != 0 {
		switch match[1] {
		case "ARG":
			value = this.parseByArgFunc(panicKey, match[2])
		case "MATH":
			value = this.parseByMathFunc(panicKey, match[2])
		default:
			panic("Not support func " + match[1] + " in " + panicKey)
		}
	}

	return value
}

func (this *RiggerConf) parseByArgFunc(panicKey string, argName string) string {
	v, ok := this.extArgs[argName]
	if !ok {
		panic("Not has arg " + argName + " for " + panicKey)
	}

	return v
}

func (this *RiggerConf) parseByMathFunc(panicKey string, express string) string {
	re := regexp.MustCompile("([0-9]+)([+\\-*/])([0-9]+)")
	match := re.FindStringSubmatch(express)

	if len(match) == 0 {
		panic("Invalid match express " + express + " for " + panicKey)
	}

	lv, _ := strconv.Atoi(match[1])
	rv, _ := strconv.Atoi(match[3])
	var value int

	switch match[2] {
	case "+":
		value = lv + rv
	case "-":
		value = lv - rv
	case "*":
		value = lv * rv
	case "/":
		value = lv / rv
	}

	return strconv.Itoa(value)
}

func (this *RiggerConf) parseTplConf() {
	path := this.rconfDir + "/tpl.json"
	if !misc.FileExist(path) {
		panic("Tpl conf not exists in " + path)
	}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, &this.tplConf)
	if nil != err {
		panic("Parse tpl conf error: " + err.Error())
	}

	for key, item := range this.tplConf {
		pkeyPrefix := "tpl " + key + " item "
		item.Tpl, _ = this.ParseValueByDefined(pkeyPrefix+" tpl", item.Tpl)
		item.Dst, _ = this.ParseValueByDefined(pkeyPrefix+" dst", item.Dst)
		item.Ln, _ = this.ParseValueByDefined(pkeyPrefix+" ln", item.Ln)
	}
}

func (this *RiggerConf) parseActionConf() {
	path := this.rconfDir + "/action.json"
	if !misc.FileExist(path) {
		panic("Action conf not exists in " + path)
	}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, this.aconf)
	if nil != err {
		panic("Parse action conf error: " + err.Error())
	}

	for i, item := range this.aconf.Mkdir {
		pkeyPrefix := "action mkdir " + strconv.Itoa(i) + " item "
		item.Dir, _ = this.ParseValueByDefined(pkeyPrefix+" dir", item.Dir)
		item.Mask, _ = this.ParseValueByDefined(pkeyPrefix+" mask", item.Mask)
	}
	for i, cmd := range this.aconf.Exec {
		pkey := "action exec " + strconv.Itoa(i) + " cmd "
		this.aconf.Exec[i], _ = this.ParseValueByDefined(pkey, cmd)
	}
}
