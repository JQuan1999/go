package ch6

import "math"

// 声明一个接口类型 包含一个Area方法
type Shaper interface {
	Area() float32
}

type Square struct {
	Side float32
}

func (sq *Square) Area() float32 {
	return sq.Side * sq.Side
}

func (sq *Square) volumn() float64 {
	return math.Pow(float64(sq.Side), 3)
}

func (sq *Square) Volumn() float64 {
	return sq.volumn()
}
