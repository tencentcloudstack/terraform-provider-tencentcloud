# Change: 优化 RabbitMQ 实例的 update 逻辑

## Why

当前 RabbitMQ VIP 实例 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 的 Update 逻辑存在以下问题：

1. **过度限制可更新字段**：大量本应支持更新的字段被标记为不可变（immutable），包括：
   - `zone_ids` - 可用区配置
   - `vpc_id`, `subnet_id` - 网络配置
   - `node_spec` - 节点规格
   - `node_num` - 节点数量
   - `storage_size` - 存储大小
   - `cluster_version` - 集群版本
   - `band_width` - 带宽
   - `enable_public_access` - 公网访问开关
   - `auto_renew_flag` - 自动续费标志
   - `time_span` - 购买时长
   - `pay_mode` - 付费模式
   - `enable_create_default_ha_mirror_queue` - 镜像队列开关

2. **用户体验差**：用户无法通过 Terraform 声明式地调整实例配置，必须手动在控制台修改后再重新导入，违背了基础设施即代码的初衷。

3. **与腾讯云 API 能力不匹配**：腾讯云 API (`ModifyRabbitMQVipInstance`) 支持修改多个配置项，但当前实现未充分利用这些能力。

腾讯云提供的 `ModifyRabbitMQVipInstance` API 支持以下参数：
- `InstanceId` - 实例 ID
- `ClusterName` - 集群名称
- `Tags` - 资源标签
- `RemoveAllTags` - 移除所有标签

此外，还有其他专门的 API 用于修改特定配置：
- `ModifyRabbitMQVipInstanceSpec` - 修改规格（节点数量、存储等）
- `ModifyRabbitMQVipInstanceNetwork` - 修改网络配置
- `ModifyRabbitMQVipInstancePublicAccess` - 修改公网访问

通过优化 Update 逻辑，用户可以：
- 通过 Terraform 声明式地修改实例配置
- 实现实例的弹性扩缩容
- 简化运维流程，减少手动操作
- 提高基础设施配置的一致性和可追溯性

## What Changes

优化 RabbitMQ VIP 实例的 Update 逻辑，支持更多字段的更新，并调用相应的 API 实现配置修改。

### 修改文件
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 优化 Update 函数逻辑

### 具体变更

#### 1. 移除不必要的 immutable 限制

当前 Update 函数中标记为不可变的字段，实际上部分可以通过 API 修改。需要根据腾讯云 API 能力调整：

**可以更新的字段**：
- `cluster_name` - 通过 `ModifyRabbitMQVipInstance` API 修改 ✓（已实现）
- `resource_tags` - 通过 `ModifyRabbitMQVipInstance` API 修改 ✓（已实现）
- `node_spec` - 通过 `ModifyRabbitMQVipInstanceSpec` API 修改（需新增）
- `node_num` - 通过 `ModifyRabbitMQVipInstanceSpec` API 修改（需新增）
- `storage_size` - 通过 `ModifyRabbitMQVipInstanceSpec` API 修改（需新增）
- `enable_public_access` - 通过 `ModifyRabbitMQVipInstancePublicAccess` API 修改（需新增）
- `band_width` - 通过 `ModifyRabbitMQVipInstancePublicAccess` API 修改（需新增）

**仍然不可变的字段**：
- `zone_ids` - 可用区在创建后不可修改
- `vpc_id`, `subnet_id` - VPC 和子网在创建后不可修改
- `cluster_version` - 集群版本需要升级操作，不属于普通 Update 范畴
- `pay_mode` - 付费模式在创建后不可修改
- `time_span` - 购买时长在创建后不可修改
- `enable_create_default_ha_mirror_queue` - 镜像队列设置在创建后不可修改
- `auto_renew_flag` - 自动续费标志可能通过单独 API 修改，需要进一步确认

#### 2. 增强 Update 函数实现

在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中：

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
		"zone_ids", "vpc_id", "subnet_id", "cluster_version",
		"pay_mode", "time_span", "enable_create_default_ha_mirror_queue",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	// 1. 修改基本配置（名称、标签）
	if d.HasChange("cluster_name") || d.HasChange("resource_tags") {
		request := tdmq.NewModifyRabbitMQVipInstanceRequest()
		request.InstanceId = &instanceId

		if d.HasChange("cluster_name") {
			if v, ok := d.GetOk("cluster_name"); ok {
				request.ClusterName = helper.String(v.(string))
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
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance basic config failed, reason:%+v", logId, err)
			return err
		}
	}

	// 2. 修改规格（节点数量、存储大小）
	if d.HasChange("node_spec") || d.HasChange("node_num") || d.HasChange("storage_size") {
		request := tdmq.NewModifyRabbitMQVipInstanceSpecRequest()
		request.InstanceId = &instanceId

		if d.HasChange("node_spec") {
			if v, ok := d.GetOk("node_spec"); ok {
				request.NodeSpec = helper.String(v.(string))
			}
		}

		if d.HasChange("node_num") {
			if v, ok := d.GetOk("node_num"); ok {
				request.NodeNum = helper.IntInt64(v.(int))
			}
		}

		if d.HasChange("storage_size") {
			if v, ok := d.GetOk("storage_size"); ok {
				request.StorageSize = helper.IntInt64(v.(int))
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstanceSpec(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance spec failed, reason:%+v", logId, err)
			return err
		}

		// 等待规格变更完成
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
			log.Printf("[CRITAL]%s wait for spec update failed, reason:%+v", logId, err)
			return err
		}
	}

	// 3. 修改公网访问配置
	if d.HasChange("enable_public_access") || d.HasChange("band_width") {
		request := tdmq.NewModifyRabbitMQVipInstancePublicAccessRequest()
		request.InstanceId = &instanceId

		if d.HasChange("enable_public_access") {
			if v, ok := d.GetOk("enable_public_access"); ok {
				request.EnablePublicAccess = helper.Bool(v.(bool))
			}
		}

		if d.HasChange("band_width") {
			if v, ok := d.GetOk("band_width"); ok {
				request.Bandwidth = helper.IntUint64(v.(int))
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstancePublicAccess(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance public access failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}
```

#### 3. 更新 Schema 中的 ForceNew 标记

需要从以下字段的 Schema 中移除 `ForceNew: true` 标记：
- `node_spec`
- `node_num`
- `storage_size`
- `enable_public_access`
- `band_width`

#### 4. 更新文档和测试

- 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 文档，明确哪些字段可以更新
- 更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`，增加 Update 相关的测试用例
- 更新网站文档 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`

## Impact

### 受影响的规范
- 修改规范：`tdmq-rabbitmq-vip-instance` - RabbitMQ 实例资源

### 受影响的代码
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 优化 Update 函数

### 向后兼容性
- ✅ 完全向后兼容，不会破坏现有配置
- ✅ 只是解锁了更多可更新字段，现有配置行为不变
- ⚠️ 之前因字段变更导致资源被重建的行为，现在会改为原地更新

### 依赖关系
- 无新增依赖
- 使用已有的腾讯云 API：`ModifyRabbitMQVipInstance`, `ModifyRabbitMQVipInstanceSpec`, `ModifyRabbitMQVipInstancePublicAccess`

### 测试影响
- 需要验收测试环境中的 RabbitMQ 实例
- 新增测试用例需要验证：
  - 修改节点数量
  - 修改存储大小
  - 修改节点规格
  - 修改公网访问配置
  - 修改带宽
  - 多字段同时修改

### 性能影响
- Update 操作可能需要更长时间（规格变更需要等待实例状态）
- 增加了状态检查的重试逻辑，但这是必需的以确保配置变更完成

## Migration Guide

### 对于已有 Terraform 配置

如果用户之前因为无法更新某些字段而删除并重新创建资源，现在可以直接更新配置：

**之前（需要删除重建）**：
```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... 其他配置
  node_num = 3
}

# 想要修改节点数量，必须先删除再创建
```

**现在（原地更新）**：
```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... 其他配置
  node_num = 5  # 直接修改即可，Terraform 会调用 API 更新
}
```

### 注意事项

1. 规格变更（node_num, storage_size）可能需要较长时间（几分钟到十几分钟），请耐心等待
2. 某些变更可能产生额外费用（如增加存储、带宽），请注意成本控制
3. 更新过程中实例状态会变为 "Running"，此期间部分功能可能受限
4. 确保有足够的配额和资源支持变更操作

## Implementation Notes

### API 调用顺序

1. 先调用 `ModifyRabbitMQVipInstance` 修改名称和标签（快速）
2. 再调用 `ModifyRabbitMQVipInstanceSpec` 修改规格（需要等待）
3. 最后调用 `ModifyRabbitMQVipInstancePublicAccess` 修改公网访问（快速）

### 错误处理

- 如果某个 API 调用失败，应该立即返回错误，不继续执行后续更新
- 使用 `resource.Retry` 处理临时性错误和 API 限流
- 规格变更后必须等待实例状态变为 "Success" 才能返回

### 状态管理

- 使用 `defer tccommon.LogElapsed()` 记录执行时间
- 使用 `defer tccommon.InconsistentCheck()` 检查状态一致性
- Update 后必须调用 Read 函数刷新状态

### 日志记录

- 在每个 API 调用时记录请求和响应（DEBUG 级别）
- 在错误发生时记录错误信息（CRITICAL 级别）
- 使用统一的 logId 追踪请求链路
