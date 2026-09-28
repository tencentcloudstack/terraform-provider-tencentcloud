## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/dlc/service_tencentcloud_dlc.go` 中新增 `DescribeDlcTCLakeMetaInstance` 方法，调用 `DescribeTCLakeMetaInstanceWithContext`，包含三层空指针保护（`response == nil || response.Response == nil || response.Response.Status == nil`）
- [x] 1.2 方法包含 `ratelimit.Check`、DEBUG 日志、CRITAL 错误日志

## 2. Data Source Schema and Read

- [x] 2.1 创建 `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance.go`，定义 `DataSourceTencentCloudDlcTCLakeMetaInstance()` schema（`status` computed TypeString、`result_output_file` optional TypeString）
- [x] 2.2 实现 `dataSourceTencentCloudDlcTCLakeMetaInstanceRead` 方法，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹服务层调用，错误用 `tccommon.RetryError(e)` 包装
- [x] 2.3 在 retry 块内检查服务层返回 nil status 时返回 `resource.NonRetryableError` 并打印 `dlc tc_lake_meta_instance read empty, skip SetId`，禁止 `d.SetId("")`
- [x] 2.4 设置 `status` 前判断非 nil，`d.SetId("tc_lake_meta_instance")` 使用固定字符串，支持 `result_output_file` 写出

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册 `tencentcloud_dlc_tc_lake_meta_instance` 映射到 `dlc.DataSourceTencentCloudDlcTCLakeMetaInstance()`
- [x] 3.2 在 `tencentcloud/provider.md` 的 DLC Data Source 列表追加资源名

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance.md`，包含一句话描述（带 DLC 云产品名称，"Use this data source to query ..."前缀）和 Example Usage，不含 Argument Reference/Attribute Reference

## 5. Unit Testing

- [x] 5.1 创建 `tencentcloud/services/dlc/data_source_tc_dlc_tc_lake_meta_instance_test.go`，使用 mock（gomonkey）方法对云 API 进行 mock，仅测试业务代码逻辑，不使用 Terraform 测试套件
- [x] 5.2 确保生成的代码在当前环境下可正确编译执行，检查所有函数返回的 error
