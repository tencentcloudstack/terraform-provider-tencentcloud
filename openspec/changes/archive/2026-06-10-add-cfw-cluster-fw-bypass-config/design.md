## Context

腾讯云防火墙（CFW）提供集群防火墙 Bypass 状态管理能力。Bypass 模式允许流量绕过防火墙（Enable=true）或经过防火墙（Enable=false）。

当前 Terraform Provider 的 CFW 服务目录（`tencentcloud/services/cfw/`）已有部分资源，本次新增 `tencentcloud_cfw_cluster_fw_bypass` CONFIG 类型资源，用于管理集群防火墙的 Bypass 状态。

涉及的云 API：
- `DescribeClusterNatCcnFwSwitchList`：查询 NAT CCN 集群模式防火墙开关列表，用于 Read 操作
- `ModifyClusterFwBypass`：修改集群防火墙 Bypass 状态，用于 Update 操作

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_cfw_cluster_fw_bypass` CONFIG 资源，支持读取和更新集群防火墙 Bypass 状态
- 资源 ID 使用 `fw_type` + `ccn_id`（+ `nat_ins_id` 当 fw_type 为 NAT_FW 时）的联合 ID
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 提供对应的单元测试（使用 gomonkey mock 云 API）
- 提供对应的 `.md` 文档

**Non-Goals:**
- 不支持创建或删除集群防火墙实例（仅管理 Bypass 配置）
- 不支持管理防火墙开关状态（Status 字段只读）

## Decisions

### 1. 资源类型选择 RESOURCE_KIND_CONFIG
由于集群防火墙实例由其他资源或控制台创建，本资源仅管理其 Bypass 配置。CONFIG 类型资源只需 Read（R）和 Update（U）操作，无需 Create/Delete。

### 2. 资源 ID 设计
使用联合 ID：`{fw_type}#{ccn_id}` 或 `{fw_type}#{ccn_id}#{nat_ins_id}`（当 fw_type 为 NAT_FW 时）。
- `fw_type`：防火墙类型（VPC_FW 或 NAT_FW），必填
- `ccn_id`：云联网实例 ID，必填
- `nat_ins_id`：NAT 防火墙实例 ID，当 fw_type 为 NAT_FW 时必填

使用 `tccommon.FILED_SP`（`#`）作为分隔符，与 Provider 其他资源保持一致。

### 3. Read 接口参数设计
`DescribeClusterNatCcnFwSwitchList` 接口支持通过 `NatType` 和 `Filters` 过滤。Read 时通过 `nat_type` 和 `filters` 参数查询，从返回的 `Data` 列表中找到匹配的实例。

### 4. 只读字段处理
`DescribeClusterNatCcnFwSwitchList` 返回的 `Data`（NatFwSwitchDetailS 列表）、`Total`、`RegionList` 等字段作为 Computed 字段存储在 state 中，供用户查看。

### 5. 单元测试策略
使用 gomonkey 对云 API 进行 mock，测试资源的 CRUD 业务逻辑，不依赖真实云环境。

## Risks / Trade-offs

- [风险] `DescribeClusterNatCcnFwSwitchList` 返回列表，需要根据 ID 字段匹配目标实例。若实例不存在，需正确处理 `d.SetId("")` → 使用 `log.Printf` 保留现场后再清空 ID
- [风险] `ModifyClusterFwBypass` 接口未标注为异步，但实际可能存在延迟 → 当前按同步处理，若后续发现需要轮询可扩展
- [权衡] CONFIG 资源无 Create/Delete，用户需要先通过控制台或其他方式创建防火墙实例，再通过本资源管理 Bypass 配置
