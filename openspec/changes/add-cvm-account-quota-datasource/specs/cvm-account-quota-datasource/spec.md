## ADDED Requirements

### Requirement: Data source registration
数据源 MUST 在 Provider 中注册为 `tencentcloud_cvm_account_quota`，使其可在 Terraform 配置中使用。

#### Scenario: Data source is accessible
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_cvm_account_quota"`
- **THEN** Terraform 能够成功识别并初始化该数据源

### Requirement: Query account quota without filters
数据源 MUST 支持不带任何过滤条件的查询，返回完整的账户配额信息。

#### Scenario: Query all quotas
- **WHEN** 用户未设置任何过滤参数
- **THEN** 数据源调用 DescribeAccountQuota API 不传递 Filters 参数
- **THEN** 返回当前地域的所有配额类型数据

### Requirement: Filter by availability zone
数据源 MUST 支持通过 `zone` 参数按可用区过滤配额信息。

#### Scenario: Filter single zone
- **WHEN** 用户设置 `zone = ["ap-guangzhou-3"]`
- **THEN** 只返回广州三区的配额信息

#### Scenario: Filter multiple zones
- **WHEN** 用户设置 `zone = ["ap-guangzhou-3", "ap-guangzhou-4"]`
- **THEN** 返回广州三区和四区的配额信息

#### Scenario: Zone format validation
- **WHEN** 用户设置 `zone` 参数
- **THEN** 接受标准的可用区格式（如 `ap-guangzhou-3`）

### Requirement: Filter by quota type
数据源 MUST 支持通过 `quota_type` 参数按配额类型过滤。

#### Scenario: Filter post-paid quota
- **WHEN** 用户设置 `quota_type = "PostPaidQuotaSet"`
- **THEN** 只返回后付费配额数据

#### Scenario: Filter pre-paid quota
- **WHEN** 用户设置 `quota_type = "PrePaidQuotaSet"`
- **THEN** 只返回预付费配额数据

#### Scenario: Filter spot quota
- **WHEN** 用户设置 `quota_type = "SpotPaidQuotaSet"`
- **THEN** 只返回竞价实例配额数据

#### Scenario: Filter image quota
- **WHEN** 用户设置 `quota_type = "ImageQuotaSet"`
- **THEN** 只返回镜像配额数据

#### Scenario: Filter disaster recover group quota
- **WHEN** 用户设置 `quota_type = "DisasterRecoverGroupQuotaSet"`
- **THEN** 只返回置放群组配额数据

#### Scenario: Valid quota type values
- **WHEN** 用户设置 `quota_type` 参数
- **THEN** 支持的值为: `PostPaidQuotaSet`, `PrePaidQuotaSet`, `SpotPaidQuotaSet`, `ImageQuotaSet`, `DisasterRecoverGroupQuotaSet`

### Requirement: Return AppId
数据源 MUST 返回用户的 AppId。

#### Scenario: AppId is present
- **WHEN** API 调用成功
- **THEN** 输出属性 `app_id` 包含用户的 AppId 整数值

### Requirement: Return post-paid quota details
数据源 MUST 返回后付费配额详细信息。

#### Scenario: Post-paid quota structure
- **WHEN** 配额数据包含后付费配额
- **THEN** `post_paid_quota_set` 列表包含以下字段:
  - `zone`: 可用区
  - `total_quota`: 总配额
  - `used_quota`: 已使用配额
  - `remaining_quota`: 剩余配额

### Requirement: Return pre-paid quota details
数据源 MUST 返回预付费配额详细信息。

#### Scenario: Pre-paid quota structure
- **WHEN** 配额数据包含预付费配额
- **THEN** `pre_paid_quota_set` 列表包含以下字段:
  - `zone`: 可用区
  - `total_quota`: 总配额
  - `used_quota`: 已使用配额
  - `remaining_quota`: 剩余配额
  - `once_quota`: 单次购买配额

### Requirement: Return spot instance quota details
数据源 MUST 返回竞价实例配额详细信息。

#### Scenario: Spot quota structure
- **WHEN** 配额数据包含竞价实例配额
- **THEN** `spot_paid_quota_set` 列表包含以下字段:
  - `zone`: 可用区
  - `total_quota`: 总配额
  - `used_quota`: 已使用配额
  - `remaining_quota`: 剩余配额

### Requirement: Return image quota details
数据源 MUST 返回镜像配额详细信息。

#### Scenario: Image quota structure
- **WHEN** 配额数据包含镜像配额
- **THEN** `image_quota_set` 列表包含以下字段:
  - `total_quota`: 总配额
  - `used_quota`: 已使用配额

### Requirement: Return disaster recover group quota details
数据源 MUST 返回置放群组配额详细信息。

#### Scenario: Disaster recover group quota structure
- **WHEN** 配额数据包含置放群组配额
- **THEN** `disaster_recover_group_quota_set` 列表包含以下字段:
  - `group_quota`: 置放群组配额
  - `current_num`: 当前已使用的置放群组数量
  - `cvm_in_host_group_quota`: 交换机类型容灾组内最大实例数
  - `cvm_in_switch_group_quota`: 专用宿主机类型容灾组内最大实例数
  - `cvm_in_rack_group_quota`: 机架类型容灾组内最大实例数

### Requirement: Return region information
数据源 MUST 返回配额所属地域信息。

#### Scenario: Region is present
- **WHEN** API 调用成功
- **THEN** `account_quota_overview` 的 `region` 字段包含地域标识（如 `ap-guangzhou`）

### Requirement: Support result output file
数据源 MUST 支持将结果保存到文件。

#### Scenario: Save to file
- **WHEN** 用户设置 `result_output_file = "./quota.json"`
- **THEN** 查询结果以 JSON 格式保存到指定文件

#### Scenario: File path is optional
- **WHEN** 用户未设置 `result_output_file`
- **THEN** 不保存文件，仅返回数据到 Terraform state

### Requirement: Handle API errors gracefully
数据源 MUST 正确处理 API 错误并返回清晰的错误信息。

#### Scenario: API call fails
- **WHEN** DescribeAccountQuota API 调用失败
- **THEN** 返回包含错误详情的 Terraform 错误
- **THEN** 错误信息包含 API 返回的错误代码和描述

#### Scenario: Network error
- **WHEN** 网络连接失败
- **THEN** 返回明确的网络错误信息

### Requirement: Combine multiple filters
数据源 MUST 支持同时使用多个过滤条件。

#### Scenario: Filter by zone and quota type
- **WHEN** 用户同时设置 `zone = ["ap-guangzhou-3"]` 和 `quota_type = "PostPaidQuotaSet"`
- **THEN** 只返回广州三区的后付费配额信息
