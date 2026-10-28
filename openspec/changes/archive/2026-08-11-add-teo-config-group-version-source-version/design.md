## Context

`tencentcloud_teo_config_group_version` 资源封装了 TEO（EdgeOne）配置组版本的创建与查询能力。当前资源只暴露 `zone_id`、`group_id`、`description`、`content` 四个入参字段，以及若干 Computed 出参字段（`version_id`、`version_number`、`group_type`、`status`、`create_time`）。

当前状态：
- 资源为 RESOURCE_KIND_GENERAL 类型，但仅有 Create/Read/Delete 三个 CRUD 方法（云 API 未提供独立的 Update 接口），现有 schema 中 `content`、`description`、`zone_id`、`group_id` 均标记为 `ForceNew: true`。
- Create 方法调用 `teov20220901.CreateConfigGroupVersionWithContext`，Read 方法通过 `TeoService.DescribeTeoConfigGroupVersionById` 调用 `DescribeConfigGroupVersionDetail` 获取详情。
- SDK 已支持 `CreateConfigGroupVersionRequest.SourceVersion`（可选 `*string`）入参，语义为"新版本所基于的来源版本 ID，未传入时默认采用当前生产环境正在生效的版本作为来源版本"。
- SDK 返回结构 `ConfigGroupVersionInfo.SourceVersion`（`*string`）已存在，但当前 Read 方法未读取该字段。

约束：
- 必须保持向后兼容：现有 TF 配置（未写 `source_version`）不能产生 plan diff，且不能破坏已有 state。
- 资源无 Update 接口，新增字段必须为 `ForceNew: true`，与现有 `content`/`description` 等字段保持一致。
- vendor 模式管理依赖，本次使用已 vendored 的 SDK 字段，无需变更 vendor。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_config_group_version` 资源中新增 `source_version` 可选参数，支持用户指定新版本所基于的来源版本 ID。
- 创建时将 `source_version` 传递给 `CreateConfigGroupVersion` API 的 `request.SourceVersion`。
- 读取时从 `ConfigGroupVersionInfo.SourceVersion` 回填到 Terraform 状态。
- 保持向后兼容：未配置 `source_version` 时沿用云平台默认行为。
- 通过单元测试覆盖 `source_version` 的创建与读取场景。

**Non-Goals:**
- 不改变 Create/Read/Delete 的现有行为（除新增 `source_version` 字段处理外）。
- 不新增 Update 方法或 Update 接口调用（云 API 未提供配置组版本的独立更新接口）。
- 不修改 `content`、`description`、`zone_id`、`group_id` 等现有字段的行为。
- 不变更 vendor 依赖版本。
- 不新增 Timeouts 子块（资源无异步 Update 操作）。

## Decisions

### Decision 1: `source_version` 字段为 Optional + ForceNew

**选择**：在 Schema 中将 `source_version` 声明为 `Optional: true, ForceNew: true, TypeString`。

**备选**：声明为 `Optional + Computed`（不强制重建）。

**理由**：
- 该资源无 Update 接口（仅 CRD），所有可变入参均为 `ForceNew`。若用户修改 `source_version`，语义上是"基于不同的来源版本派生新版本"，必须重建资源，因此 `ForceNew` 与现有字段行为一致。
- 不使用 `Computed`：`source_version` 是用户主动指定的入参，而非云端自动生成。Read 阶段虽会回填，但遵循现有 `description` 等字段的 `ForceNew`（非 Computed）模式即可，避免引入 `Optional+Computed` 导致的 plan diff 复杂性。

### Decision 2: Create 方法中仅当 `d.GetOk` 为 true 时设置 `request.SourceVersion`

**选择**：在 Create 方法中使用 `if v, ok := d.GetOk("source_version"); ok { request.SourceVersion = helper.String(v.(string)) }`，未配置时不设置该字段（保持 nil）。

**理由**：
- 与现有 `description` 字段处理方式一致。
- SDK 注释明确说明"该字段可选，未传入时默认采用当前生产环境正在生效的版本作为来源版本"，保持 nil 即可触发云平台默认行为，确保向后兼容。

### Decision 3: Read 方法中对 `SourceVersion` 进行 nil 检查后回填

**选择**：在 Read 方法中，于 `respData.ConfigGroupVersionInfo != nil` 分支内新增 `if respData.ConfigGroupVersionInfo.SourceVersion != nil { _ = d.Set("source_version", respData.ConfigGroupVersionInfo.SourceVersion) }`。

**理由**：
- 遵循项目规范"在调用 setXX() 设置字段前，请先判断 Response 中的字段是否为 nil"。
- 老版本创建的资源或未指定来源版本时，云端可能返回空字符串或 nil，nil 检查避免空指针。

### Decision 4: 字段放置位置

**选择**：将 `source_version` 字段放置在 `content` 字段之后、`version_id`（Computed 字段）之前，保持入参在前、Computed 出参在后的分组风格。

**理由**：
- 与现有 schema 中入参（`zone_id`/`group_id`/`description`/`content`）在前、Computed 出参（`version_id` 等）在后的布局一致。

### Decision 5: 单元测试使用 gomonkey mock

**选择**：在 `resource_tc_teo_config_group_version_test.go` 中补充测试用例时，使用 gomonkey mock 云 API（`CreateConfigGroupVersionWithContext` 与 `DescribeConfigGroupVersionDetail`），进行业务代码逻辑的单元测试。

**理由**：
- 遵循项目规范："对于新增的 terraform 资源，生成 *_test.go 时，不要使用 terraform 的测试套件，而是使用 mock（gomonkey）的方法对云 API 进行 mock 处理"。
- 现有测试文件使用的是 terraform 测试套件（`resource.Test`），但本次属于新增资源参数场景，按规范应使用 gomonkey mock 方式补充单元测试。

## Risks / Trade-offs

- **Risk**：已有 state 中不包含 `source_version` 字段，升级后首次 `terraform plan` 是否产生 diff → **Mitigation**：新增字段为 `Optional`，未配置时不会产生 diff；Read 回填的值若与用户未配置的零值一致也不会触发重建。
- **Risk**：云端对未指定 `source_version` 的版本返回空字符串而非 nil，导致 Read 回填空字符串后与用户配置产生 diff → **Mitigation**：nil 检查仅在非 nil 时回填；空字符串场景下 terraform 会将其视为已配置的零值，但因字段为 `Optional` 且未配置时 state 为空字符串，diff 行为可接受。
- **Trade-off**：`source_version` 设为 `ForceNew` 意味着修改来源版本将重建资源（删除旧版本记录并创建新版本） → 这是合理行为，因云 API 无 Update 接口，版本本身不可变。

## Migration Plan

- 新增字段为纯加法（Optional 追加），无 state 迁移需求。
- 存量资源：Terraform state 中不含 `source_version`，升级后 `terraform plan` 对未在 HCL 配置中声明该字段的资源不会产生 diff（Optional 行为）。
- 文档更新：在 `resource_tc_teo_config_group_version.md` 中补充 `source_version` 字段示例。
- 回滚：若需回退，只需移除 schema 中 `source_version` 字段及 Create/Read 中的相关处理代码；state 中已回填的值不会影响其他字段。

## Open Questions

- 无
