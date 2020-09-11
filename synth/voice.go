package synth

type Voice struct {
	Audiable
	Freq int
	Enabled bool
	MidiChannel int
}

func NewVoice () *Voice {
	osc := NewOsc()
	return &Voice { &osc, 0, false, 0 }
}

func (v *Voice) Enable () {
	v.Enabled = true
}

func (v *Voice) Disable () {
	v.Enabled = false
}

func (v *Voice) SetChannel (mch int) {
	v.MidiChannel = mch
}

func (v *Voice) SetFreq (freq int) {
	v.Audiable.(*Oscillator).SetFreq(freq)
	v.Freq = freq
}

func AddVoiceToChannel (dest map[int][]*Voice, voice *Voice) {
  if _, ok := dest[voice.MidiChannel]; !ok {
    dest[voice.MidiChannel] = make([]*Voice, POLYPHONY)
  }

  dest[voice.MidiChannel] = append(dest[voice.MidiChannel], voice)
}
