## Why

腾讯云配置审计（Config）产品提供了账号组（Aggregator）能力，用于对多账号/成员账号进行统一合规管理。当前 Terraform Provider 缺少查询账号组列表的数据源，用户无法通过 Terraform 查询现有账号组信息以供其他资源引用或编排。新增 `tencentcloud_config_list_aggregators` 数据源可填补该空白，让用户能够以声明式方式获取账号组列表。

## What Changes

- 新增数据源 `tencentcloud_config_list_aggregators`（RESOURCE_KIND_DATASOURCE），调用云 API `ListAggregators`（config v20220802）查询账号组列表。
- 数据源暴露账号组列表字段：`total`、`items`（包含 `name`、`description`、`owner_uin`、`create_time`、`account_count`、`type`、`account_group_id`、`aggregator_status`、`member_name`）。
- 在 `tencentcloud/provider.go` 中注册数据源。
- 在 `tencentcloud/provider.md` 中新增数据源文档条目（通过 `make doc` 生成）。
- 新增单元测试文件，使用 gomonkey mock 云 API。

## Capabilities

### New Capabilities
- `config-list-aggregators-datasource`: 提供 Config 账号组列表查询数据源 `tencentcloud_config_list_aggregators`，封装 `ListAggregators` 云 API，返回账号组总数与详情列表。

### Modified Capabilities

（无）

## Impact

- **新增文件**:
  - `tencentcloud/services/config/data_source_tc_config_list_aggregators.go`（数据源实现）
  - `tencentcloud/services/config/data_source_tc_config_list_aggregators_test.go`（单元测试，gomonkey mock）
  - `tencentcloud/services/config/data_source_tc_config_list_aggregators.md`（文档）
- **修改文件**:
  - `tencentcloud/services/config/service_tencentcloud_config.go`（新增 `DescribeConfigListAggregatorsByFilter` 方法）
  - `tencentcloud/provider.go`（注册数据源）
  - `tencentcloud/provider.md`（文档，由 `make doc` 生成）
- **云 API 依赖**: 复用已有 vendor 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`，使用 `ListAggregators` 接口（账号组列表）。
- **向后兼容**: 仅新增数据源，不修改任何现有资源 schema，向后完全兼容。
