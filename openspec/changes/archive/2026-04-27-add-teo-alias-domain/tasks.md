# Tasks: tencentcloud_teo_alias_domain

## Implementation Steps

- [x] **Task 1**: 在 `service_tencentcloud_teo.go` 中追加 `DescribeTeoAliasDomainById` 方法
  - 使用 `DescribeAliasDomains` 接口，Filter `alias-name` 精确过滤
  - 分页循环匹配 ZoneId + AliasName
  - 返回 `*teov20220901.AliasDomain`

- [x] **Task 2**: 新增 `resource_tc_teo_alias_domain.go` 资源主文件
  - 定义 `ResourceTencentCloudTeoAliasDomain()` 返回 schema.Resource
  - Schema 字段与 CreateAliasDomain 接口入参严格对齐（zone_id, alias_name, target_name, cert_type, cert_id, paused）
  - 追加只读字段（status, forbid_mode, created_on, modified_on）
  - 实现 Create：调用 CreateAliasDomain，若 paused=true 额外调用 ModifyAliasDomainStatus
  - 实现 Read：调用 service.DescribeTeoAliasDomainById，不覆盖 cert_type/cert_id（API 不返回）
  - 实现 Update：检测 target_name/cert_type/cert_id 变更调 ModifyAliasDomain；检测 paused 变更调 ModifyAliasDomainStatus
  - 实现 Delete：调用 DeleteAliasDomain（AliasNames=[aliasName]）
  - 设置 Importer（schema.ImportStatePassthrough）

- [x] **Task 3**: 新增 `resource_tc_teo_alias_domain.md` 文档
  - Example Usage HCL 示例（create + update 两段）
  - Import 说明（zone_id#alias_name）

- [x] **Task 4**: 新增 `resource_tc_teo_alias_domain_test.go` 单元测试
  - `TestAccTencentCloudTeoAliasDomainResource_basic`（含 create、update、import 三个 Step）
  - 参考 `resource_tc_config_compliance_pack_test.go` 风格

- [x] **Task 5**: 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_alias_domain` 资源
