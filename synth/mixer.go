package synth

import (
  "fmt"
  "math"
  "sync"
)

type Mixer struct {
  voices map[int][]*Voice
  mux    sync.Mutex
}

func NewMixer (vm map[int][]*Voice) *Mixer {
  return &Mixer {
    voices: vm,
  }
}

func (mixer *Mixer) ProcessAudio (out [][]float32) {
  for i := range out[0] {
    val := mixer.mix()
    out[0][i] = val
    out[1][i] = val
  }
}

func (mixer *Mixer) MidiNote (channel, freq int, message bool) error {
  if message {
    v, err := mixer.findFreeVoice(channel)
    if err != nil {
      return err
    }

    if v != nil {
      mixer.mux.Lock()
      v.SetFreq(freq)
      v.Enable()
      mixer.mux.Unlock()
    }
  } else {
    v, err := mixer.findEnabledVoice(channel, freq)
    if err != nil {
      return err
    }
    if v != nil {
      mixer.mux.Lock()
      v.Disable()
      mixer.mux.Unlock()
    }
  }

  return nil
}

func (mixer *Mixer) findFreeVoice (channel int) (*Voice, error) {
  if _, ok := mixer.voices[channel]; !ok {
    return nil, fmt.Errorf("Error findFreeVoice@%d \n", channel)
  }

  for _, v := range mixer.voices[channel] {
    if v != nil && !v.Enabled {
      return v, nil
    }
  }

  mixer.mux.Lock()
  mixer.voices[channel] = append(
    mixer.voices[channel][len(mixer.voices[channel])-1:],
    mixer.voices[channel][:len(mixer.voices[channel])-1]...
  )
  mixer.mux.Unlock()

  return mixer.voices[channel][0], nil
}

func (mixer *Mixer) findEnabledVoice (channel, freq int) (*Voice, error) {
  if _, ok := mixer.voices[channel]; !ok {
    return nil, fmt.Errorf("Error findFreeVoice@%d \n", channel)
  }

  for _, v := range mixer.voices[channel] {
    if v != nil && v.Freq == freq {
      return v, nil
    }
  }

  return nil, nil
}


func (mixer *Mixer) mix () float32 {
  samples   := make([]float64, 0)

  for _, channel := range mixer.voices {
    for _, voice := range channel {
      if voice != nil && voice.Enabled {
        mixer.mux.Lock()
        sample := voice.GetSample()
        voice.NextPhase()
        mixer.mux.Unlock()
        samples = append(samples, sample)
      }
    }
  }

  poly      := float64(len(samples))
  threshold := 1.0 / poly

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
