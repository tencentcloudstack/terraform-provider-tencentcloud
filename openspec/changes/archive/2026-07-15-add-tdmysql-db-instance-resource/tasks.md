## 1. SDK Client Helper 基础设施

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 的 import 块中新增 tdmysql SDK 包导入：`tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"`
- [x] 1.2 在 `tencentcloud/connectivity/client.go` 的 `TencentCloudClient` 结构体中新增连接字段 `tdmysqlv20211122Conn *tdmysqlv20211122.Client`
- [x] 1.3 在 `tencentcloud/connectivity/client.go` 中新增 `UseTdmysqlV20211122Client()` 方法，返回 `*tdmysqlv20211122.Client`，参考现有 `UseTdcpgClient`/`UseIgtmV20231024Client` 模式实现（含 lazy init 与 `LogRoundTripper`）

## 2. 资源 Schema 定义与 CRUD 函数实现

- [x] 2.1 创建资源文件 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.go`，定义 `ResourceTencentCloudTdmysqlDbInstance()` 函数，包含完整 schema（Create 入参字段、Create 出参 Computed 字段、Read 出参 Computed 字段）和 CRUD + Import 支持
- [x] 2.2 实现 Create 函数 `resourceTencentCloudTdmysqlDbInstanceCreate`：组装 `CreateDBInstancesRequest`（将所有配置的 Create 入参字段映射到 API 参数），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `CreateDBInstancesWithContext`；调用成功后检查返回值是否为空（response/Response/InstanceIds 为空则返回 `NonRetryableError`）；获取 `FlowId` 后调用 `DescribeFlow` 轮询直至 `Status=success`（使用 `resource.Retry` 包裹，非 success 返回 `resource.RetryableError`）；流程成功后取 `InstanceIds` 首元素作为 `instance_id`，设置 `d.SetId()`，并 set `instance_ids` 和 `flow_id` 字段
- [x] 2.3 实现 Read 函数 `resourceTencentCloudTdmysqlDbInstanceRead`：以 `d.Id()` 作为 `InstanceId` 调用 `DescribeDBInstanceDetailWithContext`（在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中调用）；若返回空则先 `log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())` 再 `d.SetId("")`；逐字段检查 nil 后 setXX() 回写所有 Read 出参字段
- [x] 2.4 实现 Update 函数 `resourceTencentCloudTdmysqlDbInstanceUpdate`：定义 immutableArgs 数组（除 `instance_name` 外的所有顶层 Create 入参字段），遍历检查 `d.HasChange()`，若有变更返回 error；若 `instance_name` 变更，组装 `ModifyInstanceNameRequest`（InstanceId 来自 d.Id()），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `ModifyInstanceNameWithContext`
- [x] 2.5 实现 Delete 函数 `resourceTencentCloudTdmysqlDbInstanceDelete`：将 `d.Id()` 包装为单元素数组作为 `InstanceIds`，组装 `IsolateDBInstanceRequest`，在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `IsolateDBInstanceWithContext`；检查 `SuccessInstanceIds` 是否包含目标 instance_id

## 3. Provider 注册与文档索引

- [x] 3.1 在 `tencentcloud/provider.go` 的 import 块中新增 tdmysql 服务包导入：`"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"`
- [x] 3.2 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中添加 `"tencentcloud_tdmysql_db_instance": tdmysql.ResourceTencentCloudTdmysqlDbInstance()` 注册
- [x] 3.3 在 `tencentcloud/provider.md` 中添加 `tencentcloud_tdmysql_db_instance` 资源条目（收尾阶段通过 `make doc` 生成，本步骤仅在 provider.md 索引中添加条目）

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.md` 文件，包含一句话描述（提及 TDSQL-C for MySQL/tdmysql 产品名，格式 "Provides a resource to ..."）、Example Usage（含 zone、vpc_id、subnet_id、spec_code、disk、storage_node_num、replications、instance_name 等配置，涉及 JSON 字段使用 `jsonencode()`）、Import 部分（说明使用 instance_id 导入）

## 5. 单元测试

- [x] 5.1 创建测试文件 `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance_test.go`，使用 gomonkey mock 方式编写：Create 测试（mock `CreateDBInstances` 和 `DescribeFlow`）、Read 测试（mock `DescribeDBInstanceDetail`）、Update 测试（mock `ModifyInstanceName`）、Delete 测试（mock `IsolateDBInstance`）
- [x] 5.2 使用 `go test -gcflags=all=-l` 运行 `tencentcloud/services/tdmysql/` 下的单元测试，确保所有测试通过

## 6. 代码正确性验证

- [x] 6.1 验证 Create 入参字段均在 `CreateDBInstancesRequest` 中存在，Read 出参字段均在 `DescribeDBInstanceDetailResponse` 中存在，Update 入参字段在 `ModifyInstanceNameRequest` 中存在，Delete 入参字段在 `IsolateDBInstanceRequest` 中存在
- [x] 6.2 验证所有函数返回的 error 均被检查，必定不出错的函数用 `_ = func()` 将 err 赋值给 `_`
- [x] 6.3 验证 Create 函数中检查云 API 返回值为空（response/Response/InstanceIds 为空）时返回 `NonRetryableError`
- [x] 6.4 验证 Read 函数中设置字段前判断 Response 字段是否为 nil，nil 则不调用 setXX()
