package main

import (
	"github.com/oakmound/oak/v4"
	"parking/scenes"
)

func main() {
	parkingScene := scenes.NewParkingScene()

	parkingScene.Draw()

	_ = oak.Init("mainScene")
}