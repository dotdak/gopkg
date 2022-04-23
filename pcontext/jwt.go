package pcontext

import (
	"context"
)

var CtxInternalUser = NewContextKey("internal user")

func GetInternalUser(ctx context.Context) map[string]string {
	if md, ok := ctx.Value(CtxInternalUser).(map[string]string); ok {
		return md
	}
	return make(map[string]string)
}
