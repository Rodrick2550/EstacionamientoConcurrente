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
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(445, -20, 465, 0)
	spritePath := "assets/R.png"
	sprite, _ := render.LoadSprite(spritePath)
	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{255, 0, 0, 255}), entities.WithRenderable(sprite), entities.WithDrawLayers([]int{1, 2}))

	return &Car{
		area:   area,
		entity: entity,
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
	for index := 0; index < len(*spot.GetDirectionsForParking()); index++ {
		directions := *spot.GetDirectionsForParking()
		fmt.Println("Direction: " + directions[index].Direction)
		fmt.Println("Point: " + fmt.Sprintf("%f", directions[index].Point))
		if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.isCollision("right", manager.GetCars()) {
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.isCollision("down", manager.GetCars()) {
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "left" {
			for c.X() > directions[index].Point {
				if !c.isCollision("left", manager.GetCars()) {
					c.ShiftX(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "up" {
			for c.Y() > directions[index].Point {
				if !c.isCollision("up", manager.GetCars()) {
					c.ShiftY(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}
}

func (c *Car) Leave(spot *ParkingSpot, manager *CarManager) {
	for index := 0; index < len(*spot.GetDirectionsForLeaving()); index++ {
		directions := *spot.GetDirectionsForLeaving()
		if directions[index].Direction == "left" {

			for c.X() > directions[index].Point {
				if !c.isCollision("left", manager.GetCars()) {
					c.ShiftX(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.isCollision("right", manager.GetCars()) {
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "up" {
			for c.Y() > directions[index].Point {
				if !c.isCollision("up", manager.GetCars()) {
					c.ShiftY(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.isCollision("down", manager.GetCars()) {
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
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

func CarCycle(car *Car, manager *CarManager, parking *Parking, doorM *sync.Mutex) {

	manager.Add(car)

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