## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/monitor/service_tencentcloud_monitor.go` 中新增 `DescribeMonitorGrafanaVersionsByFilter` 方法
- [x] 1.2 方法调用 `me.client.UseMonitorClient().DescribeGrafanaVersions(request)` 获取版本列表
- [x] 1.3 处理 API 响应，返回 `[]*monitor.GrafanaVersion`，正确处理空响应返回 nil

## 2. 数据源核心功能实现

- [x] 2.1 创建 `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions.go` 文件
- [x] 2.2 实现 `DataSourceTencentCloudMonitorGrafanaVersions()` schema 定义，包含 `versions` 列表（含 `alias`、`version` 字段）和 `result_output_file` 可选参数
- [x] 2.3 实现 `dataSourceTencentCloudMonitorGrafanaVersionsRead` 函数，使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 2.4 在 retry 块内调用 service 层方法，使用 `tccommon.ReadRetryTimeout` 超时和 `tccommon.RetryError` 包装错误
- [x] 2.5 在 retry 块内检查空响应（response 为空或 Versions 列表为空），返回 `NonRetryableError` 并打印 `[DATASOURCE] read empty, skip SetId` 日志
- [x] 2.6 将 API 响应映射到 schema 字段（遍历 Versions，设置 alias 和 version）
- [x] 2.7 使用 `helper.DataResourceIdsHash` 设置数据源 ID
- [x] 2.8 处理 `result_output_file` 参数，使用 `tccommon.WriteToFile` 导出结果

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_monitor_grafana_versions` 数据源，使用 `tcmg.DataSourceTencentCloudMonitorGrafanaVersions()`

## 4. 文档与测试

- [x] 4.1 创建 `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions.md` 文档文件，包含一句话描述（带云产品名称 monitor）、Example Usage
- [x] 4.2 创建 `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions_test.go` 单元测试文件，使用 gomonkey mock `DescribeMonitorGrafanaVersionsByFilter` 方法，验证字段映射逻辑

## 5. 验证与收尾

- [x] 5.1 运行 `gofmt` 格式化新增/修改的 Go 代码
- [x] 5.2 运行 `make doc` 生成 website/docs/ 下的文档
- [x] 5.3 创建 `.changelog` 目录下的 changelog 文件
