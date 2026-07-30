## Why

当前 TencentCloud Provider 完全基于 `terraform-plugin-sdk/v2` 实现。SDKv2 已进入维护模式，HashiCorp 推荐新功能使用 `terraform-plugin-framework`，因为 framework 提供：

- 原生 Protocol v6 支持，更好的嵌套对象/集合/动态类型表达能力（例如 nested attribute、`types.Dynamic`）。
- 更严格的类型系统与 plan modifier、validator 机制，能减少历史上 SDKv2 中靠 `DiffSuppressFunc`/`CustomizeDiff` 维护的脆弱代码。
- 对 ephemeral resources、provider-defined functions、list resources 等新能力的官方支持，这些 SDKv2 不会再补齐。

但 80+ 服务、上千个资源/数据源全量重写不现实。需要一种渐进式方案：**新能力使用 framework 实现，存量资源继续运行在 SDKv2，两者通过 mux server 在同一 Provider 二进制中共存**，这是社区主流 Provider 在面对 SDKv2 → framework 迁移时的通用做法。

项目已在 [main.go](/Users/arunma/project/terraform-provider-tencentcloud/main.go) 和 [framework_provider.go](/Users/arunma/project/terraform-provider-tencentcloud/tencentcloud/framework_provider.go) 落了 `tf5muxserver` 脚手架，但 framework provider 仍是空壳：`Schema` 为空、`Configure` 不做任何事、`Resources/DataSources` 返回 `nil`。本变更将其补齐为可生产使用的双栈架构。

## What Changes

- 在 `FrameworkProvider` 中实现与 SDKv2 等价的 provider schema（`secret_id`、`secret_key`、`security_token`、`region`、`protocol`、`domain`、`assume_role`、`shared_credentials_dir`、`profile`、`cam_role_name` 等所有现有字段），并复用 SDKv2 已实现的环境变量/配置文件回退逻辑。
- 实现 `FrameworkProvider.Configure`：构造 `*connectivity.TencentCloudClient`（与 SDKv2 完全相同的 meta 对象），通过 `resp.DataSourceData` / `resp.ResourceData` / `resp.EphemeralResourceData` 注入给 framework 侧的资源/数据源使用，确保两栈共享同一个 client、同一份凭证与重试逻辑。
- 建立 framework 侧的目录与命名规范：`tencentcloud/services/<service>/` 下与 SDKv2 文件并列新增 `resource_tc_<service>_<name>_framework.go`、`data_source_tc_<service>_<name>_framework.go`，并提供资源/数据源注册聚合点。
- 提供一套 framework 侧的公共基础设施：`tencentcloud/internal/tcfwhelper`（types 转换、retry、超时、错误归一化）、`tencentcloud/internal/tcfwprovider`（持有 client 的 ProviderMeta 类型与共享 meta 桥接）。
- 升级 mux：保持 `tf5muxserver` 不变（兼容 Terraform 0.13+），但要求所有 framework 资源声明 schema 时遵守 mux 的限制（同一资源类型不能同时在两栈注册；schema 字段名不能冲突）。在 CI 中加入 `tf5muxserver` 的 schema 冲突校验。
- 提供资源迁移指引与一个示范资源：选择 1 个低风险新资源（无存量 state）作为首个 framework 资源落地，验证整条链路（注册、Configure、CRUD、acceptance test、文档生成）。
- 文档与开发流程：更新 `CONTRIBUTING` 类指引（新增资源默认走 framework；存量资源不主动迁移），`make doc`、`make lint`、`make test` 必须同时覆盖两栈。
- **非破坏**：不变更任何已有资源的 schema、ID 格式、state 结构；不修改 [main.go](/Users/arunma/project/terraform-provider-tencentcloud/main.go) 对外暴露的 provider address。

## Capabilities

### New Capabilities
- `provider-framework-runtime`: framework provider 的 schema、Configure 流程、ProviderData 注入、与 SDKv2 共享 `connectivity.TencentCloudClient` 的契约。
- `provider-muxing`: SDKv2 与 framework 双栈通过 `tf5muxserver` 合并为单一 ProviderServer 的注册、冲突检测与发布约束。
- `provider-framework-conventions`: framework 侧资源/数据源的目录布局、文件命名、注册聚合、公共工具包（fwhelper/fwtransport）使用规范。

### Modified Capabilities
<!-- 当前 openspec/specs/ 下没有覆盖 provider 入口/runtime 的既有 spec，因此本次只新增能力，不修改既有 spec。 -->

## Impact

- **入口代码**：[main.go](/Users/arunma/project/terraform-provider-tencentcloud/main.go) 保持 mux 入口不变；[tencentcloud/framework_provider.go](/Users/arunma/project/terraform-provider-tencentcloud/tencentcloud/framework_provider.go) 大幅扩充（schema/Configure/Resources/DataSources）。
- **依赖**：新增/确认 `github.com/hashicorp/terraform-plugin-framework`、`terraform-plugin-framework-validators`、`terraform-plugin-mux`、`terraform-plugin-go` 在 `go.mod` 与 `vendor/` 中存在并锁定版本；需要 `go mod vendor` 同步。
- **共享层**：`tencentcloud/connectivity` 不变；新增 `tencentcloud/internal/fwhelper`、`tencentcloud/internal/fwtransport` 包。
- **服务包**：`tencentcloud/services/<service>/` 下未来可并列出现 SDKv2 与 framework 两套文件；本变更只落首个示范资源，其它服务零改动。
- **协议**：继续使用 `tf5muxserver`（Protocol v5），不强制升级到 v6，避免破坏 Terraform 0.13/0.14 兼容性承诺。
- **测试**：framework 资源 acceptance test 必须使用 `ProtoV5ProviderFactories`（而非 `Providers`）；现有 SDKv2 测试不受影响，但测试基类需要新增一个 framework 友好的 factory 帮助函数。
- **文档**：`make doc` / tfplugindocs 流程需验证对两栈资源都能生成文档；`website/docs/` 命名规则不变。
- **风险**：mux 下若 framework 与 SDKv2 注册了同名资源/数据源，启动即失败 → 通过 CI 检查与 PR 模板提示规避；framework 与 SDKv2 的 schema 类型不同（如 `types.String` vs `*string`），开发者上手成本上升 → 通过 `fwhelper` 与示范资源降低门槛。
