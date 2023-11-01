package models

import (
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

func NewParking(spots []*ParkingSpot) *Parking {
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