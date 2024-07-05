package expectedkeepers

import (
	"context"

	"github.com/samricotta/vote/x/crs"
)

// type BankKeeper interface {
// 	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
// 	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
// }

type CrsKeeper interface {
	CreateDecision(ctx context.Context, decision crs.Decision) error
}
