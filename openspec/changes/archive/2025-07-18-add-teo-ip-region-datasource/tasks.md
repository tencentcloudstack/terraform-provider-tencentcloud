## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoIPRegionByFilter` 方法，构建 `DescribeIPRegionRequest`，设置 `IPs` 字段，调用 `DescribeIPRegion` API，返回 `IPRegionInfo` 列表

## 2. Data Source Schema & Read

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_region.go`，定义 `DataSourceTencentCloudTeoIPRegion()` 函数，包含 schema 定义（`i_ps`、`ip_region_info`、`result_output_file`）和 Read 函数
- [x] 2.2 实现 `dataSourceTencentCloudTeoIPRegionRead` 函数：从 schema 获取 `i_ps` 参数，调用服务层方法（带 retry），将结果映射到 `ip_region_info`，设置数据源 ID，处理 `result_output_file` 输出
- [x] 2.3 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_region_extension.go` 扩展文件（仅包含 `package teo`）

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 的数据源 map 中添加 `"tencentcloud_teo_ip_region": teo.DataSourceTencentCloudTeoIPRegion()` 注册项
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_ip_region` 数据源的文档条目

## 4. Documentation

- [x] 4.1 在 `gendoc/teo/` 目录下创建 `data_source_tc_teo_ip_region.md` 文档文件，包含一句话描述、Example Usage（使用 jsonencode 处理 JSON 场景）

## 5. Unit Tests

- [x] 5.1 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_region_test.go`，使用 gomonkey mock `DescribeIPRegion` API 调用，验证 Read 函数正确映射 `ip_region_info` 中的 `ip` 和 `is_edge_one_ip` 字段
