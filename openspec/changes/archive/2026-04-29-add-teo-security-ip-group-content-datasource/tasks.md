## 1. 数据源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/data_source_tc_teo_security_ip_group_content.go`，定义 `DataSourceTencentCloudTeoSecurityIPGroupContent` 数据源 Schema（zone_id Required TypeString, group_id Required TypeInt, ip_total_count Computed TypeInt, ip_list Computed TypeList of TypeString, result_output_file Optional TypeString）
- [x] 1.2 实现 `dataSourceTencentCloudTeoSecurityIPGroupContentRead` 函数，调用 `DescribeSecurityIPGroupContent` API，内部自动分页获取所有 IP 数据，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，对 response 字段做 nil 判断后再 Set，使用 `helper.BuildToken()` 生成数据源 ID

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_security_ip_group_content` 数据源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_security_ip_group_content` 数据源条目

## 3. 数据源文档

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_security_ip_group_content.md`，包含一句话描述（Use this data source to query ...）、Example Usage 部分

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_security_ip_group_content_test.go`，使用 gomonkey mock `DescribeSecurityIPGroupContent` API 调用，测试 Read 操作的正常场景
- [x] 4.2 使用 `go test -gcflags=all=-l` 运行单元测试并确保通过
