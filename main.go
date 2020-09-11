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

	fmt.Println("Dakota lip Suffer")
	misc.Song()
}
