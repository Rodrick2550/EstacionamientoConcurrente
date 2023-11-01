package models

import (
	"fmt"
	"github.com/oakmound/oak/v4/render"
	"image/color"
	"parking/utils"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
    initDoorPoint = 185.00
    endDoorPoint  = 145.00
    speed         = 10
)

type Car struct {
    area   floatgeom.Rect2
    entity *entities.Entity
    mu     sync.Mutex
    manager *CarManager 
}

type CarManager struct {
    Cars  []*Car
    Mutex sync.Mutex
}

func NewCarManager() *CarManager {
    return &CarManager{
        Cars: make([]*Car, 0),
    }
}

var manager = NewCarManager()

func (cm *CarManager) Add(car *Car) {
    cm.Mutex.Lock()
    defer cm.Mutex.Unlock()
    cm.Cars = append(cm.Cars, car)
}

func (cm *CarManager) Remove(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	for i, c := range cm.Cars {
		if c == car {
			cm.Cars = append(cm.Cars[:i], cm.Cars[i+1:]...)
			break
		}
	}
}

func (cm *CarManager) GetCars() []*Car {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	return cm.Cars
}


func NewCar(ctx *scene.Context) *Car {

    area := floatgeom.NewRect2(445, -20, 465, 0)
	spritePath := "assets/R.png"
	sprite, _ := render.LoadSprite(spritePath)
	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{255, 0, 0, 255}), entities.WithRenderable(sprite), entities.WithDrawLayers([]int{1, 2}))

    return &Car{
        area:   area,
        entity: entity,
        manager: manager, 
    }
}
func (c *Car) Enqueue(manager *CarManager) {

	for c.Y() < 145 {
		if !c.isCollision("down", manager.GetCars()) {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}

}

func (c *Car) JoinDoor(manager *CarManager) {
	for c.Y() < initDoorPoint {
		if !c.isCollision("down", manager.GetCars()) {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ExitDoor(manager *CarManager) {
	for c.Y() > endDoorPoint {
		if !c.isCollision("up", manager.GetCars()) {
			c.ShiftY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) Park(spot *ParkingSpot, manager *CarManager) {
	directions := *spot.GetDirectionsForParking()

	for _, dir := range directions {
		fmt.Println("Direction: " + dir.Direction)
		fmt.Println("Point: " + fmt.Sprintf("%f", dir.Point))

		switch dir.Direction {
		case "right":
			c.moveUntilCondition("right", dir.Point, func() bool {
				return !c.isCollision("right", manager.GetCars())
			})
		case "down":
			c.moveUntilCondition("down", dir.Point, func() bool {
				return !c.isCollision("down", manager.GetCars())
			})
		case "left":
			c.moveUntilCondition("left", dir.Point, func() bool {
				return !c.isCollision("left", manager.GetCars())
			})
		case "up":
			c.moveUntilCondition("up", dir.Point, func() bool {
				return !c.isCollision("up", manager.GetCars())
			})
		}
	}
}


func (c *Car) Leave(spot *ParkingSpot, manager *CarManager) {
	directions := *spot.GetDirectionsForLeaving()

	for _, dir := range directions {
		switch dir.Direction {
		case "left":
			c.moveUntilCondition("left", dir.Point, func() bool {
				return !c.isCollision("left", manager.GetCars())
			})
		case "right":
			c.moveUntilCondition("right", dir.Point, func() bool {
				return !c.isCollision("right", manager.GetCars())
			})
		case "up":
			c.moveUntilCondition("up", dir.Point, func() bool {
				return !c.isCollision("up", manager.GetCars())
			})
		case "down":
			c.moveUntilCondition("down", dir.Point, func() bool {
				return !c.isCollision("down", manager.GetCars())
			})
		}
	}
}

func (c *Car) moveUntilCondition(direction string, point float64, condition func() bool) {
	switch direction {
	case "left":
		for c.X() > point && condition() {
			c.ShiftX(-1)
			time.Sleep(speed * time.Millisecond)
		}
	case "right":
		for c.X() < point && condition() {
			c.ShiftX(1)
			time.Sleep(speed * time.Millisecond)
		}
	case "up":
		for c.Y() > point && condition() {
			c.ShiftY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	case "down":
		for c.Y() < point && condition() {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}


func (c *Car) LeaveSpot(manager *CarManager) {
	spotX := c.X()
	for c.X() > spotX-30 {
		if !c.isCollision("left", manager.GetCars()) {
			c.ShiftX(-1)
			time.Sleep(speed * time.Millisecond)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (c *Car) GoAway(manager *CarManager) {
	for c.Y() > -20 {
		if !c.isCollision("up", manager.GetCars()) {
			c.ShiftY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}


func CarCycle(car *Car, parking *Parking, doorM *sync.Mutex) {
    car.manager.Add(car) 
    
	car.Enqueue(manager)

	spotAvailable := parking.GetParkingSpotAvailable()

	doorM.Lock()

	car.JoinDoor(manager)

	doorM.Unlock()

	car.Park(spotAvailable, manager)

	time.Sleep(time.Millisecond * time.Duration(utils.Number(40000, 50000)))

	car.LeaveSpot(manager)

	parking.ReleaseParkingSpot(spotAvailable)

	car.Leave(spotAvailable, manager)

	doorM.Lock()

	car.ExitDoor(manager)

	doorM.Unlock()

	car.GoAway(manager)

	car.Remove()

	manager.Remove(car)
}

func (c *Car) ShiftY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) ShiftX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) isCollision(direction string, cars []*Car) bool {
	distance := 25.0
	for _, car := range cars {
		if direction == "left" {
			if c.X() > car.X() && c.X()-car.X() < distance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "right" {
			if c.X() < car.X() && car.X()-c.X() < distance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "up" {
			if c.Y() > car.Y() && c.Y()-car.Y() < distance && c.X() == car.X() {
				return true
			}
		} else if direction == "down" {
			if c.Y() < car.Y() && car.Y()-c.Y() < distance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}
