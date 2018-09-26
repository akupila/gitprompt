package gitprompt

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const (
	tAttribute rune = '@'
	tColor          = '#'
	tReset          = '_'
	tData           = '%'
	tGroupOp        = '['
	tGroupCl        = ']'
	tEsc            = '\\'
)

var attrs = map[rune]uint8{
	'b': 1, // bold
	'f': 2, // faint
	'i': 3, // italic
}

var resetAttrs = map[rune]uint8{
	'B': 1, // bold
	'F': 2, // faint
	'I': 3, // italic
}

var colors = map[rune]uint8{
	'k': 30, // black
	'r': 31, // red
	'g': 32, // green
	'y': 33, // yellow
	'b': 34, // blue
	'm': 35, // magenta
	'c': 36, // cyan
	'w': 37, // white
	'K': 90, // highlight black
	'R': 91, // highlight red
	'G': 92, // highlight green
	'Y': 93, // highlight yellow
	'B': 94, // highlight blue
	'M': 95, // highlight magenta
	'C': 96, // highlight cyan
	'W': 97, // highlight white
}

const (
	head      rune = 'h'
	untracked      = 'u'
	modified       = 'm'
	staged         = 's'
	conflicts      = 'c'
	ahead          = 'a'
	behind         = 'b'
)

type group struct {
	buf bytes.Buffer

	parent *group
	format formatter

	hasData  bool
	hasValue bool
	width    int
}

// Print prints the status according to the format.
//
// The integer returned is the print width of the string.
func Print(s *GitStatus, format string) (string, int) {
	if s == nil {
		return "", 0
	}

	in := make(chan rune)
	go func() {
		r := bufio.NewReader(strings.NewReader(format))
		for {
			ch, _, err := r.ReadRune()
			if err != nil {
				close(in)
				return
			}
			in <- ch
		}
	}()

	return buildOutput(s, in)
}

func buildOutput(s *GitStatus, in chan rune) (string, int) {
	root := &group{}
	g := root

	col := false
	att := false
	dat := false
	esc := false

	for ch := range in {
		if esc {
			esc = false
			g.addRune(ch)
			continue
		}

		if col {
			setColor(g, ch)
			col = false
			continue
		}

		if att {
			setAttribute(g, ch)
			att = false
			continue
		}

		if dat {
			setData(g, s, ch)
			dat = false
			continue
		}

		switch ch {
		case tEsc:
			esc = true
		case tColor:
			col = true
		case tAttribute:
			att = true
		case tData:
			dat = true
		case tGroupOp:
			g = &group{
				parent: g,
				format: g.format,
			}
			g.format.clearAttributes()
			g.format.clearColor()
		case tGroupCl:
			if g.writeTo(&g.parent.buf) {
				g.parent.format = g.format
				g.parent.format.setColor(0)
				g.parent.format.clearAttributes()
				g.parent.width += g.width
			}
			g = g.parent
		default:
			g.addRune(ch)
		}
	}

	// trailing characters
	if col {
		g.addRune(tColor)
	}
	if att {
		g.addRune(tAttribute)
	}
	if dat {
		g.addRune(tData)
	}

	g.format.clearColor()
	g.format.clearAttributes()
	g.format.printANSI(&g.buf)

	return root.buf.String(), root.width
}

func setColor(g *group, ch rune) {
	if ch == tReset {
		// Reset color code.
		g.format.clearColor()
		return
	}
	code, ok := colors[ch]
	if ok {
		g.format.setColor(code)
		return
	}
	g.addRune(tColor)
	g.addRune(ch)
}

func setAttribute(g *group, ch rune) {
	if ch == tReset {
		// Reset attribute.
		g.format.clearAttributes()
		return
	}
	code, ok := attrs[ch]
	if ok {
		g.format.setAttribute(code)
		return
	}
	code, ok = resetAttrs[ch]
	if ok {
		g.format.clearAttribute(code)
		return
	}
	g.addRune(tAttribute)
	g.addRune(ch)
}

func setData(g *group, s *GitStatus, ch rune) {
	switch ch {
	case head:
		g.hasData = true
		g.hasValue = true
		if s.Branch != "" {
			g.addString(s.Branch)
		} else {
			g.addString(s.Sha[:7])
		}
	case modified:
		g.addInt(s.Modified)
		g.hasData = true
		if s.Modified > 0 {
			g.hasValue = true
		}
	case untracked:
		g.addInt(s.Untracked)
		g.hasData = true
		if s.Untracked > 0 {
			g.hasValue = true
		}
	case staged:
		g.addInt(s.Staged)
		g.hasData = true
		if s.Staged > 0 {
			g.hasValue = true
		}
	case conflicts:
		g.addInt(s.Conflicts)
		g.hasData = true
		if s.Conflicts > 0 {
			g.hasValue = true
		}
	case ahead:
		g.addInt(s.Ahead)
		g.hasData = true
		if s.Ahead > 0 {
			g.hasValue = true
		}
	case behind:
		g.addInt(s.Behind)
		g.hasData = true
		if s.Behind > 0 {
			g.hasValue = true
		}
	default:
		g.addRune(tData)
		g.addRune(ch)
	}
}

func (g *group) writeTo(b io.Writer) bool {
	if g.hasData && !g.hasValue {
		return false
	}
	if _, err := g.buf.WriteTo(b); err != nil {
		log.Panic(err)
	}
	return true
}

func (g *group) addRune(r rune) {
	if !unicode.IsSpace(r) {
		g.format.printANSI(&g.buf)
	}
	g.width++
	g.buf.WriteRune(r)
}

func (g *group) addString(s string) {
	g.format.printANSI(&g.buf)
	g.width += len(s)
	g.buf.WriteString(s)
}

func (g *group) addInt(i int) {
	g.addString(strconv.Itoa(i))
}
