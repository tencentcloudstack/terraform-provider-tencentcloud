## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoContentQuotaByFilter` 方法，封装 `DescribeContentQuota` API 调用，支持 `ZoneId` 参数，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 处理重试，返回 `purgeQuota []*teo.Quota, prefetchQuota []*teo.Quota, err error`

## 2. 数据源 Schema 与 Read 函数实现

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota.go`，实现 `DataSourceTencentCloudTeoContentQuota()` 函数定义 schema，包含 `zone_id`（Required）、`purge_quota`（Computed, TypeList 嵌套对象）、`prefetch_quota`（Computed, TypeList 嵌套对象）、`result_output_file`（Optional）
- [x] 2.2 实现 `dataSourceTencentCloudTeoContentQuotaRead` 函数，调用 service 层方法获取配额数据，将 PurgeQuota/PrefetchQuota flatten 到 schema 中，使用 `helper.DataResourceIdsHash` 生成数据源 ID，处理 nil 响应字段

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册 `tencentcloud_teo_content_quota` 数据源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_content_quota` 数据源条目

## 4. 文档与测试

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota.md` 文档文件，包含一句话描述（提及 TEO 产品名）、Example Usage 部分
- [x] 4.2 创建 `tencentcloud/services/teo/data_source_tc_teo_content_quota_test.go`，使用 gomonkey mock 方式编写单元测试，mock `DescribeContentQuota` API 调用，验证 Read 函数的业务逻辑
- [x] 4.3 使用 `go test -gcflags=all=-l` 运行单元测试并确保通过

## 5. 收尾验证

- [ ] 5.1 执行 `gofmt` 格式化所有新增/修改的 Go 文件
- [ ] 5.2 执行 `make doc` 生成 website/docs/ 下的文档
- [ ] 5.3 创建 `.changelog` 目录下的 changelog 文件
