package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/testutil/nullify"
	"github.com/alice/checkers/x/checkers/types"
)

func TestSystemInfoQuery(t *testing.T) {
	keeper, ctx := keepertest.CheckersKeeper(t)
	item := createTestSystemInfo(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetSystemInfoRequest
		response *types.QueryGetSystemInfoResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSystemInfoRequest{},
			response: &types.QueryGetSystemInfoResponse{SystemInfo: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.SystemInfo(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
