# Design: 优化 RabbitMQ 实例的 update 逻辑

## 概述

本设计文档详细说明如何优化 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑，使其支持更多参数的动态修改，实现用户对 RabbitMQ 实例的完整管理能力。

## 背景与问题

### 当前实现

当前实现中，`resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 方法将以下参数标记为不可修改（`immutableArgs`）：

```go
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num",
    "storage_size", "enable_create_default_ha_mirror_queue",
    "auto_renew_flag", "time_span", "pay_mode", "cluster_version",
    "band_width", "enable_public_access",
}
```

这意味着，任何尝试修改这些参数的操作都会直接返回错误：

```go
for _, v := range immutableArgs {
    if d.HasChange(v) {
        return fmt.Errorf("argument `%s` cannot be changed", v)
    }
}
```

### 问题分析

这种实现方式存在以下问题：

1. **功能限制**：腾讯云的 `ModifyRabbitMQVipInstance` API 实际上支持修改多个参数，但当前实现限制了用户使用这些功能
2. **用户体验差**：用户无法通过 Terraform 管理实例的生命周期，无法实现弹性伸缩、规格升级等自动化运维场景
3. **与 API 不一致**：Provider 的实现与底层 API 的能力不匹配

## 解决方案

### 总体思路

通过分析腾讯云 `ModifyRabbitMQVipInstance` API 的能力，将支持修改的参数从 `immutableArgs` 列表中移除，并在 Update 方法中实现相应的 API 调用逻辑。

### 参数分类

根据 API 支持情况，将参数分为以下几类：

#### 1. 可修改参数（从 `immutableArgs` 移除）

这些参数可以通过 `ModifyRabbitMQVipInstance` API 修改：

- `node_spec` - 节点规格
- `node_num` - 节点数量
- `storage_size` - 存储规格
- `band_width` - 带宽
- `enable_public_access` - 公网访问开关
- `auto_renew_flag` - 自动续费标识
- `cluster_version` - 集群版本

#### 2. 真正不可修改的参数（保留在 `immutableArgs`）

这些参数在实例创建后确实无法修改：

- `zone_ids` - 可用区
- `vpc_id` - VPC ID
- `subnet_id` - 子网 ID
- `enable_create_default_ha_mirror_queue` - 是否创建默认镜像队列
- `time_span` - 购买时长
- `pay_mode` - 付费模式

### 实现方案

#### 1. 修改 `immutableArgs` 列表

将可修改的参数从 `immutableArgs` 中移除：

```go
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id",
    "enable_create_default_ha_mirror_queue",
    "time_span", "pay_mode",
}
```

#### 2. 增强 `Update` 方法

为每个可修改的参数添加相应的 API 调用逻辑：

```go
func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId      = tccommon.GetLogId(tccommon.ContextNil)
        request    = tdmq.NewModifyRabbitMQVipInstanceRequest()
        instanceId = d.Id()
        needUpdate bool
    )

    // 真正不可修改的参数
    immutableArgs := []string{
        "zone_ids", "vpc_id", "subnet_id",
        "enable_create_default_ha_mirror_queue",
        "time_span", "pay_mode",
    }

    // 检查不可修改参数
    for _, v := range immutableArgs {
        if d.HasChange(v) {
            return fmt.Errorf("argument `%s` cannot be changed", v)
        }
    }

    request.InstanceId = &instanceId

    // 处理可修改参数
    if d.HasChange("node_spec") {
        if v, ok := d.GetOk("node_spec"); ok {
            request.NodeSpec = helper.String(v.(string))
            needUpdate = true
        }
    }

    if d.HasChange("node_num") {
        if v, ok := d.GetOkExists("node_num"); ok {
            request.NodeNum = helper.IntInt64(v.(int))
            needUpdate = true
        }
    }

    if d.HasChange("storage_size") {
        if v, ok := d.GetOkExists("storage_size"); ok {
            request.StorageSize = helper.IntInt64(v.(int))
            needUpdate = true
        }
    }

    if d.HasChange("band_width") {
        if v, ok := d.GetOkExists("band_width"); ok {
            request.Bandwidth = helper.IntUint64(v.(int))
            needUpdate = true
        }
    }

    if d.HasChange("enable_public_access") {
        if v, ok := d.GetOkExists("enable_public_access"); ok {
            request.EnablePublicAccess = helper.Bool(v.(bool))
            needUpdate = true
        }
    }

    if d.HasChange("auto_renew_flag") {
        if v, ok := d.GetOkExists("auto_renew_flag"); ok {
            request.AutoRenewFlag = helper.Bool(v.(bool))
            needUpdate = true
        }
    }

    if d.HasChange("cluster_version") {
        if v, ok := d.GetOk("cluster_version"); ok {
            request.ClusterVersion = helper.String(v.(string))
            needUpdate = true
        }
    }

    // 原有的 cluster_name 和 resource_tags 修改逻辑保持不变
    if d.HasChange("cluster_name") {
        if v, ok := d.GetOk("cluster_name"); ok {
            request.ClusterName = helper.String(v.(string))
            needUpdate = true
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
            needUpdate = true
        } else {
            request.RemoveAllTags = helper.Bool(true)
            needUpdate = true
        }
    }

    // 调用 API
    if needUpdate {
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

        // 对于某些需要等待异步操作完成的参数（如节点规格、节点数量等），添加等待逻辑
        if d.HasChange("node_spec") || d.HasChange("node_num") || d.HasChange("storage_size") {
            ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
            service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
            paramMap := make(map[string]interface{})
            tmpSet := make([]*tdmq.Filter, 0)
            filter := tdmq.Filter{}
            filter.Name = helper.String("instanceIds")
            filter.Values = helper.Strings([]string{instanceId})
            tmpSet = append(tmpSet, &filter)
            paramMap["filters"] = tmpSet

            err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
                result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
                if e != nil {
                    return tccommon.RetryError(e)
                }

                if result == nil {
                    return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s does not exist", instanceId))
                }

                if len(result) != 1 {
                    return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s id error", instanceId))
                }

                if *result[0].Status == svctdmq.RabbitMQVipInstanceRunning {
                    return resource.RetryableError(fmt.Errorf("rabbitmq_vip_instance status is updating"))
                } else if *result[0].Status == svctdmq.RabbitMQVipInstanceSuccess {
                    return nil
                } else {
                    return resource.NonRetryableError(fmt.Errorf("rabbitmq_vip_instance status illegal"))
                }
            })

            if err != nil {
                log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
                return err
            }
        }
    }

    return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}
```

#### 3. 添加 Timeouts 支持

为资源添加自定义超时配置支持，以便在涉及长时间运行的操作（如规格变更）时，用户可以自定义等待时间：

```go
Schema: map[string]*schema.Schema{
    // ... 其他字段 ...

    "timeouts": {
        Type:        schema.TypeList,
        MaxItems:    1,
        Optional:    true,
        Description: "Configuration timeouts for different operations.",
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "create": {
                    Type:        schema.TypeInt,
                    Optional:    true,
                    Default:     60,  // 默认 60 分钟
                    Description: "Timeout for creating the RabbitMQ instance.",
                },
                "update": {
                    Type:        schema.TypeInt,
                    Optional:    true,
                    Default:     60,  // 默认 60 分钟
                    Description: "Timeout for updating the RabbitMQ instance.",
                },
                "delete": {
                    Type:        schema.TypeInt,
                    Optional:    true,
                    Default:     30,  // 默认 30 分钟
                    Description: "Timeout for deleting the RabbitMQ instance.",
                },
            },
        },
    },
}
```

#### 4. 更新文档

更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 文档，明确说明哪些参数可以修改，哪些不能修改，并添加自定义超时配置的示例。

## 实施步骤

### Phase 1: 核心功能实现

1. **修改 Update 方法**：调整 `immutableArgs` 列表，并添加各参数的修改逻辑
2. **添加等待逻辑**：对于需要异步等待的操作（如节点规格、节点数量、存储规格变更），添加状态等待逻辑
3. **测试验证**：编写单元测试和集成测试，验证各种参数修改场景

### Phase 2: 增强功能

1. **添加 Timeouts 支持**：实现自定义超时配置功能
2. **更新文档**：更新资源文档和示例，说明新增功能和用法

### Phase 3: 优化和稳定性

1. **性能优化**：优化 API 调用和重试逻辑
2. **错误处理增强**：完善错误处理和用户提示信息
3. **日志和监控**：添加详细的日志输出，便于问题排查

## 风险与缓解

### 风险

1. **API 兼容性**：如果腾讯云 API 发生变更，可能导致某些参数无法修改
2. **异步操作等待**：某些修改操作可能是异步的，等待时间较长
3. **回滚机制**：如果修改失败，可能需要手动恢复

### 缓解措施

1. **API 测试**：充分测试 API 能力，确保支持修改的参数确实可修改
2. **合理超时**：为异步操作设置合理的超时时间和重试机制
3. **详细错误提示**：提供清晰的错误信息，帮助用户理解问题所在
4. **文档说明**：在文档中明确说明哪些参数可以修改，哪些不能修改

## 测试策略

### 单元测试

1. 测试各个参数的修改逻辑
2. 测试不可修改参数的错误处理
3. 测试异步操作的等待逻辑

### 集成测试

1. 测试完整的修改流程，从参数变更到状态更新
2. 测试多种参数组合的修改场景
3. 测试错误场景和回滚机制

## 性能考虑

1. **批量修改**：如果用户同时修改多个参数，应优化为单次 API 调用
2. **状态等待**：对于异步操作，应采用合理的轮询间隔和最大重试次数
3. **并发控制**：确保不会因为多个并发修改导致状态不一致

## 未来改进

1. **Drift Detection**：增强配置漂移检测能力
2. **Rollback Support**：支持修改失败时的自动回滚
3. **Preview Mode**：支持预览模式，让用户在实际修改前了解将发生的变化
