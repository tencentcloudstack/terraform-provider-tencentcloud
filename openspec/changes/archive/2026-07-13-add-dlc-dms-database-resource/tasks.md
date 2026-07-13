## 1. 资源代码实现

- [x] 1.1 创建文件 `tencentcloud/services/dlc/resource_tc_dlc_dms_database.go`，包名 `dlc`，导入 `dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"`、`tccommon`、`helper`、`resource`、`schema` 等，严格参考 `resource_tc_igtm_strategy.go` 代码风格（不开头加注释）。
- [x] 1.2 实现 `ResourceTencentCloudDlcDmsDatabase()` 的 Schema 定义：顶层 `name`(Required,ForceNew)、`schema_name`(Required,ForceNew)、`datasource_connection_name`(Required,ForceNew)、`location`(Optional)、`delete_data`(Optional,TypeBool)、`cascade`(Optional,TypeBool)、`asset`(Optional,TypeList,MaxItems=1)；`asset` 嵌套块含 `id`(Computed)、`name`、`guid`(Computed)、`catalog`、`description`、`owner`、`owner_account`、`perm_values`(TypeList)、`params`(TypeList)、`biz_params`(TypeList)、`data_version`(Computed)、`create_time`(Computed)、`modified_time`(Computed)、`datasource_id`(Computed)；`perm_values`/`params`/`biz_params` 元素含 `key`(Required)、`value`(Optional)。声明 `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。
- [x] 1.3 实现 `resourceTencentCloudDlcDmsDatabaseCreate`：构建 `CreateDMSDatabaseRequest`（设置 Name/SchemaName/DatasourceConnectionName/Location/Asset），用 `resource.Retry(tccommon.WriteRetryTimeout,...)` 包裹 `UseDlcClient().CreateDMSDatabaseWithContext`；retry 内错误用 `tccommon.RetryError(e)` 包装；成功后检查 `response==nil||response.Response==nil` 返回 `NonRetryableError`；打印 logId 与 name；设置复合 ID `name#schema_name#datasource_connection_name`（`tccommon.FILED_SP` 分隔）；调用 Read。
- [x] 1.4 实现 `resourceTencentCloudDlcDmsDatabaseRead`：从 `d.Id()` 按 `tccommon.FILED_SP` 拆出 name/schema_name/datasource_connection_name；构建 `DescribeDMSDatabaseRequest`（设置 Name/SchemaName/DatasourceConnectionName，不设 Pattern）；用 `resource.Retry(tccommon.ReadRetryTimeout,...)` 包裹调用；retry 内若 `response==nil||response.Response==nil` 返回 `NonRetryableError`（不清空 id）；retry 外若资源不存在先 `log.Printf("[CRUD] dlc_dms_database id=%s", d.Id())` 再 `d.SetId("")`；每个 `d.Set` 前判断对应字段非 nil；将 `asset` 嵌套结构（含 3 个 KVPair 列表）正确反序列化到 state。
- [x] 1.5 实现 `resourceTencentCloudDlcDmsDatabaseUpdate`：定义 `immutableArgs = []string{"name", "schema_name", "datasource_connection_name"}`，遍历若 `d.HasChange(v)` 返回 error；当 `location` 或 `asset` 变更时构建 `AlterDMSDatabaseRequest`（`CurrentName` 取当前 name，设置 SchemaName/DatasourceConnectionName/Location/Asset），用 `resource.Retry(tccommon.WriteRetryTimeout,...)` 包裹 `AlterDMSDatabaseWithContext`；成功后调用 Read。
- [x] 1.6 实现 `resourceTencentCloudDlcDmsDatabaseDelete`：从 `d.Id()` 拆出 name/datasource_connection_name；构建 `DropDMSDatabaseRequest`（Name/DatasourceConnectionName/DeleteData/Cascade，后两者默认 false）；用 `resource.Retry(tccommon.WriteRetryTimeout,...)` 包裹 `DropDMSDatabaseWithContext`；检查函数返回 error 并处理。
- [x] 1.7 检查所有函数返回的 error：必定不出错的用 `_ = func()` 将 err 赋值给 `_`，避免未使用变量错误。

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中（与现有 `tencentcloud_dlc_*` 资源同组，按字母顺序）添加 `"tencentcloud_dlc_dms_database": dlc.ResourceTencentCloudDlcDmsDatabase()`。

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/dlc/resource_tc_dlc_dms_database_test.go`，包名 `dlc_test`，使用 gomonkey mock 云 API（`CreateDMSDatabase`/`DescribeDMSDatabase`/`AlterDMSDatabase`/`DropDMSDatabase`），不使用 terraform 测试套件。
- [x] 3.2 编写测试用例覆盖 Create→Read→Update→Delete 业务逻辑（含 asset 嵌套结构、复合 ID、immutableArgs 检查、nil 响应处理），使用 `go test -gcflags=all=-l` 跑通所涉及的测试文件。

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/dlc/resource_tc_dlc_dms_database.md`：一句话描述带上 DLC 产品名（"Provides a resource to ..."）；Example Usage 部分（涉及 json 字符串值时用 `jsonencode()`）；Import 部分说明使用联合 id `name#schema_name#datasource_connection_name`；不添加 Argument Reference / Attribute Reference 部分。

## 5. 验证

- [ ] 5.1 由收尾阶段通过 `make doc` 生成 `website/docs/r/dlc_dms_database.html.markdown`，并在 `tencentcloud/provider.md` 自动添加资源索引条目。
- [ ] 5.2 由收尾阶段执行 `gofmt` 格式化新增/修改的 Go 文件，并生成 `.changelog` 文件。
