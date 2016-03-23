package main

import (
	"flag"
	"fmt"
	"os"
	"rigger/pkg/commands"
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

	if len(os.Args) < 2 {
		panic("Usage " + os.Args[0] + " cmd flags")
	}

	cmd := commands.GetCommand(os.Args[1])

	var rconfDir string

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&rconfDir, "rconfDir", "", "rigger conf dir")
	cmd.SetFlags(fs)
	fs.Parse(os.Args[2:])

	if rconfDir == "" {
		panic("must have flag rconfDir")
	}
	cmd.VerifyFlags()

	args := parseArgs(fs.Args())
	rconf.Init(rconfDir, args)
	cmd.Run()
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
