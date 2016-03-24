package main

import (
	"flag"
	"fmt"
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
	rconf.Init(rconfDir, args)
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
