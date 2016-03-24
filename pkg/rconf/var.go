package rconf

import (
	"andals/gobox/misc"
	"encoding/json"
	//     "fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var varConf map[string]string

func init() {
	varConf = make(map[string]string)
}

func parseVarConf(path string) {
	if !misc.FileExist(path) {
		panic("Var conf not exists in " + path)
	}

	var varJsonConf map[string]interface{}

	jsonBytes, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonBytes, &varJsonConf)
	if nil != err {
		panic("Parse var conf error: " + err.Error())
	}

	// Run first loop since map is out-of-order so that read undefined var

	waitValues := make(map[string]string)
	for key, item := range varJsonConf {
		str, succ := toString(item)
		if !succ {
			panic("Key " + key + " item's type not support")
		}
		waitValues[key] = str
	}
	for len(waitValues) > 0 {
		for key, value := range waitValues {
			value, undefined := parseValueByDefined(value)
			if undefined == "" {
				varConf[key] = parseValueByFunc(key, value)
				delete(waitValues, key)
			} else {
				_, ok := waitValues[undefined]
				if !ok {
					panic("Var " + key + " item " + undefined + " is undefined")
				}
			}
		}
	}
}

func toString(item interface{}) (string, bool) {
	var r string
	succ := true

	switch item.(type) {
	case string:
		r = item.(string)
	case map[string]interface{}:
		mv := item.(map[string]interface{})
		k := os.Getenv("USER")
		v, ok := mv[k]
		if !ok {
			v = mv["default"]
		}
		r = v.(string)
	default:
		succ = false
	}

	return strings.TrimSpace(r), succ
}

func parseValueByFunc(key string, value string) string {
	re := regexp.MustCompile("__([A-Z]+)__\\(([^)]+)\\)")
	match := re.FindStringSubmatch(value)

	if len(match) != 0 {
		switch match[1] {
		case "ARG":
			value = parseByArgFunc(key, match[2])
		case "MATH":
			value = parseByMathFunc(key, match[2])
		}
	}

	return value
}

func parseByArgFunc(key string, argName string) string {
	v, ok := extArgs[argName]
	if !ok {
		panic("Not has arg " + argName + " for " + key)
	}

	return v
}

func parseByMathFunc(key string, express string) string {
	re := regexp.MustCompile("([0-9]+)([+\\-*/])([0-9]+)")
	match := re.FindStringSubmatch(express)

	if len(match) == 0 {
		panic("Invalid match express " + express + " for " + key)
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
