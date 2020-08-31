package synth

import (
	"github.com/gordonklaus/portaudio"
)

type Oscillator struct {
	*portaudio.Stream

	table Wavetable
	phase int
	freq  int
}

func NewOsc () Oscillator {
	return Oscillator {
		table: NewWaveTable(),
		phase: 0,
		freq:  440,
	}
}

func (osc *Oscillator) NextPhase () {
	osc.phase = (osc.phase + osc.freq) % SAMPLE_RATE
}

func (osc *Oscillator) GetSample () float32 {
	return osc.table[osc.phase]
}


func (osc *Oscillator) ProcessAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = osc.GetSample()
		out[1][i] = osc.GetSample()
		osc.NextPhase()
	}
}
