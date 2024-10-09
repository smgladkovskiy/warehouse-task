package now

import "time"

type WithNowGeneratorable interface {
	GetNowGen() Generatorable
	SetNowGen(nowFunc Generatorable)
}

type WithNowGenerator struct {
	nowFunc Generatorable
}

var generator = defaultGenerator{}

func (wng *WithNowGenerator) Now() time.Time {
	if wng.nowFunc == nil {
		wng.nowFunc = generator
	}

	return wng.nowFunc.Now()
}

func (wng *WithNowGenerator) NowP() *time.Time {
	if wng.nowFunc == nil {
		wng.nowFunc = generator
	}

	return wng.nowFunc.NowP()
}

func (wng *WithNowGenerator) SetNowGen(nowFunc Generatorable) {
	if nowFunc == nil {
		nowFunc = generator
	}

	wng.nowFunc = nowFunc
}

func (wng *WithNowGenerator) GetNowGen() Generatorable {
	if wng.nowFunc == nil {
		wng.nowFunc = generator
	}

	return wng.nowFunc
}
