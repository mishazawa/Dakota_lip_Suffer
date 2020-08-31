package synth

import (
	"math"
)

type Wavetable [SAMPLE_RATE]float32

func NewWaveTable () Wavetable {
	var table Wavetable
	for n := 0; n < SAMPLE_RATE; n += 1 {
		table[n] = sinePoint(n)
	}
	return table
}


func sinePoint (n int) float32 {
	return float32(math.Sin(float64(n) * 2.0 * math.Pi / float64(SAMPLE_RATE)))
}
