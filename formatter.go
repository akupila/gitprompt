package gitprompt

import (
	"bytes"
	"strconv"
	"strings"
)

type formatter struct {
	color        uint8
	currentColor uint8
	attr         uint8
	currentAttr  uint8
}

func (f *formatter) setColor(c uint8) {
	f.color = c
}

func (f *formatter) clearColor() {
	f.color = 0
}

func (f *formatter) setAttribute(a uint8) {
	f.attr |= (1 << a)
}

func (f *formatter) clearAttribute(a uint8) {
	f.attr &= ^(1 << a)
}

func (f *formatter) attributeSet(a uint8) bool {
	return (f.attr & (1 << a)) != 0
}

func (f *formatter) clearAttributes() {
	f.attr = 0
}

func (f *formatter) printANSI(b *bytes.Buffer) {
	if f.color == f.currentColor && f.attr == f.currentAttr {
		return
	}
	b.WriteString("\x1b[")
	if f.color == 0 && f.attr == 0 {
		// reset all
		b.WriteString("0m")
		f.currentColor = 0
		f.currentAttr = 0
		return
	}
	mm := []string{}
	aAdded, aRemoved := attrDiff(f.currentAttr, f.attr)
	if len(aRemoved) > 0 {
		mm = append(mm, "0")
		var i uint8 = 1
		for ; i < 8; i++ {
			if f.attributeSet(i) {
				mm = append(mm, strconv.Itoa(int(i)))
			}
		}
	} else if len(aAdded) > 0 {
		for _, a := range aAdded {
			mm = append(mm, strconv.Itoa(int(a)))
		}
	}
	if f.color != f.currentColor || len(aRemoved) > 0 {
		mm = append(mm, strconv.Itoa(int(f.color)))
	}
	b.WriteString(strings.Join(mm, ";"))
	b.WriteString("m")
	f.currentColor = f.color
	f.currentAttr = f.attr

}

func attrDiff(a, b uint8) ([]uint8, []uint8) {
	added := []uint8{}
	removed := []uint8{}
	var i uint8
	for ; i < 8; i++ {
		inA := (a & (1 << i)) != 0
		inB := (b & (1 << i)) != 0
		if inA && inB {
			continue
		}
		if inB {
			added = append(added, i)
			continue
		}
		if inA {
			removed = append(removed, i)
		}
	}
	return added, removed
}
