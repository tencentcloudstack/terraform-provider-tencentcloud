## 1. 资源 schema 与 CRUD 实现

- [x] 1.1 创建 `tencentcloud/services/dlc/resource_tc_dlc_dms_table.go`，定义 `ResourceTencentCloudDlcDmsTable()` 返回 `*schema.Resource`，包含 Create/Read/Update/Delete 及 Importer（State: schema.ImportStatePassthrough）
- [x] 1.2 定义顶层 schema 字段：`asset`（TypeList，Asset 嵌套）、`type`（TypeString，Optional）、`db_name`（TypeString，Required）、`name`（TypeString，Required）、`datasource_connection_name`（TypeString，Optional）、`storage_size`（TypeInt，Optional）、`record_count`（TypeInt，Optional）、`life_time`（TypeInt，Optional）、`data_update_time`/`struct_update_time`/`last_access_time`/`view_original_text`/`view_expanded_text`（TypeString，Optional）、`sds`（TypeList，DMSSds 嵌套）、`columns`（TypeList，DMSColumn 嵌套）、`partition_keys`（TypeList，DMSColumn 嵌套）、`partitions`（TypeList，DMSPartition 嵌套）、`delete_data`（TypeBool，Optional）、`env_props`（TypeList，KVPair 嵌套）、`schema_name`（TypeString，Computed）、`retention`（TypeInt，Computed）
- [x] 1.3 定义嵌套 schema：`asset`（含 id/name/guid/catalog/description/owner/owner_account/perm_values/params/biz_params/data_version/create_time/modified_time/datasource_id，其中 perm_values/params/biz_params 为 KVPair 列表）、`sds`（含 location/input_format/output_format/num_buckets/compressed/stored_as_sub_directories/serde_lib/serde_name/bucket_cols/serde_params/params/sort_cols/dms_cols/sort_columns，其中 serde_params/params 为 KVPair 列表、sort_cols 为 DMSColumnOrder、dms_cols 为 DMSColumn 列表、sort_columns 为 DMSColumnOrder 列表）、`columns`/`partition_keys`（DMSColumn：name/description/type/position/params/biz_params/is_partition）、`partitions`（DMSPartition：database_name/schema_name/table_name/data_version/name/values/storage_size/record_count/create_time/modified_time/last_access_time/params/sds/datasource_connection_name）、`env_props`（KVPair：key/value）
- [x] 1.4 实现 `resourceTencentCloudDlcDmsTableCreate`：构造 `CreateDMSTableRequest` 填入 Create 入参字段；`resource.Retry(tccommon.WriteRetryTimeout)` 内仅调用 `CreateDMSTableWithContext`，错误用 `tccommon.RetryError` 包装；retry 外检查 response/Response 非空，空则返回 NonRetryableError；成功后 `d.SetId(strings.Join([]string{db_name, name}, tccommon.FILED_SP))`，并调用 Read 回查
- [x] 1.5 实现 `resourceTencentCloudDlcDmsTableRead`：从 `d.Id()` 按 `FILED_SP` 拆分 db_name/name；调用 service 层 `DescribeDmsTableById`（内部 ReadRetryTimeout）；response 为空时先 `log.Printf("[CRUD] dlc_dms_table id=%s", d.Id())` 再 `d.SetId("")`；非空时对每个 Response 字段做 nil 判断后 set（包括 Computed 的 schema_name、retention）；嵌套结构按 SDK 结构转换为 []map[string]interface{}
- [x] 1.6 实现 `resourceTencentCloudDlcDmsTableUpdate`：构造 `AlterDMSTableRequest`；对 `db_name`/`name` 用 `d.GetChange` 取旧值填入 `CurrentDbName`/`CurrentName`，新值填入 `DbName`/`Name`；其余 AlterDMSTable 支持字段从 schema 读取填入；`schema_name`/`retention` 不加入请求；`resource.Retry(tccommon.WriteRetryTimeout)` 内仅调用 `AlterDMSTableWithContext`，错误用 `tccommon.RetryError` 包装；成功后若 db_name/name 变更则更新 `d.SetId`，再调用 Read 回查
- [x] 1.7 实现 `resourceTencentCloudDlcDmsTableDelete`：从 `d.Id()` 拆分 db_name/name；构造 `DropDMSTableRequest` 填入 db_name/name/delete_data/env_props/datasource_connection_name；`resource.Retry(tccommon.WriteRetryTimeout)` 内仅调用 `DropDMSTableWithContext`，错误用 `tccommon.RetryError` 包装
- [x] 1.8 在 service 层文件（`tencentcloud/services/dlc/service_tencentcloud_dlc.go`）中新增 `DescribeDmsTableById` 方法：内部 `resource.Retry(tccommon.ReadRetryTimeout)` 调用 `DescribeDMSTableWithContext`，request 填入 DbName/Name；retry 块内检查 response/Response 为空返回 NonRetryableError；返回 `*DescribeDMSTableResponseParams`

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的资源 map 中新增 `"tencentcloud_dlc_dms_table": dlc.ResourceTencentCloudDlcDmsTable()`
- [x] 2.2 在 `tencentcloud/provider.md` 中新增 `tencentcloud_dlc_dms_table` 资源条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/dlc/resource_tc_dlc_dms_table.md`：一句话描述带上 DLC 产品名（"Provides a resource to create a DLC DMS table"）、Example Usage（涉及 json 字符串用 jsonencode）、Import 部分说明使用 `db_name#name` 复合 ID；不添加 Argument Reference 与 Attribute Reference（由 make doc 自动生成）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/dlc/resource_tc_dlc_dms_table_test.go`，使用 gomonkey mock `UseDlcClient` 的 `CreateDMSTableWithContext`/`DescribeDMSTableWithContext`/`AlterDMSTableWithContext`/`DropDMSTableWithContext`，覆盖 Create/Read/Update/Delete 业务逻辑
- [x] 4.2 执行 `go test -gcflags=all=-l ./tencentcloud/services/dlc/ -run TestResourceTencentCloudDlcDmsTable` 确保单元测试通过

## 5. 验证与收尾

- [x] 5.1 检查所有函数返回的 error 均被处理，必定不出错的用 `_ = func()` 赋值给 `_`
- [x] 5.2 核对 Create/Update/Delete 入参字段均在对应云 API 接口中存在，schema 字段类型与 vendor models.go 一致
- [ ] 5.3 通过 tfpacer-finalize skill 执行 gofmt、make doc 生成 website/docs 文档、生成 changelog 文件并推送
