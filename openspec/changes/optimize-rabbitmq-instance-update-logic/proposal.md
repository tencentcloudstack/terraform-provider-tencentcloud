# Change: 优化 RabbitMQ 实例的 Update 逻辑

## Why

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `Update` 方法存在严重的功能限制，将多个可以通过腾讯云 API 修改的参数标记为不可修改（immutable）。这导致用户无法通过 Terraform 实现常见的运维场景，如：

1. **实例扩容**：无法通过修改 `node_num` 来增加节点数量
2. **规格升级**：无法通过修改 `node_spec` 来升级节点规格
3. **存储扩容**：无法通过修改 `storage_size` 来扩容存储
4. **配置调整**：无法调整 `band_width`（带宽）、`enable_public_access`（公网访问）等配置参数
5. **版本升级**：无法升级 `cluster_version`（集群版本）

这些问题违背了基础设施即代码（IaC）的原则，迫使用户通过控制台手动操作，导致：
- **配置漂移**：手动修改后 Terraform 状态与实际配置不一致
- **运维复杂度增加**：破坏了 IaC 的一致性和可追溯性
- **自动化能力受限**：无法实现自动化的弹性扩缩容和资源优化

根据腾讯云 TDMQ RabbitMQ 官方文档，`ModifyRabbitMQVipInstance` API 支持修改多个实例配置参数。因此，需要优化 Update 逻辑，支持更多参数的动态修改。

## What Changes

- 优化 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中的 `Update` 方法
- 将以下参数从 `immutableArgs` 列表中移除，并支持通过 API 修改：
  - `node_spec`: 节点规格
  - `node_num`: 节点数量
  - `storage_size`: 存储规格
  - `band_width`: 带宽
  - `enable_public_access`: 是否启用公网访问
  - `cluster_version`: 集群版本
  - `auto_renew_flag`: 自动续费标识
- 保留以下参数为不可修改（创建后不能修改）：
  - `zone_ids`: 可用区
  - `vpc_id`: VPC ID
  - `subnet_id`: 子网 ID
  - `enable_create_default_ha_mirror_queue`: 是否创建默认 HA 镜像队列
  - `time_span`: 购买时长
  - `pay_mode`: 付费模式
- 为修改 `cluster_version` 添加异步任务等待逻辑，因为版本升级需要时间
- 添加 `timeouts` 配置支持，允许用户配置 update 操作的超时时间
- 更新资源文档，明确说明哪些参数可以修改、哪些不能修改

## Impact

- **新增能力**: RabbitMQ 实例的动态扩容、规格升级、存储扩容、带宽调整、公网访问开关、版本升级等
- **受影响的服务**: TDMQ (tencentcloud/services/trabbit)
- **新增文件**:
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` (修改)
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` (更新文档)
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` (更新测试)
- **向后兼容**: 完全向后兼容，仅添加新的修改能力，不破坏现有配置
