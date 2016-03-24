package rconf

import (
	//     "fmt"
	"andals/gobox/shell"
	"regexp"
	"strings"
)

var extArgs map[string]string

func init() {
	extArgs = make(map[string]string)
}

func Init(rconfDir string, args map[string]string) {
	extArgs = args

	parseVarConf(rconfDir + "/var.json")
	parseTplConf(rconfDir + "/tpl.json")
	parseActionConf(rconfDir + "/action.json")
}

func parseValueByDefined(value string) (string, string) {
	re := regexp.MustCompile("\\${([^}]+)}")
	matches := re.FindAllStringSubmatch(value, -1)

	var undefined string
	if len(matches) != 0 {
		var rs []string
		for _, item := range matches {
			rs = append(rs, item[0])
			k := item[1]
			v, ok := varConf[k]
			if !ok {
				v = shell.RunCmd("echo $" + k).Output
				v = strings.TrimSpace(v)
				if v == "" {
					undefined = k
					break
				}
			}
			rs = append(rs, v)
		}

		if undefined == "" {
			value = strings.NewReplacer(rs...).Replace(value)
		}
	}

	return value, undefined
}

func parseValueByDefinedWithPanic(panicKey string, value string) string {
	value, undefined := parseValueByDefined(value)
	if undefined != "" {
		panic("Parse " + panicKey + " has undefined field " + undefined)
	}

	return value
}
