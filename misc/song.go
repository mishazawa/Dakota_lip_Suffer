package misc

import (
	"fmt"
	"sync"
	"math"

	"github.com/mishazawa/Dakota_lip_Suffer/synth"
	"github.com/mishazawa/heartache/parser/events"
)



func Song () {
	midiEvents := make(chan *events.MidiEvent)
	eotEvents  := make(chan bool)

	msong, err := parseSong(midiEvents, eotEvents)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", msong)

	voices := make(map[int][]*synth.Voice)

	for v := 0; v < synth.POLYPHONY; v += 1 {
		voice := synth.NewVoice()
		synth.AddVoiceToChannel(voices, voice)
	}

	mixer := synth.NewMixer(voices)

	out, err := synth.NewOutput(mixer)

	if err != nil {
		panic(err)
	}

	defer out.Close()

	if err := out.Start(); err != nil {
		panic(err)
	}

	for ti, track := range msong.Tracks {
		go func (track *SynthTrack, i int) {
			for _, events := range track.Events {
				var wg sync.WaitGroup
				wg.Add(len(events))
				for _, callback := range events {
					go func (track *SynthTrack, i int, cb EventCallback) {
						defer wg.Done()
						cb(track)
					}(track, i, callback)
				}
				wg.Wait()
			}
		}(track, ti)
	}


	tracksEnded := 0

  for {
    select {
    case <-eotEvents:
    	tracksEnded += 1
      fmt.Printf("received %d eot \n", tracksEnded)
      if tracksEnded == len(msong.Tracks) {
      	return
      }
    case mevent := <-midiEvents:
			if mevent.Status & 0xf0 == 0x90 || mevent.Status & 0xf0 == 0x80 {
	      channel := int(mevent.Status & 0x0f)
	      disable := mevent.Status & 0xf0 == 0x80 || mevent.Data[1] == 0x00
	      freq := mtof(mevent.Data[0])
	      err := out.Mixer.MidiNote(channel, freq, !disable)
	      if err != nil {
		      panic(err)
	      }
			}
    }
  }
}

func mtof (key byte) int {
	return int(440 * math.Pow(2, float64(int(key) - 69) / 12))
}
