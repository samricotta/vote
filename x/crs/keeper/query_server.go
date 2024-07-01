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

// Counter defines the handler for the Query/Counter RPC method.
// func (qs queryServer) Counter(ctx context.Context, req *crs.QueryCounterRequest) (*crs.QueryCounterResponse, error) {
// 	if _, err := qs.k.addressCodec.StringToBytes(req.Address); err != nil {
// 		return nil, fmt.Errorf("invalid sender address: %w", err)
// 	}

// 	counter, err := qs.k.Counter.Get(ctx, req.Address)
// 	if err != nil {
// 		if errors.Is(err, collections.ErrNotFound) {
// 			return &crs.QueryCounterResponse{Counter: 0}, nil
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &crs.QueryCounterResponse{Counter: counter}, nil
// }

// Counters defines the handler for the Query/Counters RPC method.
// func (qs queryServer) Counters(ctx context.Context, req *crs.QueryCountersRequest) (*crs.QueryCountersResponse, error) {
// 	counters, pageRes, err := query.CollectionPaginate(
// 		ctx,
// 		qs.k.Counter,
// 		req.Pagination,
// 		func(key string, value uint64) (*crs.Counter, error) {
// 			return &crs.Counter{
// 				Address: key,
// 				Count:   value,
// 			}, nil
// 		})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &crs.QueryCountersResponse{Counters: counters, Pagination: pageRes}, nil
// }

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
