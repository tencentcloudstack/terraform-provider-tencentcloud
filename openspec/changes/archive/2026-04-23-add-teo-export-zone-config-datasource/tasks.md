## 1. 数据源代码实现

- [x] 1.1 创建数据源文件 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config.go`，实现 `DataSourceTencentCloudTeoExportZoneConfig()` 函数，定义 Schema（zone_id: Required StringType, types: Optional List of StringType, content: Computed StringType, result_output_file: Optional StringType）和 Read 函数
- [x] 1.2 实现 `dataSourceTencentCloudTeoExportZoneConfigRead()` 函数：构建请求参数，直接调用 service.ExportZoneConfigByFilter()（service 层已包含 retry 逻辑），解析 response Content 设置到 content 字段，设置 d.SetId(helper.BuildToken())

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的数据源注册部分添加 `"tencentcloud_teo_export_zone_config": teo.DataSourceTencentCloudTeoExportZoneConfig()`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_export_zone_config` 数据源条目

## 3. 文档示例

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config.md` 文件，包含一句话描述、Example Usage 示例（含 zone_id 和 types 参数的用法）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config_test.go`，使用 gomonkey mock 方式编写单元测试，覆盖正常导出、指定 types 导出等场景

## 5. 验证

- [x] 5.1 运行 `go test` 验证单元测试通过
