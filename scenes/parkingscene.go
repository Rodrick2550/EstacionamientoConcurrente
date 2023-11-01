package scenes

import (
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
	"image/color"
	"math/rand"
	"parking/models"
	"sync"
	"time"
)

var (
	spots = []*models.ParkingSpot{

		// first row
		models.NewParkingSpot(410, 210, 440, 240, 1, 1),
		models.NewParkingSpot(410, 255, 440, 285, 1, 2),
		models.NewParkingSpot(410, 300, 440, 330, 1, 3),
		models.NewParkingSpot(410, 345, 440, 375, 1, 4),

		// second row
		models.NewParkingSpot(320, 210, 350, 240, 2, 5),
		models.NewParkingSpot(320, 255, 350, 285, 2, 6),
		models.NewParkingSpot(320, 300, 350, 330, 2, 7),
		models.NewParkingSpot(320, 345, 350, 375, 2, 8),

		// third row
		models.NewParkingSpot(230, 210, 260, 240, 3, 9),
		models.NewParkingSpot(230, 255, 260, 285, 3, 10),
		models.NewParkingSpot(230, 300, 260, 330, 3, 11),
		models.NewParkingSpot(230, 345, 260, 375, 3, 12),

		// fourth row
		models.NewParkingSpot(140, 210, 170, 240, 4, 13),
		models.NewParkingSpot(140, 255, 170, 285, 4, 14),
		models.NewParkingSpot(140, 300, 170, 330, 4, 15),
		models.NewParkingSpot(140, 345, 170, 375, 4, 16),

		// fifth row
		models.NewParkingSpot(50, 210, 80, 240, 5, 17),
		models.NewParkingSpot(50, 255, 80, 285, 5, 18),
		models.NewParkingSpot(50, 300, 80, 330, 5, 19),
		models.NewParkingSpot(50, 345, 80, 375, 5, 20),
	}
	parking    = models.NewParking(spots)
	doorMutex  sync.Mutex
	carManager = models.NewCarManager()
)

type ParkingScene struct {
}

func NewParkingScene() *ParkingScene {
	return &ParkingScene{}
}

func (ps *ParkingScene) Start() {
	isFirstTime := true

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for i := 0; i < 100; i++ {
					go carCycle(ctx)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})
}

func setUpScene(ctx *scene.Context) {

	parkingArea := floatgeom.NewRect2(20, 180, 500, 405)
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithColor(color.RGBA{100, 100, 100, 1}))

	parkingDoor := floatgeom.NewRect2(440, 170, 500, 180)
	entities.New(ctx, entities.WithRect(parkingDoor), entities.WithColor(color.RGBA{200, 0, 0, 1}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}

func carCycle(ctx *scene.Context) {
	car := models.NewCar(ctx)

	carManager.AddCar(car)

	car.Enqueue(carManager)

	spotAvailable := parking.GetParkingSpotAvailable()

	doorMutex.Lock()

	car.JoinDoor(carManager)

	doorMutex.Unlock()

	car.Park(spotAvailable, carManager)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(40000, 50000)))

	car.LeaveSpot(carManager)

	parking.ReleaseParkingSpot(spotAvailable)

	car.Leave(spotAvailable, carManager)

	doorMutex.Lock()

	car.ExitDoor(carManager)

	doorMutex.Unlock()

	car.GoAway(carManager)

	car.Remove()

	carManager.RemoveCar(car)
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}