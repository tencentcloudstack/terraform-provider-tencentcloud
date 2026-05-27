## Why

为 GA2（全球加速2.0）云产品新增 Terraform 数据源 `tencentcloud_ga2_accelerate_regions`，使用户能够通过 Terraform 查询可选的加速区域信息，便于在编排加速实例时选择合适的加速区域。

## What Changes

- 新增数据源 `tencentcloud_ga2_accelerate_regions`，调用 `DescribeAccelerateRegions` 云API查询可选加速区域列表
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册该数据源
- 新增数据源文档 `data_source_tc_ga2_accelerate_regions.md`
- 新增单元测试文件 `data_source_tc_ga2_accelerate_regions_test.go`

## Capabilities

### New Capabilities

- `ga2-accelerate-regions-datasource`: 提供查询 GA2 可选加速区域的数据源能力，返回加速区域列表（包含地域名称、可用性、地域标识、地区名称、是否中国地域、支持的ISP类型、是否腾讯地域等信息）

### Modified Capabilities

（无）

## Impact

- 新增文件: `tencentcloud/services/ga2/data_source_tc_ga2_accelerate_regions.go`
- 新增文件: `tencentcloud/services/ga2/data_source_tc_ga2_accelerate_regions_test.go`
- 新增文件: `tencentcloud/services/ga2/data_source_tc_ga2_accelerate_regions.md`
- 修改文件: `tencentcloud/provider.go`（注册数据源）
- 修改文件: `tencentcloud/provider.md`（添加数据源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`（已在 vendor 中）
