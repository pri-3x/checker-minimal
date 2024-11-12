package checkers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alice/checkers"
	"github.com/alice/checkers/rules"
)

const (
	alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob   = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
)

func GetValidGame() checkers.StoredGame {
	return checkers.StoredGame{
		Board:      rules.New().String(),
		Turn:       "b",
		Black:      alice,
		Red:        bob,
		StartTime:  100, // Unix timestamp
		EndTime:    0,
		GameStatus: "ACTIVE",
	}
}

func TestValidateBasic(t *testing.T) {
	testCases := []struct {
		name    string
		game    checkers.StoredGame
		wantErr bool
	}{
		{
			name:    "valid game",
			game:    GetValidGame(),
			wantErr: false,
		},
		{
			name: "invalid black address",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "b",
				Black:      "invalid",
				Red:        bob,
				StartTime:  100,
				EndTime:    0,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "invalid red address",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "b",
				Black:      alice,
				Red:        "invalid",
				StartTime:  100,
				EndTime:    0,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "zero game start time",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "b",
				Black:      alice,
				Red:        bob,
				StartTime:  0,
				EndTime:    0,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "non-zero game end time",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "b",
				Black:      alice,
				Red:        bob,
				StartTime:  100,
				EndTime:    200,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "invalid game status",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "b",
				Black:      alice,
				Red:        bob,
				StartTime:  100,
				EndTime:    0,
				GameStatus: "INVALID",
			},
			wantErr: true,
		},
		{
			name: "invalid board string",
			game: checkers.StoredGame{
				Board:      "invalid board",
				Turn:       "b",
				Black:      alice,
				Red:        bob,
				StartTime:  100,
				EndTime:    0,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "invalid turn",
			game: checkers.StoredGame{
				Board:      rules.New().String(),
				Turn:       "invalid",
				Black:      alice,
				Red:        bob,
				StartTime:  100,
				EndTime:    0,
				GameStatus: "ACTIVE",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.game.Validate()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetBlackAddress(t *testing.T) {
	game := GetValidGame()
	black, err := game.GetBlackAddress()
	require.NoError(t, err)
	require.Equal(t, alice, black.String())
}

func TestGetRedAddress(t *testing.T) {
	game := GetValidGame()
	red, err := game.GetRedAddress()
	require.NoError(t, err)
	require.Equal(t, bob, red.String())
}

func TestParseGame(t *testing.T) {
	game := GetValidGame()
	parsed, err := game.ParseGame()
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, game.Board, parsed.String())
	require.Equal(t, rules.StringPieces[game.Turn].Player, parsed.Turn)
}
