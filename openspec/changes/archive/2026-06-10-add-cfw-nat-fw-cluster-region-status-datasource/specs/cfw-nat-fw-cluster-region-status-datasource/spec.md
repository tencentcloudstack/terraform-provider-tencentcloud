## ADDED Requirements

### Requirement: 数据源查询 NAT 防火墙引流集群地域状态
系统 SHALL 提供 `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源，通过调用 `DescribeNatFwClusterRegionStatus` 接口查询 CFW NAT 防火墙引流集群的地域状态。

#### Scenario: 不带过滤条件查询
- **WHEN** 用户配置数据源且不指定 `nat_cluster_region_status_query_list`
- **THEN** 系统调用 `DescribeNatFwClusterRegionStatus` 接口（不传入查询列表），返回所有地域的集群状态，并将 `total` 和 `region_fw_status` 写入 state

#### Scenario: 带过滤条件查询
- **WHEN** 用户配置数据源并指定 `nat_cluster_region_status_query_list`（包含 `ccn_id`、`nat_ins_id`、`asset_type`、`routing_mode` 等字段）
- **THEN** 系统将查询列表传入 `DescribeNatFwClusterRegionStatus` 接口，返回符合条件的地域集群状态，并将结果写入 state

### Requirement: 数据源 Schema 定义
系统 SHALL 按照以下 Schema 定义数据源参数：
- 入参 `nat_cluster_region_status_query_list`：Optional，List of Object，每个对象包含：
  - `ccn_id`：Optional，String，云联网 ID
  - `nat_ins_id`：Optional，String，NAT 网关 ID
  - `asset_type`：Optional，String，资产类型（nat_ccn 或 nat）
  - `routing_mode`：Optional，Int，引流路由方法（0-多路由表模式，1-策略路由模式）
- 出参 `total`：Computed，Int，返回地域数量
- 出参 `region_fw_status`：Computed，List of Object，每个对象包含：
  - `nat_ins_id`：Computed，String，NAT 网关 ID
  - `ccn_id`：Computed，String，云联网 ID
  - `region`：Computed，String，地域
  - `status`：Computed，String，地域集群状态
  - `cidr`：Computed，String，引流网络 CIDR
  - `routing_mode`：Computed，Int，引流路由方法

#### Scenario: Schema 字段正确映射
- **WHEN** 云 API 返回 `RegionFwStatus` 列表
- **THEN** 系统将列表中每个元素的 `NatInsId`、`CcnId`、`Region`、`Status`、`Cidr`、`RoutingMode` 字段正确映射到 `region_fw_status` 列表中对应的 schema 字段

### Requirement: 错误处理与重试
系统 SHALL 在调用云 API 时使用 `tccommon.ReadRetryTimeout` 作为超时时间，并对可重试错误进行重试处理。

#### Scenario: 云 API 调用失败时重试
- **WHEN** 调用 `DescribeNatFwClusterRegionStatus` 接口返回可重试错误
- **THEN** 系统使用 `tccommon.RetryError()` 包装错误并重试，直到超时

#### Scenario: 云 API 返回空数据时不清空 state
- **WHEN** 调用 `DescribeNatFwClusterRegionStatus` 接口返回空响应
- **THEN** 系统在 retry 块内返回 `NonRetryableError`，不直接调用 `d.SetId("")`，避免 state 数据丢失

### Requirement: 注册数据源
系统 SHALL 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册 `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源。

#### Scenario: 数据源注册成功
- **WHEN** Terraform 初始化 provider
- **THEN** `tencentcloud_cfw_nat_fw_cluster_region_status` 数据源可被正常使用

### Requirement: 单元测试
系统 SHALL 为数据源提供使用 gomonkey mock 的单元测试，覆盖主要业务逻辑。

#### Scenario: 单元测试通过
- **WHEN** 使用 `go test -gcflags=all=-l` 运行单元测试
- **THEN** 所有测试用例通过，无编译错误
