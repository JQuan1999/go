package ch5

import (
	"fmt"
)

type Base struct{}

func (Base) Magic() {
	fmt.Println("base magic")
}

func (Base) MagicStr() string {
	return "base magic"
}

func (b Base) MoreMagic() {
	b.Magic()
	b.Magic()
}

type Voodoo struct {
	Base
}

func (Voodoo) Magic() {
	fmt.Println("voodoo magic")
}

func (Voodoo) MagicStr() string {
	return "voodoo magic"
}
