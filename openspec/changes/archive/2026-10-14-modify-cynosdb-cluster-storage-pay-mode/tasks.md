## 1. immutableArgs 限制移除

- [x] 1.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go` 的 Update 函数中，从 `immutableArgs` 列表（约第 647-658 行）中移除 `"storage_pay_mode"`

## 2. 服务层实现

- [x] 2.1 在 `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go` 中新增 `TransferStoragePrepayToPostpay` 方法
  - 构造 `cynosdb.NewTransferStoragePrepayToPostpayRequest()`，设置 `ClusterId`
  - 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，失败时用 `tccommon.RetryError(err)` 包装
  - 在 retry 块内调用 `ratelimit.Check(request.GetAction())`
  - 成功后返回 `response.Response.DealNames`（`[]*string`）

## 3. Update 函数 storage_pay_mode 变更处理

- [x] 3.1 在 `resourceTencentCloudCynosdbClusterUpdate` 函数中新增 `d.HasChange("storage_pay_mode")` 处理分支，建议置于 `storage_limit` 处理之后
- [x] 3.2 使用 `d.GetChange("storage_pay_mode")` 获取新旧值，校验仅允许旧值 `1`（预付费）→ 新值 `0`（后付费），反向返回错误 `storage_pay_mode only supports transfer from prepaid(1) to postpaid(0)`
- [x] 3.3 校验通过后调用 `cynosdbService.TransferStoragePrepayToPostpay(ctx, clusterId)` 获取 `dealNames`，检查 `len(dealNames) < 1` 时返回错误并打印 clusterId
- [x] 3.4 取 `dealNames[0]` 作为 DealName，构造 `DescribeResourcesByDealNameRequest`，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 轮询，参考 create 模块（第 249-277 行）实现
  - 错误码 `InvalidParameterValue.DealNameNotFound` 作为 `RetryableError` 重试
  - 其他错误用 `tccommon.RetryError` 包装
  - 成功返回 `BillingResourceInfos` 后停止轮询

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.md` 中 `storage_pay_mode` 参数说明，标注支持从预付费(`1`)转为后付费(`0`)，为单向操作
- [x] 4.2 在 `website/docs/` 下的文档通过收尾阶段的 `make doc` 命令自动生成（禁止手动编写 website/ 文件）

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_test.go` 中补充单元测试用例（使用 gomonkey mock 云 API，不使用 terraform 测试套件）
  - mock `TransferStoragePrepayToPostpay` 与 `DescribeResourcesByDealName` 接口
  - 测试预付费转后付费成功场景
  - 测试拒绝后付费转预付费场景

## 6. 验证（收尾阶段执行）

- [ ] 6.1 执行 `gofmt` 格式化变更的 Go 代码（由 tfpacer-finalize skill 统一执行）
- [ ] 6.2 执行 `make doc` 生成 `website/docs/` 文档（由 tfpacer-finalize skill 统一执行）
- [ ] 6.3 创建 `.changelog/` 下的变更日志文件（由 tfpacer-finalize skill 统一执行）

## 依赖关系

- 任务 1.1、2.1 为基础，需先完成
- 任务 3 依赖任务 1.1 和 2.1 完成
- 任务 4 依赖任务 3 完成（确认参数行为后更新文档）
- 任务 5 依赖任务 3 完成
- 任务 6 在所有代码修改完成后由收尾阶段统一执行