package alien_invastion

import "math/rand"

var alienSerialCount = 0

// Alien is a struct that represents an alien. Alien always attached in struct City
type Alien struct {
	Number int
	Steps  int
	Alive  bool
}

func NewAlien() *Alien {
	defer func() {
		alienSerialCount++
	}()
	return &Alien{Number: alienSerialCount, Steps: 0, Alive: true}
}

// Move is a method that moves alien to a random neighbor city
// All battle result against two aliens is done in City.AlienMigrate
func (a *Alien) Move(from *City) (to *City, step int) {
	a.Steps++
	if a.Alive == false {
		return from, a.Steps
	}

	if from.IsIsolatedOrDestroyed() {
		//Still alive, but no longer able to move
		return from, a.Steps
	}

	// Let's pick one of the cities that are not destroyed
	var candidates []*City
	for _, city := range from.Neighborhoods {
		if city == nil {
			continue
		}
		if city.Exists {
			candidates = append(candidates, city)
		}
	}

	// randomly pick one of the candidates
	var index = rand.Intn(len(candidates))
	from.AlienMigrate(candidates[index])

	return candidates[index], a.Steps
}
