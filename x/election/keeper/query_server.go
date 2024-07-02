package keeper

import (
	// "context"
	// "errors"

	// "cosmossdk.io/collections"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"github.com/samricotta/vote/x/election"
)

var _ election.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) election.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// // Counter defines the handler for the Query/Counter RPC method.
// func (qs queryServer) Counter(ctx context.Context, req *election.QueryCountersRequest) (*election.QueryCounterResponse, error) {
// 	if _, err := qs.k.addressCodec.StringToBytes(req.Address); err != nil {
// 		return nil, fmt.Errorf("invalid sender address: %w", err)
// 	}

// 	counter, err := qs.k.Counters.Get(ctx, req.Address)
// 	if err != nil {
// 		if errors.Is(err, collections.ErrNotFound) {
// 			return &election.QueryCountersRequest{Counter: 0}, nil
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &election.QueryCounterResponse{Counter: counter}, nil
// }

// // Counters defines the handler for the Query/Counters RPC method.
// func (qs queryServer) Counters(ctx context.Context, req *election.QueryCountersRequest) (*election.QueryCountersResponse, error) {
// 	counters, pageRes, err := query.CollectionPaginate(
// 		ctx,
// 		qs.k.Counter,
// 		req.Pagination,
// 		func(key string, value uint64) (*election.Counter, error) {
// 			return &election.Counter{
// 				Address: key,
// 				Count:   value,
// 			}, nil
// 		})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &election.QueryCountersResponse{Counters: counters, Pagination: pageRes}, nil
// }

// Params defines the handler for the Query/Params RPC method.
// func (qs queryServer) Params(ctx context.Context, req *election.QueryParamsRequest) (*election.QueryParamsResponse, error) {
// 	params, err := qs.k.Params.Get(ctx)
// 	if err != nil {
// 		if errors.Is(err, collections.ErrNotFound) {
// 			return &election.QueryParamsResponse{Params: election.Params{}}, nil
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &election.QueryParamsResponse{Params: params}, nil
// }
