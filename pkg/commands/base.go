package commands

import (
	"flag"
)

const (
	COMMAND_NAME_CONF = "conf"
)

type ICommand interface {
	SetFlags(fs *flag.FlagSet)
	VerifyFlags()
	Run()
}

func GetCommand(name string) ICommand {
	switch name {
	case COMMAND_NAME_CONF:
		return &confCommand{new(baseCommand)}
	default:
		panic("invalid command " + name)
	}
}

type baseCommand struct {
}

func (this *baseCommand) SetFlags(fs *flag.FlagSet) {
}

func (this *baseCommand) VerifyFlags() {
}
