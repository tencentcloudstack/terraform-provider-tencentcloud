## Why

用户需要通过 Terraform 查询 IP 是否属于腾讯云 EdgeOne (TEO) 节点，以便在基础设施即代码中实现基于 IP 归属的条件逻辑。当前 TEO 产品没有提供 `tencentcloud_teo_ip_region` 数据源，用户无法在 Terraform 配置中直接获取 IP 归属信息。

## What Changes

- 新增数据源 `tencentcloud_teo_ip_region`，调用 TEO 的 `DescribeIPRegion` 接口查询 IP 归属信息
- 数据源入参 `i_ps`：待查询的 IP 列表（支持 IPv4 和 IPv6，最大 100 条）
- 数据源出参 `ip_region_info`：IP 归属信息列表，包含 IP 地址和是否属于 EdgeOne 节点的标识
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 新增数据源文档 `data_source_tc_teo_ip_region.md`

## Capabilities

### New Capabilities
- `teo-ip-region-datasource`: 新增 TEO IP 归属查询数据源，支持通过 IP 列表查询其是否属于 EdgeOne 节点

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_ip_region.go`、`tencentcloud/services/teo/data_source_tc_teo_ip_region_extension.go`、`tencentcloud/services/teo/data_source_tc_teo_ip_region_test.go`
- 修改文件：`tencentcloud/provider.go`（注册数据源）、`tencentcloud/provider.md`（文档）
- 新增文档：`gendoc/teo/data_source_tc_teo_ip_region.md`
- 依赖云 API：`teo.v20220901.DescribeIPRegion`
