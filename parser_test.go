package gitprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestParseValues(t *testing.T) {
	tests := []struct {
		name     string
		setup    string
		expected *GitStatus
	}{
		{
			name:     "not git repo",
			expected: nil,
		},
		{
			name: "dirty",
			setup: `
				git init
				touch test
			`,
			expected: &GitStatus{
				Untracked: 1,
			},
		},
		{
			name: "staged",
			setup: `
				git init
				touch test
				git add test
			`,
			expected: &GitStatus{
				Staged: 1,
			},
		},
		{
			name: "modified",
			setup: `
				git init
				echo "hello" >> test
				git add test
				git commit -m 'initial'
				echo "world" >> test
			`,
			expected: &GitStatus{
				Modified: 1,
			},
		},
		{
			name: "deleted",
			setup: `
				git init
				echo "hello" >> test
				git add test
				git commit -m 'initial'
				rm test
			`,
			expected: &GitStatus{
				Modified: 1,
			},
		},
		{
			name: "conflicts",
			setup: `
				git init
				git commit --allow-empty -m 'initial'
				git checkout -b other
				git checkout master
				echo foo >> test
				git add test
				git commit -m 'first'
				git checkout other
				echo bar >> test
				git add test
				git commit -m 'first'
				git rebase master || true
			`,
			expected: &GitStatus{
				Conflicts: 1,
			},
		},
		{
			name: "ahead",
			setup: `
				git init
				git remote add origin $REMOTE
				git commit --allow-empty -m 'first'
				git push -u origin HEAD
				git commit --allow-empty -m 'second'
			`,
			expected: &GitStatus{
				Ahead: 1,
			},
		},
		{
			name: "behind",
			setup: `
				git init
				git remote add origin $REMOTE
				git commit --allow-empty -m 'first'
				git commit --allow-empty -m 'second'
				git push -u origin HEAD
				git reset --hard HEAD^
			`,
			expected: &GitStatus{
				Behind: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir, cleanupDir := setupTestDir(t)
			defer cleanupDir()

			if test.setup != "" {
				remote, cleanupRemote := setupRemote(t, dir)
				defer cleanupRemote()
				commands := "export REMOTE=" + remote + "\n" + test.setup
				setupCommands(t, dir, commands)
			}

			actual, err := Parse()
			if err != nil {
				t.Errorf("Received unexpected error: %v", err)
				return
			}
			if test.expected == nil {
				if actual != nil {
					t.Errorf("Expected nil return, got %v", actual)
				}
				return
			}
			assertInt(t, "Untracked", test.expected.Untracked, actual.Untracked)
			assertInt(t, "Modified", test.expected.Modified, actual.Modified)
			assertInt(t, "Staged", test.expected.Staged, actual.Staged)
			assertInt(t, "Conflicts", test.expected.Conflicts, actual.Conflicts)
			assertInt(t, "Ahead", test.expected.Ahead, actual.Ahead)
			assertInt(t, "Behind", test.expected.Behind, actual.Behind)
		})
	}
}

func TestParseHead(t *testing.T) {
	dir, done := setupTestDir(t)
	defer done()

	setupCommands(t, dir, `
		git init
	`)
	s, _ := Parse()
	assertString(t, "branch", "master", s.Branch)

	setupCommands(t, dir, `
		git commit --allow-empty -m 'initial'
	`)
	s, _ = Parse()
	assertString(t, "branch", "master", s.Branch)

	setupCommands(t, dir, `
		git checkout -b other
	`)
	s, _ = Parse()
	assertString(t, "branch", "other", s.Branch)

	setupCommands(t, dir, `
		git commit --allow-empty -m 'second'
		git checkout HEAD^
	`)
	s, _ = Parse()
	assertString(t, "branch", "", s.Branch)
	if len(s.Sha) != 40 {
		t.Errorf("Expected 40 char hash, got %v (%s)", len(s.Sha), s.Sha)
	}
}

func TestExecGitErr(t *testing.T) {
	path := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", path)

	_, err := Parse()
	if err == nil {
		t.Errorf("Expected error when git not found on $PATH")
	}
}

func setupTestDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "gitprompt-test")
	if err != nil {
		t.Fatalf("Create temp dir: %v", err)
	}

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("Could not change dir: %v", err)
	}

	return dir, func() {
		if err = os.RemoveAll(dir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clean up test dir: %v\n", err)
		}
	}
}

func setupRemote(t *testing.T, dir string) (string, func()) {
	// Set up remote dir
	remote, err := ioutil.TempDir("", "gitprompt-test-remote.git")
	if err != nil {
		t.Fatalf("create temp dir for remote: %v", err)
	}
	remoteCmd := exec.Command("git", "init", "--bare")
	remoteCmd.Dir = remote
	if err = remoteCmd.Run(); err != nil {
		t.Fatalf("Set up remote: %v", err)
	}

	return remote, func() {
		if err = os.RemoveAll(remote); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clean up remote dir: %v\n", err)
		}
	}
}

func setupCommands(t *testing.T, dir, commands string) {
	commands = "set -eo pipefail\n" + commands
	script := path.Join(dir, "setup.sh")
	if err := ioutil.WriteFile(script, []byte(commands), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", "setup.sh")
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = bufio.NewWriter(&stderr)
	cmd.Stdout = bufio.NewWriter(&stdout)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Log("STDOUT:", stdout.String())
		t.Log("STDERR:", stderr.String())
		t.Fatalf("Setup command failed: %v", err)
	}
	if err := os.Remove(script); err != nil {
		t.Fatal(err)
	}
}

func assertString(t *testing.T, name, expected, actual string) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("%s does not match\n\tExpected: %q\n\tActual:   %q", name, expected, actual)
}

func assertInt(t *testing.T, name string, expected, actual int) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("%s does not match\n\tExpected: %v\n\tActual:   %v", name, expected, actual)
}
