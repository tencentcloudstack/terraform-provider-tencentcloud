## Context

腾讯云 tdmysql（TDSQL-C MySQL，SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122`）目前未在 Terraform Provider 中提供实例资源管理能力。本次新增 `tencentcloud_tdmysql_db_instance`（RESOURCE_KIND_GENERAL）资源，实现实例的创建、查询、改名、隔离。

接口语义关键点（基于 vendor 目录下的云 API 校验）：

- `CreateDBInstances`：批量创建实例，入参 35 个，出参 `InstanceIds`（`[]*string`）与 `FlowId`（`*int64`）。创建为异步流程，返回 FlowId。接口注释明确"批量创建实例功能"。
- `DescribeDBInstanceDetail`：查询实例详情，入参 `InstanceId`，出参约 60 个字段（含 `InstanceIds`/`FlowId` 之外的全部只读属性，以及创建时传入的部分参数回显）。
- `ModifyInstanceName`：仅修改实例名称，入参 `InstanceId` + `InstanceName`。
- `IsolateDBInstance`：批量隔离实例，入参 `InstanceIds`（`[]*string`），出参 `SuccessInstanceIds`/`FailedInstanceIds`。tdmysql 的销毁语义为"隔离"。
- `DescribeFlow`：查询异步任务流程状态，入参 `FlowId`，出参 `Status`（running/success/paused/failed）。用于创建后轮询。

CRUD 参数一致性核对结论：
- Create（`CreateDBInstances`）入参与 Read（`DescribeDBInstanceDetail`）出参中，可写的"创建参数"在 Read 中均有回显（zone、vpc_id、subnet_id、spec_code、disk、storage_node_num、replications、full_replications、create_version、instance_name、init_params、resource_tags、storage_node_cpu、storage_node_mem、pay_mode、vport、instance_type、storage_type、az_mode、instance_mode、template_id、sql_mode、auto_scale_config、encryption_enable 等）。
- Update 仅 `ModifyInstanceName` 可用 → 除 `instance_name` 外的顶层创建参数均为不可变（Immutable）。
- Delete 用 `IsolateDBInstance`。

参考样板：代码风格严格对齐 `tencentcloud_igtm_strategy`（RESOURCE_KIND_GENERAL，含 import、复合/单 ID、retry 块、immutableArgs、service 层）。服务层与 client 访问对齐 `tencentcloud_mongodb` / `tencentcloud_igtm`。

## Goals / Non-Goals

**Goals:**
- Schema 字段名与 `CreateDBInstances` 入参 1:1 映射（snake_case），类型按 SDK 指针类型正确转换（string/int64/bool/float64/list/嵌套 struct）。
- 资源 ID = `CreateDBInstances` 返回的 `InstanceIds[0]`（单实例 ID，支持 import）。
- Create 完成后通过 `DescribeFlow` 轮询 FlowId 直到 `success`，再用 `DescribeDBInstanceDetail` 回填状态。
- Update 仅处理 `instance_name` 变更（调 `ModifyInstanceName`）；其余顶层字段加入 `immutableArgs`，命中变更返回 error。
- Delete 调 `IsolateDBInstance`，校验目标实例在 `SuccessInstanceIds` 中。
- 全部 SDK 调用包裹 `resource.Retry(Write/ReadRetryTimeout, ...)`，retry 块内仅做接口调用与 nil 检查；set id / set 字段放到 retry 块外。
- 在 `connectivity/client.go` 新增 `UseTdmysqlV20211122Client()` 与 `tdmysqlv20211122Conn` 字段、import。
- 在 `provider.go` 注册资源并在 `provider.md` 声明。
- 单元测试使用 gomonkey mock 云 API（非 terraform 验收测试套件），用 `go test -gcflags=all=-l` 跑通。

**Non-Goals:**
- 不实现 tdmysql 的其他资源/数据源（如节点、参数模板、备份等）——仅本实例资源。
- 不暴露 `instance_count` > 1 的多实例批量管理语义给用户作为多资源：`CreateDBInstances` 虽支持批量，但 Terraform 资源语义为单资源，故 `instance_count` 作为可选入参透传，state id 取返回列表的第一项。
- 不实现 `CancelIsolateDBInstances`（取消隔离）/恢复逻辑——超出 CRUD 范围。
- 不修改任何既有资源/数据源/service 方法。
- 不生成 `_extension.go` 文件。
- 不在资源 go 文件开头添加注释。
- `website/docs/` 与 `.changelog/` 文件不在本阶段生成，统一由收尾阶段（tfpacer-finalize / `make doc`）处理。

## Decisions

### D1 — 资源 ID 与批量创建语义
`CreateDBInstances` 返回 `InstanceIds`（列表）+ `FlowId`。Terraform 资源模型为单资源，故：
- 资源 ID = `*response.Response.InstanceIds[0]`。
- `instance_count` 作为可选入参（默认由 API 决定）透传给 `CreateDBInstances.InstanceCount`，不在 Terraform 层做"展开成多个 resource"。
- `instance_ids`（完整列表）与 `flow_id` 作为 Computed 属性回显，便于用户排障。

**理由**：与 Provider 内其他批量创建但单资源管理的资源一致；避免引入 `count` 元编程带来的状态复杂度。

### D2 — Create 后异步轮询
`CreateDBInstances` 是异步接口（返回 FlowId）。创建流程：
1. 调 `CreateDBInstances`，取 `InstanceIds[0]` 作为待设 id（先不 SetId，避免空 id 触发状态混乱）。
2. 校验返回值非空（`Response`/`InstanceIds` 为空 → `NonRetryableError`）。
3. 若返回 `FlowId` 非 nil，用 `resource.StateChangeConf` 或 `resource.Retry` 轮询 `DescribeFlow`，`success` 视为完成，`failed` 返回错误，`running` 继续。
4. `d.SetId(instanceIds[0])`。
5. 调 Read 回填状态。

**理由**：用户需求"对于异步接口，调用完后要调用 Read 接口轮询直到接口生效"。`DescribeFlow` 提供了标准的流程状态查询能力，比直接轮询 `DescribeDBInstanceDetail`（实例状态枚举复杂）更直接可靠。

### D3 — Update 的 immutableArgs 策略
`ModifyInstanceName` 仅能改名称。Update 实现：
```go
immutableArgs := []string{
    "zone", "vpc_id", "subnet_id", "spec_code", "disk", "storage_node_num",
    "replications", "instance_count", "full_replications", "create_version",
    "resource_tags", "init_params", "time_unit", "time_span", "storage_node_cpu",
    "storage_node_mem", "pay_mode", "mc_num", "vport", "zones", "auto_voucher",
    "voucher_ids", "instance_type", "storage_type", "az_mode", "instance_mode",
    "template_id", "sql_mode", "auto_scale_config", "security_group_ids",
    "user_name", "password", "encryption_enable",
}
```
遍历 `immutableArgs`，若 `d.HasChange(v)` 返回 `fmt.Errorf("argument `%s` cannot be changed", v)`。
若 `d.HasChange("instance_name")` → 调 `ModifyInstanceName`（`InstanceId` = `d.Id()`，`InstanceName` = 新值）。

**理由**：遵循"若资源只有部分更新接口，则不可变字段加入 immutableArgs"的项目规范。`instance_name` 是唯一可通过 `ModifyInstanceName` 修改的顶层字段。

### D4 — Read 的 nil 检查与清空 id 顺序
Read 中：
- retry 块内调 `DescribeDBInstanceDetail`，检查 `result == nil || result.Response == nil` → `NonRetryableError`（避免短暂波动清空 id）。
- retry 块外：若 `response.Response.InstanceId == nil`（实例不存在），**先** `log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())`，**再** `d.SetId("")`，返回 nil。
- 回填字段前逐个判断 `xx != nil` 才 `d.Set(...)`。

**理由**：遵循项目规范"打印日志保留现场再清空 id"与"set 前判断 nil"。

### D5 — Delete 的成功校验
`IsolateDBInstance` 返回 `SuccessInstanceIds`/`FailedInstanceIds`。Delete：
- 请求 `InstanceIds = []*string{&instanceId}`。
- retry 块内调用并检查 `result.Response.SuccessInstanceIds` 是否包含目标 id；包含则成功，否则 `NonRetryableError`。
- 成功的 SetId 等操作放到 retry 块外。

**理由**：隔离是批量接口，需确认目标实例确实被隔离成功。

### D6 — 复杂字段 Schema 结构
- `resource_tags`：`TypeList`，嵌套 Resource `{tag_key: TypeString, tag_value: TypeString}`。
- `init_params`：`TypeList`，嵌套 Resource `{param: TypeString, value: TypeString}`。
- `auto_scale_config`：`TypeList`（`MaxItems: 1`），嵌套 Resource `{range_min: TypeFloat, range_max: TypeFloat}`。
- `zones`、`voucher_ids`、`security_group_ids`：`TypeList`，`Elem: schema.TypeString`。
- `node`（只读，`[]*NodeInfo`）：`TypeList`，嵌套 Resource 含 NodeInfo 的字段（ip、type、node_id、port、zone、host、cpu、mem、data_disk 等）。
- `analysis_relation_infos`（只读，`[]*AnalysisRelationInfo`）：`TypeList` 嵌套。
- `analysis_instance_info`、`maintenance_window`：`TypeList`（`MaxItems: 1`）嵌套。

**理由**：与 igtm 等资源处理嵌套结构的方式一致；`MaxItems: 1` 用于单对象字段。

### D7 — client 访问方法
在 `connectivity/client.go`：
- import：`tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"`
- `TencentCloudClient` 结构体新增字段 `tdmysqlv20211122Conn *tdmysqlv20211122.Client`。
- 新增方法 `UseTdmysqlV20211122Client() *tdmysqlv20211122.Client`，懒加载，`NewClientProfile(300)` + `LogRoundTripper`，与 `UseIgtmV20231024Client` 同款。

**理由**：所有服务 SDK 客户端均通过 `connectivity.TencentCloudClient` 上的 `UseXxxClient()` 方法获取，统一懒加载与日志拦截。

### D8 — 单元测试策略
测试文件 `resource_tc_tdmysql_db_instance_test.go`，包名 `tdmysql_test`，使用 gomonkey：
- mock `mockMeta.GetAPIV3Conn().UseTdmysqlV20211122Client()` 返回 `&tdmysqlv20211122.Client{}`。
- mock `CreateDBInstancesWithContext` / `DescribeFlowWithContext` / `DescribeDBInstanceDetailWithContext` / `ModifyInstanceNameWithContext` / `IsolateDBInstanceWithContext`。
- 用 `schema.TestResourceDataRaw` 构造 `*schema.ResourceData`，调用 `res.Create/Read/Update/Delete(d, meta)`，`assert` 校验。
- 用 `go test -gcflags=all=-l` 跑通涉及文件。

**理由**：新增资源用 mock 单测（项目规范），不依赖 terraform 验收测试套件与真实云环境。

## Risks / Trade-offs

- **Risk**: `CreateDBInstances` 支持批量创建（`instance_count`），但 Terraform 仅取第一个实例 id 作为单资源 → 多实例场景下其余实例不在 state 中 → Mitigation: `instance_count` 默认不强制为 1（透传给 API），`instance_ids` Computed 属性回显完整列表，文档说明"资源管理列表中第一个实例，批量创建请使用 instance_count 并关注 instance_ids"。
- **Risk**: 创建异步轮询可能因 `DescribeFlow` 返回 `paused` 状态而无限等待 → Mitigation: 轮询使用带超时的 `StateChangeConf`/`resource.Retry`（`tccommon.ReadRetryTimeout` 量级或更长），`paused`/`failed` 视为终态失败。
- **Trade-off**: 不实现"真正销毁"（`DeleteDBInstance` 类接口）。tdmysql 的 `IsolateDBInstance` 是隔离语义，隔离后资源在云端仍保留一段时间 → 与 igtm 等资源"调隔离接口即视为删除"的惯例一致，文档说明销毁语义为隔离。
- **Trade-off**: Update 仅支持改名，其他规格变配（如扩容磁盘、变更副本数）tdmysql 暂未在本需求提供的接口中覆盖 → 变更这些字段会报错而非自动重建，符合 immutableArgs 规范；后续若有对应 Modify 接口可扩展。
