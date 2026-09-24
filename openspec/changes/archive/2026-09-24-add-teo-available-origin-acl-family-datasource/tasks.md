## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoAvailableOriginAclFamilyByFilter` 方法，封装云 API `DescribeAvailableOriginACLFamily`，内部以 `Limit=100`（API 最大值）循环分页累加 `Offset` 获取所有数据，调用通过 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹并使用 `tccommon.RetryError(e)` 包装错误

## 2. 数据源 Schema 与 Read 函数实现

- [x] 2.1 新建 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.go`，定义 `DataSourceTencentCloudTeoAvailableOriginAclFamily()` schema：`zone_id`（必填 string）、`filters`（选填 TypeList，子字段 `name`/`values`，对应云 API `Filter`）、`origin_acl_family_info_set`（Computed TypeList，元素展开为 `version`/`active_time`/`entire_addresses`(含 `ipv4`/`ipv6`)/`origin_acl_family`）、`total_count`（Computed int）、`result_output_file`（选填 string）
- [x] 2.2 在同文件实现 `dataSourceTencentCloudTeoAvailableOriginAclFamilyRead` 函数：组装 `paramMap`（ZoneId、Filters），retry 块内调用 service 方法，若返回空则返回 `NonRetryableError` 并 `log.Printf("[DATASOURCE] read empty, skip SetId")`，将 `OriginACLFamilyInfos` 映射到 `origin_acl_family_info_set`（仅在非 nil 时 set），`d.SetId(helper.BuildToken())`，处理 `result_output_file`

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册数据源 `tencentcloud_teo_available_origin_acl_family`，参考 teo 分组中 `tencentcloud_teo_origin_acl` 的注册方式
- [x] 3.2 在 `tencentcloud/provider.md` 中为数据源 `tencentcloud_teo_available_origin_acl_family` 添加文档索引注册，参考 `tencentcloud_teo_origin_acl` 的注册方式

## 4. 文档与测试

- [x] 4.1 新建 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.md`，参考 `data_source_tc_teo_origin_acl.md` 格式编写数据源文档主体：一句话描述（含 TEO）、Example Usage、Import 不需要，文档生成由 `make doc` 完成
- [x] 4.2 新建 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family_test.go`，使用 gomonkey mock 云 API `DescribeAvailableOriginACLFamily`，对 service 层 `DescribeTeoAvailableOriginAclFamilyByFilter` 与数据源 schema 校验进行单元测试

## 5. 验证

- [x] 5.1 验证 provider.go 与 provider.md 注册项完整、数据源可被 Terraform 识别
- [x] 5.2 验证生成的单元测试代码在当前环境下可正确构建执行