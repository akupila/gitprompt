package main

import (
	"fmt"
	"strings"

	"github.com/akupila/gitprompt"
)

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
	return strings.TrimSpace(fmt.Sprintf(`
How to format output

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
	@I	Clear italic
`, defaultFormat, gitprompt.Print(exampleStatus, defaultFormat)))
}
