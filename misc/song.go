package misc

import (
	"time"
	"github.com/mishazawa/Dakota_lip_Suffer/synth"
)

func Song (out *synth.Output) {
	out.Mixer.AddVoice(synth.NewVoice(440))
	time.Sleep(2 * time.Second)
	out.Mixer.AddVoice(synth.NewVoice(220))
	time.Sleep(3 * time.Second)
	out.Mixer.AddVoice(synth.NewVoice(7878))
	time.Sleep(3 * time.Second)
	out.Mixer.AddVoice(synth.NewVoice(666))
	time.Sleep(3 * time.Second)
}
