## 1. 基础设施（client 访问方法）

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 中新增 import：`tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"`
- [x] 1.2 在 `TencentCloudClient` 结构体中新增字段 `tdmysqlv20211122Conn *tdmysqlv20211122.Client`
- [x] 1.3 新增 `UseTdmysqlV20211122Client()` 方法（懒加载、`NewClientProfile(300)`、`LogRoundTripper`），风格对齐 `UseIgtmV20231024Client`

## 2. 服务层实现

- [x] 2.1 创建 `tencentcloud/services/tdmysql/service_tencentcloud_tdmysql.go`，定义 `NewTdmysqlService` 与 `TdmysqlService` 结构体（含 `client *connectivity.TencentCloudClient`）
- [x] 2.2 实现 `DescribeTdmysqlDbInstanceById(ctx, instanceId)`：调 `DescribeDBInstanceDetail`，`ReadRetryTimeout` retry，retry 内 nil 检查返回 `NonRetryableError`，含 `ratelimit.Check`
- [x] 2.3 实现 `IsolateTdmysqlDbInstance(ctx, instanceId)`：调 `IsolateDBInstance`，`WriteRetryTimeout` retry，校验 `SuccessInstanceIds` 包含目标 id
- [x] 2.4 实现 `DescribeTdmysqlFlow(ctx, flowId)`：调 `DescribeFlow`，`ReadRetryTimeout` retry，返回 Status 字符串

## 3. 资源 Schema 实现

- [x] 3.1 创建 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.go`，定义 `ResourceTencentCloudTdmysqlDbInstance()` 返回 `*schema.Resource`（含 Create/Read/Update/Delete、Importer、Schema、Timeouts 块）
- [x] 3.2 定义 Create 输入参数 schema（zone/vpc_id/subnet_id/spec_code/disk/storage_node_num/replications/instance_count/full_replications/create_version/instance_name/resource_tags/init_params/time_unit/time_span/storage_node_cpu/storage_node_mem/pay_mode/mc_num/vport/zones/auto_voucher/voucher_ids/instance_type/storage_type/az_mode/instance_mode/template_id/sql_mode/auto_scale_config/security_group_ids/user_name/password/encryption_enable），类型按 SDK 指针类型映射
- [x] 3.3 定义嵌套 schema：resource_tags(tag_key/tag_value)、init_params(param/value)、auto_scale_config(MaxItems=1, range_min/range_max 为 TypeFloat)、zones/voucher_ids/security_group_ids(TypeList+TypeString)、node/analysis_relation_infos/analysis_instance_info/maintenance_window（只读嵌套）
- [x] 3.4 定义 Computed 只读属性（instance_id/vip/status/create_time/update_time/char_set/region/status_desc/renew_flag/expire_at/isolated_at/disk_usage/binlog_status/standby_flag/binlog_type/timing_modify_instance_flag/columnar_*/dumper_*/template_name/analysis_mode/is_switch_full_replications_enable/instance_category/is_support_columnar/encryption_kms_region/instance_ids/flow_id 等）

## 4. 资源 CRUD 实现

- [x] 4.1 实现 `resourceTencentCloudTdmysqlDbInstanceCreate`：填充请求参数、`WriteRetryTimeout` retry 调 `CreateDBInstances`、校验 Response/InstanceIds 非空（空返回 NonRetryableError）、记录 logId 与 id、轮询 `DescribeFlow` 直到 success、`d.SetId(InstanceIds[0])`、调 Read
- [x] 4.2 实现 `resourceTencentCloudTdmysqlDbInstanceRead`：用 `d.Id()` 构造请求、调 service `DescribeTdmysqlDbInstanceById`、空时先日志后 `d.SetId("")`、逐字段 nil 判断后 set
- [x] 4.3 实现 `resourceTencentCloudTdmysqlDbInstanceUpdate`：遍历 immutableArgs 命中变更返回 error；`instance_name` 变更时调 `ModifyInstanceName`（`WriteRetryTimeout` retry），成功后调 Read
- [x] 4.4 实现 `resourceTencentCloudTdmysqlDbInstanceDelete`：用 `d.Id()` 构造 `IsolateDBInstance` 请求、调 service `IsolateTdmysqlDbInstance`、`WriteRetryTimeout` retry
- [x] 4.5 所有 CRUD 函数添加 `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(d, meta)`，日志中资源名统一使用 `tdmysql_db_instance`

## 5. Provider 注册

- [x] 5.1 在 `tencentcloud/provider.go` 中 import `github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql`（别名 `tdmysql`）
- [x] 5.2 在 `provider.go` 的 ResourcesMap 中注册 `"tencentcloud_tdmysql_db_instance": tdmysql.ResourceTencentCloudTdmysqlDbInstance()`
- [x] 5.3 在 `tencentcloud/provider.md` 资源列表中声明 `tencentcloud_tdmysql_db_instance`（用于 `make doc`）

## 6. 单元测试（gomonkey mock）

- [x] 6.1 创建 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance_test.go`，包名 `tdmysql_test`，定义 mockMeta 结构体与 `newMockMeta` 辅助函数
- [x] 6.2 实现 `TestTdmysqlDbInstance_Create_Success`：mock `UseTdmysqlV20211122Client`/`CreateDBInstancesWithContext`/`DescribeFlowWithContext`/`DescribeDBInstanceDetailWithContext`，验证 `d.Id()` 正确
- [x] 6.3 实现 `TestTdmysqlDbInstance_Read_Success`：mock `DescribeDBInstanceDetailWithContext` 返回详情，验证字段回填
- [x] 6.4 实现 `TestTdmysqlDbInstance_Update_Success`：mock `ModifyInstanceNameWithContext`，验证改名成功
- [x] 6.5 实现 `TestTdmysqlDbInstance_Delete_Success`：mock `IsolateDBInstanceWithContext` 返回含目标 id 的 `SuccessInstanceIds`，验证隔离成功
- [x] 6.6 实现 Read 资源不存在场景：mock 返回空，验证先日志后 `d.SetId("")` 且无 error
- [x] 6.7 使用 `go test -gcflags=all=-l` 跑通该测试文件

## 7. 资源文档

- [x] 7.1 创建 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.md`：一句话描述（带上 TDSQL-C MySQL 产品名，格式 "Provides a resource to ..."）
- [x] 7.2 添加 Example Usage 部分（含必要字段示例，json 字符串场景用 `jsonencode()`）
- [x] 7.3 添加 Import 部分（说明使用实例 ID），不添加 Argument Reference/Attribute Reference（由工具生成）
