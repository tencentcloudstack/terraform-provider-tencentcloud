## Context

本变更新增一个 terraform-plugin-framework **action** 资源 `tencentcloud_antiddos_unblock_resources`，用于调用 AntiDDoS（DDoS 防护）的 `UnblockResources` 接口申请解封被封堵的资源（公网 IP 列表）。该 action 属于一次性操作，操作完成后不记录任何云侧状态。

**参考实现**：严格参考 `tencentcloud/framework/registry.go` 中已注册的两个 action —— `teo.NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）与 `bdrc.NewBdrcRunCopyPairTasks`（`tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks.go`）的代码风格。

**云 API SDK**：`UnblockResources` 接口位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903` 包中。现有 `UseAntiddosClient()` 方法仅支持 `v20200309` 版本，因此需在 `tencentcloud/connectivity/client.go` 中新增 `UseAntiddosV20250903Client()` 方法。

**vendor 中云 API 关键字段定义**（从 `vendor/.../antiddos/v20250903/models.go` 按行区间摘录）：

`UnblockResourcesRequest`（models.go:153-158）:
```go
type UnblockResourcesRequest struct {
	*tchttp.BaseRequest
	// 申请解封的资源列表。支持按照公网 IP 解封。
	// 入参限制：列表长度最大限制 10。
	Resources []*string `json:"Resources,omitnil,omitempty" name:"Resources"`
}
```

`UnblockResourcesResponse`（models.go:185-188）:
```go
type UnblockResourcesResponse struct {
	*tchttp.BaseResponse
	Response *UnblockResourcesResponseParams `json:"Response"`
}
// UnblockResourcesResponseParams 仅含 RequestId *string，无业务返回字段。
```

Client 方法（client.go:117-142）:
```go
func (c *Client) UnblockResources(request *UnblockResourcesRequest) (response *UnblockResourcesResponse, err error)
func (c *Client) UnblockResourcesWithContext(ctx context.Context, request *UnblockResourcesRequest) (response *UnblockResourcesResponse, err error)
```

接口注释说明：`UnblockResources` —— 申请解封资源，可通过 `DescribeDDoSBlockRecords` 接口获取资源的封堵解封状态。该接口为**同步接口**（response 无 FlowId/TaskId/JobId 等任务标识），无需异步轮询。

## Goals / Non-Goals

**Goals:**
- 提供一个 framework action 资源，让用户通过 Terraform 触发 AntiDDOS 资源解封操作。
- action 资源的 schema 包含一个必填参数 `resources`（公网 IP 列表，`types.List` of `types.String`）。
- 严格遵循 `NewTeoConfirmOriginAclUpdate` / `NewBdrcRunCopyPairTasks` 的代码风格：实现 `action.ActionWithConfigure` 接口，通过 `fw.ActionWithConfigure` 嵌入获取 client，在 `Invoke` 方法中校验入参并调用服务层。
- 新增 `UseAntiddosV20250903Client()` 连接方法，支持 `antiddos/v20250903` SDK。
- 服务层 `UnblockResources` 方法使用 `resource.Retry(WriteRetryTimeout, ...)` + `ratelimit.Check` + `tccommon.RetryError` 的重试机制。
- 提供单元测试（gomonkey mock 云 API）与示例文档。

**Non-Goals:**
- 不实现 Read/Update/Delete（action 资源无这些生命周期方法，仅有 `Invoke`）。
- 不做异步轮询（接口为同步接口，response 无任务标识）。
- 不修改现有 `tencentcloud/provider.go`（framework action 通过 `registry.go` 的 `actionFactories` 注册，无需改动 provider.go）。

## Decisions

### Decision 1: 资源类型为 framework action，而非传统 SDKv2 operation 资源
用户明确要求使用 terraform framework plugin 模式开发，属于 action 资源。framework action 与传统 SDKv2 operation 资源的本质区别：
- SDKv2 operation 资源实现 `schema.Resource` 的 Create/Read/Delete，通过 `provider.go` 的 `ResourcesMap` 注册。
- framework action 实现 `action.Action` 接口的 `Metadata`/`Schema`/`Invoke`，通过 `registry.go` 的 `actionFactories` 注册。

参考 `NewBdrcRunCopyPairTasks`：它有 `copy_pair_ids`（List[String]）这样的列表入参，本资源的 `resources` 同为 List[String]，结构完全对应。

### Decision 2: 文件落地位置
- action 实现文件：`tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.go`
- action 示例文档：`tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.md`
- action 单元测试：`tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources_test.go`

依据当前仓库实际状态：`restructure-framework-types-and-naming` 变更尚未落地，`registry.go` 中已注册的全部 4 个 framework action（`teo`/`bdrc`/`dlc`/`mongodb`）均实际位于 `tencentcloud/services/<product>/`，服务方法与 action 同包。为保证可编译且与真实代码库约定一致，本 action 落在 `tencentcloud/services/antiddos/`（复用既有 `AntiddosService`，包名 `antiddos`），registry 中 import `tencentcloud/services/antiddos`。`restructure` 变更后续落地时再统一迁移到 `framework/antiddos/`。

### Decision 3: 服务层方法位置
`UnblockResources` 服务方法放置在现有 `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go` 中（已有 `AntiddosService` 结构体，可复用），新增方法：
```go
func (me *AntiddosService) UnblockResources(ctx context.Context, resources []*string) (errRet error)
```
该方法内部使用 `me.client.UseAntiddosV20250903Client().UnblockResourcesWithContext(ctx, request)`，并包裹 `resource.Retry(tccommon.WriteRetryTimeout, ...)`。

framework action 的 `Invoke` 中通过 `antiddos.NewAntiddosService(a.Client())` 获取服务层实例，调用该方法。

### Decision 4: 新增 UseAntiddosV20250903Client() 连接方法
在 `tencentcloud/connectivity/client.go` 中：
- import 区新增 `antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"`
- struct 字段区新增 `antiddosV20250903Conn *antiddosv20250903.Client`
- 新增方法（参考 `UseBdrcV20260330Client` 风格）：
```go
func (me *TencentCloudClient) UseAntiddosV20250903Client() *antiddosv20250903.Client {
	if me.antiddosV20250903Conn != nil {
		return me.antiddosV20250903Conn
	}
	cpf := me.NewClientProfile(300)
	me.antiddosV20250903Conn, _ = antiddosv20250903.NewClient(me.Credential, me.Region, cpf)
	me.antiddosV20250903Conn.WithHttpTransport(&LogRoundTripper{})
	return me.antiddosV20250903Conn
}
```

### Decision 5: registry 注册方式
在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片末尾新增一行 `antiddos.NewAntiddosUnblockResources,`，并在 import 区新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"`。

### Decision 6: Schema 定义
action schema 仅含一个必填属性 `resources`（`schema.ListAttribute`，`ElementType: types.StringType`），对应云 API 的 `request.Resources []*string`。参考 `NewBdrcRunCopyPairTasks` 中 `copy_pair_ids` 的定义方式。

### Decision 7: Invoke 方法校验逻辑
参考 `NewBdrcRunCopyPairTasks` 的 Invoke 方法：
1. `req.Config.Get(ctx, &data)` 解析入参。
2. 校验 `resources` 非空（null/unknown/空列表均报错 "Missing resources"）。
3. 校验 client 非空（报错 "Provider not configured"）。
4. 将 `types.List` 转换为 `[]*string`，调用服务层 `UnblockResources(ctx, resources)`。
5. 错误通过 `resp.Diagnostics.AddError` 返回。

### Decision 8: 测试风格
严格参考 `tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks_test.go` 的测试风格，使用 gomonkey mock 云 API。测试用例包括：
- `TestAntiddosUnblockResources_Invoke_Success`：成功调用，校验 request 参数正确。
- `TestAntiddosUnblockResources_Invoke_APIError`：API 报错，校验 diagnostics 含错误。
- `TestAntiddosUnblockResources_Invoke_MissingInput`：缺少必填参数，校验在 API 调用前报错。
- `TestAntiddosUnblockResources_Invoke_ClientNotConfigured`：client 为 nil，校验报错。
- `TestAntiddosUnblockResources_MetadataAndSchema`：校验 TypeName 和 schema 定义。

### Decision 9: 示例文档格式
参考 `tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks.md`，包含：
- 一句话描述（带上 AntiDDoS 产品名）
- `~> **NOTE:**` action 版本提示
- Example Usage（hcl `action` 块）

## Risks / Trade-offs

- [action 资源需要 Terraform >= 1.14] → 在示例文档中添加 NOTE 提示用户 Terraform 版本要求。
- [UnblockResources 接口有调用配额限制（解封次数配额）] → 在示例文档中说明，超出配额时 API 会返回 LimitExceeded 错误，由用户自行处理。
- [v20250903 SDK 未提交到 vendor tracked 状态] → 实施时需确保 `go mod vendor` 后 vendor 目录包含该包；当前 vendor 中已存在该包文件（untracked），需在 go.mod 中确认依赖。