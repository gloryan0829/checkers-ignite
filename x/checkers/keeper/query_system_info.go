package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SystemInfo(goCtx context.Context, req *types.QueryGetSystemInfoRequest) (*types.QueryGetSystemInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetSystemInfo(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found ::: ")
	}

	return &types.QueryGetSystemInfoResponse{SystemInfo: val}, nil
}
