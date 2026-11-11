## Why

DBCustom 节点安全组信息是查询节点安全配置的关键数据。当前 Terraform Provider 缺少 `tencentcloud_dbdc_db_custom_node_security_groups` 数据源，用户无法通过 Terraform 直接查询 DBCustom 节点的安全组绑定信息。新增此数据源可以补齐 DBDC 产品数据源矩阵，满足用户查询节点安全组的需求。

## What Changes

- 新增 RESOURCE_KIND_DATASOURCE 数据源 `tencentcloud_dbdc_db_custom_node_security_groups`
- 调用云 API `DescribeDBCustomNodeSecurityGroups`，传入 `node_id` 查询指定节点的安全组列表
- 将安全组信息展开为 Terraform schema 字段，支持安全组 ID、名称、入/出站规则等属性
- 在 `tencentcloud/provider.go` 中注册新数据源
- 新增对应的 `.md` 使用样例文档

## Capabilities

### New Capabilities
- `dbdc-db-custom-node-security-groups-datasource`: 提供 DBDC 自定义节点安全组查询数据源，通过节点 ID 查询该节点绑定的安全组及其规则详情

### Modified Capabilities
<!-- None - this is a new capability, no existing specs are modified -->

## Impact

- **新增文件**: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups.go`
- **新增文件**: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups_test.go`
- **新增文件**: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups.md`
- **修改文件**: `tencentcloud/provider.go` - 注册新数据源
- **SDK 依赖**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` (已在 vendor 中)
- **API 接口**: `DescribeDBCustomNodeSecurityGroups` (非异步接口，无需轮询)