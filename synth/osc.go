package synth

type Audiable interface {
	GetSample() float64
	NextPhase()
}

type Oscillator struct {
	table Wavetable
	phase int
	freq  int
}

func NewOsc () Oscillator {
	return Oscillator {
		table: NewWaveTable(),
		phase: 0,
		freq:  0,
	}
}

func (osc *Oscillator) NextPhase () {
	osc.phase = (osc.phase + osc.freq) % SAMPLE_RATE
}

func (osc *Oscillator) GetSample () float64 {
	return osc.table[osc.phase]
}

func (osc *Oscillator) SetFreq (freq int) {
	osc.freq = freq
}
