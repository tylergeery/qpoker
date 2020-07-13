package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndRemoveFromTable(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: int64(i + 10)})
	}

	table := NewTable(players[:4])
	err := table.AddPlayer(players[4])
	assert.Error(t, err)

	err = table.RemovePlayer(players[4].ID)
	assert.Error(t, err)

	err = table.RemovePlayer(players[2].ID)
	assert.NoError(t, err)
	err = table.AddPlayer(players[4])
	assert.NoError(t, err)
}
