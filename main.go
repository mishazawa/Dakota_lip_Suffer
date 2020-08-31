package main

import (
	"fmt"
	"math"

	"github.com/gordonklaus/portaudio"

	"github.com/mishazawa/Dakota_lip_Suffer/synth"
)

func main() {
	err := portaudio.Initialize()
	if err != nil {
		panic(err)
	}

	defer portaudio.Terminate()

	osc := synth.NewOsc()
	stream, err := portaudio.OpenDefaultStream(0, 2, synth.SAMPLE_RATE, 0, osc.ProcessAudio)

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	if err := stream.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Dakota lip Suffer")
	for {

	}
}
