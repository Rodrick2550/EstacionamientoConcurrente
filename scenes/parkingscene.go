package scenes

import (
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
	"parking/models"
	"parking/utils"
	"sync"
	"time"
)

type MainScene struct {
}

func NewParkingScene() *MainScene {
	return &MainScene{}
}

func (ps *MainScene) Draw() {
	firstTime := true
	manager := models.NewCarManager()
	doorM := sync.Mutex{}

	_ = oak.AddScene("mainScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			parking := models.NewParking(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !firstTime {
					return 0
				}
				firstTime = false

				for {
					car := models.NewCar(ctx)
					go models.CarCycle(car, manager, parking, &doorM)

					time.Sleep(time.Millisecond * time.Duration(utils.Number(1000, 2000)))
				}

				return 0
			})
		},
	})
}