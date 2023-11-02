package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
	"image/color"
	"sync"
)

type CarQueue struct {
	cars []Car
}

type Parking struct {
	spots         []*ParkingSpot
	queueCars     *CarQueue
	mu            sync.Mutex
	availableCond *sync.Cond
}

func NewParking(ctx *scene.Context) *Parking {
	spots := []*ParkingSpot{

		NewParkingSpot(410, 210, 440, 240, 1, 1),
		NewParkingSpot(410, 255, 440, 285, 1, 2),
		NewParkingSpot(410, 300, 440, 330, 1, 3),
		NewParkingSpot(410, 345, 440, 375, 1, 4),

		NewParkingSpot(320, 210, 350, 240, 2, 5),
		NewParkingSpot(320, 255, 350, 285, 2, 6),
		NewParkingSpot(320, 300, 350, 330, 2, 7),
		NewParkingSpot(320, 345, 350, 375, 2, 8),

		NewParkingSpot(230, 210, 260, 240, 3, 9),
		NewParkingSpot(230, 255, 260, 285, 3, 10),
		NewParkingSpot(230, 300, 260, 330, 3, 11),
		NewParkingSpot(230, 345, 260, 375, 3, 12),

		NewParkingSpot(140, 210, 170, 240, 4, 13),
		NewParkingSpot(140, 255, 170, 285, 4, 14),
		NewParkingSpot(140, 300, 170, 330, 4, 15),
		NewParkingSpot(140, 345, 170, 375, 4, 16),

		NewParkingSpot(50, 210, 80, 240, 5, 17),
		NewParkingSpot(50, 255, 80, 285, 5, 18),
		NewParkingSpot(50, 300, 80, 330, 5, 19),
		NewParkingSpot(50, 345, 80, 375, 5, 20),
	}

	setUpScene(ctx, spots)
	queue := NewCarQueue()
	p := &Parking{
		spots:     spots,
		queueCars: queue,
	}
	p.availableCond = sync.NewCond(&p.mu)
	return p
}

func (p *Parking) GetSpots() []*ParkingSpot {
	return p.spots
}

func (p *Parking) GetParkingSpotAvailable() *ParkingSpot {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}

func (p *Parking) ReleaseParkingSpot(spot *ParkingSpot) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsAvailable(true)
	p.availableCond.Signal()
}

func (p *Parking) GetQueueCars() *CarQueue {
	return p.queueCars
}

func NewCarQueue() *CarQueue {
	return &CarQueue{
		cars: make([]Car, 0),
	}
}

func setUpScene(ctx *scene.Context, spots []*ParkingSpot) {

	parkingArea := floatgeom.NewRect2(20, 180, 500, 405)
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithColor(color.RGBA{100, 100, 100, 1}))

	parkingDoor := floatgeom.NewRect2(440, 170, 500, 180)
	entities.New(ctx, entities.WithRect(parkingDoor), entities.WithColor(color.RGBA{200, 0, 0, 1}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}
