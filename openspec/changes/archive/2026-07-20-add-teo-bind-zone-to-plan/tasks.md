## 1. 资源实现

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation.go` 中定义 `ResourceTencentCloudTeoBindZoneToPlan()` 资源，schema 包含 `zone_id`、`plan_id` 两个 `Required + ForceNew + TypeString` 字段，并定义 `Create`、`Read`、`Delete`（不定义 `Update`，不定义 `Importer`）
- [x] 1.2 实现 `resourceTencentCloudTeoBindZoneToPlanCreate`：构造 `teov20220901.NewBindZoneToPlanRequest()`，填充 `ZoneId` 与 `PlanId`；使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `BindZoneToPlanWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装返回；retry 块外、错误处理后调用 `d.SetId(helper.BuildToken())`，再调用 Read
- [x] 1.3 实现 `resourceTencentCloudTeoBindZoneToPlanRead` 与 `resourceTencentCloudTeoBindZoneToPlanDelete`：均为 no-op，直接返回 nil

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `tencentcloud_teo_bind_zone_to_plan` 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中登记 `tencentcloud_teo_bind_zone_to_plan` 资源名

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation.md`，包含一句话描述（带上 TEO）与 Example Usage（展示 `zone_id`、`plan_id`），不含 Import 部分（operation 资源无）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation_test.go`，使用 gomonkey mock `BindZoneToPlanWithContext`，覆盖 Create 成功用例（断言 `ZoneId`、`PlanId` 正确传入，资源 id 非空）
- [x] 4.2 新增 Create API 错误用例（mock 返回 `ResourceNotFound` 等非可重试错误，断言返回 error）
- [x] 4.3 新增 no-op Read、no-op Delete 用例
- [x] 4.4 新增 schema 校验用例（断言 Create/Read/Delete 非空、Update 为 nil、`zone_id`/`plan_id` 为 Required+ForceNew）

## 5. 验证

- [x] 5.1 运行 `go test ./tencentcloud/services/teo/ -run TestBindZoneToPlanOperation -v -count=1 -gcflags="all=-l"` 确认所有单测通过
- [ ] 5.2 执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改）
