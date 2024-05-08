package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ctxResourceDataKey struct{}
type ctxProviderMetaKey struct{}

// NewResourceLifeCycleHandleFuncContext 创建一个资源生命周期处理方法上下文
func NewResourceLifeCycleHandleFuncContext(
	parent context.Context,
	logID string,
	d *schema.ResourceData,
	meta interface{},
) context.Context {
	ctx := context.WithValue(parent, LogIdKey, logID)
	ctx = context.WithValue(ctx, ctxResourceDataKey{}, d)
	ctx = context.WithValue(ctx, ctxProviderMetaKey{}, meta)
	return ctx
}

// ResourceDataFromContext 从上下文获取资源数据
func ResourceDataFromContext(ctx context.Context) *schema.ResourceData {
	if d, ok := ctx.Value(ctxResourceDataKey{}).(*schema.ResourceData); ok {
		return d
	}
	return nil
}

// ProviderMetaFromContext 从上下文获取 provider meta
func ProviderMetaFromContext(ctx context.Context) interface{} {
	if meta, ok := ctx.Value(ctxProviderMetaKey{}).(ProviderMeta); ok {
		return meta
	}
	return nil
}
