## 1. 云 API 核查（前置）

- [x] 1.1 确认 vendor 中 `ModifyLoadBalancerSla` 接口存在，入参 `Force *bool` 存在（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go`）
- [x] 1.2 确认 `DescribeLoadBalancers` 返回的 `LoadBalancer` 结构不含 `Force` 字段，确定 `force` 为非 Computed
- [x] 1.3 确认 clb 包不存在 `ModifyTags` 接口、`tencentcloud_clb_instance` 已存在 `tags` 字段，排除 `tags` 变更

## 2. Schema 与 CRUD 代码修改

- [x] 2.1 在 `tencentcloud/services/clb/resource_tc_clb_instance.go` 的 `ResourceTencentCloudClbInstance` schema 中新增 `force` 字段：`schema.TypeBool`、`Optional: true`、`Default: false`，Description 说明"是否强制升级，默认否；仅在 sla_type 变更调用 ModifyLoadBalancerSla 时生效；共享型升级为性能容量型后不可回退"
- [x] 2.2 在 `resourceTencentCloudClbInstanceUpdate` 的 `if d.HasChange("sla_type")` 块内，构造 `slaRequest` 后、retry 调用 `ModifyLoadBalancerSla` 前，增加 `slaRequest.Force = helper.Bool(d.Get("force").(bool))`，保持现有 retry 与异步任务等待逻辑不变
- [x] 2.3 确认 `resourceTencentCloudClbInstanceRead` 不回填 `force`（非 Computed），不新增 `d.Set("force", ...)`
- [x] 2.4 确认 Create 阶段不调用 `ModifyLoadBalancerSla`，无需处理 `force`

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/clb/resource_tc_clb_instance_test.go` 中补充测试用例：验证 `sla_type` 变更且 `force=true` 时，`ModifyLoadBalancerSlaRequest.Force` 被正确透传为 true（使用 gomonkey mock 云 API）
- [x] 3.2 补充测试用例：验证未声明 `force` 时，`Force` 取默认值 false，行为与变更前一致

## 4. 文档

- [x] 4.1 更新 `tencentcloud/services/clb/resource_tc_clb_instance.md`，补充 `force` 参数的 Example Usage 与说明（强制升级、不可回退、仅在 sla_type 变更时生效）
- [ ] 4.2 通过收尾阶段 `make doc` 命令生成 `website/docs/` 下的资源文档（禁止手动编写 website/ 文件）

## 5. 验证（收尾阶段执行，非本步骤）

- [ ] 5.1 执行 `gofmt` 格式化变更的 Go 代码（由 tfpacer-finalize skill 执行）
- [ ] 5.2 执行 `make doc` 生成文档（由 tfpacer-finalize skill 执行）
- [ ] 5.3 生成 `.changelog` 文件（由 tfpacer-finalize skill 执行）
