# 任务清单：fix-sqlserver-account-custom-timeouts

## 1. 新增 Timeouts 块 + 引入 time 包

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] 在 import 中新增 `"time"`
- [x] 在 `schema.Resource` 中新增 `Timeouts` 块：
  ```go
  Timeouts: &schema.ResourceTimeout{
      Create: schema.DefaultTimeout(30 * time.Minute),
      Delete: schema.DefaultTimeout(10 * time.Minute),
  },
  ```

---

## 2. Create 模块替换超时时间

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] 将 `CreateSqlserverAccount` 的 retry 超时从 `tccommon.WriteRetryTimeout` 改为 `d.Timeout(schema.TimeoutCreate)`

---

## 3. Delete 模块替换超时时间

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] 将确认账号存在的 retry 超时从 `tccommon.ReadRetryTimeout` 改为 `d.Timeout(schema.TimeoutDelete)`
- [x] 将 `DeleteSqlserverAccount` 的 retry 超时从 `tccommon.WriteRetryTimeout` 改为 `d.Timeout(schema.TimeoutDelete)`
- [x] 将等待删除完成的 retry 超时从 `tccommon.ReadRetryTimeout` 改为 `d.Timeout(schema.TimeoutDelete)`
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`

---

## 4. 编译验证

- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 总结

- **预计工作量**：极小（约 10 分钟）
- **风险等级**：极低（纯增量，不影响已有逻辑）
- **破坏性变更**：无
- **状态**: 已完成

---

## 5. 拆分 Create 模块的创建操作与等待操作（post-apply 补充）

**背景**: `CreateSqlserverAccount` service 方法内部同时包含 API 调用和 `WaitForTaskFinish`（硬编码 30 分钟轮询），导致 resource 层的 `d.Timeout(schema.TimeoutCreate)` 无法真正控制等待时间。需拆分为两个独立步骤。

### service 层

**文件**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

- [x] 新增 `CreateSqlserverAccountReturnFlowId` 方法：只调用 `CreateAccount` API，返回 `(flowId int64, errRet error)`，不调用 `WaitForTaskFinish`
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`

### resource 层

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] Create 模块改为两个独立 retry 块：
  1. `resource.Retry(tccommon.WriteRetryTimeout)` 调用 `CreateSqlserverAccountReturnFlowId` 获取 flowId
  2. `resource.Retry(d.Timeout(schema.TimeoutCreate))` 调用 `sqlserverService.WaitForTaskFinish` 等待 flowId 完成
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 6. 用直接轮询替换 WaitForTaskFinish（post-apply 补充）

**背景**: `WaitForTaskFinish` 内部自带 `resource.Retry(6*WriteRetryTimeout)` 轮询，外层的 `d.Timeout(Create)` 只会调用它一次，实际超时不受用户控制。需直接在 resource 层用 `d.Timeout` 控制 `DescribeFlowStatus` 的轮询。

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] 将等待块从 `WaitForTaskFinish` 改为直接构建 `DescribeFlowStatusRequest` 并轮询：
  - `status == SQLSERVER_TASK_RUNNING` → `RetryableError`
  - `status == SQLSERVER_TASK_FAIL` → `NonRetryableError`
  - 其他（成功）→ `nil`
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 7. 修正 Delete 模块各步骤的超时分配（post-apply 补充）

**背景**: Delete 三个 retry 块都用了 `d.Timeout(Delete)`，但实际上只有最后的轮询等待才需要自定义超时；前两步（确认存在、调删除 API）是快速操作，应使用全局常量。

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] Block 1（确认账号存在）：改为 `tccommon.ReadRetryTimeout`
- [x] Block 2（调 DeleteSqlserverAccount API）：改为 `tccommon.WriteRetryTimeout`
- [x] Block 3（轮询等账号消失）：保持 `d.Timeout(schema.TimeoutDelete)`
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 8. 拆分 Delete 第二步：API 调用与等待分离（post-apply 补充）

**背景**: `DeleteSqlserverAccount` service 方法内部同样捆绑了 `DeleteAccount` API 调用 + `ReadRetryTimeout(3m)` 轮询等待，与 Create 问题完全对称。需拆分，使 resource 层的 `d.Timeout(Delete)` 真正控制等待时间。

### service 层

**文件**: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

- [x] 新增 `DeleteSqlserverAccountOnly` 方法：只调用 `DeleteAccount` API，不包含轮询等待，出错时处理 `ResourceNotFound.InstanceNotFound`

### resource 层

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] Delete 改为三个 retry 块：
  1. `Retry(tccommon.ReadRetryTimeout)` 确认账号存在
  2. `Retry(tccommon.WriteRetryTimeout)` 调 `DeleteSqlserverAccountOnly` 只发 API
  3. `Retry(d.Timeout(schema.TimeoutDelete))` 轮询 `DescribeSqlserverAccountById`，等账号消失（`!has`）；原 Block 4 合并进来
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 9. 修复 Delete Block3 缺少异步状态判断（post-apply 补充）

**背景**: Block3 轮询只判断了 `!has`，丢失了原 `DeleteSqlserverAccount` 中对 `status == -1`（删除中）和其他异常状态的区分语义。

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] Block3 补充完整状态判断：
  - `!has` → `nil`（删除完成）
  - `status == -1` → `RetryableError`（删除中，继续等）
  - 其他 status → `NonRetryableError`（异常状态）
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过

---

## 10. 将 Delete Block3 拆为异步等待 + 最终验证两个独立块（post-apply 补充）

**背景**: 用户要求将原功能2（删除+等待）拆分后变成4个功能块：
1. 删除前查询（不变）
2. 调删除 API（不变）
3. 删除后异步等待：`d.Timeout(Delete)` 轮询等 status 从 -1 变为消失
4. 删除后验证查询：`ReadRetryTimeout` 最终确认账号已消失

当前 Block3 把等待和验证合并在一起，需要拆开。

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go`

- [x] 将当前 Block3 拆为两个独立块：
  - **新 Block3**（`d.Timeout(Delete)`）：等 `status == -1` → `RetryableError`；`!has` → `nil`；其他 → `NonRetryableError`
  - **新 Block4**（`ReadRetryTimeout`）：最终确认 `!has`，如果还存在则 `RetryableError`
- [x] 执行 `go fmt ./tencentcloud/services/sqlserver/`
- [x] `go build ./tencentcloud/services/sqlserver/` 确认编译通过
