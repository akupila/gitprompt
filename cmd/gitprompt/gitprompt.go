package main

import (
	"flag"
	"os"

	"github.com/akupila/gitprompt"
)

const defaultFormat = "#B([@b#R%h][#y ›%s][#m ↓%b][#m ↑%a][#r x%c][#g +%m][#y %u]#B)"

type formatFlag struct {
	set   bool
	value string
}

func (f *formatFlag) Set(v string) error {
	f.set = true
	f.value = v
	return nil
}

func (f *formatFlag) String() string {
	if f.set {
		return f.value
	}

	if envVar := os.Getenv("GITPROMPT_FORMAT"); envVar != "" {
		return envVar
	}

	return defaultFormat
}

var format formatFlag

func main() {
	flag.Var(&format, "format", formatHelp())
	flag.Parse()
	gitprompt.Exec(format.String())
}
