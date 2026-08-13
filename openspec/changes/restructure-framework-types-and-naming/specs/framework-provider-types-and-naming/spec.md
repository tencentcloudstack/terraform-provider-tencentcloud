## ADDED Requirements

### Requirement: Framework Helper Package Naming

framework 侧的通用 helper 包 SHALL 位于 `tencentcloud/framework/internal/helper/`，包名 SHALL 为 `helper`；该包 SHALL 受 Go `internal/` 可见性规则约束，仅 `tencentcloud/framework/...` 子树可 import。SHALL NOT 在仓内任何位置出现 `tcfwhelper` 或 `frameworkhelper` 标识符（包括 import 路径与符号引用）。

#### Scenario: 通用 helper 引用
- **WHEN** 任何 `tencentcloud/framework/` 子树下的 `.go` 文件需要使用 framework 侧的 types 转换 / retry / timeouts 等通用工具
- **THEN** 该文件 import 路径 MUST 为 `github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/internal/helper`
- **AND** 该文件 MUST 通过 `helper.<Symbol>` 形式调用，不得使用其他别名

#### Scenario: internal 可见性硬约束
- **WHEN** 任何位于 `tencentcloud/framework/` 子树**之外**的 `.go` 文件尝试 import `tencentcloud/framework/internal/helper`
- **THEN** Go 编译器 MUST 拒绝该 import（由 `internal/` 包可见性规则保证）

#### Scenario: 全仓无遗留命名
- **WHEN** 在仓内执行 `grep -rE 'tcfwhelper|internal/frameworkhelper|frameworkhelper\.' -- 'tencentcloud/' main.go`
- **THEN** 命中条数 MUST 为 0
- **AND** `openspec/changes/<active-changes>/` 与历史 `archive/` 目录中的旧引用 MAY 保留（仅历史文档）

### Requirement: Shared Meta Package Naming

跨 SDKv2 与 framework 两侧共享 `*connectivity.TencentCloudClient` 的桥接包 SHALL 位于 `tencentcloud/internal/sharedmeta/`，包名 SHALL 为 `sharedmeta`，且必须导出 `ProviderMeta` / `SetSharedMeta` / `GetSharedMeta` / `ResetSharedMetaForTest` 四个符号以保持与原 `tcfwprovider` 同名 API 的语义一致。

#### Scenario: SDKv2 注入共享 client
- **WHEN** SDKv2 `providerConfigure` 完成 `tcClient` 构造
- **THEN** 同一函数 MUST 在 `return &tcClient, nil` 之前调用 `sharedmeta.SetSharedMeta(tcClient.apiV3Conn)`

#### Scenario: framework 读取共享 client
- **WHEN** framework `Provider.Configure` 被调用
- **THEN** 该方法 MUST 通过 `sharedmeta.GetSharedMeta()` 读取 client
- **AND** 返回 nil 时 MUST 在 `resp.Diagnostics` 追加 Error 诊断
- **AND** 返回非 nil 时 MUST 把 `&sharedmeta.ProviderMeta{Client: client}` 写入 `resp.ResourceData`、`resp.DataSourceData`、`resp.EphemeralResourceData`、`resp.ActionData`

#### Scenario: 全仓无遗留命名
- **WHEN** 在仓内执行 `grep -r tcfwprovider -- 'tencentcloud/**/*.go'`
- **THEN** 命中条数 MUST 为 0

### Requirement: Framework Package Single Home

framework 一切（provider 入口、registry、业务实现、测试）SHALL 统一位于 `tencentcloud/framework/` 目录树下；SHALL NOT 再存在 `tencentcloud/provider/framework/` 或 `tencentcloud/services/tcprovider/` 这两个旧目录。`tencentcloud/framework/` 顶层 Go 包名 SHALL 为 `framework`。

#### Scenario: 顶层入口位置
- **WHEN** 检查 `tencentcloud/framework/` 目录的直接子项（不递归）
- **THEN** MUST 至少存在 `provider.go` / `registry.go` / `provider_test.go` / `testhelpers_test.go` / `README.md`
- **AND** `provider.go` / `registry.go` 的 package 声明 MUST 为 `package framework`

#### Scenario: 旧目录全部清退
- **WHEN** 在仓内检查 `tencentcloud/provider/framework/` 与 `tencentcloud/services/tcprovider/`
- **THEN** 这两个目录 MUST 都不存在（目录与其下文件均已删除）

#### Scenario: import 路径清退
- **WHEN** 在仓内执行 `grep -rEn 'tencentcloud/provider/framework|tencentcloud/services/tcprovider' -- 'tencentcloud/' main.go`
- **THEN** 命中条数 MUST 为 0

### Requirement: Framework Product Single-Level Layout

`tencentcloud/framework/` 下的业务实现 SHALL 采用**单层产品目录**布局：每一个云产品（或元产品 `meta`）SHALL 对应一个 `tencentcloud/framework/<product>/` 目录，且该目录下 SHALL 直接放置该产品的全部 plugin-framework 实现文件（resource / datasource / function / ephemeral / list / action 不再各自下沉到子目录）。SHALL NOT 在 `tencentcloud/framework/<product>/` 下创建 `resources/` / `datasources/` / `functions/` / `ephemerals/` / `lists/` / `actions/` 这 6 个类型子目录中的任何一个。

#### Scenario: 业务实现位置
- **WHEN** 在 `tencentcloud/framework/` 下新增任何 plugin-framework 类型实例
- **THEN** 该实例 MUST 位于 `tencentcloud/framework/<product>/<file>.go` 形式的**单层**路径下
- **AND** `<product>` MUST 是真实云产品的 lowercase 名称（如 `cvm` / `vpc` / `cbs` / `ssm`）或元产品 `meta`
- **AND** 该文件 MUST NOT 落到 `<product>/<type>/<file>.go` 这种两层路径

#### Scenario: 类型子目录禁用
- **WHEN** 检查任意 `tencentcloud/framework/<product>/` 目录的直接子项
- **THEN** 该目录 MUST NOT 包含名为 `resources` / `datasources` / `functions` / `ephemerals` / `lists` / `actions` 的子目录

#### Scenario: 产品归属约束
- **WHEN** 新增的 reference 能明确归属于一个真实云产品（如 CVM 的 reboot action、SSM 的 secret_version ephemeral）
- **THEN** 该实现 MUST 落到对应产品目录（如 `tencentcloud/framework/cvm/` / `tencentcloud/framework/ssm/`）
- **AND** 跨产品 / 不归属任何具体云产品的 reference（如 provider runtime / region list / parse_resource_id）MUST 落到 `tencentcloud/framework/meta/` 下

### Requirement: Framework File Naming Convention

`tencentcloud/framework/<product>/` 下的实现文件 SHALL 采用**类型前缀法**命名：`<type>_tc_<product>_<name>.go`，其中 `<type>` SHALL 取自下表的 6 个固定值，`<product>` SHALL 与所在目录名一致，`<name>` SHALL 为业务领域名。SHALL NOT 使用类型后缀法（如 `*_resource.go` / `*_ephemeral_resource.go` / `*_action.go` 等）；SHALL NOT 出现 `example_` / `demo_` / `sample_` 等占位前缀。

| 类型 | 文件名前缀 |
|---|---|
| resource | `resource_tc_` |
| datasource | `data_source_tc_` |
| function | `function_tc_` |
| ephemeral | `ephemeral_tc_` |
| list | `list_tc_` |
| action | `action_tc_` |

#### Scenario: 实现文件命名
- **WHEN** 在 `tencentcloud/framework/<product>/` 下新增任一 plugin-framework 类型的 `.go` 实现文件
- **THEN** 该文件名 MUST 形如 `<type>_tc_<product>_<name>.go`，其中 `<type>` ∈ {`resource`, `data_source`, `function`, `ephemeral`, `list`, `action`}
- **AND** 文件名中的 `<product>` 段 MUST 与其所在目录名完全一致

#### Scenario: 测试与文档文件命名
- **WHEN** 为某实现文件 `<type>_tc_<product>_<name>.go` 创建测试或文档
- **THEN** 测试文件名 MUST 为 `<type>_tc_<product>_<name>_test.go`
- **AND** 文档文件名 MUST 为 `<type>_tc_<product>_<name>.md`（gendoc 兼容）

#### Scenario: 后缀法禁用
- **WHEN** 在仓内执行 `find tencentcloud/framework -type f -name '*_resource.go' -o -name '*_data_source.go' -o -name '*_ephemeral_resource.go' -o -name '*_function.go' -o -name '*_list_resource.go' -o -name '*_action.go'`
- **THEN** 命中条数 MUST 为 0
- **AND** （该规则仅约束 `tencentcloud/framework/` 子树；`tencentcloud/services/` 下的 SDKv2 文件不受影响）

### Requirement: Framework Sub-package Naming Convention

`tencentcloud/framework/<product>/` 下的 Go 包名 SHALL 直接采用 `<product>`（产品名本身、全小写、不含分隔符）。SHALL NOT 使用 `<product><type>` 拼接（如 `metaresources` / `metaephemerals` / `cvmactions`），SHALL NOT 使用 `fw<product>` / `<product>fw` 等带前/后缀的形式。同名包靠 import path 区分，天然不需要 alias。

#### Scenario: 包名规则
- **WHEN** 检查 `tencentcloud/framework/<product>/` 目录下任一 `.go` 文件
- **THEN** 其 `package` 声明 MUST 为 `package <product>`（例：`package meta` / `package ssm` / `package cvm`）
- **AND** 测试文件的 `package` 声明 MUST 为 `package <product>` 或 `package <product>_test`（外部测试包惯例）

#### Scenario: 旧拼接包名清退
- **WHEN** 在仓内执行 `grep -rEn 'package\s+(meta(resources|datasources|functions|ephemerals|lists)|cvmactions|ssmephemerals|[a-z]+(resources|datasources|functions|ephemerals|lists|actions))\b' -- 'tencentcloud/framework/'`
- **THEN** 命中条数 MUST 为 0

### Requirement: Framework Registry Direct Import Without Alias

`tencentcloud/framework/registry.go` SHALL 是 framework 6 类型工厂的**唯一聚合点**，SHALL 直接 import 各产品包并通过 `<product>.NewXxx` 构造工厂 slice，暴露 6 个聚合方法：`Resources` / `DataSources` / `Functions` / `EphemeralResources` / `ListResources` / `Actions`。registry 对各产品包的 import SHALL NOT 使用任何 alias（不同 import path 天然唯一，无需 alias）；registry 之外 SHALL NOT 存在中间聚合层。

#### Scenario: registry import 无 alias
- **WHEN** 检查 `tencentcloud/framework/registry.go` 的 import 块
- **THEN** 对每一个产品包的 import 行 MUST 形如 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/<product>"`（无 alias 前缀）
- **AND** 6 个聚合方法的方法体 MUST 直接通过 `<product>.NewXxx` 引用各产品包内导出的工厂函数

#### Scenario: 6 个聚合方法均在 registry
- **WHEN** 检查 `tencentcloud/framework/registry.go`
- **THEN** MUST 同时定义 `Resources` / `DataSources` / `Functions` / `EphemeralResources` / `ListResources` / `Actions` 6 个方法
- **AND** 每个方法 MUST 返回 `[]func() <iface>{ <product>.NewXxx, ... }` 形式的 slice，方法体内不调用其他聚合函数

#### Scenario: 无中间聚合层
- **WHEN** 在仓内执行 `grep -rEn 'func\s+Framework(Resources|DataSources|Functions|EphemeralResources|ListResources|Actions)\s*\(' -- 'tencentcloud/'`
- **THEN** 命中条数 MUST 为 0

### Requirement: Framework Top-Level Files Stay At Framework Root

`tencentcloud/framework/` 顶层 SHALL 保留 provider 入口、registry、provider 测试、测试辅助 4 类文件，且其 `package` 声明 SHALL 为 `package framework`。SHALL NOT 把这些顶层文件下沉到任何产品目录。

#### Scenario: 顶层入口位置
- **WHEN** 检查 `tencentcloud/framework/` 目录的直接子项（不递归）
- **THEN** MUST 至少存在 `provider.go` / `registry.go` / `provider_test.go` / `testhelpers_test.go`（或语义等价文件）以及 `README.md`
- **AND** `provider.go` / `registry.go` 的 package 声明 MUST 为 `package framework`

#### Scenario: 顶层不放业务实现
- **WHEN** 检查 `tencentcloud/framework/` 目录的直接子项 `.go` 文件（不递归）
- **THEN** MUST NOT 出现 `resource_tc_*.go` / `data_source_tc_*.go` / `function_tc_*.go` / `ephemeral_tc_*.go` / `list_tc_*.go` / `action_tc_*.go` 形式的文件（任何业务实现都属于某个产品包）

### Requirement: Reference Implementations Per Type Under Single-Level Layout

framework 6 种类型 SHALL 各包含至少一份 reference 实现，落地路径与命名 SHALL 严格遵循单层布局 + 类型前缀法。默认 SHALL 通过 registry 注册到 provider schema（L2）；MAY 降级为 L0 占位仅在接口起点跟业务 reference 意图不匹配且全量实现超出 apply scope 的情形下适用。所有实现 MUST NOT 触发任何真实云 API 调用；命名采用业务领域名（不含 `example` 前缀）。

#### Scenario: resource reference
- **WHEN** 检查 `tencentcloud/framework/meta/` 目录
- **THEN** MUST 至少存在一个文件名形如 `resource_tc_meta_<name>.go`、实现 `resource.Resource` 接口的 `.go` 文件（建议 `resource_tc_meta_local_note.go`）
- **AND** 其 `Metadata.TypeName` MUST 形如 `tencentcloud_local_note`
- **AND** 其 CRUD 实现 MUST 仅操作进程内 in-memory 状态，不得调用 SDK / HTTP

#### Scenario: datasource reference
- **WHEN** 检查 `tencentcloud/framework/meta/` 目录
- **THEN** MUST 至少存在一个文件名形如 `data_source_tc_meta_<name>.go`、实现 `datasource.DataSource` 接口的 `.go` 文件（建议 `data_source_tc_meta_provider_runtime.go`）
- **AND** 该实现 MUST 是从原 `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go` 迁入并改名后的版本，schema / Read 逻辑 MUST 与迁移前**逐字段一致**

#### Scenario: function reference
- **WHEN** 检查 `tencentcloud/framework/meta/` 目录
- **THEN** MUST 至少存在一个文件名形如 `function_tc_meta_<name>.go`、实现 `function.Function` 接口的 `.go` 文件（建议 `function_tc_meta_parse_resource_id.go`）
- **AND** 其 `Metadata.Name` MUST 形如 `parse_resource_id`
- **AND** 其 `Run` 实现 MUST 是纯字符串处理，无 IO

#### Scenario: ephemeral reference
- **WHEN** 检查 `tencentcloud/framework/meta/` 目录
- **THEN** MUST 至少存在一个文件名形如 `ephemeral_tc_meta_<name>.go`、实现 `ephemeral.EphemeralResource` 接口的 `.go` 文件（建议 `ephemeral_tc_meta_temp_credential.go`）
- **AND** 其 `Metadata.TypeName` MUST 形如 `tencentcloud_temp_credential`
- **AND** 其 `Open` 实现 MUST 仅返回本地构造的占位凭证（不得调用 STS / CAM API）

#### Scenario: list reference【L0 降级占位】
- **WHEN** 检查 `tencentcloud/framework/meta/` 目录
- **THEN** MUST 至少存在一个文件名形如 `list_tc_meta_<name>.go`（建议 `list_tc_meta_region.go`）的占位 `.go` 文件
- **AND** 该文件顶部 doc 注释 MUST 明确说明本文件是 L0 占位（仅提供静态 region 数据 + helper，未实现 `list.ListResource` 接口）以及暂不实现的原因（framework v1.19 的 ListResource 要求同名 managed resource + ResourceIdentity，超出 apply scope）
- **AND** 该文件提供的静态区域数据 MUST 至少包含 5 条，且每条 ID/Name 非空
- **AND** 该包 MUST NOT 被 `tencentcloud/framework/registry.go` 通过聚合方法注册到 provider schema（L0 占位不接入 registry）

#### Scenario: action reference（按产品归属落地）
- **WHEN** 检查 `tencentcloud/framework/cvm/` 目录
- **THEN** MUST 至少存在一个文件名形如 `action_tc_cvm_<name>.go`、实现 `action.Action` 接口的 `.go` 文件（建议 `action_tc_cvm_reboot_instance.go`）
- **AND** 其 `Metadata.TypeName` MUST 形如 `tencentcloud_reboot_instance`
- **AND** 其 `Invoke` 实现 MUST 仅做参数校验与日志记录，不得调用 CVM API
- **AND** `Invoke` 返回的错误诊断仅来自参数校验，永不来自网络/SDK 错误

### Requirement: Framework Provider Action Support

framework provider 实现 SHALL 暴露 `Actions(ctx context.Context) []func() action.Action` 方法，且 `tencentcloud/framework/registry.go` SHALL 通过直接 import 产品包提供其返回值。

#### Scenario: provider 接口完整
- **WHEN** 编译 `tencentcloud/framework/provider.go`
- **THEN** 其类型 MUST 同时满足 `provider.Provider` / `provider.ProviderWithFunctions` / `provider.ProviderWithEphemeralResources` / `provider.ProviderWithListResources` / `provider.ProviderWithActions` 全部接口
- **AND** `Actions` 方法 MUST 直接返回 `r.Actions()`，其中 `r` 为本 provider 持有的 registry 实例

#### Scenario: action mux 协议可用
- **WHEN** 启动 `main.go` 的 mux server
- **THEN** mux server MUST 能在不报错的前提下注册 framework 这一栈的 action 集合
- **AND** `tf5muxserver` 或 `tf6muxserver` 的版本 MUST 支持 action 协议；若不支持则 `main.go` MUST 显式升级到支持的 mux 实现

### Requirement: Plugin Framework Minimum Version

`go.mod` 中的 `github.com/hashicorp/terraform-plugin-framework` 版本 SHALL ≥ `v1.19.0`，且 `vendor/github.com/hashicorp/terraform-plugin-framework/action/` 目录 MUST 存在并包含 `action.go`。

> 实测：`v1.19.0` 上游已包含完整 `action` 包；`terraform-plugin-mux v0.23.1` 已发布 `tf5muxserver/mux_server_InvokeAction.go` 支持 action 协议。本次 change 不进行依赖升级。

#### Scenario: 依赖锁定
- **WHEN** 读取 `go.mod`
- **THEN** 其中 `github.com/hashicorp/terraform-plugin-framework` 的语义版本 MUST 满足 `>= 1.19.0`

#### Scenario: vendor 完整
- **WHEN** 构建产物的 `vendor/` 目录
- **THEN** `vendor/github.com/hashicorp/terraform-plugin-framework/action/` 目录 MUST 存在且包含 `action.go`
- **AND** `go mod verify` MUST 返回 `all modules verified`

### Requirement: Framework Test Co-location

framework provider 的单元测试与辅助测试 SHALL 位于 `tencentcloud/framework/` 目录下，与生产代码同包；SHALL NOT 位于 `tencentcloud/` 顶层目录或旧的 `tencentcloud/provider/framework/` 目录。

#### Scenario: 测试文件位置
- **WHEN** 检查 `tencentcloud/framework/`
- **THEN** MUST 存在 `provider_test.go` 与 `testhelpers_test.go`（或语义等价文件）
- **AND** 这两个文件的 package 声明 MUST 为 `package framework`

#### Scenario: 顶层目录清理
- **WHEN** 检查 `tencentcloud/`（不递归）
- **THEN** MUST NOT 存在 `framework_provider_test.go` 或 `framework_provider_testhelpers_test.go`

### Requirement: No Import Alias Required

所有 framework provider 的引用方 SHALL 直接使用包名 `framework` 而非别名 `fwprovider`。各产品子包（`meta` / `cvm` / `ssm` 等）在 registry 中同样直接 import 无需 alias（不同 import path 天然唯一）。

#### Scenario: 顶层 main / 测试 / 工厂
- **WHEN** 检查 `main.go` / `tencentcloud/framework/acctest/factories.go` / `tencentcloud/framework/*_test.go` 中对 `tencentcloud/framework` 的 import
- **THEN** 该 import MUST NOT 使用任何别名（直接 `"github.com/.../tencentcloud/framework"` 即可）
- **AND** 文件中对该包符号的引用 MUST 形如 `framework.NewProvider(...)`

#### Scenario: registry import 产品包无 alias
- **WHEN** 检查 `tencentcloud/framework/registry.go` 中对产品包的 import
- **THEN** 各 import 行 MUST NOT 使用 alias，直接形如 `"github.com/.../tencentcloud/framework/meta"` 即可

#### Scenario: 全仓无遗留别名
- **WHEN** 在仓内执行 `grep -rE 'fwprovider\s+"' -- 'tencentcloud/**/*.go' main.go`
- **THEN** 命中条数 MUST 为 0
