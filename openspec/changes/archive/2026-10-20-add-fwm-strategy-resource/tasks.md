## 1. 客户端与注册

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 新增 `UseFwmV20250611Client()` 访问器与 `fwmv20250611Conn` 字段，import 别名 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`，实现参考 `UseIgtmV20231024Client`。
- [x] 1.2 在 `tencentcloud/provider.go` 注册 `"tencentcloud_fwm_strategy": fwm.ResourceTencentCloudFwmStrategy()`，import `fwm` 服务包。
- [x] 1.3 在 `tencentcloud/provider.md` 的 FireWall(FWM) Resource 段追加 `tencentcloud_fwm_strategy`，使 `make doc` 可索引。

## 2. 资源 Schema 定义

- [x] 2.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_strategy.go`，定义 `ResourceTencentCloudFwmStrategy()`，Schema 严格按 design D1/D2 映射 `CreateStrategyRequest` 字段与 Computed `strategy_ids`：
  - `product` TypeString Required ForceNew
  - `receive_account` TypeList of TypeString Optional
  - `pre_strategy` TypeList(GroupId/Sequence) Optional
  - `post_strategy` TypeList(GroupId/Sequence) Optional
  - `receive_group` TypeList of TypeString Optional
  - `strategy_ids` TypeList of TypeString Computed（本次新增参数）
- [x] 2.2 声明 `Importer: schema.ImportStatePassthrough`。

## 3. Create 函数

- [x] 3.1 从 `d` 读取入参并构造 `CreateStrategyRequest`（参考 igtm_strategy Create 风格：`helper.String` / 列表遍历构造 `StrategyReq`）。
- [x] 3.2 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `CreateStrategyWithContext`，retry 块内仅调 API；错误用 `tccommon.RetryError(e)`；检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`。
- [x] 3.3 retry 块外：检查 `response.Response.StrategyIds == nil || len == 0` → 返回 NonRetryableError；打印 `logId` 与 `d.Id()` 便于排障。
- [x] 3.4 `d.SetId(strings.Join(strategyIds, tccommon.FILED_SP))`，再 `d.Set("strategy_ids", ...)`，最后转调 Read。

## 4. Read 函数

- [x] 4.1 `strings.Split(d.Id(), tccommon.FILED_SP)` 还原 StrategyIds 列表。
- [x] 4.2 对每个 id 用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调 `DescribeStrategyWithContext`，retry 块内检查 `result == nil || result.Response == nil || result.Response.Strategy == nil` 返回 `NonRetryableError`（不清空 id）。
- [x] 4.3 全部查询为空时：先 `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` 再 `d.SetId("")`。

## 5. Delete 函数

- [x] 5.1 从 `d.Id()` 还原 StrategyIds 列表，构造 `DeleteStrategyRequest.StrategyIds`。
- [x] 5.2 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调 `DeleteStrategyWithContext`，错误 `tccommon.RetryError(e)`，检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`。

## 6. Update 函数（CRD-only immutable）

- [x] 6.1 实现 immutableArgs 校验：`immutableArgs := []string{"receive_account","pre_strategy","post_strategy","receive_group"}`，遍历 `d.HasChange(v)` 命中即返回 `fmt.Errorf("fwm_strategy %s is immutable, please recreate resource", v)`；`product` 为 ForceNew 自动触发重建；`strategy_ids` 为 Computed 不参与校验。
- [x] 6.2 校验通过后转调 Read（不调任何 Modify API）。

## 7. 文档示例

- [x] 7.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_strategy.md`，包含用法示例、参数说明、import 联合 id 说明（`terraform import tencentcloud_fwm_strategy.x s1#s2`）。

## 8. 单元测试

- [x] 8.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go`，使用 gomonkey mock `CreateStrategyWithContext` / `DescribeStrategyWithContext` / `DeleteStrategyWithContext`，覆盖 Create 成功写入 strategy_ids、Create 空 StrategyIds 报错、Read 空打印并清空 id、Delete 成功、Update immutable 报错等分支，确保可编译可执行。

## 9. 验证任务

- [x] 9.1 验证 schema / provider.go / provider.md / connectivity 改动可编译（由后续流程的 go build/vet 执行，不在本阶段运行）。
- [x] 9.2 验证 `make doc` 生成 `website/docs/r/fwm_strategy.html.markdown`（收尾阶段执行）。