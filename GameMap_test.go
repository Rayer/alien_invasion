package alien_invastion

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDirectionFromString(t *testing.T) {
	type args struct {
		from string
	}
	tests := []struct {
		name string
		args args
		want Direction
	}{
		{
			name: "Good - North",
			args: args{
				from: "north",
			},
			want: North,
		},
		{
			name: "Good - East",
			args: args{
				from: "east",
			},
			want: East,
		},
		{
			name: "Good - South",
			args: args{
				from: "south",
			},
			want: South,
		},
		{
			name: "Good - West",
			args: args{
				from: "west",
			},
			want: West,
		},
		{
			name: "Good - With mixed case",
			args: args{
				from: "NorTh",
			},
			want: North,
		},
		{
			name: "Bad - Empty",
			args: args{
				from: "",
			},
			want: Invalid,
		},
		{
			name: "Bad - Unknown",
			args: args{
				from: "unknown",
			},
			want: Invalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DirectionFromString(tt.args.from), "DirectionFromString(%v)", tt.args.from)
		})
	}
}

func TestDirection_GetOpposite(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want Direction
	}{
		{
			name: "Good - North",
			d:    North,
			want: South,
		},
		{
			name: "Good - East",
			d:    East,
			want: West,
		},
		{
			name: "Good - South",
			d:    South,
			want: North,
		},
		{
			name: "Good - West",
			d:    West,
			want: East,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.d.GetOpposite(), "GetOpposite()")
		})
	}
}

func TestGameMap_DestroyCity(t *testing.T) {
	type fields struct {
		cities map[string]*City
	}
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		validate func(*testing.T, *GameMap)
	}{
		{
			name: "Good - Destroy city",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: true,
					},
				},
			},
			args: args{
				name: "city1",
			},
			wantErr: assert.NoError,
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Falsef(t, gameMap.cities["city1"].Exists, "Exists")
			},
		},
		{
			name: "Bad - City not found",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: true,
					},
				},
			},
			args: args{
				name: "city2",
			},
			wantErr: assert.Error,
			validate: func(t *testing.T, gameMap *GameMap) {

			},
		},
		{
			name: "Bad - already destroyed",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: false,
					},
				},
			},
			args: args{
				name: "city1",
			},
			wantErr: assert.Error,
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Falsef(t, gameMap.cities["city1"].Exists, "Exists")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GameMap{
				cities: tt.fields.cities,
			}

			tt.wantErr(t, m.destroyCity(tt.args.name), fmt.Sprintf("DestroyCity(%v)", tt.args.name))
			tt.validate(t, m)
		})
	}
}

func TestGameMap_GetExistCity(t *testing.T) {
	type fields struct {
		cities map[string]*City
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *City
	}{
		{
			name: "Good - City found",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: true,
					},
				},
			},
			args: args{
				name: "city1",
			},
			want: &City{
				Name:   "city1",
				Exists: true,
			},
		},
		{
			name: "Good - Should not return destroyed city",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: false,
					},
				},
			},
			args: args{
				name: "city1",
			},
			want: nil,
		},
		{
			name: "Bad - City not found",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: true,
					},
				},
			},
			args: args{
				name: "city2",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GameMap{
				cities: tt.fields.cities,
			}
			assert.Equalf(t, tt.want, m.GetExistCity(tt.args.name), "GetExistCity(%v)", tt.args.name)
		})
	}
}

func TestGameMap_UpdateCityWithNeighborhood(t *testing.T) {
	type fields struct {
		cities map[string]*City
	}
	type args struct {
		name                 string
		direction            Direction
		neighborhoodCityName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		// Patch uses to alter map, in case `fields` not working
		patch    func(t *testing.T, gameMap *GameMap)
		validate func(*testing.T, *GameMap)
	}{
		{
			name: "Good - Updated empty side",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:          "city1",
						Neighborhoods: make(Neighborhoods, DirectionSize),
						Exists:        true,
					},
				},
			},
			args: args{
				name:                 "city2",
				direction:            East,
				neighborhoodCityName: "city1",
			},
			wantErr: assert.NoError,
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Equal(t, "city2", gameMap.cities["city1"].Neighborhoods[West].Name)
				assert.Equal(t, "city1", gameMap.cities["city2"].Neighborhoods[East].Name)
			},
		},
		{
			name: "Good - Update neighborhood and it exists",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:          "city1",
						Neighborhoods: make(Neighborhoods, DirectionSize),
						Exists:        true,
					},
				},
			},
			patch: func(t *testing.T, gameMap *GameMap) {
				_ = gameMap.UpdateCityWithNeighborhood("city1", West, "city2")
			},
			args: args{
				name:                 "city2",
				direction:            East,
				neighborhoodCityName: "city1",
			},
			wantErr: assert.NoError,
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Equal(t, "city2", gameMap.cities["city1"].Neighborhoods[West].Name)
				assert.Equal(t, "city1", gameMap.cities["city2"].Neighborhoods[East].Name)
			},
		},
		{
			name: "Bad - Update neighborhood conflict",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:          "city1",
						Neighborhoods: make(Neighborhoods, DirectionSize),
						Exists:        true,
					},
				},
			},
			patch: func(t *testing.T, gameMap *GameMap) {
				_ = gameMap.UpdateCityWithNeighborhood("city1", West, "city3")
			},
			args: args{
				name:                 "city2",
				direction:            East,
				neighborhoodCityName: "city1",
			},
			wantErr: assert.Error,
			validate: func(t *testing.T, gameMap *GameMap) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GameMap{
				cities: tt.fields.cities,
			}
			if tt.patch != nil {
				tt.patch(t, m)
			}
			tt.wantErr(t, m.UpdateCityWithNeighborhood(tt.args.name, tt.args.direction, tt.args.neighborhoodCityName), fmt.Sprintf("UpdateCityWithNeighborhood(%v, %v, %v)", tt.args.name, tt.args.direction, tt.args.neighborhoodCityName))
			tt.validate(t, m)
		})
	}
}

func TestGameMap_UpsertCity(t *testing.T) {
	type fields struct {
		cities map[string]*City
	}
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		patch    func(t *testing.T, gameMap *GameMap)
		validate func(t *testing.T, gameMap *GameMap)
	}{
		// This test case doesn't have BAD condition
		{
			name: "Good - Update with empty map",
			fields: fields{
				cities: make(map[string]*City),
			},
			args: args{
				name: "city1",
			},
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Equal(t, 1, len(gameMap.cities))
			},
		},
		{
			name: "Good - Update with existing map",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:   "city1",
						Exists: true,
					},
				},
			},
			args: args{
				name: "city1",
			},
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Equal(t, 1, len(gameMap.cities))
			},
		},
		{
			name: "Good - Don't override exist neighborhood",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						Name:          "city1",
						Neighborhoods: make(Neighborhoods, DirectionSize),
						Exists:        true,
					},
				},
			},
			args: args{
				name: "city1",
			},
			patch: func(t *testing.T, gameMap *GameMap) {
				_ = gameMap.UpdateCityWithNeighborhood("city1", West, "city2")
			},
			validate: func(t *testing.T, gameMap *GameMap) {
				assert.Equal(t, 2, len(gameMap.cities))
				assert.Equal(t, "city2", gameMap.cities["city1"].Neighborhoods[West].Name)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GameMap{
				cities: tt.fields.cities,
			}
			if tt.patch != nil {
				tt.patch(t, m)
			}
			//ALWAYS return value
			assert.NotNil(t, m.UpsertCity(tt.args.name), "UpsertCity(%v)", tt.args.name)
			tt.validate(t, m)
		})
	}
}

func TestNewCity(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *City
	}{
		{
			name: "Sanity Case",
			args: args{
				name: "city1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city := newCity(tt.args.name)
			assert.NotNil(t, city.Neighborhoods)
			assert.True(t, city.Exists)
		})
	}
}

func TestNewGameMap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Sanity Case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewGameMap()
			assert.NotNil(t, m.cities)
		})
	}
}

func TestCity_IsIsolatedOrDestroyed(t *testing.T) {
	type fields struct {
		Name          string
		Neighborhoods Neighborhoods
		Exists        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Not Isolated, city exists and have alive neighborhood",
			fields: fields{
				Name: "city1",
				Neighborhoods: Neighborhoods{
					&City{
						Name:   "city2",
						Exists: true,
					},
					nil,
					nil,
					nil,
				},
				Exists: true,
			},

			want: false,
		},
		{
			name: "Destroyed City",
			fields: fields{
				Name: "city1",
				Neighborhoods: Neighborhoods{
					&City{
						Name:   "city2",
						Exists: true,
					},
					nil,
					nil,
					nil,
				},
				Exists: false,
			},
			want: true,
		},
		{
			name: "No neighborhoods",
			fields: fields{
				Name: "city1",
				Neighborhoods: Neighborhoods{
					nil,
					nil,
					nil,
					nil,
				},
				Exists: false,
			},
			want: true,
		},
		{
			name: "All neighborhoods destroyed",
			fields: fields{
				Name: "city1",
				Neighborhoods: Neighborhoods{
					&City{
						Name:   "city2",
						Exists: false,
					},
					nil,
					nil,
					nil,
				},
				Exists: true,
			},
			want: true,
		},
		{
			name: "Being destroyed and all neighborhoods destroyed",
			fields: fields{
				Name: "city1",
				Neighborhoods: Neighborhoods{
					&City{
						Name:   "city2",
						Exists: false,
					},
					nil,
					nil,
					nil,
				},
				Exists: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &City{
				Name:          tt.fields.Name,
				Neighborhoods: tt.fields.Neighborhoods,
				Exists:        tt.fields.Exists,
			}
			assert.Equalf(t, tt.want, c.IsIsolatedOrDestroyed(), "IsIsolatedOrDestroyed()")
		})
	}
}

func TestCity_AlienMigrate(t *testing.T) {

	type fields struct {
		Name          string
		Neighborhoods Neighborhoods
		Exists        bool
		AlienInCity   *Alien
	}
	type args struct {
		to    *City
		patch func(t *testing.T, from *City, alien *Alien)
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		validate func(t *testing.T, from, to *City)
	}{
		{
			name: "Migrate to empty city",
			fields: fields{
				Name:        "city1",
				Exists:      true,
				AlienInCity: NewAlien(),
			},
			args: args{
				to: &City{
					Name:   "city2",
					Exists: true,
				},
			},
			validate: func(t *testing.T, from, to *City) {
				assert.Nil(t, from.AlienInCity)
				assert.NotNil(t, from.AlienInCity)
				assert.True(t, from.Exists)
				assert.True(t, to.Exists)
			},
		},
		{
			name: "Migrate to occupied city",
			fields: fields{
				Name:        "city1",
				Exists:      true,
				AlienInCity: NewAlien(),
			},
			args: args{
				to: &City{
					Name:        "city2",
					Exists:      true,
					AlienInCity: NewAlien(),
				},
			},
			validate: func(t *testing.T, from, to *City) {
				assert.Nil(t, from.AlienInCity)
				assert.NotNil(t, from.AlienInCity)
				assert.True(t, from.Exists)
				//Target city should have been destroyed
				assert.False(t, to.Exists)
			},
		},
		{
			name: "Should not migrate to a destroyed city",
			fields: fields{
				Name:   "city1",
				Exists: true,
				AlienInCity: &Alien{
					Number: 0,
					Steps:  0,
					Alive:  true,
				},
			},
			args: args{
				to: &City{
					Name:   "city2",
					Exists: false,
				},
			},
			validate: func(t *testing.T, from, to *City) {
				//Alien should stay in current city
				assert.NotNil(t, to.AlienInCity)
				assert.True(t, from.Exists)
				assert.False(t, to.Exists)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &City{
				Name:          tt.fields.Name,
				Neighborhoods: tt.fields.Neighborhoods,
				Exists:        tt.fields.Exists,
				AlienInCity:   tt.fields.AlienInCity,
			}
			//connect from and to
			c.Neighborhoods = Neighborhoods{tt.args.to, nil, nil, nil}
			tt.args.to.Neighborhoods = Neighborhoods{nil, nil, c, nil}
			c.AlienMigrate(tt.args.to)
		})
	}
}

func TestGameMap_DumpMap(t *testing.T) {
	parser := StreamParser{}
	tests := []struct {
		name    string
		gameMap *GameMap
		// It is not a good way to compare strings... but let's not overkill
		validate func(t *testing.T, dumpedString string)
	}{
		{
			name: "Just load from map",
			gameMap: func() *GameMap {
				m, _ := parser.ParseFile("test_resources/standard_input1.txt")
				return m
			}(),
			validate: func(t *testing.T, dumpedString string) {
				//Should be 5 entities(lines), so 4 \n
				assert.Equal(t, 4, strings.Count(dumpedString, "\n"))
			},
		},
		{
			name: "Don't print destroyed cities",
			gameMap: func() *GameMap {
				m, _ := parser.ParseFile("test_resources/standard_input1.txt")
				m.GetExistCity("Qu-ux").Exists = false
				return m
			}(),
			validate: func(t *testing.T, dumpedString string) {
				assert.False(t, strings.Contains(dumpedString, "Qu-ux"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dumped := tt.gameMap.DumpMap()
			tt.validate(t, dumped)
		})
	}
}

func TestGameMap_CityCount(t *testing.T) {
	parser := StreamParser{}
	tests := []struct {
		name    string
		gameMap *GameMap
		want    int
	}{
		{
			name: "standard input map",
			gameMap: func() *GameMap {
				m, _ := parser.ParseFile("test_resources/standard_input1.txt")
				return m
			}(),
			want: 5,
		},
		{
			name: "destroyed one city",
			gameMap: func() *GameMap {
				m, _ := parser.ParseFile("test_resources/standard_input1.txt")
				m.GetExistCity("Qu-ux").Exists = false
				return m
			}(),
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.gameMap.ExistCityCount(), "ExistCityCount()")
		})
	}
}

func TestGameMap_AssignAliens(t *testing.T) {
	parser := StreamParser{}
	type args struct {
		gameMap    *GameMap
		alienCount int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Cities 5 and Aliens 5",
			args: args{
				gameMap: func() *GameMap {
					m, _ := parser.ParseFile("test_resources/standard_input1.txt")
					return m
				}(),
				alienCount: 5,
			},
			wantErr: false,
		},
		{
			name: "Cities 5 and Aliens 4",
			args: args{
				gameMap: func() *GameMap {
					m, _ := parser.ParseFile("test_resources/standard_input1.txt")
					return m
				}(),
				alienCount: 4,
			},
			wantErr: false,
		},
		{
			name: "Cities 5 and Aliens 6",
			args: args{
				gameMap: func() *GameMap {
					m, _ := parser.ParseFile("test_resources/standard_input1.txt")
					return m
				}(),
				alienCount: 6,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var aliens []*Alien
			for i := 0; i < tt.args.alienCount; i++ {
				aliens = append(aliens, NewAlien())
			}

			assert.Equal(t, tt.wantErr, nil != tt.args.gameMap.AssignAliens(aliens))
		})
	}
}
