## 1. Connectivity 层接入 BDRC 客户端

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 的 import 块中新增 `bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"`
- [x] 1.2 在 `TencentCloudClient` 结构体中新增字段 `bdrcv20260330Conn *bdrcv20260330.Client`
- [x] 1.3 新增方法 `UseBdrcV20260330Client()`，参考 `UseTeoV20220901Client()` 实现：懒初始化、`NewClient(me.Credential, me.Region, cpf)`（超时 300）、挂载 `LogRoundTripper`、缓存返回

## 2. Service 层封装 RunCopyPairTasks

- [x] 2.1 新建 `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`，定义 `BdrcService` 结构体（持有 `*connectivity.TencentCloudClient`）与 `NewBdrcService(client)` 构造函数
- [x] 2.2 实现 `BdrcService.RunCopyPairTasks(ctx, copyPairIds []*string, copyPairType string) error`：构造 `bdrcv20260330.NewRunCopyPairTasksRequest()`，填充 `CopyPairIds` 与 `CopyPairType`，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `RunCopyPairTasksWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；defer 中记录失败日志，成功记录 DEBUG 日志
- [x] 2.3 确保 import 包含 `tccommon`、`connectivity`、`ratelimit`、`resource`（SDKv2 helper）、`bdrcv20260330`（注：未使用 `helper` 包，因 `CopyPairType` 直接通过 `&copyPairType` 取址，避免引入未使用 import）

## 3. Framework Action 资源实现

- [x] 3.1 新建 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.go`，参考 `action_tc_teo_confirm_origin_acl_update.go` 结构：`var _ action.ActionWithConfigure = &BdrcCopyPairTasks{}`、工厂 `NewBdrcCopyPairTasks() action.Action`、结构体 `BdrcCopyPairTasks` 内嵌 `fw.ActionWithConfigure`、模型 `BdrcCopyPairTasksModel`（`CopyPairIds types.List`、`CopyPairType types.String`）
- [x] 3.2 实现 `Metadata`：`resp.TypeName = "tencentcloud_bdrc_copy_pair_tasks"`
- [x] 3.3 实现 `Schema`：`copy_pair_ids`（Required, `schema.ListAttribute` of String）、`copy_pair_type`（Required, `schema.StringAttribute`），含 Description
- [x] 3.4 实现 `Invoke`：解析 `req.Config` 到模型；校验 `copy_pair_ids` 为空/null 与 `copy_pair_type` 为空/null 时通过 `resp.Diagnostics.AddError` 返回；校验 `a.Client() == nil` 时返回；将 `types.List` 转为 `[]*string`，调用 `service.RunCopyPairTasks(ctx, ids, copyPairType)`，失败时 `resp.Diagnostics.AddError` 上报
- [x] 3.5 不实现 Read/Update/Delete、不设置 id、不回写任何输出属性（一次性操作）

## 4. Action 注册

- [x] 4.1 在 `tencentcloud/framework/registry.go` 的 import 块中新增 `bdrc` 服务包导入（若不存在）
- [x] 4.2 在 `actionFactories` 切片末尾追加 `bdrc.NewBdrcCopyPairTasks`

## 5. 单元测试

- [x] 5.1 新建 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks_test.go`，使用 gomonkey mock 云 API（不使用 Terraform 测试套件），定义 ptrString 辅助与 schema/config 构造辅助
- [x] 5.2 新增 `TestBdrcCopyPairTasks_Invoke_Success`：mock `UseBdrcV20260330Client` 与 `RunCopyPairTasksWithContext` 返回成功响应，断言 request 的 `CopyPairIds`/`CopyPairType` 符合预期、无错误 diagnostic
- [x] 5.3 新增 `TestBdrcCopyPairTasks_Invoke_APIError`：mock `RunCopyPairTasksWithContext` 返回错误，断言产生错误 diagnostic 且包含错误信息
- [x] 5.4 新增 `TestBdrcCopyPairTasks_Invoke_MissingInput`：校验 `copy_pair_ids` 或 `copy_pair_type` 为空时直接返回错误 diagnostic、未调用 API

## 6. 文档

- [x] 6.1 新建 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.md`：一句话描述带上 BDRC 云产品名称；Example Usage 使用 `action "tencentcloud_bdrc_copy_pair_tasks" "example" { config { copy_pair_ids = [...] copy_pair_type = "DISK" } }` 语法；含 NOTE 提示 Terraform 1.14+ 支持；无 Import 段、无手动 Argument/Attribute Reference

## 7. 收尾

- [ ] 7.1 执行收尾阶段 `tfpacer-finalize` skill：gofmt 格式化变更的 Go 代码、`make doc` 生成 website/docs 文档、生成 changelog 文件并统一 amend 推送
