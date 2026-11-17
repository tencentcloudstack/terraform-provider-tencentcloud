## 1. 数据源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_group_references.go`，定义 `DataSourceTencentCloudTeoIPGroupReferences` 数据源 Schema（zone_id Required TypeString, group_id Required TypeInt, references Computed TypeList of schema.Resource[含 zone_id/entity_type/entity_id/entity_name/sub_entity_type/sub_entity_id/sub_entity_name 均为 TypeString], result_output_file Optional TypeString）
- [x] 1.2 实现 `dataSourceTencentCloudTeoIPGroupReferencesRead` 函数：调用 `DescribeIPGroupReferences` API，内部自动分页（Limit=200，offset 递增）获取所有 references；使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，retry 块内仅调用接口；对 response 字段做 nil 判断后再 Set；retry 块内若返回空（response==nil / response.Response==nil / len(References)==0）则返回 NonRetryableError 而非清空 id，外层 retry 失败路径打印 `[DATASOURCE] read empty, skip SetId`；使用 `helper.BuildToken()` 生成数据源 ID；资源描述统一使用 `teo_ip_group_references`

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_ip_group_references` 数据源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_ip_group_references` 数据源条目

## 3. 数据源文档

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_group_references.md`，包含一句话描述（Use this data source to query ...，带上 TEO/EdgeOne 产品名称）、Example Usage 部分（使用 zone_id 和 group_id 查询），不添加 Argument Reference / Attribute Reference 部分

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_ip_group_references_test.go`，使用 gomonkey mock `DescribeIPGroupReferences` API 调用，测试 Read 操作的正常场景（含分页与字段映射校验），不使用 terraform 测试套件

## 5. 验证

- [x] 5.1 检查新增参数在云 API 接口中的存在性（入参 ZoneId/GroupId 与出参 References/TotalCount 均已在 `DescribeIPGroupReferencesRequest`/`DescribeIPGroupReferencesResponse` 中验证）
- [x] 5.2 确保 `data_source_tc_teo_ip_group_references.go` 中所有函数返回的 error 均被检查，无未使用变量；代码可在当前环境下正确构建执行
