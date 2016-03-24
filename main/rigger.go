package main

import (
	"andals/gobox/misc"
	"andals/gobox/shell"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"rigger/pkg/rconf"
	"strings"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	var rconfDir string

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&rconfDir, "rconfDir", "", "rigger conf dir")
	fs.Parse(os.Args[1:])

	if rconfDir == "" {
		panic("must have flag rconfDir")
	}

	args := parseArgs(fs.Args())
	rconf.Parse(rconfDir, args)

	genConf()
}

func parseArgs(args []string) map[string]string {
	result := make(map[string]string)

	for _, str := range args {
		item := strings.Split(str, "=")
		if len(item) == 2 {
			result[item[0]] = item[1]
		}
	}

	return result
}

func genConf() {
	for key, item := range rconf.GetTplConf() {
		if !misc.FileExist(item.Tpl) {
			panic("Gen conf " + key + " tpl " + item.Tpl + " not exists")
		}
		tplBytes, _ := ioutil.ReadFile(item.Tpl)
		dstString := rconf.ParseValueByDefinedWithPanic(key+" tpl ", string(tplBytes))
		ioutil.WriteFile(item.Dst, []byte(dstString), 0644)

		var cmd string
		var cmdPrefix string
		if item.Sudo {
			cmdPrefix += "sudo "
		}
		cmd += cmdPrefix + "rm -f " + item.Ln + "; "
		cmd += cmdPrefix + "ln -s " + item.Dst + " " + item.Ln
		shell.RunCmdBindTerminal(cmd)
	}
}
