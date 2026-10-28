## Context

TencentCloud Provider 当前由单一的 `terraform-plugin-sdk/v2` Provider 构成，入口在 [main.go](/Users/arunma/project/terraform-provider-tencentcloud/main.go)：

```go
primary := tencentcloud.Provider() // SDKv2
providers := []func() tfprotov5.ProviderServer{
    primary.GRPCProvider,
    providerserver.NewProtocol5(tencentcloud.New(primary)), // framework 空壳
}
muxServer, _ := tf5muxserver.NewMuxServer(ctx, providers...)
tf5server.Serve("registry.terraform.io/tencentcloudstack/tencentcloud", muxServer.ProviderServer, ...)
```

`FrameworkProvider` 已存在于 [framework_provider.go](/Users/arunma/project/terraform-provider-tencentcloud/tencentcloud/framework_provider.go)，但 `Schema` 为空、`Configure` 为空、`Resources/DataSources` 返回 `nil`，等同于一个会被 mux 调用却什么都不暴露的占位 provider。

约束：
- 80+ 服务、数千个资源/数据源，零回归是硬性要求；不可改 schema、ID、state。
- 兼容 Terraform 0.13+ → 必须留在 Protocol v5（用 `tf5muxserver`，而非 `tf6muxserver`）。
- vendor 模式管理依赖，所有新依赖必须 `go mod vendor` 同步。
- 凭证逻辑（环境变量、shared credentials、CAM role、assume_role、profile）只能有一份事实来源。
- 社区主流多云 Provider 在双栈迁移时的通用思路：framework 资源 ProviderData 复用 SDKv2 已构造好的 client 对象，不在 framework 侧重复解析凭证。

## Goals / Non-Goals

**Goals:**
- 让 `FrameworkProvider` 成为生产可用的 framework provider：完整 schema、Configure 流程、ProviderData 注入。
- SDKv2 与 framework 共享同一个 `*connectivity.TencentCloudClient`（同一份凭证、同一份 SDK client 缓存、同一份 retry/UA），用户在配置中只需写一次 `provider "tencentcloud" {}` 块。
- 建立明确的开发约定：新资源默认走 framework；存量资源停留在 SDKv2，不主动迁移。
- 落地一个示范 framework 资源（无存量 state 的低风险资源），验证从注册、Configure、CRUD、acceptance test、`make doc` 的全链路。
- CI 中加入 mux 启动校验，避免同名资源在两栈重复注册导致运行时崩溃。

**Non-Goals:**
- 不迁移任何已有 SDKv2 资源到 framework（迁移由独立 change 推进，每次一个资源）。
- 不升级到 Protocol v6 / `tf6muxserver`（保留向后兼容）。
- 不重写 `tencentcloud/connectivity` 包；framework 侧通过适配器复用。
- 不为 framework 资源单独实现凭证解析；仍由 SDKv2 provider 解析后传递。
- 不引入 ephemeral resources / provider-defined functions 的具体实现（仅保留扩展点）。

## Decisions

### D1：mux 协议保持 Protocol v5（`tf5muxserver`）

**选择**：继续使用 `tf5muxserver`，framework 用 `providerserver.NewProtocol5`。

**理由**：
- `terraform-provider-tencentcloud` 公开承诺兼容 Terraform 0.13+，Protocol v6 要求 Terraform 1.0+。
- v5 已能承载 framework 绝大多数能力；nested attribute 在 v5 下也有降级表达。
- 切到 v6 是后续独立变更，不应与本次 mux 落地耦合。

**备选**：`tf6muxserver` + `tf5to6server` 包裹 SDKv2。被否，因 v6 切换是单独的兼容性变更，不应混在本次。

### D2：凭证与 client 由 SDKv2 provider 解析，framework 复用

**选择**：framework provider 的 `Configure` 不自行实现凭证解析。在 SDKv2 `Provider().ConfigureContextFunc` 完成 client 构建后，将 `*connectivity.TencentCloudClient` 通过进程级共享指针传递；framework provider 通过共享的 `ProviderMeta` 类型获得同一指针。

具体做法（项目自行设计，命名与包路径贴合腾讯云 provider 现有风格）：
- 新增 `tencentcloud/internal/tcfwprovider/meta.go`：定义 `type ProviderMeta struct { Client *connectivity.TencentCloudClient }`。
- 新增 `tencentcloud/internal/tcfwprovider/shared_meta.go`：提供进程级 `var sharedMeta atomic.Pointer[connectivity.TencentCloudClient]` 以及 `SetSharedMeta`/`GetSharedMeta`/`ResetSharedMetaForTest` 函数，由 SDKv2 的 `providerConfigure` 在成功返回前写入。
- `FrameworkProvider.Schema` 与 SDKv2 schema 字段镜像（同名、同类型语义），但 `Configure` 中**不重新解析**，而是从 `tcfwprovider.GetSharedMeta()` 读出 SDKv2 已经构造好的 client，写入 `resp.ResourceData = &tcfwprovider.ProviderMeta{Client: c}`。
- 同样 `resp.DataSourceData` / `resp.EphemeralResourceData` 写同一指针。

**理由**：
- 避免双份解析导致两栈凭证不一致（例如一栈读到了 `assume_role`，另一栈没读到）。
- 凭证、SDK client cache、UA 头、retry 策略只有一份维护成本。
- 用户的 `provider "tencentcloud" {}` 配置块由 mux 同时投递给两栈，两栈先后 Configure；SDKv2 先 Configure 完成，framework 拿到同一份结果。

**备选**：
- 让 framework 也独立实现 schema 解析 → 否，维护成本高，凭证逻辑漂移风险大。
- 把凭证解析下沉到一个新公共包，两栈都从公共包读 → 工作量过大，超出本变更范围；列入未来 follow-up。

**风险**：mux 调用顺序未定义。需要在 framework `Configure` 中处理 `sharedClient == nil` 情况：返回明确诊断"SDKv2 provider not configured yet"。实际上 `tf5muxserver` 会对每个底层 server 都派发 `ConfigureProvider`，且 SDKv2 注册在前，但仍需防御式判空。

### D3：framework provider schema 与 SDKv2 字段镜像（仅声明，不解析）

**选择**：`FrameworkProvider.Schema` 声明所有 SDKv2 provider 块字段（`secret_id`、`secret_key`、`security_token`、`region`、`protocol`、`domain`、`assume_role` 等），全部 `Optional: true`。Configure 中读取后**仅做日志/告警**，真正生效以 SDKv2 解析结果为准。

**理由**：mux 校验 schema 时要求两栈对 provider 块字段集合一致，否则用户在 HCL 中写的字段会被其中一栈视为 unknown。

**备选**：framework 侧 schema 完全为空 → 否，Terraform 在解析 provider block 时会按 mux 合并的 schema 校验，缺字段会报错。

### D4：资源/数据源注册采用聚合函数，按服务分包

**选择**：
- `tencentcloud/services/<service>/framework.go`（每服务一个）导出 `func FrameworkResources() []func() resource.Resource` 和 `func FrameworkDataSources() []func() datasource.DataSource`。
- `tencentcloud/framework_provider.go` 的 `Resources()` / `DataSources()` 汇总调用所有服务包的聚合函数。
- framework 资源文件命名：`resource_tc_<service>_<name>_framework.go`，避免与 SDKv2 同名文件冲突。
- 资源类型名（TypeName）保持 `tencentcloud_<service>_<name>`，确保用户视角无差异；同名资源不能同时在两栈注册（CI 校验）。

**理由**：
- 与现有按服务分包的目录结构一致，开发者无需学习新布局。
- 显式聚合便于 grep 与代码审查；不引入运行时反射注册。
- `_framework.go` 后缀让 ripgrep / IDE 在同一服务目录下能快速区分两栈实现。

### D5：公共工具包 `tencentcloud/internal/fwhelper`

**选择**：新增 `fwhelper` 包，提供：
- `StringValueOrNull(s string) types.String` / `Int64ValueOrNull` 等 SDK ↔ framework 类型转换。
- `RetryFramework(ctx, timeout, fn)`：包装 `helper.Retry`，把 `*resource.RetryError` 桥接到 framework 的 `diag.Diagnostics`。
- `WrapSDKError(err) diag.Diagnostic`：把 `*sdkErrors.TencentCloudSDKError` 翻译成结构化诊断。
- `TimeoutsAttribute()`：复用 `terraform-plugin-framework-timeouts`，统一 Create/Read/Update/Delete 超时块声明。

**理由**：framework 没有 SDKv2 中 `helper.Retry`、`tccommon.LogElapsed`、`tccommon.InconsistentCheck` 的等价物；如果不沉淀公共层，每个新资源都会重复造轮子，且不一致。

### D6：CI 校验

**选择**：在 `make lint` 之外新增 `make check-mux`：
1. 构建二进制，`./terraform-provider-tencentcloud -dump-schema`（或在测试中调用 `tf5muxserver.NewMuxServer` 直接拿到 schema）→ 校验无 panic、无重复 type name。
2. 单元测试 `tencentcloud/framework_provider_test.go` 调用 `New(...)` 后遍历 `Resources()` / `DataSources()`，与 SDKv2 `Provider().ResourcesMap` 求交集，断言交集为空。

**理由**：mux 同名资源是启动期 panic，必须在 PR 阶段拦截。

## Risks / Trade-offs

- **风险：mux 调用顺序未定义** → 通过 `sharedClient` 防御式判空 + 明确 diagnostics；并在单元测试中用 `tf5muxserver` 模拟两栈 Configure 顺序，确保 framework 后 Configure 时能拿到 client。
- **风险：framework 与 SDKv2 schema 类型语义差异**（如 `Computed + Optional` 在 framework 下需要显式 `UseStateForUnknown` plan modifier） → 在 `fwhelper` 中提供 `OptionalComputedString()` 等帮手函数，并写入 contributor 文档。
- **风险：vendor 不同步导致 CI 红** → 在本变更中明确执行 `go mod vendor` 并在 tasks.md 中作为独立步骤；提交时一并校验 `vendor/modules.txt`。
- **风险：开发者把同一个 type name 同时注册到两栈** → CI 的 `make check-mux` 拦截；PR 模板里增加 checkbox。
- **Trade-off：维护两套基础设施** → 长期来看不可避免；本变更通过 `fwhelper` + 示范资源把成本降到最低，新资源全部走 framework，存量随业务自然演进。
- **Trade-off：暂不升 v6** → 暂时无法使用 framework 的部分纯 v6 能力（如 `MoveResourceState` 全功能），可在后续 change 中升级。

## Migration Plan

1. **阶段 0（本变更）**：补齐 framework provider runtime + 落地 1 个示范资源；不动任何存量资源。
2. **阶段 1（后续 change）**：开放 framework 给业务团队，所有新资源走 framework；编写 contributor 指引。
3. **阶段 2（后续 change）**：评估 Protocol v6 升级时机（Terraform 0.13 用户占比 < 阈值时切 `tf6muxserver` + `tf5to6server`）。
4. **阶段 3（按需）**：针对个别复杂资源（嵌套对象、动态类型痛点）单独立项迁移，每次一个资源 + 完整测试。

**回滚**：
- 本变更不改任何已有资源的 schema/state，回滚等价于把 `FrameworkProvider.Resources/DataSources` 改回返回空切片，并删除示范资源；用户层无感知。
- 回滚不需要 state 迁移，因为示范资源是新增、无存量。

## Open Questions

- 示范资源的具体选择：建议挑选一个尚未实现、API 简单且无嵌套对象的新资源（候选由后续 tasks 阶段确认，例如 `tencentcloud_tag_attachment_v2` 之类的轻量资源）。是否在本变更中确定，还是留到 apply 阶段？→ **决策**：留到 apply 阶段挑选，proposal 不绑定具体资源。
- 是否需要在本变更中就建立 `terraform-plugin-framework-validators` 与 `terraform-plugin-framework-timeouts` 的版本锁定？→ **决策**：是，在 tasks 中作为依赖步骤。
- `make doc` 当前基于 SDKv2 schema 反射；framework 资源是否走 `tfplugindocs`？→ 倾向使用 `tfplugindocs`（HashiCorp 官方），但需确认现有自定义文档生成器是否兼容。在 tasks 中加一个调研步骤。
