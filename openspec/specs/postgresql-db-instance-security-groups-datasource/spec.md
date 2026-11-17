# postgresql-db-instance-security-groups-datasource Specification

## Purpose
TBD - created by archiving change add-postgresql-db-instance-security-groups-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data source registration
数据源 MUST 在 Provider 中注册为 `tencentcloud_postgresql_db_instance_security_groups`，使其可在 Terraform 配置中使用。

#### Scenario: Data source is accessible
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_postgresql_db_instance_security_groups"`
- **THEN** Terraform 能够成功识别并初始化该数据源

### Requirement: Query security groups by db_instance_id
数据源 MUST 支持通过 `db_instance_id` 参数查询实例关联的安全组。

#### Scenario: Query with db_instance_id
- **WHEN** 用户设置 `db_instance_id = "postgres-xxxxx"`
- **THEN** 数据源调用 DescribeDBInstanceSecurityGroups API 并将 `DBInstanceId` 设置为该值
- **THEN** 返回该实例关联的安全组列表

#### Scenario: db_instance_id is optional
- **WHEN** 用户未设置 `db_instance_id`
- **THEN** 不向 API 传递 `DBInstanceId` 参数

### Requirement: Query security groups by read_only_group_id
数据源 MUST 支持通过 `read_only_group_id` 参数查询只读组关联的安全组。

#### Scenario: Query with read_only_group_id
- **WHEN** 用户设置 `read_only_group_id = "pgro-xxxxx"`
- **THEN** 数据源调用 DescribeDBInstanceSecurityGroups API 并将 `ReadOnlyGroupId` 设置为该值
- **THEN** 返回该只读组关联的安全组列表

#### Scenario: read_only_group_id is optional
- **WHEN** 用户未设置 `read_only_group_id`
- **THEN** 不向 API 传递 `ReadOnlyGroupId` 参数

### Requirement: At least one query parameter required
数据源 MUST 依赖云 API 的入参约束：`DBInstanceId`、`ReadOnlyGroupId` 至少需传一个。

#### Scenario: No query parameters provided
- **WHEN** 用户同时未设置 `db_instance_id` 和 `read_only_group_id`
- **THEN** 云 API 返回参数校验错误 `InvalidParameter.ParameterCheckError`

#### Scenario: Both parameters provided
- **WHEN** 用户同时设置了 `db_instance_id` 和 `read_only_group_id`
- **THEN** 数据源将两个参数均传递给 API

### Requirement: Return security group set
数据源 MUST 返回 `security_group_set` 列表，将云 API 返回的 `SecurityGroupSet` 展开为顶层列表，每个元素包含安全组的详细属性。

#### Scenario: Security group set structure
- **WHEN** API 调用成功并返回安全组列表
- **THEN** `security_group_set` 列表包含以下字段:
  - `project_id`: 项目 ID（整数）
  - `create_time`: 创建时间（字符串）
  - `security_group_id`: 安全组 ID
  - `security_group_name`: 安全组名称
  - `security_group_description`: 安全组备注
  - `inbound`: 入站规则列表
  - `outbound`: 出站规则列表

#### Scenario: Empty security group set
- **WHEN** API 调用成功但返回空安全组列表
- **THEN** 数据源返回重试耗尽错误（遵循 DATASOURCE 空响应处理规范，不直接清空 id）

### Requirement: Return inbound policy rules
数据源 MUST 返回每个安全组的入站规则列表 `inbound`，元素为 `PolicyRule` 结构。

#### Scenario: Inbound rule structure
- **WHEN** 安全组包含入站规则
- **THEN** `inbound` 列表每个元素包含以下字段:
  - `action`: 策略，ACCEPT 或 DROP
  - `cidr_ip`: 来源或目的 IP 或 IP 段
  - `port_range`: 端口
  - `ip_protocol`: 网络协议，支持 UDP、TCP 等
  - `description`: 规则描述

#### Scenario: Inbound is nil
- **WHEN** 安全组的入站规则字段为 nil
- **THEN** 不调用 set 方法设置 `inbound`

### Requirement: Return outbound policy rules
数据源 MUST 返回每个安全组的出站规则列表 `outbound`，元素为 `PolicyRule` 结构。

#### Scenario: Outbound rule structure
- **WHEN** 安全组包含出站规则
- **THEN** `outbound` 列表每个元素包含以下字段:
  - `action`: 策略，ACCEPT 或 DROP
  - `cidr_ip`: 来源或目的 IP 或 IP 段
  - `port_range`: 端口
  - `ip_protocol`: 网络协议，支持 UDP、TCP 等
  - `description`: 规则描述

#### Scenario: Outbound is nil
- **WHEN** 安全组的出站规则字段为 nil
- **THEN** 不调用 set 方法设置 `outbound`

### Requirement: Support result output file
数据源 MUST 支持将结果保存到文件。

#### Scenario: Save to file
- **WHEN** 用户设置 `result_output_file = "./security_groups.json"`
- **THEN** 查询结果以 JSON 格式保存到指定文件

#### Scenario: File path is optional
- **WHEN** 用户未设置 `result_output_file`
- **THEN** 不保存文件，仅返回数据到 Terraform state

### Requirement: Handle API errors gracefully
数据源 MUST 正确处理 API 错误并返回清晰的错误信息。

#### Scenario: API call fails
- **WHEN** DescribeDBInstanceSecurityGroups API 调用失败
- **THEN** 使用 tccommon.RetryError() 包装错误返回
- **THEN** 错误信息包含 API 返回的错误代码和描述

#### Scenario: API returns empty response
- **WHEN** 云 API 返回空响应（`response == nil` 或 `response.Response == nil` 或 `len(SecurityGroupSet) == 0`）
- **THEN** 在 retry 块内返回 NonRetryableError，不直接 d.SetId("")，避免清空 state
- **THEN** 在外层 retry 失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

### Requirement: Retry with read timeout
数据源 MUST 使用 `tccommon.ReadRetryTimeout` 作为超时时间进行重试。

#### Scenario: Retry on transient failure
- **WHEN** API 调用因短暂故障失败
- **THEN** 在 tccommon.ReadRetryTimeout 时间内进行重试
- **THEN** 重试耗尽后返回最终错误

### Requirement: Nil check before setting fields
数据源在设置每个字段前 MUST 检查云 API Response 中对应字段是否为 nil。

#### Scenario: Skip nil fields
- **WHEN** Response 中某字段为 nil
- **THEN** 不调用对应的 set 方法

