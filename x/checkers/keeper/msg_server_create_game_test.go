package keeper_test

import (
	"github.com/alice/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateGame(t *testing.T) {
	_, msgServer, context := setupMsgServer(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})

	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "",
	}, *createResponse)
}
