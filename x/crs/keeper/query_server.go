package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/samricotta/vote/x/crs"
)

var _ crs.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) crs.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Decision implements crs.QueryServer.
func (qs queryServer) Decision(context.Context, *crs.QueryDecisionRequest) (*crs.QueryDecisionResponse, error) {
	panic("unimplemented")
}

// Params defines the handler for the Query/Params RPC method.
func (qs queryServer) Params(ctx context.Context, req *crs.QueryParamsRequest) (*crs.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &crs.QueryParamsResponse{Params: crs.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &crs.QueryParamsResponse{Params: params}, nil
}
