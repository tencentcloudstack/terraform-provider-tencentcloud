## 1. 连接层：新增 UseAntiddosV20250903Client

- [x] 1.1 在 `tencentcloud/connectivity/client.go` import 区新增 `antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"`
- [x] 1.2 在 `TencentCloudClient` struct 字段区新增 `antiddosV20250903Conn *antiddosv20250903.Client`
- [x] 1.3 新增 `UseAntiddosV20250903Client()` 方法，参考 `UseBdrcV20260330Client()` 的实现风格（懒加载 + NewClientProfile(300) + LogRoundTripper）

## 2. 服务层：新增 UnblockResources 服务方法

- [x] 2.1 在 `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go` 新增 `UnblockResources(ctx context.Context, resources []*string) (errRet error)` 方法
- [x] 2.2 方法内部构造 `antiddosv20250903.NewUnblockResourcesRequest()`，设置 `request.Resources = resources`
- [x] 2.3 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 `me.client.UseAntiddosV20250903Client().UnblockResourcesWithContext(ctx, request)`，API 错误通过 `tccommon.RetryError(e)` 包装，调用前执行 `ratelimit.Check(request.GetAction())`
- [x] 2.4 在 defer 中打印 `[CRITAL]` 错误日志（参考 `BdrcService.RunCopyPairTasks`）

## 3. Action 实现

- [x] 3.1 创建目录 `tencentcloud/framework/antiddos/`
- [x] 3.2 创建 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.go`，package antiddos
- [x] 3.3 定义 `AntiddosUnblockResources` struct，嵌入 `fw.ActionWithConfigure`；定义 `AntiddosUnblockResourcesModel` struct，含 `Resources types.List` 字段（tfsdk: "resources"）
- [x] 3.4 实现 `NewAntiddosUnblockResources() action.Action` 工厂函数
- [x] 3.5 实现 `Metadata` 方法，设置 `resp.TypeName = "tencentcloud_antiddos_unblock_resources"`
- [x] 3.6 实现 `Schema` 方法，定义必填 `resources` 属性（`schema.ListAttribute`，`ElementType: types.StringType`），描述带上 AntiDDoS 产品名
- [x] 3.7 实现 `Invoke` 方法：解析入参 → 校验 resources 非空 → 校验 client 非空 → 转换 types.List 为 []*string → 调用 `NewAntiddosService(a.Client()).UnblockResources(ctx, resources)`，严格参考 `NewBdrcRunCopyPairTasks.Invoke`

> 注：当前仓库内 `restructure-framework-types-and-naming` 变更尚未落地，`registry.go` 中全部 4 个 framework action 均实际位于 `tencentcloud/services/<product>/`，服务方法同包。为保证可编译且与真实代码库约定一致，本 action 落在 `tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.go`，`registry.go` 中 import `tencentcloud/services/antiddos`。

## 4. Registry 注册

- [x] 4.1 在 `tencentcloud/framework/registry.go` import 区新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/antiddos"`
- [x] 4.2 在 `actionFactories` 切片末尾新增 `antiddos.NewAntiddosUnblockResources,`

## 5. 示例文档

- [x] 5.1 创建 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.md`，格式参考 `action_tc_bdrc_run_copy_pair_tasks.md`：一句话描述（带 AntiDDoS）+ `~> NOTE` action 版本提示 + Example Usage（hcl `action` 块）

## 6. 单元测试

- [x] 6.1 创建 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources_test.go`，package antiddos_test，严格参考 `action_tc_bdrc_run_copy_pair_tasks_test.go` 风格
- [x] 6.2 实现 helper 函数：`antiddosUnblockResourcesSchema`、`newAntiddosUnblockResourcesInvokeRequest`、`newAntiddosUnblockResourcesInvokeRequestNullResources`、`setAntiddosUnblockResourcesClient`、`ptrStringAntiddosUnblockResources`
- [x] 6.3 实现 `TestAntiddosUnblockResources_Invoke_Success`：mock `UseAntiddosV20250903Client` 与 `UnblockResourcesWithContext`，校验 request.Resources 正确
- [x] 6.4 实现 `TestAntiddosUnblockResources_Invoke_APIError`：mock API 返回错误，校验 diagnostics 含错误摘要
- [x] 6.5 实现 `TestAntiddosUnblockResources_Invoke_MissingInput`：校验空 resources 在 API 调用前报错 "Missing resources"
- [x] 6.6 实现 `TestAntiddosUnblockResources_Invoke_ClientNotConfigured`：校验 nil client 报错 "Provider not configured"
- [x] 6.7 实现 `TestAntiddosUnblockResources_MetadataAndSchema`：校验 TypeName 与 schema 属性定义

## 7. 文档生成

- [ ] 7.1 执行 `make doc` 命令生成 `website/docs/r/` 下的 action 文档（禁止手动编写 website/ 文件）

> 注：`make doc` 由收尾阶段（tfpacer-finalize skill）统一执行，不在本 apply 阶段运行。

## 8. 验证

- [ ] 8.1 执行 `gofmt` 格式化新增/修改的 .go 文件
- [ ] 8.2 确认 `go build ./tencentcloud/...` 无编译错误（由后续流程执行，不在本阶段手动运行）

> 注：`gofmt` 与 `go build` 由收尾阶段执行；本阶段已通过人工核对 `git diff` 与 vendor 云 API 字段一致（`UnblockResourcesRequest.Resources []*string` ↔ schema `resources` List[String] ↔ 服务层 `request.Resources = resources`），确保可编译。