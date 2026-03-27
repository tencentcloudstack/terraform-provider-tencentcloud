# Change: 优化 RabbitMQ 实例的 update 逻辑

## Why

当前 RabbitMQ 实例资源 `tencentcloud_tdmq_rabbitmq_vip_instance` 的 update 逻辑存在以下问题:

1. **过度限制的字段不可变性**: 当前代码中有 12 个字段被标记为不可变,但实际上部分字段应该是可以通过 API 更新的
2. **缺少对可变字段的支持**: 腾讯云 TDMQ RabbitMQ API 支持动态修改多个实例属性,但 Provider 未充分利用这些能力
3. **用户体验差**: 用户需要删除并重建实例才能更改某些配置,导致不必要的资源中断和数据迁移成本
4. **与 API 能力不匹配**: 腾讯云 API 提供了 `ModifyRabbitMQVipInstance` 接口,支持修改更多属性,但当前实现未充分利用

通过优化 update 逻辑,可以:
- 支持更多字段的动态更新,减少资源重建
- 提高用户体验和运维效率
- 降低因配置变更导致的服务中断风险
- 与底层 API 能力保持一致

## What Changes

### 问题分析

当前 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中 update 函数标记了以下字段为不可变:

```go
immutableArgs := []string{
    "zone_ids",              // ✅ 确实不可变 - 集群创建后不能修改可用区
    "vpc_id",                // ✅ 确实不可变 - 网络配置通常不可变
    "subnet_id",             // ✅ 确实不可变 - 网络配置通常不可变
    "node_spec",             // ❓ 可考虑支持 - API 可能支持规格变更
    "node_num",              // ❓ 可考虑支持 - API 可能支持节点数调整
    "storage_size",          // ❓ 可考虑支持 - API 可能支持存储扩容
    "enable_create_default_ha_mirror_queue", // ❓ 可考虑支持 - 镜像队列配置可能可改
    "auto_renew_flag",       // ❌ 应该可变 - 自动续费标志应该是可修改的
    "time_span",             // ✅ 确实不可变 - 购买时长不能直接修改
    "pay_mode",              // ✅ 确实不可变 - 付费模式转换可能有限制
    "cluster_version",       // ✅ 确实不可变 - 版本升级通常需要专用流程
    "band_width",            // ❓ 可考虑支持 - 带宽调整应该是可修改的
    "enable_public_access",  // ❌ 应该可变 - 公网访问开关应该是可修改的
}
```

### 优化方案

#### 1. 支持可变字段的更新

新增以下字段的 update 支持:

- `auto_renew_flag` - 自动续费标志
- `enable_public_access` - 公网访问开关
- `band_width` - 公网带宽

#### 2. 调用对应的 API 方法

使用腾讯云 API 提供的修改接口:
- `ModifyRabbitMQVipInstance` - 修改实例基本信息
- 可能需要调用 `ModifyInstances` - 修改实例规格/节点数/存储

#### 3. 添加状态等待机制

对于影响实例状态的更新操作(如规格变更),需要添加状态等待逻辑,确保更新完成后再返回。

### 修改文件

#### 主要修改: `/repo/tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`

**修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数**:

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId      = tccommon.GetLogId(tccommon.ContextNil)
        ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
        service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
        instanceId = d.Id()
    )

    // 真正不可变的字段
    immutableArgs := []string{
        "zone_ids", "vpc_id", "subnet_id",
        "time_span", "pay_mode", "cluster_version",
    }

    for _, v := range immutableArgs {
        if d.HasChange(v) {
            return fmt.Errorf("argument `%s` cannot be changed", v)
        }
    }

    // 处理可以立即更新的字段
    if d.HasChange("auto_renew_flag") || d.HasChange("enable_public_access") ||
       d.HasChange("band_width") || d.HasChange("cluster_name") ||
       d.HasChange("resource_tags") {

        request := tdmq.NewModifyRabbitMQVipInstanceRequest()
        request.InstanceId = &instanceId

        if d.HasChange("cluster_name") {
            if v, ok := d.GetOk("cluster_name"); ok {
                request.ClusterName = helper.String(v.(string))
            }
        }

        if d.HasChange("auto_renew_flag") {
            if v, ok := d.GetOkExists("auto_renew_flag"); ok {
                request.AutoRenewFlag = helper.Bool(v.(bool))
            }
        }

        if d.HasChange("band_width") {
            if v, ok := d.GetOkExists("band_width"); ok {
                request.Bandwidth = helper.IntUint64(v.(int))
            }
        }

        if d.HasChange("enable_public_access") {
            if v, ok := d.GetOkExists("enable_public_access"); ok {
                request.EnablePublicAccess = helper.Bool(v.(bool))
            }
        }

        if d.HasChange("resource_tags") {
            if v, ok := d.GetOk("resource_tags"); ok {
                for _, item := range v.([]interface{}) {
                    dMap := item.(map[string]interface{})
                    tag := tdmq.Tag{}
                    if v, ok := dMap["tag_key"]; ok {
                        tag.TagKey = helper.String(v.(string))
                    }
                    if v, ok := dMap["tag_value"]; ok {
                        tag.TagValue = helper.String(v.(string))
                    }
                    request.Tags = append(request.Tags, &tag)
                }
            } else {
                request.RemoveAllTags = helper.Bool(true)
            }
        }

        err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
            result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstance(request)
            if e != nil {
                return tccommon.RetryError(e)
            } else {
                log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
                    logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
            }
            return nil
        })

        if err != nil {
            log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
            return err
        }
    }

    // 处理需要异步更新的字段(规格、节点数、存储)
    // 注意: 这些更新可能需要调用不同的 API,并且需要等待实例状态稳定
    if d.HasChange("node_spec") || d.HasChange("node_num") || d.HasChange("storage_size") {
        // TODO: 调用规格变更 API,如果 API 支持的话
        // 目前保留原有行为,返回错误提示用户需要重建
        if d.HasChange("node_spec") {
            return fmt.Errorf("argument `node_spec` requires instance recreation, please use terraform taint")
        }
        if d.HasChange("node_num") {
            return fmt.Errorf("argument `node_num` requires instance recreation, please use terraform taint")
        }
        if d.HasChange("storage_size") {
            return fmt.Errorf("argument `storage_size` requires instance recreation, please use terraform taint")
        }
    }

    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}
```

#### Schema 调整

移除以下字段的 `ForceNew` 标记(如果有的话):
- `auto_renew_flag`
- `enable_public_access`
- `band_width`

### 测试文件更新

更新 `/repo/tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`:

新增测试用例:
- `TestAccTencentCloudTdmqRabbitmqVipInstance_updateAutoRenewFlag` - 测试自动续费标志更新
- `TestAccTencentCloudTdmqRabbitmqVipInstance_updatePublicAccess` - 测试公网访问开关更新
- `TestAccTencentCloudTdmqRabbitmqVipInstance_updateBandWidth` - 测试带宽更新

## Impact

### 受影响的规范
- 更新规范: `tdmq-rabbitmq-vip-instance` - RabbitMQ 实例资源

### 受影响的代码
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 核心修改
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` - 测试更新

### 向后兼容性
- ✅ 完全向后兼容,不破坏现有配置
- ✅ 保留真正不可变字段的限制
- ✅ 现有的 Terraform 配置无需修改即可继续工作

### 用户体验改善
- ✅ 支持动态更新自动续费标志,无需重建实例
- ✅ 支持动态切换公网访问开关,无需重建实例
- ✅ 支持动态调整带宽,无需重建实例
- ✅ 更清晰的错误提示,告知哪些字段需要重建

### 测试影响
- 需要验收测试环境中的 RabbitMQ 实例
- 测试需要验证新支持的可变字段能够正确更新
- 需要验证不可变字段仍然被正确限制

### 风险评估
- ⚠️ 低风险: 新增功能,不影响现有行为
- ⚠️ 需要验证: 确保腾讯云 API 确实支持这些字段的更新
- ⚠️ 需要测试: 验证更新操作不会影响实例稳定性

### 性能影响
- ✅ 无性能影响: 仅修改 update 逻辑,不改变 read/create/delete 流程
- ✅ 减少不必要的实例重建,整体提升性能
