### Requirement: 读取集群防火墙 Bypass 配置
资源 SHALL 通过 `DescribeClusterNatCcnFwSwitchList` 接口读取集群防火墙的当前 Bypass 状态及相关信息，并将结果存储到 Terraform state 中。

#### Scenario: 成功读取 Bypass 配置
- **WHEN** 用户执行 `terraform plan` 或 `terraform refresh`，且防火墙实例存在
- **THEN** 资源 SHALL 调用 `DescribeClusterNatCcnFwSwitchList` 接口，从返回的 `Data` 列表中找到匹配实例，并将 `total`、`data`、`region_list` 等字段写入 state

#### Scenario: 防火墙实例不存在时处理
- **WHEN** `DescribeClusterNatCcnFwSwitchList` 返回空列表或找不到匹配实例
- **THEN** 资源 SHALL 打印日志保留现场，然后调用 `d.SetId("")` 将资源标记为不存在

### Requirement: 更新集群防火墙 Bypass 状态
资源 SHALL 通过 `ModifyClusterFwBypass` 接口更新集群防火墙的 Bypass 开关状态。

#### Scenario: 开启 Bypass 模式
- **WHEN** 用户将 `enable` 设置为 `true` 并执行 `terraform apply`
- **THEN** 资源 SHALL 调用 `ModifyClusterFwBypass` 接口，传入 `fw_type`、`ccn_id`、`enable=true`（以及 fw_type 为 NAT_FW 时的 `nat_ins_id`），使流量绕过防火墙

#### Scenario: 关闭 Bypass 模式
- **WHEN** 用户将 `enable` 设置为 `false` 并执行 `terraform apply`
- **THEN** 资源 SHALL 调用 `ModifyClusterFwBypass` 接口，传入 `enable=false`，使流量经过防火墙

#### Scenario: 更新后读取最新状态
- **WHEN** `ModifyClusterFwBypass` 调用成功
- **THEN** 资源 SHALL 调用 Read 方法刷新 state 中的最新配置

### Requirement: 资源 ID 管理
资源 SHALL 使用联合 ID 唯一标识一个集群防火墙 Bypass 配置。

#### Scenario: VPC_FW 类型资源 ID
- **WHEN** `fw_type` 为 `VPC_FW`
- **THEN** 资源 ID SHALL 为 `{fw_type}#{ccn_id}` 格式

#### Scenario: NAT_FW 类型资源 ID
- **WHEN** `fw_type` 为 `NAT_FW`
- **THEN** 资源 ID SHALL 为 `{fw_type}#{ccn_id}#{nat_ins_id}` 格式，且 `nat_ins_id` 为必填字段

### Requirement: 资源参数 Schema 定义
资源 SHALL 定义以下参数：
- `fw_type`（Required, ForceNew）：防火墙类型，取值 "VPC_FW" 或 "NAT_FW"
- `ccn_id`（Required, ForceNew）：云联网实例 ID
- `enable`（Required）：Bypass 开关，true-开启 Bypass，false-关闭 Bypass
- `nat_ins_id`（Optional, ForceNew）：NAT 防火墙实例 ID，fw_type 为 NAT_FW 时必填
- `nat_type`（Optional）：NAT 防火墙类型筛选，用于 Read 查询
- `filters`（Optional）：过滤条件列表，用于 Read 查询
- `total`（Computed）：符合条件的总记录数
- `data`（Computed）：NAT 防火墙开关详情列表
- `region_list`（Computed）：地域列表

#### Scenario: 必填参数校验
- **WHEN** 用户未提供 `fw_type` 或 `ccn_id`
- **THEN** Terraform SHALL 报错，提示必填参数缺失

#### Scenario: NAT_FW 类型时 nat_ins_id 处理
- **WHEN** `fw_type` 为 `NAT_FW` 且用户提供了 `nat_ins_id`
- **THEN** 资源 SHALL 在调用 `ModifyClusterFwBypass` 时传入 `nat_ins_id` 参数

### Requirement: 注册到 Provider
资源 SHALL 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册，使其可被 Terraform 识别和使用。

#### Scenario: Provider 注册
- **WHEN** 用户在 Terraform 配置中使用 `tencentcloud_cfw_cluster_fw_bypass` 资源
- **THEN** Terraform SHALL 能够识别并加载该资源的 schema 和 CRUD 函数

### Requirement: 单元测试覆盖
资源 SHALL 提供使用 gomonkey mock 云 API 的单元测试，覆盖 Read 和 Update 的主要业务逻辑。

#### Scenario: Read 方法单元测试
- **WHEN** 运行 `go test -gcflags=all=-l` 执行单元测试
- **THEN** 测试 SHALL 通过，覆盖 Read 方法的正常路径和实例不存在路径

#### Scenario: Update 方法单元测试
- **WHEN** 运行 `go test -gcflags=all=-l` 执行单元测试
- **THEN** 测试 SHALL 通过，覆盖 Update 方法的正常路径
