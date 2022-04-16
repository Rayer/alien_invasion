package alien_invastion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlien_Move(t *testing.T) {
	type fields struct {
		Number int
		Steps  int
		Alive  bool
	}
	type args struct {
		from *City
		to   *City
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantMoved bool
		wantStep  int
	}{
		{
			name: "Will move to other city with correct steps",
			fields: fields{
				Number: 0,
				Steps:  10,
				Alive:  true,
			},
			args: args{
				from: newCity("city1"),
				to:   newCity("city2"),
			},
			wantMoved: true,
			wantStep:  11,
		},
		{
			name: "Will not move due to next city is destroyed",
			fields: fields{
				Number: 0,
				Steps:  10,
				Alive:  true,
			},
			args: args{
				from: newCity("city1"),
				to: func() *City {
					ret := newCity("city2")
					ret.Exists = false
					return ret
				}(),
			},
			wantMoved: false,
			wantStep:  11,
		},
		{
			name: "Will not move to other city since alien is dead",
			fields: fields{
				Number: 0,
				Steps:  10,
				Alive:  false,
			},
			args: args{
				from: newCity("city1"),
				to:   newCity("city2"),
			},
			wantMoved: false,
			wantStep:  11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Alien{
				Number: tt.fields.Number,
				Steps:  tt.fields.Steps,
				Alive:  tt.fields.Alive,
			}
			//Connect cities together
			tt.args.from.Neighborhoods[North] = tt.args.to
			tt.args.to.Neighborhoods[South] = tt.args.from

			gotTo, gotStep := a.Move(tt.args.from)
			assert.Equalf(t, tt.wantMoved, gotTo == tt.args.to, "Move(%v)", tt.args.from)
			assert.Equalf(t, tt.wantStep, gotStep, "Move(%v)", tt.args.from)
		})
	}
}
