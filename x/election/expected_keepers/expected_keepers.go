package expectedkeepers

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/samricotta/vote/x/crs/keeper"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type CrsKeeper interface {
	keeper.Keeper
}
