## Context

`tencentcloud_teo_config_group_version` 是管理 EdgeOne (TEO) 配置组版本的 RESOURCE_KIND_GENERAL 资源，其复合 ID 由 `ZoneId#GroupId#VersionId` 拼接而成。当前资源在 Create 时已从 `CreateConfigGroupVersion` 响应中获取 `VersionId` 并拼入复合 ID；但在 Read 阶段从 `DescribeConfigGroupVersionDetail` 读取并回填时，`version_id` 未作为独立的、规范的顶层出参向用户暴露。

当前状态：
- `resource_tc_teo_config_group_version.go` 的 schema 已声明 `version_id` 为 `Computed, TypeString`。
- Read 方法已通过 `respData.ConfigGroupVersionInfo.VersionId` 调用 `d.Set("version_id", ...)`（仅在非 nil 时设置）。
- 云 API `DescribeConfigGroupVersionDetail` 返回结构 `Response.ConfigGroupVersionInfo.VersionId`（`*string`，格式如 `ver-2kplomhisdcb`），SDK 已 vendored。

约束：
- Terraform Provider 必须保持向后兼容：新增 Computed 出参不会破坏现有 TF 配置和 state（Computed 字段由 Read 回填，不产生 plan diff）。
- 资源 ID 为复合 ID（`ZoneId#GroupId#VersionId`），`version_id` 仅作为只读出参，不能设为 ForceNew 或 Required。

## Goals / Non-Goals

**Goals:**
- 确保 `version_id` 作为规范的顶层 Computed 出参暴露给用户，便于跨资源引用与 state 查询。
- Read 方法从 `DescribeConfigGroupVersionDetail` 的 `Response.ConfigGroupVersionInfo.VersionId` 正确回填该字段，并对 nil 做安全处理。
- 通过单元测试覆盖 `version_id` 出参的读取逻辑（使用 gomonkey mock 云 API）。

**Non-Goals:**
- 不改变 Create/Delete 行为（该资源无 Update 接口，所有顶层字段除 id 外均为 ForceNew）。
- 不修改复合 ID 的拼接格式（仍为 `ZoneId#GroupId#VersionId`）。
- 不新增独立的 `tencentcloud_teo_config_group_version` 数据源资源。
- 不引入新的外部依赖或变更 vendor。

## Decisions

### Decision 1: `version_id` 作为 Computed 只读出参，不改 ForceNew

**选择**：`version_id` 保持 `Computed: true, TypeString`，不设置 `Optional`/`Required`/`ForceNew`。

**备选**：将 `version_id` 改为 `Required + ForceNew` 供用户指定。

**理由**：
- `version_id` 由云 API 在创建时分配，是云端生成资源标识，用户无法预先指定。
- 该字段已是复合 ID 的组成部分，设为 Computed 出参即可满足用户读取与引用需求，无需改变其生命周期语义。
- 保持 Computed 与现有 schema 及向后兼容性一致，不会产生 plan diff。

### Decision 2: Read 方法在非 nil 时设置，沿用现有 nil 安全模式

**选择**：在 `resourceTencentCloudTeoConfigGroupVersionRead` 中，沿用现有 `if respData.ConfigGroupVersionInfo.VersionId != nil { _ = d.Set("version_id", respData.ConfigGroupVersionInfo.VersionId) }` 模式。

**理由**：
- 符合项目规范：在调用 `setXX()` 前判断 Response 字段是否为 nil，nil 时不设置。
- 符合资源 Read 方法中"云 API 返回空时先打印 `[CRUD]` 日志保留现场再 `d.SetId("")`"的约束（本变更不涉及该路径，但保持一致风格）。

### Decision 3: 单元测试使用 gomonkey mock 云 API

**选择**：由于本资源为新资源（按需求定义为新增参数），单测在 `resource_tc_teo_config_group_version_test.go` 中使用 gomonkey mock `DescribeConfigGroupVersionDetail` 接口，只测试业务代码逻辑，不使用 terraform 测试套件。

**备选**：沿用现有 terraform 测试套件补充用例。

**理由**：
- 项目规范要求：新增 terraform 资源时使用 mock（gomonkey）方法对云 API 进行 mock 处理，只做业务代码逻辑的单元测试。
- mock 方式可避免依赖真实云环境与凭据，保证单测在 CI 中可稳定构建执行。

### Decision 4: 文档更新通过 `.md` 文件驱动，禁止手改 `website/`

**选择**：仅在 `resource_tc_teo_config_group_version.md` 中补充 `version_id` 出参说明，`website/docs/` 由收尾阶段 `make doc` 自动生成。

**理由**：
- 项目规范禁止直接新增/修改 `website/` 目录文件，只能在收尾阶段通过 `make doc` 生成。
- `.md` 文件作为 `make doc` 的输入源，需同步补充新增字段。

## Risks / Trade-offs

- **Risk**：已有 state 中 `version_id` 已有值，本次规范其为出参后行为不变 → **Mitigation**：保持 `Computed` 语义，Read 回填逻辑不变，无 state drift。
- **Risk**：Read 阶段云 API 短暂波动导致 `ConfigGroupVersionInfo` 为 nil 时 `version_id` 不被设置 → **Mitigation**：沿用现有 nil 安全判断，不主动清空 id；已有 retry 机制会在外层重试。
- **Trade-off**：`version_id` 仅可读不可写，用户无法手动指定 → 可接受，该值为云端生成的资源标识。

## Migration Plan

- 新增字段属性为纯加法（Computed 出参规范化），无 state 迁移需求。
- 存量资源：Terraform state 中已有 `version_id` 值，升级后 `terraform plan` 对未在 HCL 配置中引用的资源不会产生 diff。
- 文档更新：在 `resource_tc_teo_config_group_version.md` 中补充 `version_id` 字段说明，`website/docs/` 由 `make doc` 重新生成。
- 回滚：若需回退，移除 schema 中 `version_id` 出参声明与 Read 中的设置即可；state 中已有值不会丢失。

## Open Questions

- 无
