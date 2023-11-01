package models

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
	entranceSpot = 185.00
	speed        = 10
)

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(445, -20, 465, 0)
	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{255, 0, 0, 255}))

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
	for c.Y() < entranceSpot {
		if !c.isCollision("down", manager.GetCars()) {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ExitDoor(manager *CarManager) {
	for c.Y() > 145 {
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
	minDistance := 30.0
	for _, car := range cars {
		if direction == "left" {
			if c.X() > car.X() && c.X()-car.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "right" {
			if c.X() < car.X() && car.X()-c.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "up" {
			if c.Y() > car.Y() && c.Y()-car.Y() < minDistance && c.X() == car.X() {
				return true
			}
		} else if direction == "down" {
			if c.Y() < car.Y() && car.Y()-c.Y() < minDistance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}