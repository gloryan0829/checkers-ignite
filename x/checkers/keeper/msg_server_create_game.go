package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("system info not found")
	}
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	newGame := rules.New()

	storedGame := types.StoredGame{
		Index:  newIndex,
		Board:  newGame.String(),
		Turn:   rules.PieceStrings[newGame.Turn],
		Black:  msg.Black,
		Red:    msg.Red,
		Winner: rules.PieceStrings[rules.NO_PLAYER],
	}

	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	k.Keeper.SetStoredGame(ctx, storedGame)

	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
		),
	)

	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
