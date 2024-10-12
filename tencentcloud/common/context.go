package common

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ctxResourceDataKey struct{}
type ctxProviderMetaKey struct{}
type ctxDataKey struct{}

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
	ctx = context.WithValue(ctx, ctxDataKey{}, &ContextData{})
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

// DataFromContext 从上下文获取 data
func DataFromContext(ctx context.Context) *ContextData {
	if data, ok := ctx.Value(ctxDataKey{}).(*ContextData); ok {
		return data
	}
	return nil
}

// ContextData 上下文临时数据
type ContextData struct {
	lock sync.RWMutex
	data map[string]interface{}
}

// Set 设置值
func (d *ContextData) Set(key string, v interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.data == nil {
		d.data = make(map[string]interface{})
	}
	d.data[key] = v
}

// Delete 删除值
func (d *ContextData) Delete(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.data, key)
}

// Get 获取键
func (d *ContextData) Get(key string) interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.data[key]
}

// GetInt 获取 int 数据键
func (d *ContextData) GetInt(key string) (ret int, ok bool) {
	ret, ok = d.Get(key).(int)
	return
}

// GetUInt 获取 uint 数据键
func (d *ContextData) GetUInt(key string) (ret uint, ok bool) {
	ret, ok = d.Get(key).(uint)
	return
}

// GetInt64 获取 int64 数据键
func (d *ContextData) GetInt64(key string) (ret int64, ok bool) {
	ret, ok = d.Get(key).(int64)
	return
}

// GetUInt64 获取 uint64 数据键
func (d *ContextData) GetUInt64(key string) (ret uint64, ok bool) {
	ret, ok = d.Get(key).(uint64)
	return
}

// GetString 获取 string 数据键
func (d *ContextData) GetString(key string) (ret string, ok bool) {
	ret, ok = d.Get(key).(string)
	return
}

// GetBool 获取 bool 数据键
func (d *ContextData) GetBool(key string) (ret bool, ok bool) {
	ret, ok = d.Get(key).(bool)
	return
}
