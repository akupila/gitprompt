# gitprompt

[![CircleCI](https://circleci.com/gh/akupila/gitprompt.svg?style=svg)](https://circleci.com/gh/akupila/gitprompt)
[![goreportcard](https://goreportcard.com/badge/github.com/akupila/gitprompt)](https://goreportcard.com/report/github.com/akupila/gitprompt)
[![codecov](https://codecov.io/gh/akupila/gitprompt/branch/master/graphs/badge.svg?branch=master)](https://codecov.io/gh/akupila/gitprompt)

gitprompt is a configurable, fast and zero-dependencies* way of getting the
current git status to be displayed in the `PROMPT`.

Displays:

- Current branch / sha1
- Untracked files
- Modified files
- Staged files
- Commits behind / ahead remote

When executed, gitprompt gets the git status of the current working directory
then prints it according to the format specified. If the current working
directory is not part of a git repository, gitprompt
exits with code `0` and no output.

`*` git is required

## Configuration

The output is configured with `-format` or the `GITPROMPT_FORMAT` environment
variable. If both are set, the flag takes precedence.

A very bare-bones format can look like this:

```
gitprompt -format="%h >%s ↓%b ↑%a @%c +%m %u "

# or

export GITPROMPT_FORMAT="%h >%s ↓%b ↑%a @%c +%m %u "
gitprompt
```

Characters that don't have a special meaning are printed as usual _(unicode
characters are fine, go crazy with emojis if that's your thing)_.

### Data

Various data from git can be displayed in the output. Data tokens are prefixed
with `%`:

| token | explanation                       |
|-------|-----------------------------------|
| `%h`  | Current branch or sha1            |
| `%s`  | Number of files staged            |
| `%b`  | Number of commits behind remote   |
| `%a`  | Number of commits ahead of remote |
| `%c`  | Number of conflicts               |
| `%m`  | Number of files modified          |
| `%u`  | Number of untracked files         |

Normally `%h` displays the current branch (`master`) but if you're detached
from `HEAD`, it will display the current sha1. Only first 7 characters of the
sha1 are displayed.

### Colors

The color can be set with color tokens, prefixed with `#`:

| token | color             |
|-------|-------------------|
| `#k`  | Black             |
| `#r`  | Red               |
| `#g`  | Green             |
| `#y`  | Yellow            |
| `#b`  | Blue              |
| `#m`  | Magenta           |
| `#c`  | Cyan              |
| `#w`  | White             |
| `#K`  | Highlight Black   |
| `#R`  | Highlight Red     |
| `#G`  | Highlight Green   |
| `#Y`  | Highlight Yellow  |
| `#B`  | Highlight Blue    |
| `#M`  | Highlight Magenta |
| `#C`  | Highlight Cyan    |
| `#W`  | Highlight White   |

The color is set until another color overrides it, or a group ends (see below).
If a color was set when gitprompt is done, it will add a color reset escape
code at the end, meaning text after gitprompt won't have the color applied.

### Text attributes

The text attributes can be set with attribute tokens, prefixed with `@`:

| token | attribute             |
|-------|-----------------------|
| `@b`  | Set bold              |
| `@B`  | Clear bold            |
| `@f`  | Set faint/dim color   |
| `@F`  | Clear faint/dim color |
| `@i`  | Set italic            |
| `@I`  | Clear italic          |

As with colors, if an attribute was set when gitprompt is done, an additional
escape code is automatically added to clear it.

### Groups

Groups can be used for adding logic to the format. A group's output is only
printed if at least one item in the group has data.

| token | action      |
|-------|-------------|
| `[`   | start group |
| `]`   | end group   |

Consider the following:

```
%h[ %a][ %m]
```

With `head=master, ahead=3, modified=0`, this will print `master 3` since there
were not modified files. Note the space that's included in the group, if
spacing should be added when the group is present, the spacing should be added
to the group itself.

Colors and attributes are also scoped to a group, meaning they won't leak
outside so there's no need to reset colors. This prints `behind` in red, `-`
without any formatting, and `ahead` in green:

```
[#r behind: %b] - [#g ahead: %a]
```

A group requires one data token to have a non-zero value. The following prints
the current branch/sha1 in cyan, then number of staged files (if not zero),
then commits behind and ahead (if both are not zero). This allows for symmetry,
if it's desired to show `master >1 ↓0 ↑2` instead of `master >1 ↑2`:

```
#c%h[ >%s][ ↓%b ↑%a]
```

### Complete example

Putting everything together, a complex format may look something like this:

```
gitprompt -format="#B(@b#R%h[#y >%s][#m ↓%b ↑%a][#r x%c][#g +%m][#y ϟ%u]#B) "
```

- `(` in highlight blue
- Current branch/sha1 in bold red
- If there are staged files, number of staged files in yellow, prefixed with `>`
- If there are commits ahead or behind, show them with arrows in magenta
- If there are conflicts, show them in red, prefixed with `x`
- If files have been modified since the previous commit, show `+3` for 3 modified files
- If files are untracked (added since last commit), show a lightning and the number in yellow
- `)` in highlight blue

> Any text printed after gitprompt will have all formatting cleared

## Installation

Installation consists of two parts: get the binary & configure your shell to
use it.

### Get binary

#### Homebrew

Ensure you have [Homebrew] installed.

```
$ brew tap akupila/gitprompt
$ brew install gitprompt
```

#### Install binary directly

```
$ curl -sfL https://install.goreleaser.com/github.com/akupila/gitprompt.sh | bash -s -- -b /usr/local/bin
```

Feel free to change the path from `/usr/local/bin`, just make sure `gitprompt`
is available on your `$PATH` (check with `gitprompt -version`).

#### Build from source

Install the go toolchain: https://golang.org/doc/install

```
$ go get github.com/akupila/gitprompt/...
```

The code has no vendored dependencies so no need to worry about that.

### Configure your shell

#### zsh

Execute `gitprompt` as part of `PROMPT`. Add this to your  `~/.zshrc`:

```
export PROMPT='$PROMPT $(gitprompt)'
```

Now reload the config (`source ~/.zshrc`) and gitprompt should show up. Feel
free to add anything else here too, just execute `gitprompt` where you want the
status, for example:

```
local ret_status="%(?:%{$fg_bold[green]%}$:%{$fg_bold[red]%}$)"
local dir="%{$fg[cyan]%}%3d"
export PROMPT='${ret_status} ${dir} $(gitprompt)'
```

Alternatively, you can add this to `RPROMPT` instead, which will make the
status appear on the right hand side of the screen. `gitprompt` will by default
add a trailing space so you you may want to customize the formatting if you
don't want a trailing space here.

#### bash

Set `PS1` in your `~/.bashrc` and reload the config (`source ~/.bashrc`).

For example:

```
export PS1='$PS1 $(gitprompt)'
```

See [bashrcgenerator] for more, just add `$(gitprompt)` where you want the git
status to appear.

### Uninstallation

1. Remove `gitprompt` from your shell config
2. Delete the binary `rm $(which gitprompt)` or `brew uninstall gitprompt`



[Homebrew]: https://brew.sh/
[bashrcgenerator]: http://bashrcgenerator.com/
