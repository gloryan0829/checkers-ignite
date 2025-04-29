package keeper_test

import (
	"context"
	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/testutil/sample"
	"github.com/alice/checkers/x/checkers/keeper"
	checkers "github.com/alice/checkers/x/checkers/module"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMsgServerWithOneGameForPlayMove(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	return server, k, context
}

func TestPlayMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgPlayMoveResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    "*",
	}, *playMoveResponse)
}

func TestPlayMoveCannotParseGame(t *testing.T) {
	msgServer, k, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	storedGame, _ := k.GetStoredGame(ctx, "1")
	storedGame.Board = "not a board"
	k.SetStoredGame(ctx, storedGame)
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, r, "game cannot be parsed: invalid board string: not a board")
	}()
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
}

func TestMsgPlayMove_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgPlayMove
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgPlayMove{
				Creator:   "invalid_address",
				GameIndex: "5",
				FromX:     0,
				FromY:     5,
				ToX:       1,
				ToY:       4,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid game index",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "invalid_index",
				FromX:     0,
				FromY:     5,
				ToX:       1,
				ToY:       4,
			},
			err: types.ErrInvalidGameIndex,
		},
		{
			name: "invalid fromX too high",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     rules.BOARD_DIM,
				FromY:     5,
				ToX:       1,
				ToY:       4,
			},
			err: types.ErrInvalidPositionIndex,
		},
		{
			name: "invalid fromY too high",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     0,
				FromY:     rules.BOARD_DIM,
				ToX:       1,
				ToY:       4,
			},
			err: types.ErrInvalidPositionIndex,
		},
		{
			name: "invalid toX too high",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     0,
				FromY:     5,
				ToX:       rules.BOARD_DIM,
				ToY:       4,
			},
			err: types.ErrInvalidPositionIndex,
		},
		{
			name: "invalid toY too high",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     0,
				FromY:     5,
				ToX:       1,
				ToY:       rules.BOARD_DIM,
			},
			err: types.ErrInvalidPositionIndex,
		},
		{
			name: "invalid no move",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     0,
				FromY:     5,
				ToX:       0,
				ToY:       5,
			},
			err: types.ErrMoveAbsent,
		},
		{
			name: "valid address",
			msg: types.MsgPlayMove{
				Creator:   sample.AccAddress(),
				GameIndex: "5",
				FromX:     0,
				FromY:     5,
				ToX:       1,
				ToY:       4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
