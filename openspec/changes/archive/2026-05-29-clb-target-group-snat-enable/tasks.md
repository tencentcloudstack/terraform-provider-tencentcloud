## 1. Schema 与 Service 层改动

- [x] 1.1 在 `tencentcloud/services/clb/resource_tc_clb_target_group.go` 的 `ResourceTencentCloudClbTargetGroup().Schema` 中新增 `snat_enable` 字段：`Type: schema.TypeBool, Optional: true, Computed: true`，描述说明 SNAT 行为及生效条件
- [x] 1.2 在 `tencentcloud/services/clb/service_tencentcloud_clb.go` 的 `ClbService.CreateTargetGroup` 函数签名末尾新增 `snatEnable *bool` 参数，并在函数体中加入 `if snatEnable != nil { request.SnatEnable = snatEnable }`

## 2. 资源 CRUD 改动

- [x] 2.1 在 `resourceTencentCloudClbTargetCreate` 中通过 `d.GetOkExists("snat_enable")` 取值，封装为 `*bool` 后作为新参数传入 `clbService.CreateTargetGroup(...)`
- [x] 2.2 在 `resourceTencentCloudClbTargetRead` 中加入 `if targetGroup.SnatEnable != nil { _ = d.Set("snat_enable", targetGroup.SnatEnable) }`，与现有字段处理风格一致
- [x] 2.3 在 `resourceTencentCloudClbTargetUpdate` 中加入 `if d.HasChange("snat_enable") { request.SnatEnable = helper.Bool(d.Get("snat_enable").(bool)); isChanged = true }` 分支
- [x] 2.4 检查 `immutableFields` 切片，确保不会把 `snat_enable` 错误地加入不可变字段（应保持可变）

## 3. 文档与示例

- [x] 3.1 更新 `tencentcloud/services/clb/resource_tc_clb_target_group.md`，在示例 HCL 中加入 `snat_enable = true` 演示，并补充字段说明段落
- [x] 3.2 运行 `make doc` 自动生成 / 更新 `website/docs/r/clb_target_group.html.markdown`（不可手写）

## 4. 测试

- [x] 4.1 在 `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` 中新增 `TestAccTencentCloudClbTargetGroup_snatEnable`，使用两步 `TestStep`：Step 1 创建时 `snat_enable = true`，Step 2 Update 为 `snat_enable = false`
- [x] 4.2 在测试中加入断言：`resource.TestCheckResourceAttr(...,"snat_enable","true"/"false")`，并验证资源 ID 在 Update 前后未变（in-place update）
- [x] 4.3 在 `resource_tc_clb_target_group_testing_test.go` 中（如需）补充非真实账号的 schema 校验用例

## 5. 变更日志

- [x] 5.1 在 `.changelog/` 下新增 `<PR_NUMBER>.txt`，使用 `enhancement` 类型，文案：`resource/tencentcloud_clb_target_group: support new argument 'snat_enable'`

## 6. 验证

- [x] 6.1 运行 `make fmt` 完成代码格式化
- [x] 6.2 运行 `make lint`（含 golangci-lint 与 tfproviderlint）通过
- [x] 6.3 运行 `go build ./...` 通过
- [x] 6.4 运行 `go vet ./...` 通过
- [ ] 6.5 在具备测试账号的环境下，使用 `TF_ACC=1 go test -run TestAccTencentCloudClbTargetGroup_snatEnable ./tencentcloud/services/clb/...` 跑通验收测试
- [x] 6.6 运行 `openspec validate clb-target-group-snat-enable --strict` 校验本变更结构
