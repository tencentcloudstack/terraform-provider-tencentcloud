# Spec: RabbitMQ 实例更新能力

## Capability Name

`rabbitmq-instance-update`

## 描述

本能力定义了 RabbitMQ VIP 实例资源的更新操作，允许用户通过 Terraform 修改实例的配置参数。

## 支持的更新参数

### 节点规格 (`node_spec`)

**描述**: 更新 RabbitMQ VIP 实例的节点规格

**类型**: String

**有效值**: 
- `rabbit-vip-basic-5` (2C4G)
- `rabbit-vip-profession-2c8g` (2C8G)
- `rabbit-vip-basic-1` (4C8G)
- `rabbit-vip-profession-4c16g` (4C16G)
- `rabbit-vip-basic-2` (8C16G)
- `rabbit-vip-profession-8c32g` (8C32G)
- `rabbit-vip-basic-4` (16C32G)
- `rabbit-vip-profession-16c64g` (16C64G)

**注意事项**:
- 规格升级可能导致短暂的服务中断，建议在业务低峰期执行
- 规格降级可能不支持，请确保目标规格满足当前数据量要求
- 某些规格可能暂时售罄，请选择其他可用规格

### 节点数量 (`node_num`)

**描述**: 更新 RabbitMQ VIP 实例的节点数量

**类型**: Int

**最小值**: 1 (单可用区)

**多可用区最小值**: 3

**注意事项**:
- 节点扩缩容可能导致短暂的服务中断
- 扩容时需要确保目标规格下有足够的资源
- 缩容时需要确保剩余节点能够承载当前负载

### 存储规格 (`storage_size`)

**描述**: 更新 RabbitMQ VIP 实例的单节点存储大小

**类型**: Int

**单位**: GB

**默认值**: 200 GB

**注意事项**:
- 存储扩容为单向操作，不支持缩容
- 扩容过程中可能会有短暂的性能波动
- 扩容后的空间会自动分配给各个节点

### 带宽 (`band_width`)

**描述**: 更新 RabbitMQ VIP 实例的公网带宽

**类型**: Int

**单位**: Mbps

**注意事项**:
- 带宽升级可能导致短暂的服务中断
- 带宽降级可能不支持，请确认目标带宽满足当前业务需求
- 带宽配置会影响公网访问的性能

### 公网访问开关 (`enable_public_access`)

**描述**: 启用或禁用 RabbitMQ VIP 实例的公网访问

**类型**: Bool

**默认值**: false

**注意事项**:
- 启用公网访问会分配一个公网 IP 和带宽
- 禁用公网访问会释放公网 IP 和带宽资源
- 启用公网访问需要确保安全组配置正确，避免安全风险

### 自动续费标识 (`auto_renew_flag`)

**描述**: 设置 RabbitMQ VIP 实例的自动续费状态

**类型**: Bool

**默认值**: true

**注意事项**:
- 仅对预付费实例有效
- 设置为 true 时，会在到期前自动续费
- 设置为 false 时，需要在到期前手动续费

### 集群版本 (`cluster_version`)

**描述**: 升级 RabbitMQ VIP 实例的集群版本

**类型**: String

**有效值**:
- `3.8.30` (默认版本)
- `3.11.8`
- `3.13.7`

**注意事项**:
- 版本升级可能导致短暂的服务中断
- 大版本升级（如 3.8.30 到 3.13.7）需要特别注意兼容性问题
- 升级前建议备份数据，并在测试环境中验证
- 某些版本可能存在已知的兼容性问题，请参考官方文档

## 不可修改的参数

以下参数在实例创建后无法修改，需要在创建时指定：

- `zone_ids` - 可用区列表
- `vpc_id` - VPC ID
- `subnet_id` - 子网 ID
- `enable_create_default_ha_mirror_queue` - 是否创建默认镜像队列
- `time_span` - 购买时长
- `pay_mode` - 付费模式

## 异步操作

某些更新操作（如节点规格升级、版本升级）是异步进行的，Terraform 会等待操作完成后才返回。

在等待过程中，Provider 会轮询实例状态，直到状态变为运行中或超时。

超时时间可以通过 `timeouts.update` 参数进行自定义，默认为 60 分钟。

## 错误处理

如果更新操作失败，Provider 会返回详细的错误信息，包括：

- 错误代码
- 错误描述
- 建议的解决方案

常见的错误原因包括：

1. **规格不支持**：目标规格暂时售罄或不支持当前实例配置
2. **资源不足**：目标可用区资源不足，无法完成扩容
3. **版本不兼容**：目标版本与当前配置不兼容
4. **并发操作冲突**：同时进行多个更新操作可能导致冲突

## 示例

### 示例 1: 升级节点规格

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids  = [1]
  vpc_id    = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  cluster_name = "example-rabbitmq"
  
  # 初始规格
  node_spec = "rabbit-vip-basic-1"  # 4C8G
  node_num  = 3
  storage_size = 200
  
  enable_create_default_ha_mirror_queue = true
  auto_renew_flag = true
  
  time_span = 1
}

# 后续升级规格
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... 其他配置保持不变
  
  # 升级到更高规格
  node_spec = "rabbit-vip-basic-2"  # 8C16G
}
```

### 示例 2: 扩容节点数量

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... 其他配置
  
  # 初始节点数量
  node_num = 3
  
  # 扩容到 5 个节点
  node_num = 5
}
```

### 示例 3: 升级集群版本

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... 其他配置
  
  # 初始版本
  cluster_version = "3.8.30"
  
  # 升级到 3.13.7
  cluster_version = "3.13.7"
  
  # 自定义超时时间（单位：分钟）
  timeouts {
    update = 90  # 90 分钟
  }
}
```

## 测试用例

### 测试用例 1: 节点规格升级

**前置条件**: 已有一个 RabbitMQ VIP 实例，规格为 `rabbit-vip-basic-1` (4C8G)

**操作**: 将规格升级为 `rabbit-vip-basic-2` (8C16G)

**预期结果**:
- 操作成功完成
- 实例规格更新为 8C16G
- 实例状态变为运行中

### 测试用例 2: 节点数量扩容

**前置条件**: 已有一个 RabbitMQ VIP 实例，节点数量为 3

**操作**: 将节点数量扩展为 5

**预期结果**:
- 操作成功完成
- 实例节点数量更新为 5
- 实例状态变为运行中

### 测试用例 3: 版本升级

**前置条件**: 已有一个 RabbitMQ VIP 实例，版本为 `3.8.30`

**操作**: 将版本升级为 `3.13.7`

**预期结果**:
- 操作成功完成
- 实例版本更新为 `3.13.7`
- 实例状态变为运行中

### 测试用例 4: 不可修改参数

**前置条件**: 已有一个 RabbitMQ VIP 实例

**操作**: 尝试修改 `zone_ids` 参数

**预期结果**:
- 操作失败
- 返回错误信息：参数 `zone_ids` 不可修改

## 相关资源

- [腾讯云 TDMQ RabbitMQ 官方文档](https://cloud.tencent.com/document/product/1493)
- [腾讯云 TDMQ RabbitMQ API 文档](https://cloud.tencent.com/document/product/1493/62966)
- [腾讯云 Terraform Provider RabbitMQ 实例资源文档](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/resources/tdmq_rabbitmq_vip_instance)

## 版本历史

- v1.0 (2026-03-29): 初始版本，支持节点规格、节点数量、存储规格、带宽、公网访问开关、自动续费标识、集群版本的更新
