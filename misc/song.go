package misc

import (
	"fmt"
	"sync"

	"github.com/mishazawa/Dakota_lip_Suffer/synth"
)



func Song (out *synth.Output) {
	msong, err := parseSong()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", msong)

	for ti, track := range msong.Tracks {
		go func (track *SynthTrack, i int) {
			for _, events := range track.Events {
				var wg sync.WaitGroup
				wg.Add(len(events))
				for _, callback := range events {
					go func (track *SynthTrack, i int) {
						defer wg.Done()
						fmt.Printf("track %d\n", i)
						callback(track)
					}(track, i)
				}
				wg.Wait()
			}
		}(track, ti)
	}


	// out.Mixer.AddVoice(synth.NewVoice(440))
	// out.Mixer.AddVoice(synth.NewVoice(220))
	// time.Sleep(3 * time.Second)
	// out.Mixer.AddVoice(synth.NewVoice(7878))
	// time.Sleep(3 * time.Second)
	// out.Mixer.AddVoice(synth.NewVoice(666))
	// time.Sleep(3 * time.Second)
}
