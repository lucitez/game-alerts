package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	cases := []struct {
		name     string
		location string
		validate func(field string, err error)
	}{
		{
			name:     "valid location",
			location: "SM Airport Field - WEST",
			validate: func(field string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, field, "west")
			},
		},
		{
			name:     "invalid location",
			location: "foo",
			validate: func(field string, err error) {
				assert.Error(t, err)
				assert.Zero(t, field)
			},
		},
		{
			name: "no location",
			validate: func(field string, err error) {
				assert.Error(t, err)
				assert.Zero(t, field)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			game := Game{
				Location: tc.location,
			}
			tc.validate(game.Field())
		})
	}
}
