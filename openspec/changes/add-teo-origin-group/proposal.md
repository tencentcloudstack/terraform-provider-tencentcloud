## Why

TEO (Tencent Edge One) 服务已经支持通过资源 `tencentcloud_teo_origin_group` 管理源站组，但缺少对应的数据源来查询源站组信息。添加数据源 `tencentcloud_teo_origin_group` 可以让用户以只读方式查询已创建的源站组配置，满足基础设施即代码中引用现有资源的需求。

## What Changes

- 新增数据源 `data_source_tc_teo_origin_group.go`，用于查询 TEO 源站组信息
- 新增数据源文档 `data_source_tc_teo_origin_group.md`
- 新增数据源测试 `data_source_tc_teo_origin_group_test.go`
- 数据源使用 `DescribeOriginGroup` API 进行只读查询

## Capabilities

### New Capabilities
- `teo-origin-group-datasource`: 提供 TEO 源站组数据源的查询功能，支持通过源站组 ID 获取源站组的详细配置信息，包括名称、类型、源站记录列表、回源 Host Header、引用实例、创建时间、更新时间等。

### Modified Capabilities
(无，不涉及现有功能的修改)

## Impact

- **新增文件**：
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group.go`
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group.md`
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group_test.go`

- **依赖服务**：
  - 使用现有的 `DescribeOriginGroup` API（Teo 20220901 版本）
  - 复用现有的 `DescribeTeoOriginGroupById` 服务方法

- **兼容性**：
  - 不影响现有资源 `tencentcloud_teo_origin_group` 的功能
  - 新增数据源，符合 Terraform Provider 最佳实践
