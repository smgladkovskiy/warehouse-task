package now

import (
	"time"
)

//go:generate mockgen -source=generator.go -destination=./generator_mock.go -package=now -mock_names Generatorable=Mock
type Generatorable interface {
	Now() time.Time
	NowP() *time.Time
}

type defaultGenerator struct{}

func (dng defaultGenerator) Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func (dng defaultGenerator) NowP() *time.Time {
	tn := dng.Now()

	return &tn
}
