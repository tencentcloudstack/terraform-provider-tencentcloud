## Why

仓里上一版 framework 布局规范（`restructure-framework-types-and-naming`）选择了 **`<product>/<type>/` 双层目录** + **`<product><type>` 拼接包名**（如 `metaresources` / `metaephemerals`） + **`*_ephemeral_resource.go` 后缀法文件名**。该方案在评审中暴露了三个问题：

1. **与 SDKv2 心智割裂**：现存 `tencentcloud/services/<product>/` 全是单层，文件名一律用 **类型前缀法**（`resource_tc_<product>_<name>.go` / `data_source_tc_<product>_<name>.go`）。新增 framework 类型时，开发者要在两套完全不同的目录/命名规则之间切换，认知成本高。
2. **包名拼接劝退新人**：`metaresources` / `metaephemerals` 这类拼接包名不符合 Go 社区惯例，主流 provider（AWS / Google / Azure 各自的 framework 子树）都是直接 `package <product>`，靠 import path 区分。
3. **后缀法 + 拼接包名的真实收益为零**：拼接包名是为了"避免 registry import alias"，但只要包名直接用 `<product>`（不同 import path、可同名）一样不需要 alias；后缀法（`*_ephemeral_resource.go`）则在前缀法（`ephemeral_tc_<product>_<name>.go`）面前没有任何可读性优势，反而又造出一套要记忆的规则。

旧 change 在仓里的 `tasks.md` 虽然标记为 `[x]`，但代码侧 `tencentcloud/framework/meta/` 实际目录**并未落地**——这是把布局错误从源头改正的最佳窗口期，越往后改成本越高。

## What Changes

- **BREAKING**（仅对未实现的旧 spec 而言，对运行时零影响）：framework 业务实现目录布局从 `<product>/<type>/` **双层**改为 `<product>/` **单层**。
- **BREAKING**（同上）：framework 子包包名规则从 `<product><type>`（如 `metaresources`）改为 `<product>`（如 `meta` / `ssm` / `cvm`），与 SDKv2 现有 `services/<product>/` 心智完全对齐；registry 直接 import 各产品包，不需 alias（不同 import path 天然唯一）。
- **BREAKING**（同上）：framework 文件名规则从"类型后缀法"（`*_resource.go` / `*_data_source.go` / `*_ephemeral_resource.go` / …）改为"类型前缀法"（`resource_tc_<product>_<name>.go` / `data_source_tc_<product>_<name>.go` / `ephemeral_tc_<product>_<name>.go` / `function_tc_<product>_<name>.go` / `list_tc_<product>_<name>.go` / `action_tc_<product>_<name>.go`），与 SDKv2 完全对齐。
- 旧 change `restructure-framework-types-and-naming` 在 `openspec/changes/` 下**整个目录硬改写**：spec / design / tasks / proposal 文档中所有"双层路径 / 拼接包名 / 后缀法"的文字一律改为新规范，避免今后回看时的误读。`tasks.md` 中的 `[x]` 标记保留（行为意图未变：reference 实现仍然要落地，只是落到新规范的路径上）。
- 同步刷新 `tencentcloud/framework/README.md`、`tencentcloud/framework/registry.go` 顶部说明（如已写入旧规范字样）。

## Capabilities

### New Capabilities

- `framework-provider-layout-and-naming`: framework provider 的目录结构、文件命名、Go 包命名、registry 聚合方式的硬约束。覆盖以下方面：
  - 单层产品目录布局（`tencentcloud/framework/<product>/`）
  - 类型前缀法文件命名（`<type>_tc_<product>_<name>.go`）
  - 直接同名 Go 包（`package <product>`）
  - registry 直接 import 产品包、无 alias
  - 跨产品 / 元产品 reference 落到 `meta/`
  - 6 种 plugin-framework 类型（resource / datasource / function / ephemeral / list / action）的 reference 实现要求

### Modified Capabilities

（无：旧 change `restructure-framework-types-and-naming` 的 spec 尚未 archive 进 `openspec/specs/`，其条款本质上仍是"提案态"。本 change 用一份全新的 New Capability 完整取代旧条款，并在旧 change 目录内做硬改写以消除误导。）

## Impact

- **代码影响**：零。旧 change 的代码搬迁动作（`framework/meta/` 树）实际未执行；本 change 落地后仍由后续 `openspec-apply-change` 一次性按新规范创建 framework 目录与 reference 实现。
- **文档影响**：
  - 新增 `openspec/changes/restructure-framework-single-level-layout/` 全套（proposal / design / specs / tasks）。
  - 硬改写 `openspec/changes/restructure-framework-types-and-naming/` 下的 `proposal.md` / `design.md` / `specs/framework-provider-types-and-naming/spec.md` / `tasks.md`：将所有"双层 / 拼接包名 / 后缀法"表述替换为"单层 / 同名包 / 前缀法"，保持其他语义不变。
- **依赖影响**：无（不动 `go.mod` / `vendor/`，旧 change 已锁定的 plugin-framework v1.19+ 与 mux v0.23.1+ 仍然适用）。
- **下游影响**：
  - 后续 `tencentcloud_ssm_secret_version`（ephemeral 类型）将按新规范落到 `tencentcloud/framework/ssm/ephemeral_tc_ssm_secret_version.go`（包名 `ssm`），而非旧规范的 `framework/ssm/ephemerals/ssm_secret_version_ephemeral_resource.go`（包名 `ssmephemerals`）。
  - 已规划的 reference 文件名同步刷新：
    - `local_note_resource.go` → `resource_tc_meta_local_note.go`
    - `provider_runtime_data_source.go` → `data_source_tc_meta_provider_runtime.go`
    - `temp_credential_ephemeral_resource.go` → `ephemeral_tc_meta_temp_credential.go`
    - `region_list_resource.go` → `list_tc_meta_region.go`
    - `parse_resource_id_function.go`（旧 change 未起名时的占位）→ `function_tc_meta_parse_resource_id.go`
    - `reboot_instance_action.go` → `action_tc_cvm_reboot_instance.go`
