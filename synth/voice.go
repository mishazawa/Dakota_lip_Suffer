package synth

type Voice struct {
	Audiable
}

func NewVoice (freq int) *Voice {
	osc := NewOsc(freq)
	return &Voice { &osc }
}
