package auth

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"qpoker/models"
)

func TestCreatingAndExtractingToken(t *testing.T) {
	player := &models.Player{ID: 444}
	token, err := CreateToken(player)

	assert.Nil(t, err)

	claims, err := ExtractToken(token)
	assert.Nil(t, err)

	val, ok := claims["player_id"]
	assert.True(t, ok)
	assert.Equal(t, player.ID, int64(val.(float64)), "Could not get player_id from claims")

	nano, ok := claims["nbf"]
	assert.True(t, ok)

	ts := int64(nano.(float64))
	assert.Greater(t, ts, time.Now().Unix()-100)
	assert.GreaterOrEqual(t, time.Now().Unix(), ts)

	playerIDFromClaims, err := GetPlayerIDFromAccessToken(token)
	assert.Nil(t, err)
	assert.Equal(t, playerIDFromClaims, player.ID)
}

func TestInvalidClaims(t *testing.T) {
	_, err := GetPlayerIDFromAccessToken("invalid")
	assert.Error(t, err)

	dur, _ := time.ParseDuration("2m")
	claims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(dur).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(signingKey)

	extracted, _ := ExtractToken(tokenString)

	_, ok := extracted["player_id"]
	assert.False(t, ok)

	_, err = GetPlayerIDFromAccessToken(tokenString)
	assert.Error(t, err)
}
