## 1. 数据源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota.go`，定义 `DataSourceTencentCloudTeoContentQuota()` 函数，包含 schema 定义（zone_id Required、purge_quota Computed、prefetch_quota Computed、result_output_file Optional）和 Read 函数
- [x] 1.2 实现 `dataSourceTencentCloudTeoContentQuotaRead` 函数：构造 DescribeContentQuotaRequest，设置 ZoneId，使用 resource.Retry(tccommon.ReadRetryTimeout) 包装 API 调用，对 PurgeQuota/PrefetchQuota 及其子字段进行 nil 检查后设置到 state，使用 helper.BuildToken() 生成 ID
- [x] 1.3 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_content_quota` 数据源
- [x] 1.4 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_content_quota` 数据源文档条目

## 2. 数据源文档

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota.md`，包含一句话描述、Example Usage（含 zone_id 参数和 purge_quota/prefetch_quota 输出）

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota_test.go`，使用 gomonkey mock DescribeContentQuota API，覆盖成功读取场景，验证 schema 字段正确填充
- [x] 3.2 使用 `go test -gcflags=all=-l` 运行单元测试确认通过

## 4. 代码正确性验证

- [x] 4.1 检查数据源 schema 中的参数与云 API DescribeContentQuota 的入参/出参一致
- [x] 4.2 检查 provider.go 和 provider.md 中的注册信息完整
