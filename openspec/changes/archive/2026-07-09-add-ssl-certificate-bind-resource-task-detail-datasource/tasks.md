## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/ssl/service_tencent_ssl_certificate.go` 中追加 `DescribeSslCertificateBindResourceTaskDetail` service 方法
- [x] 1.2 实现 request 参数构造（TaskId、ResourceTypes、Regions、Limit=100、Offset 分页循环）
- [x] 1.3 实现异步任务轮询逻辑：在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内检查 `Status`，为 0 时返回 `tccommon.RetryError` 重试，为 1/2 时返回结果
- [x] 1.4 实现分页合并：累加 Offset 直到所有云资源详情获取完毕，返回合并后的完整 `DescribeCertificateBindResourceTaskDetailResponseParams`

## 2. 数据源 Schema 与 Read 函数实现

- [x] 2.1 创建 `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail.go` 文件
- [x] 2.2 实现 `DataSourceTencentCloudSslCertificateBindResourceTaskDetail` Schema 定义，包含输入参数（task_id、resource_types、regions、result_output_file）和输出属性（status、cache_time 及 16 类云资源详情列表）
- [x] 2.3 实现 16 类云资源详情的嵌套 Schema（CLB、CDN、WAF、DDoS、Live、VOD、TKE、APIGateway、TCB、TEO、TSE、COS、TDMQ、MQTT、GAAP、SCF），完整映射 vendor SDK 结构体字段
- [x] 2.4 实现 `dataSourceTencentCloudSslCertificateBindResourceTaskDetailRead` 函数：构造请求参数、调用 service 方法、处理 retry 块内空响应返回 NonRetryableError、设置各字段到 state
- [x] 2.5 在设置字段前检查 nil，使用 `helper.BuildToken()` 生成数据源 ID，支持 `result_output_file` 输出

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_ssl_certificate_bind_resource_task_detail` 数据源

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail_test.go` 文件
- [x] 4.2 使用 gomonkey mock 云 API（`DescribeCertificateBindResourceTaskDetail`），实现业务逻辑单元测试
- [x] 4.3 测试正常查询场景（Status=1 成功）
- [x] 4.4 测试异步轮询场景（Status=0 重试后成功）
- [x] 4.5 测试空响应返回 NonRetryableError 场景
- [x] 4.6 使用 `go test -gcflags=all=-l` 跑通涉及的单元测试文件

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail.md` 文档文件
- [x] 5.2 包含一句话描述（带 SSL 云产品名称）、Example Usage、Import 部分（不适用，RESOURCE_KIND_DATASOURCE 无 Import）
- [x] 5.3 不添加 Argument Reference 和 Attribute Reference 部分（由工具自动生成）

## 6. 验证与收尾

- [x] 6.1 运行 `gofmt` 格式化代码
- [x] 6.2 运行 `make doc` 生成 `website/docs/` 下的 markdown 文档
- [x] 6.3 生成 `.changelog` 目录下的 changelog 文件
