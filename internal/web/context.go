package web

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/internal"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/utils"
)

// * Sama kayak req.user kalo di express
func GetUserFromContext(ctx context.Context) (*utils.JWTClaims, bool) {
	claims, ok := ctx.Value(internal.UserClaimsKey).(*utils.JWTClaims)
	return claims, ok
}
