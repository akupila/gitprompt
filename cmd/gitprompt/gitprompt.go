package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akupila/gitprompt"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	goversion = "unknown"
)

const defaultFormat = "#B([@b#R%h][#y ›%s][#m ↓%b][#m ↑%a][#r x%c][#g +%m][#y %u]#B) "

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

var exampleStatus = &gitprompt.GitStatus{
	Branch:    "master",
	Sha:       "0455b83f923a40f0b485665c44aa068bc25029f5",
	Untracked: 1,
	Modified:  2,
	Staged:    3,
	Conflicts: 4,
	Ahead:     5,
	Behind:    6,
}

var formatHelp = func() string {
	example, _ := gitprompt.Print(exampleStatus, defaultFormat)
	return fmt.Sprintf(`Define output format.

Default format is: %q
Example result:    %s

Data:
	%%h	Current branch or SHA1
	%%s	Number of files staged
	%%b	Number of commits behind remote
	%%a	Number of commits ahead of remote
	%%c	Number of conflicts
	%%m	Number of files modified
	%%u	Number of untracked files

Colors:
	#k	Black
	#r	Red
	#g	Green
	#y	Yellow
	#b	Blue
	#m	Magenta
	#c	Cyan
	#w	White
	#K	Highlight Black
	#R	Highlight Red
	#G	Highlight Green
	#Y	Highlight Yellow
	#B	Highlight Blue
	#M	Highlight Magenta
	#C	Highlight Cyan
	#W	Highlight White

Text attributes:
	@b	Set bold
	@B	Clear bold
	@f	Set faint/dim color
	@F	Clear faint/dim color
	@i	Set italic
	@I	Clear italic`, defaultFormat, example)
}

func main() {
	v := flag.Bool("version", false, "Print version inforformation.")
	zsh := flag.Bool("zsh", false, "Print zsh width control characters")
	flag.Var(&format, "format", formatHelp())
	flag.Parse()

	if *v {
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Build date: %s\n", date)
		fmt.Printf("Go version: %s\n", goversion)
		os.Exit(0)
	}

	gitprompt.Exec(format.String(), *zsh)
}
