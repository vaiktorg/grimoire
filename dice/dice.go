package dice

import (
	"github.com/google/uuid"
	"github.com/gtank/isaac"
)

type Dice struct {
	rnd         isaac.ISAAC
	valtillSeed int
}

func (d *Dice) D4() int      { return d.getDiceVal(4) }
func (d *Dice) D6() int      { return d.getDiceVal(6) }
func (d *Dice) D8() int      { return d.getDiceVal(8) }
func (d *Dice) D10() int     { return d.getDiceVal(10) }
func (d *Dice) D12() int     { return d.getDiceVal(12) }
func (d *Dice) D20() int     { return d.getDiceVal(20) }
func (d *Dice) DN(n int) int { return d.getDiceVal(n) }

func (d *Dice) getDiceVal(dn int) int {
	if d.valtillSeed == 0 {
		d.seed()
		d.valtillSeed = 256
	}

	return int(d.rnd.Rand())%dn + 1
}

func (d *Dice) seed() {
	d.rnd.Seed(uuid.New().String())
	d.valtillSeed--
}
