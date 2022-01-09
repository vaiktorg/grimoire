package documentstore

import (
	"fmt"
	"io"
	"strings"
)

func (d *Dir) Print(lvl int, w io.Writer) {
	printElem := func(times int, name string) {
		_, _ = fmt.Fprintf(w, "%s"+"%s\n", strings.Repeat(" ", times), name)
	}

	indent := 4
	idx := 0
	printElem(lvl, d.Meta.Name+"/")

	for _, f := range d.Files {
		printElem(lvl+indent, "- "+f.Meta.Name)
		idx++
	}

	for _, d := range d.Dirs {
		d.Print(lvl+indent, w)
	}
}
