package misc

import (
	"fmt"
	"time"
	"github.com/mishazawa/heartache/parser"
	"github.com/mishazawa/heartache/parser/events"
)

type EventCallback func(*SynthTrack)

type SynthTrack struct {
	Tempo  int
	Events map[int][]EventCallback
}

func (t *SynthTrack) SetTempo (bpm int) {
	t.Tempo = 60000 / bpm
	fmt.Printf("SetTempo %d \n", t.Tempo)
}

func (t *SynthTrack) EmitEnd () {
	fmt.Printf("emit end\n")
}

func (t *SynthTrack) EmitMidiEvent (ev *events.MidiEvent) {
	fmt.Printf("emit midi\n")
}

func newSynthTrack () *SynthTrack {
	return &SynthTrack{
		Tempo:  500, // ms
		Events: make(map[int][]EventCallback),
	}
}

type MidiSong struct {
	TimeDiv int
	Tracks  []*SynthTrack
}

func parseSong () (*MidiSong, error) {
	song, err := parser.ParseFile("./test_data/A Sacred Lot.mid")
	if err != nil {
		return nil, err
	}

	midiSong := &MidiSong{int(song.TimeDiv), make([]*SynthTrack, 0)}

	for _, track := range song.Tracks {
		st := newSynthTrack()
		eventTime := 0

		for _, event := range track.Events {
			eventDelta := dtoms(st.Tempo, int(*event.GetDelta()), int(song.TimeDiv))
			if _, ok := st.Events[eventTime]; !ok {
				st.Events[eventTime] = make([]EventCallback, 0)
			}

			switch event.(type) {
			case *events.SetTempoEvent:
				tempoEvent := event.(*events.SetTempoEvent)
				st.Events[eventTime] = append(st.Events[eventTime], func (t *SynthTrack) {
					time.Sleep(time.Duration(eventDelta) * time.Millisecond)
					t.SetTempo(int(tempoEvent.Tempo / 1000))
				})
			case *events.MidiEvent:
				midiEvent := event.(*events.MidiEvent)

				st.Events[eventTime] = append(st.Events[eventTime], func (t *SynthTrack) {
					time.Sleep(time.Duration(eventDelta) * time.Millisecond)
					t.EmitMidiEvent(midiEvent)
				})
			case *events.EndOfTrackEvent:
				st.Events[eventTime] = append(st.Events[eventTime], func (t *SynthTrack) {
					time.Sleep(time.Duration(eventDelta) * time.Millisecond)
					t.EmitEnd()
				})
			default:
			}
			eventTime += eventDelta
		}
		midiSong.Tracks = append(midiSong.Tracks, st)
	}

	return midiSong, nil

}

func dtoms (tempo, ticks, division int) int {
	return ticks * (tempo / division) / 1000
}
