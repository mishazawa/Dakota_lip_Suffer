package synth

import (
	"math"
)

type Mixer struct {
	voices []*Voice
}

func NewMixer () Mixer {
	return Mixer {
		voices: []*Voice {},
	}
}

func (mixer *Mixer) AddVoice (voice *Voice) {
	mixer.voices = append(mixer.voices, voice)
}

func (mixer *Mixer) ProcessAudio(out [][]float32) {
	for i := range out[0] {
		val := mixer.mix()
		out[0][i] = val
		out[1][i] = val
	}
}

func (mixer *Mixer) mix () float32 {
	if len(mixer.voices) == 0 {
		return 0.0
	}

	poly 			:= float64(len(mixer.voices))
	threshold := 1.0 / poly
	samples   := make([]float64, len(mixer.voices))


	for i, v := range mixer.voices {
		samples[i] = v.GetSample()
		v.NextPhase()
	}

	res := 0.0

	for _, s := range samples {
		res += mixer.compress(s, threshold, poly)
	}

	return mixer.limit(res, LIMIT_CLIP)
}

func (mixer *Mixer) limit (sample, threshold float64) float32 {
	return float32(sample * threshold)
}

func (mixer *Mixer) compress (sample, threshold, ratio float64) float64 {
	if math.Abs(sample) > threshold {
		return sample / ratio
	}
	return sample
}
