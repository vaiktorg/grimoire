package bitflag

import "fmt"

// BitFlag Largest uint32 4294967295
type BitFlag struct {
	bit    int
	bitmap []string
}

func (b *BitFlag) AddBitSet(bitmap []string) {
	b.bitmap = bitmap
}

func (b *BitFlag) Set(flag int)      { b.bit = b.bit | flag }
func (b *BitFlag) Clear(flag int)    { b.bit = b.bit &^ flag }
func (b *BitFlag) Toggle(flag int)   { b.bit = b.bit ^ flag }
func (b *BitFlag) Has(flag int) bool { return b.bit&flag != 0 }
func (b *BitFlag) Val() int          { return b.bit }

func (b *BitFlag) String() string {
	return fmt.Sprintf("%T", b)
}
