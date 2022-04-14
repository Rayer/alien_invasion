package alien_invastion

import (
	"fmt"
	"math/rand"
	"strings"
)

type Direction int

const (
	North Direction = iota
	West
	South
	East

	DirectionSize // must be last, if we need more direction please add in front of it.
	Invalid
)

func (d Direction) GetOpposite() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	default:
		return Invalid
	}
}

func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case South:
		return "south"
	case East:
		return "east"
	case West:
		return "west"
	default:
		return "invalid"
	}
}

func DirectionFromString(from string) Direction {
	switch strings.ToLower(from) {
	case "north":
		return North
	case "west":
		return West
	case "south":
		return South
	case "east":
		return East
	default:
		// This will make panic or error
		return Invalid
	}
}

type Neighborhoods []*City

func (n Neighborhoods) String() string {
	var lines []string
	for direction, city := range n {
		if city != nil && city.Exists {
			lines = append(lines, fmt.Sprintf("%v=%v", Direction(direction), city.Name))
		}
	}
	return strings.Join(lines, " ")
}

type City struct {
	Name          string
	Neighborhoods Neighborhoods
	Exists        bool
	AlienInCity   *Alien
}

//Should be private, and external call will only use name as index
func newCity(name string) *City {
	return &City{Name: name, Neighborhoods: make(Neighborhoods, DirectionSize), Exists: true}
}

func (c *City) AlienMigrate(to *City) {
	alien := c.AlienInCity
	if alien == nil || c.Exists == false {
		return
	}

	if to.AlienInCity != nil {
		fmt.Printf("City %s have been destroyed by alien %v and %v!\n", to.Name, alien.Number, to.AlienInCity.Number)
		to.Exists = false
	}
	to.AlienInCity = alien
	c.AlienInCity = nil
}

func (c *City) IsIsolatedOrDestroyed() bool {
	if c.Exists == false {
		return true
	}
	for _, n := range c.Neighborhoods {
		if n != nil && n.Exists {
			return false
		}
	}
	return true
}

type GameMap struct {
	cities map[string]*City
}

func NewGameMap() *GameMap {
	return &GameMap{
		cities: make(map[string]*City),
	}
}

func (m *GameMap) UpdateCityWithNeighborhood(name string, direction Direction, neighborhoodCityName string) error {
	city := m.UpsertCity(name)
	//No nil check because m.UpsertCity will be always exists
	neighborhoodCity := m.UpsertCity(neighborhoodCityName)
	city.Neighborhoods[direction] = neighborhoodCity
	if neighborhoodCity.Neighborhoods[direction.GetOpposite()] == nil {
		neighborhoodCity.Neighborhoods[direction.GetOpposite()] = city
	} else {
		if neighborhoodCity.Neighborhoods[direction.GetOpposite()] != city {
			// For example, A's north is B, but B's south is not A
			return fmt.Errorf("%s's %s is %s, but %s's %s is %s (conflict)", city.Name, direction, neighborhoodCity.Name, neighborhoodCity.Name, direction.GetOpposite(), neighborhoodCity.Neighborhoods[direction.GetOpposite()].Name)
		}
	}
	return nil
}

func (m *GameMap) UpsertCity(name string) *City {
	var city *City
	if c, exists := m.cities[name]; !exists {
		city = newCity(name)
		m.cities[name] = city
	} else {
		city = c
	}
	return city
}

func (m *GameMap) GetExistCity(name string) *City {
	if c, cityInMap := m.cities[name]; cityInMap {
		if c.Exists {
			return c
		}
	}
	return nil
}

func (m *GameMap) destroyCity(name string) error {
	if c, cityInMap := m.cities[name]; cityInMap {
		if c.Exists {
			c.Exists = false
		} else {
			return fmt.Errorf("city %s have been already destroyed", name)
		}
	} else {
		return fmt.Errorf("city %s doesn't exist", name)
	}
	return nil
}

// AssignAliens Assign aliens to cities
func (m *GameMap) AssignAliens(aliens []*Alien) error {
	for _, alien := range aliens {
		var candidates []*City
		for _, city := range m.cities {
			if city.Exists && city.AlienInCity == nil {
				candidates = append(candidates, city)
			}
		}
		if len(candidates) == 0 {
			return fmt.Errorf("not enough exist cities available to assign aliens")
		}
		candidates[rand.Intn(len(candidates))].AlienInCity = alien
	}
	return nil
}

func (m *GameMap) Update() (willContinue bool) {
	for _, city := range m.cities {
		if city.AlienInCity != nil {
			_, steps := city.AlienInCity.Move(city)
			if steps > 10000 {
				return false
			}
		}
	}
	return true
}

func (m *GameMap) DumpMap() string {
	var result []string
	for _, city := range m.cities {
		if city.Exists {
			result = append(result, fmt.Sprintf("%s %s", city.Name, city.Neighborhoods.String()))
		}
	}
	return strings.Join(result, "\n")
}

func (m *GameMap) ExistCityCount() int {
	ret := 0
	for _, c := range m.cities {
		if c.Exists {
			ret++
		}
	}
	return ret
}
