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

	genConfByTpl()
	runAction()
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

func genConfByTpl() {
	for key, item := range rconf.GetTplConf() {
		if !misc.FileExist(item.Tpl) {
			panic("Gen conf " + key + " tpl " + item.Tpl + " not exists")
		}
		tplBytes, _ := ioutil.ReadFile(item.Tpl)
		dstString := rconf.ParseValueByDefinedWithPanic(key+" tpl ", string(tplBytes))
		ioutil.WriteFile(item.Dst, []byte(dstString), 0644)

		cmd := ""
		cmdPrefix := ""
		if item.Sudo {
			cmdPrefix += "sudo "
		}
		cmd += cmdPrefix + "rm -f " + item.Ln + "; "
		cmd += cmdPrefix + "ln -s " + item.Dst + " " + item.Ln

		shell.RunCmdBindTerminal(cmd)
	}
}

func runAction() {
	aconf := rconf.GetActionConf()

	for _, item := range aconf.Mkdir {
		cmd := ""
		cmdPrefix := ""
		if item.Sudo {
			cmdPrefix += "sudo "
		}
		if !misc.DirExist(item.Dir) {
			cmd += cmdPrefix + "mkdir -p " + item.Dir + "; "
		}
		cmd += cmdPrefix + "chmod " + item.Mask + " " + item.Dir

		shell.RunCmdBindTerminal(cmd)
	}

	for _, cmd := range aconf.Exec {
		shell.RunCmdBindTerminal(cmd)
	}
}
