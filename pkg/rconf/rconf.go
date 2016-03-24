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

func Parse(rconfDir string, args map[string]string) {
	extArgs = args

	parseVarConf(rconfDir + "/var.json")
	parseTplConf(rconfDir + "/tpl.json")
	parseActionConf(rconfDir + "/action.json")
}

/**
* @return value, undefined field
 */
func parseValueByDefined(value string) (string, string) {
	re := regexp.MustCompile("\\${([^}]+)}")
	matches := re.FindAllStringSubmatch(value, -1)

	if len(matches) != 0 {
		var rs []string
		for _, item := range matches {
			k := item[1]
			v, ok := varConf[k]
			if !ok {
				_, ok = waitValues[k]
				if ok { //in waitValues since value has not been parsed
					return value, k
				}
				v = shell.RunCmd("echo $" + k).Output
				v = strings.TrimSpace(v)
				if v == "" {
					return value, k
				}
			}
			rs = append(rs, item[0])
			rs = append(rs, v)
		}

		value = strings.NewReplacer(rs...).Replace(value)
	}

	return value, ""
}

func ParseValueByDefinedWithPanic(panicKey string, value string) string {
	value, undefined := parseValueByDefined(value)
	if undefined != "" {
		panic("Parse " + panicKey + " has undefined field " + undefined)
	}

	return value
}
