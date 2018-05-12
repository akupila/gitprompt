package gitprompt

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// Parse parses the status for the repository from git. Returns nil if the
// current directory is not part of a git repository.
func Parse() (*GitStatus, error) {
	status := &GitStatus{}

	stat, err := runGitCommand("git", "status", "--branch", "--porcelain=2")
	if err != nil {
		if strings.HasPrefix(err.Error(), "fatal:") {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(stat, "\n")
	for _, line := range lines {
		switch line[0] {
		case '#':
			parseHeader(line, status)
		case '?':
			status.Untracked++
		case 'u':
			status.Conflicts++
		case '1', '2':
			parts := strings.Split(line, " ")
			if parts[1][0] != '.' {
				status.Staged++
			}
			if parts[1][1] != '.' {
				status.Modified++
			}
		}
	}

	return status, nil
}

func parseHeader(h string, s *GitStatus) {
	if strings.HasPrefix(h, "# branch.oid") {
		hash := h[13:]
		if hash != "(initial)" {
			s.Sha = hash
		}
		return
	}
	if strings.HasPrefix(h, "# branch.head") {
		branch := h[14:]
		if branch != "(detached)" {
			s.Branch = branch
		}
		return
	}
	if strings.HasPrefix(h, "# branch.ab") {
		parts := strings.Split(h, " ")
		s.Ahead, _ = strconv.Atoi(strings.TrimPrefix(parts[2], "+"))
		s.Behind, _ = strconv.Atoi(strings.TrimPrefix(parts[3], "-"))
		return
	}
}

func runGitCommand(cmd string, args ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.Command(cmd, args...)
	command.Stdout = bufio.NewWriter(&stdout)
	command.Stderr = bufio.NewWriter(&stderr)
	if err := command.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(stderr.String())
		}
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
