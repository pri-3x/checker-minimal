package checkers

import (
	fmt "fmt"

	"cosmossdk.io/errors"
	"github.com/alice/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const GameStatusActive = "ACTIVE"

func (storedGame *StoredGame) GetBlackAddress() (black sdk.AccAddress, err error) {
	black, errBlack := sdk.AccAddressFromBech32(storedGame.Black)
	return black, errors.Wrapf(errBlack, ErrInvalidBlack.Error(), storedGame.Black)
}

func (storedGame *StoredGame) GetRedAddress() (red sdk.AccAddress, err error) {
	red, errRed := sdk.AccAddressFromBech32(storedGame.Red)
	return red, errors.Wrapf(errRed, ErrInvalidRed.Error(), storedGame.Red)
}

func (storedGame *StoredGame) ParseGame() (game *rules.Game, err error) {
	board, errBoard := rules.Parse(storedGame.Board)
	if errBoard != nil {
		return nil, errors.Wrapf(errBoard, ErrGameNotParseable.Error())
	}
	board.Turn = rules.StringPieces[storedGame.Turn].Player
	if board.Turn.Color == "" {
		return nil, errors.Wrapf(fmt.Errorf("turn: %s", storedGame.Turn), ErrGameNotParseable.Error())
	}
	return board, nil
}

// Validate performs basic validation of a StoredGame
func (storedGame StoredGame) Validate() error {
	// Validate board and game rules
	_, err := storedGame.ParseGame()
	if err != nil {
		return err
	}

	// Validate black player address
	_, err = storedGame.GetBlackAddress()
	if err != nil {
		return err
	}

	// Validate red player address
	_, err = storedGame.GetRedAddress()
	if err != nil {
		return err
	}

	// Validate game start time
	if storedGame.StartTime == 0 {
		return errors.Wrapf(ErrInvalidGame, "game start time cannot be zero")
	}

	// Validate game status
	if storedGame.GameStatus != GameStatusActive {
		return errors.Wrapf(ErrInvalidGame, "game status must be ACTIVE")
	}

	// Validate game end time is zero (since we're not implementing game completion)
	if storedGame.EndTime != 0 {
		return errors.Wrapf(ErrInvalidGame, "game end time must be zero")
	}

	return nil
}
