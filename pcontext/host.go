package pcontext

import (
	"context"
)

var CtxFullPath = NewContextKey("full path url")
var CtxHeaderForward = NewContextKey("header forward")

func ExtractFullPath(ctx context.Context) string {
	if md, ok := ctx.Value(CtxHeaderForward).(map[string]string); ok {
		return md[RequestURL]
	}
	return ""
}

func ExtractQueryParams(ctx context.Context) string {
	if md, ok := ctx.Value(CtxHeaderForward).(map[string]string); ok {
		return md[RequestParams]
	}
	return ""
}
