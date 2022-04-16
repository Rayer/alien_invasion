package alien_invastion

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileStreamParser_ParseFile(t *testing.T) {
	f := &StreamParser{}
	got, err := f.ParseFile("test_resources/standard_input1.txt")
	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}
	assert.True(t, got != nil)
}

func TestStreamParser_ParseFile(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name          string
		args          args
		wantSize      int
		wantErrorSize int
	}{
		{
			name: "Happy Path",
			args: args{
				filepath: "test_resources/standard_input1.txt",
			},
			wantSize:      5,
			wantErrorSize: 0,
		},
		{
			name: "File not found",
			args: args{
				filepath: "test_resources/not_found.txt",
			},
			wantSize:      0,
			wantErrorSize: 1,
		},
		{
			name: "File with errors",
			args: args{
				filepath: "test_resources/standard_input_err.txt",
			},
			wantSize:      5,
			wantErrorSize: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamParser{}
			gotRet, gotErrors := s.ParseFile(tt.args.filepath)
			assert.Equalf(t, tt.wantErrorSize, len(gotErrors), "ParseFile(%v)", tt.args.filepath)
			var cities int
			if gotRet != nil {
				cities = len(gotRet.cities)
			} else {
				cities = 0
			}
			assert.Equalf(t, tt.wantSize, cities, "ParseFile(%v)", tt.args.filepath)
		})
	}
}

//go:embed test_resources/standard_input1.txt
var happyPathString string

//go:embed test_resources/standard_input_err.txt
var errMapString string

func TestStreamParser_ParseString(t *testing.T) {
	type args struct {
		str string
	}

	tests := []struct {
		name          string
		args          args
		wantSize      int
		wantErrorSize int
	}{
		{
			name: "Happy Path",
			args: args{
				str: happyPathString,
			},
			wantSize:      5,
			wantErrorSize: 0,
		},
		{
			name: "File with errors",
			args: args{
				str: errMapString,
			},
			wantSize:      5,
			wantErrorSize: 1,
		},
		{
			name: "Empty string",
			args: args{
				str: "",
			},
			wantSize:      0,
			wantErrorSize: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamParser{}
			gotRet, gotErrors := s.ParseString(tt.args.str)
			assert.Equalf(t, tt.wantErrorSize, len(gotErrors), "ParseString(%v)", tt.args.str)
			var cities int
			if gotRet != nil {
				cities = len(gotRet.cities)
			} else {
				cities = 0
			}
			assert.Equalf(t, tt.wantSize, cities, "ParseString(%v)", tt.args.str)
		})
	}
}

func TestStreamParser_parseSingleLine(t *testing.T) {
	type args struct {
		line    string
		gameMap *GameMap
	}
	tests := []struct {
		name          string
		args          args
		wantSize      int
		wantErrorSize int
	}{
		{
			name: "Happy Path with empty map",
			args: args{
				line:    "Foo north=Bar west=Baz south=Qu-ux",
				gameMap: NewGameMap(),
			},
			wantSize:      4,
			wantErrorSize: 0,
		},
		{
			name: "Happy Path with non-empty map",
			args: args{
				line:    "Foo north=Bar west=Baz south=Qu-ux",
				gameMap: func() *GameMap { m := NewGameMap(); m.UpsertCity("Baz"); return m }(),
			},
			wantSize:      4,
			wantErrorSize: 0,
		},
		{
			name: "Happy Path with invalid pair",
			args: args{
				line:    "Foo north=Bar west=Baz south=Qu-ux east=EAC=aina",
				gameMap: func() *GameMap { m := NewGameMap(); m.UpsertCity("Baz"); return m }(),
			},
			wantSize:      4,
			wantErrorSize: 0,
		},
		{
			name: "Conflict path",
			args: args{
				line: "Foo north=Bar west=Baz south=Qu-ux",
				gameMap: func() *GameMap {
					m := NewGameMap()
					// Foo east -> Bar, but Bar west -> Baz, so this should fail
					_ = m.UpdateCityWithNeighborhood("Baz", East, "Bar")
					return m
				}(),
			},
			wantSize:      4,
			wantErrorSize: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamParser{}
			gotErrors := s.parseSingleLine(tt.args.line, tt.args.gameMap)
			assert.Equalf(t, tt.wantErrorSize, len(gotErrors), "ParseLine(%v)", tt.args.line)
			assert.Equalf(t, tt.wantSize, len(tt.args.gameMap.cities), "ParseLine(%v)", tt.args.line)
		})
	}
}
