# bdrc-security-group-mapping-resource Specification

## Purpose
TBD - created by archiving change add-bdrc-security-group-mapping. Update Purpose after archive.
## Requirements
### Requirement: 资源注册与客户端访问器
系统 SHALL 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `tencentcloud_bdrc_security_group_mapping`，并在 `tencentcloud/connectivity/client.go` 中新增 `UseBdrcV20260330Client()` 客户端访问器及 `bdrcv20260330Conn` 连接字段，供 bdrc 资源调用 bdrc v20260330 SDK。

#### Scenario: provider 中注册资源
- **WHEN** 用户在 Terraform 配置中声明 `resource "tencentcloud_bdrc_security_group_mapping" "example"`
- **THEN** provider SHALL 能够识别该资源类型并调用对应的 CRUD 函数

#### Scenario: 客户端访问器懒加载
- **WHEN** bdrc 资源首次调用 `UseBdrcV20260330Client()`
- **THEN** 系统 SHALL 创建 `bdrcv20260330.Client` 实例（使用 `NewClientProfile(300)` 与 `LogRoundTripper`）并缓存于 `bdrcv20260330Conn`，后续调用直接返回缓存实例

### Requirement: 创建安全组映射
系统 SHALL 通过调用 `CreateSecurityGroupMapping` 接口创建安全组映射，入参映射 `src_security_group_id`→`SrcSecurityGroupId`、`target_security_group_id`→`TargetSecurityGroupId`、`site_pair_id`→`SitePairId`。由于该接口为异步接口且不返回映射 ID，系统 MUST 在调用成功后轮询 `DescribeSecurityGroupMappings` 直到命中新建映射，获取其 `SecurityGroupMappingId`。

#### Scenario: 成功创建并轮询到映射 ID
- **WHEN** 用户通过 Terraform 创建 `tencentcloud_bdrc_security_group_mapping`，提供有效的 `src_security_group_id`、`target_security_group_id`、`site_pair_id`
- **THEN** 系统 SHALL 调用 `CreateSecurityGroupMapping`，随后轮询 `DescribeSecurityGroupMappings`（使用 `src-security-group-id` 与 `target-security-group-id` 过滤），当返回集合中存在源/目标安全组 ID 均匹配的记录时，取其 `SecurityGroupMappingId`
- **AND** 系统 SHALL 将资源 ID 设置为 `site_pair_id#security_group_mapping_id`（分隔符为 `tccommon.FILED_SP`）
- **AND** 系统 SHALL 调用 Read 回填 `security_group_mapping_id`、`source_security_group_id`、`life_state` 等 Computed 字段

#### Scenario: Create 接口返回空或失败
- **WHEN** `CreateSecurityGroupMapping` 调用返回 `nil` 响应或 `Response` 为 nil
- **THEN** 系统 SHALL 返回 `NonRetryableError`，不得写入空 ID

#### Scenario: 轮询超时未命中
- **WHEN** Create 调用成功但 `DescribeSecurityGroupMappings` 在重试期限内始终未命中新建映射
- **THEN** 系统 SHALL 以重试耗尽的形式失败，不得设置空 ID

### Requirement: 读取安全组映射
系统 SHALL 通过调用 `DescribeSecurityGroupMappings` 查询单条安全组映射。Read 时从 `d.Id()` 解析 `site_pair_id` 与 `security_group_mapping_id`，以 `SitePairId` 作为查询入参（`Limit` 取云 API 注释最大值 500），在返回的 `SecurityGroupMappingSet` 中按 `SecurityGroupMappingId` 精确匹配，将匹配项字段平铺回填到资源顶层 schema（不创建列表型嵌套 schema 层）。

#### Scenario: 成功读取并回填字段
- **WHEN** 资源存在且 `DescribeSecurityGroupMappings` 返回包含目标 `SecurityGroupMappingId` 的记录
- **THEN** 系统 SHALL 回填 `site_pair_id`、`src_security_group_id`、`target_security_group_id`、`security_group_mapping_id`、`source_security_group_id`、`life_state`（仅当对应 API 出参字段非 nil 时调用 `d.Set`）

#### Scenario: 映射不存在
- **WHEN** `DescribeSecurityGroupMappings` 返回空（`response == nil` / `response.Response == nil` / `len(SecurityGroupMappingSet)==0`）或未匹配到目标 ID
- **THEN** 系统 SHALL 先打印 `log.Printf("[CRUD] bdrc security_group_mapping id=%s", d.Id())` 保留现场，再调用 `d.SetId("")`

### Requirement: 更新安全组映射（不可变约束）
由于云 API 不提供安全组映射的更新接口，系统 SHALL 将所有业务顶层字段（`src_security_group_id`、`target_security_group_id`、`site_pair_id`）视为不可变。Update 函数 MUST 检测这些字段变更并返回 error，强制用户删除重建。

#### Scenario: 修改不可变字段
- **WHEN** 用户在 Terraform 配置中修改 `src_security_group_id`、`target_security_group_id` 或 `site_pair_id`
- **THEN** Update 函数 SHALL 返回 `fmt.Errorf("argument `%s` cannot be changed, please delete and recreate", v)`，不调用任何云 API

#### Scenario: 无变更
- **WHEN** 无任何不可变字段发生变更
- **THEN** Update 函数 SHALL 直接调用 Read 回填 state

### Requirement: 删除安全组映射
系统 SHALL 通过调用 `DeleteSecurityGroupMapping` 删除安全组映射，入参 `SitePairId` 从 `d.Id()` 解析，`SecurityGroupMappingIds` 为单元素列表（包含从 `d.Id()` 解析的 `security_group_mapping_id`）。

#### Scenario: 成功删除
- **WHEN** 用户执行 `terraform destroy` 删除资源
- **THEN** 系统 SHALL 调用 `DeleteSecurityGroupMapping`（`SitePairId` + `SecurityGroupMappingIds=[mappingId]`），调用使用 `resource.Retry(tccommon.WriteRetryTimeout)` 包装，失败用 `tccommon.RetryError` 包装

### Requirement: 导入安全组映射
系统 SHALL 支持通过联合 ID 导入资源，导入时使用 `schema.ImportStatePassthrough`，用户需提供 `site_pair_id#security_group_mapping_id` 格式的联合 ID。

#### Scenario: terraform import
- **WHEN** 用户执行 `terraform import tencentcloud_bdrc_security_group_mapping.example sitePairId#securityGroupMappingId`
- **THEN** 系统 SHALL 解析联合 ID 并调用 Read 加载资源状态

### Requirement: 资源文档
系统 SHALL 在 `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping.md` 中提供资源文档，包含一句话描述（带上云产品名称 BDRC）、Example Usage、Import 部分（说明需使用联合 ID），不手动添加 Argument Reference / Attribute Reference（由工具自动生成）。

#### Scenario: 文档完整性
- **WHEN** 资源代码生成完成
- **THEN** 对应 `.md` 文件 SHALL 存在并包含一句话描述、Example Usage、Import 说明

