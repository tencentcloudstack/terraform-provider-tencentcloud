## 1. Action 实现

- [x] 1.1 在 `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan.go` 中定义 `TeoBindZoneToPlanAction` 实现 `action.Action` 与 `action.ActionWithConfigure`，schema 包含 `zone_id`、`plan_id` 两个 `Required + TypeString` 字段（`actionschema.StringAttribute`），不定义 Computed/Sensitive
- [x] 1.2 实现 `Metadata` 返回 `tencentcloud_teo_bind_zone_to_plan`；实现 `Schema` 返回 action schema；实现 `Configure` 从 `*sharedmeta.ProviderMeta` 注入共享 `*connectivity.TencentCloudClient`
- [x] 1.3 实现 `Invoke`：读取 config，校验 `zone_id`、`plan_id` 非空，构造 `teov20220901.NewBindZoneToPlanRequest()` 填充 `ZoneId` 与 `PlanId`；使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `BindZoneToPlanWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装返回；成功后不写 state（action 无 state）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册 `teo.NewTeoBindZoneToPlanAction`，并添加 `tencentcloud/services/teo` 包导入
- [x] 2.2 在 `tencentcloud/provider.md` 中 TEO 产品下新增 `Action` 子节，登记 `tencentcloud_teo_bind_zone_to_plan`

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan.md`，包含一句话描述（带上 TEO）与 Example Usage（展示 `zone_id`、`plan_id`），不含 Import 部分（action 无）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan_test.go`，使用 gomonkey mock `BindZoneToPlanWithContext`，覆盖 Invoke 成功用例（断言 `ZoneId`、`PlanId` 正确传入）
- [x] 4.2 新增 Invoke API 错误用例（mock 返回 `ResourceNotFound` 等非可重试错误，断言返回 diagnostics error）
- [x] 4.3 新增 Invoke 缺 `zone_id` / 缺 `plan_id` 用例（断言返回 diagnostics error）
- [x] 4.4 新增 Metadata、Schema、Factory、ImplementsInterfaces、Configure(nil/wrong-type) 用例

## 5. 验证

- [x] 5.1 运行 `go test ./tencentcloud/services/teo/ -run TestTeoBindZoneToPlanAction -v -count=1 -gcflags="all=-l"` 确认所有单测通过
- [ ] 5.2 执行 `make doc`，根据 provider 规范重新生成 `website/docs/a/` 下的 markdown 文档（禁止手改）