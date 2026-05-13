## 1. 资源实现

- [x] 1.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation.go` 文件
- [x] 1.2 定义 Schema：`product`、`instance_ids`、`task_id`、`operate`（必填，全部 ForceNew）
- [x] 1.3 定义 Schema：`order_auth_time`、`task_sub_method`（可选，全部 ForceNew）
- [x] 1.4 实现 `ResourceTencentCloudCvmRepairTaskControlOperation()` 返回 `*schema.Resource`
- [x] 1.5 实现 Create 函数：构造 `cvm.NewRepairTaskControlRequest`、字段映射、调用 API、SetId
- [x] 1.6 在 Create 中使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 与 `tccommon.RetryError`
- [x] 1.7 实现 Read 函数：仅 defer 日志，return nil
- [x] 1.8 实现 Delete 函数：仅 defer 日志，return nil
- [x] 1.9 添加 `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(d, meta)` 到 CRUD 函数
- [x] 1.10 错误日志使用 `log.Printf("[CRITAL]...")` 格式

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中按字母序在 CVM 区域注册 `tencentcloud_cvm_repair_task_control_operation`
- [x] 2.2 在 `tencentcloud/provider.md` 的资源列表中按字母序添加 `tencentcloud_cvm_repair_task_control_operation`

## 3. 文档模板

- [x] 3.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation.md`
- [x] 3.2 添加资源功能描述（中英文均可，遵循项目惯例）
- [x] 3.3 添加示例：立即授权（仅必填参数）
- [x] 3.4 添加示例：预约授权（含 `order_auth_time`）
- [x] 3.5 添加示例：弃盘迁移（含 `task_sub_method = "LossyLocal"`），并明确警告会清空本地盘数据
- [x] 3.6 添加 Import 说明（注明该资源为 operation 资源，通常不需要 import）

## 4. 测试文件

- [x] 4.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation_test.go`
- [x] 4.2 实现 `TestAccTencentCloudCvmRepairTaskControlOperationResource_basic` 验收测试
- [x] 4.3 测试中校验资源 ID 非空、属性正确写入 state
- [x] 4.4 添加 `TestAccTencentCloudCvmRepairTaskControlOperationResource_withOrderTime` 验证预约授权（可选）

## 5. 代码质量检查

- [x] 5.1 运行 `go build ./...` 确保编译通过
- [x] 5.2 运行 `gofmt -w` 格式化代码
- [x] 5.3 运行 `go vet ./tencentcloud/...` 检查代码问题
- [x] 5.4 检查 import 顺序与项目其他 operation 资源一致

## 6. 文档生成

- [x] 6.1 运行 `make doc` 生成 `website/docs/r/cvm_repair_task_control_operation.html.markdown`
- [x] 6.2 检查生成文档：参数描述完整、示例可运行、无格式错误

## 7. 验收测试（需云凭证）

- [ ] 7.1 设置 `TF_ACC=1`、`TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY` 环境变量
- [ ] 7.2 准备一个处于"待授权"状态的真实维修任务用于测试
- [ ] 7.3 运行 `go test -v -run TestAccTencentCloudCvmRepairTaskControlOperationResource ./tencentcloud/services/cvm/`
- [ ] 7.4 确认测试通过，无遗留资源

## 8. 提交准备

- [x] 8.1 在 `.changelog/` 目录新建 changelog 文件（参考最新 PR 的格式与编号）
- [x] 8.2 changelog 内容：`new-resource: tencentcloud_cvm_repair_task_control_operation`
- [x] 8.3 确认 git diff 仅包含本次新增/修改文件
- [x] 8.4 准备 PR 描述：说明新增的资源功能、API 映射、测试方式、注意事项（不可逆 / 弃盘迁移警告）
