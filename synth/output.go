package synth

import (
	"github.com/gordonklaus/portaudio"
)

type Output struct {
	*portaudio.Stream
	*Mixer
}

func NewOutput (mixer *Mixer) (*Output, error) {
	var err error

	output := Output {
		Stream: nil,
		Mixer: mixer,
	}
	output.Stream, err = portaudio.OpenDefaultStream(0, 2, SAMPLE_RATE, 0, output.Mixer.ProcessAudio)

	return &output, err
}


func Init () error {
	return portaudio.Initialize()
}

func Terminate () {
	portaudio.Terminate()
}
