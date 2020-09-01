package main

import (
	"fmt"
	"github.com/mishazawa/Dakota_lip_Suffer/synth"
	"github.com/mishazawa/Dakota_lip_Suffer/misc"
)

func main() {
	err := synth.Init()

	if err != nil {
		panic(err)
	}

	defer synth.Terminate()

	out, err := synth.NewOutput()

	if err != nil {
		panic(err)
	}

	defer out.Close()

	if err := out.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Dakota lip Suffer")
	misc.Song(out)
}
