## Context

CFW（云防火墙）提供 NAT 防火墙引流集群功能，用户需要通过 Terraform 数据源查询各地域的集群部署状态和引流网络配置。当前 Terraform Provider 中缺少该数据源，需要新增 `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源。

云 API `DescribeNatFwClusterRegionStatus` 支持通过 `NatClusterRegionStatusQueryList` 参数列表过滤查询，返回地域数量（`Total`）和地域防火墙集群状态列表（`RegionFwStatus`）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源，调用 `DescribeNatFwClusterRegionStatus` 接口
- 支持通过 `nat_cluster_region_status_query_list` 参数（包含 `ccn_id`、`nat_ins_id`、`asset_type`、`routing_mode` 字段）过滤查询
- 返回 `total`（地域数量）和 `region_fw_status`（地域防火墙集群状态列表，包含 `nat_ins_id`、`ccn_id`、`region`、`status`、`cidr`、`routing_mode` 字段）
- 在 `provider.go` 和 `provider.md` 中注册数据源
- 提供对应的 `.md` 文档和单元测试

**Non-Goals:**
- 不支持创建、更新、删除 NAT 防火墙引流集群地域状态（只读数据源）
- 不实现分页（该接口无分页参数）

## Decisions

### 决策1：数据源实现方式

采用标准 RESOURCE_KIND_DATASOURCE 模式，参考 `tencentcloud_igtm_instance_list` 数据源的实现风格：
- 在 `data_source_tc_cfw_nat_fw_cluster_region_status.go` 中实现数据源逻辑
- 使用 `tccommon.ReadRetryTimeout` 作为超时时间，添加 retry 处理
- 在 Read 方法中调用云 API，将结果写入 state

### 决策2：Schema 设计

- 入参 `nat_cluster_region_status_query_list` 为 Optional 的 List of Object，每个对象包含 `ccn_id`、`nat_ins_id`、`asset_type`、`routing_mode` 字段
- 出参 `total` 为 Computed 的 int，`region_fw_status` 为 Computed 的 List of Object，每个对象包含 `nat_ins_id`、`ccn_id`、`region`、`status`、`cidr`、`routing_mode` 字段
- 将 `region_fw_status` 列表中每个元素的字段平铺到 schema 中，不额外嵌套

### 决策3：ID 生成

数据源使用 `helper.BuildUUID()` 生成唯一 ID，避免与其他数据源冲突。

## Risks / Trade-offs

- [风险] 云 API 返回的 `RegionFwStatus` 可能为 nil → 在 Read 方法中判断 nil 后再设置字段，避免 panic
- [风险] 云 API 短暂波动导致返回空数据 → 在 retry 块内若返回空则返回 `NonRetryableError`，避免清空 state

## Migration Plan

无需迁移，纯新增数据源。

## Open Questions

无。
