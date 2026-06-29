> **Supersedes**: 本 spec 完整取代旧 change `restructure-framework-types-and-naming` 中 `spec.md` 的"Framework Product/Type Two-Level Layout"与"Framework Sub-package Naming Convention"条款。旧条款的双层目录 / 拼接包名 / 后缀法文件名已全部被本 spec 的单层布局 / 同名包 / 前缀法否定。

## ADDED Requirements

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

`tencentcloud/framework/<product>/` 下的 Go 包名 SHALL 直接采用 `<product>`（产品名本身、全小写、不含分隔符）。SHALL NOT 使用 `<product><type>` 拼接（如 `metaresources` / `metaephemerals` / `cvmactions`），SHALL NOT 使用 `fw<product>` / `<product>fw` 等带前/后缀的形式。

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
