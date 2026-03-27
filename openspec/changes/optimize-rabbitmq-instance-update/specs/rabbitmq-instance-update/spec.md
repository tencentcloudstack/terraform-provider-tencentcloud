# RabbitMQ 实例 Update 逻辑优化规范

## 概述

本规范定义了 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑优化方案,旨在支持更多字段的动态更新,提升用户体验。

## 字段分类

### 不可变字段(Immutable Fields)

以下字段在实例创建后不可修改,修改这些字段会触发错误:

| 字段名 | 类型 | 原因 |
|--------|------|------|
| `zone_ids` | `[]int` | 集群可用区在创建后无法更改 |
| `vpc_id` | `string` | 私有网络配置在创建后无法更改 |
| `subnet_id` | `string` | 子网配置在创建后无法更改 |
| `time_span` | `int` | 购买时长无法直接修改(可能需要续费操作) |
| `pay_mode` | `int` | 付费模式转换有特殊限制和流程 |
| `cluster_version` | `string` | 集群版本升级需要专用的升级流程 |

### 可变字段(Mutable Fields)

以下字段支持动态更新,修改这些字段会调用 `ModifyRabbitMQVipInstance` API:

| 字段名 | 类型 | API 字段 | 说明 |
|--------|------|----------|------|
| `cluster_name` | `string` | `ClusterName` | 集群名称 |
| `auto_renew_flag` | `bool` | `AutoRenewFlag` | 自动续费标志 |
| `enable_public_access` | `bool` | `EnablePublicAccess` | 公网访问开关 |
| `band_width` | `int` | `Bandwidth` | 公网带宽(Mbps) |
| `resource_tags` | `list` | `Tags`/`RemoveAllTags` | 资源标签 |

### 条件可变字段(Conditionally Mutable Fields)

以下字段目前暂不支持动态更新,但可能需要进一步评估 API 能力:

| 字段名 | 类型 | 当前状态 | 建议 |
|--------|------|----------|------|
| `node_spec` | `string` | ❌ 不可变 | 评估 API 是否支持规格变更 |
| `node_num` | `int` | ❌ 不可变 | 评估 API 是否支持节点数调整 |
| `storage_size` | `int` | ❌ 不可变 | 评估 API 是否支持存储扩容 |
| `enable_create_default_ha_mirror_queue` | `bool` | ❌ 不可变 | 评估是否支持镜像队列配置修改 |

## Update 流程

### 1. 参数验证

```go
// 验证不可变字段
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id",
    "time_span", "pay_mode", "cluster_version",
}

for _, v := range immutableArgs {
    if d.HasChange(v) {
        return fmt.Errorf("argument `%s` cannot be changed", v)
    }
}
```

### 2. 构建更新请求

检测可变字段的变化,构建 `ModifyRabbitMQVipInstanceRequest`:

```go
request := tdmq.NewModifyRabbitMQVipInstanceRequest()
request.InstanceId = &instanceId

// 设置需要更新的字段
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

if d.HasChange("enable_public_access") {
    if v, ok := d.GetOkExists("enable_public_access"); ok {
        request.EnablePublicAccess = helper.Bool(v.(bool))
    }
}

if d.HasChange("band_width") {
    if v, ok := d.GetOkExists("band_width"); ok {
        request.Bandwidth = helper.IntUint64(v.(int))
    }
}

if d.HasChange("resource_tags") {
    if v, ok := d.GetOk("resource_tags"); ok {
        // 设置新标签
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
        // 清除所有标签
        request.RemoveAllTags = helper.Bool(true)
    }
}
```

### 3. 执行 API 调用

使用重试机制执行 API 调用:

```go
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
```

### 4. 刷新资源状态

调用 Read 函数刷新资源状态:

```go
return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
```

## 错误处理

### 不可变字段错误

当用户尝试修改不可变字段时,返回清晰的错误信息:

```go
return fmt.Errorf("argument `%s` cannot be changed. To update this field, you need to recreate the instance using `terraform taint` or `terraform apply -replace`.", fieldName)
```

### API 调用失败错误

捕获 API 调用失败并返回详细的错误信息:

```go
if err != nil {
    log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
    return fmt.Errorf("failed to update RabbitMQ instance %s: %v", instanceId, err)
}
```

## 测试场景

### 1. 更新自动续费标志

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  cluster_name       = "test-rabbitmq"
  zone_ids           = [1]
  vpc_id             = "vpc-xxxxxx"
  subnet_id          = "subnet-xxxxxx"
  auto_renew_flag    = true
  # ... 其他配置
}

# 更新为不自动续费
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  cluster_name       = "test-rabbitmq"
  zone_ids           = [1]
  vpc_id             = "vpc-xxxxxx"
  subnet_id          = "subnet-xxxxxx"
  auto_renew_flag    = false
  # ... 其他配置
}
```

### 2. 切换公网访问

```hcl
# 启用公网访问
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  enable_public_access = true
  band_width          = 10
}

# 禁用公网访问
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  enable_public_access = false
}
```

### 3. 调整带宽

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  enable_public_access = true
  band_width          = 20
}
```

### 4. 更新标签

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  resource_tags = [
    {
      tag_key   = "Environment"
      tag_value = "Production"
    },
    {
      tag_key   = "Owner"
      tag_value = "Team-A"
    },
  ]
}
```

### 5. 尝试修改不可变字段(应失败)

```hcl
# 尝试修改可用区(应返回错误)
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids = [2]  # 从 [1] 改为 [2]
}
```

## Schema 更新

### 移除 ForceNew 标记

对于以下字段,从 Schema 中移除 `ForceNew: true` 标记:

```go
"auto_renew_flag": {
    Type:        schema.TypeBool,
    Optional:    true,
    // 移除: ForceNew: true
    Description: "Automatic renewal, the default is true.",
},

"enable_public_access": {
    Type:        schema.TypeBool,
    Optional:    true,
    // 移除: ForceNew: true
    Description: "Whether to enable public network access. Default is false.",
},

"band_width": {
    Type:        schema.TypeInt,
    Optional:    true,
    Computed:    true,
    // 移除: ForceNew: true
    Description: "Public network bandwidth in Mbps.",
},
```

## API 参考文档

- [ModifyRabbitMQVipInstance API](https://cloud.tencent.com/document/api/1491/XXXXXX)
- [DescribeRabbitMQVipInstances API](https://cloud.tencent.com/document/api/1491/XXXXXX)
