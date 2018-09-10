package gitprompt

import (
	"fmt"
	"os"
)

// GitStatus is the parsed status for the current state in git.
type GitStatus struct {
	Sha       string
	Branch    string
	Untracked int
	Modified  int
	Staged    int
	Conflicts int
	Ahead     int
	Behind    int
}

// Exec executes gitprompt. It first parses the git status, then outputs the
// data according to the format.
// Exits with a non-zero exit code in case git returned an error. Exits with a
// blank string if the current directory is not part of a git repository.
func Exec(format string, printZSH bool) {
	s, err := Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, num := Print(s, format)
	fmt.Fprint(os.Stdout, out)
	if printZSH {
		fmt.Fprintf(os.Stdout, "%%%dG", num)
	}
}
