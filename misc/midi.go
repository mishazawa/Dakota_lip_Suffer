package misc

import (
	"fmt"
	"time"
	"sync"
	"github.com/mishazawa/heartache/parser"
	"github.com/mishazawa/heartache/parser/events"
)

type EventCallback func(*SynthTrack)

type SynthTrack struct {
	Tempo       int
	Events      map[int][]EventCallback
	Destination chan<- *events.MidiEvent
	Eot         chan<- bool
	mux					sync.Mutex
}

func (t *SynthTrack) SetTempo (bpm int) {
	t.mux.Lock()
	t.Tempo = 60000 / bpm
	t.mux.Unlock()
	fmt.Printf("SetTempo %d \n", t.Tempo)
}

func (t *SynthTrack) EmitEnd () {
	t.Eot <- true
}

func (t *SynthTrack) EmitMidiEvent (ev *events.MidiEvent) {
	t.Destination <- ev
}

func newSynthTrack (dest chan<- *events.MidiEvent, eot chan<- bool) *SynthTrack {
	return &SynthTrack{
		Tempo:  500, // ms
		Events: make(map[int][]EventCallback),
		Destination: dest,
		Eot: eot,
	}
}

type MidiSong struct {
	TimeDiv int
	Tracks  []*SynthTrack
}

func parseSong (dest chan<- *events.MidiEvent, eot chan<- bool) (*MidiSong, error) {
	song, err := parser.ParseFile("./test_data/A Sacred Lot.mid")
	if err != nil {
		return nil, err
	}

	midiSong := &MidiSong{int(song.TimeDiv), make([]*SynthTrack, 0)}

	for _, track := range song.Tracks {
		st := newSynthTrack(dest, eot)
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
