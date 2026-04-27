## Context

TEO（EdgeOne）是腾讯云的边缘安全加速平台。当前 Terraform Provider 中已存在大量 TEO 资源（如 `tencentcloud_teo_l7_acc_rule` 等），均位于 `tencentcloud/services/teo/` 目录下。本次需要新增 `tencentcloud_teo_security_js_injection_rule` 资源，用于管理 TEO 站点级别的 JavaScript 注入安全规则。

云 API 结构特点：
- **资源粒度**：一个 ZoneId 下可管理多条 JS 注入规则，Create/Modify 接口接受规则列表，Delete 接口按规则 ID 列表删除
- **无异步操作**：4 个 CRUD 接口均为同步接口，无需轮询等待
- **分页查询**：DescribeSecurityJSInjectionRule 支持 Limit（最大100）/ Offset 分页

参考资源：`tencentcloud_igtm_strategy`（通用资源模式，包含复合 ID、嵌套 TypeList、CRUD 全流程）

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_security_js_injection_rule` 资源的完整 CRUD 生命周期管理
- 支持通过 `zone_id` 标识站点，通过 `js_injection_rules` 管理规则列表
- 支持资源导入（Import）
- 编写 gomonkey mock 方式的单元测试，验证业务逻辑正确性
- 生成资源文档和 provider 注册

**Non-Goals:**
- 不实现 datasource（当前需求仅为资源）
- 不修改现有 TEO 资源的 schema 或行为
- 不实现异步轮询机制（API 为同步接口）

## Decisions

### 1. 资源 ID 设计：使用 `zone_id` 作为资源 ID

**决策**：资源 ID 使用 `zone_id`，因为该资源的所有 CRUD 操作都以 `zone_id` 为核心标识。一个 `zone_id` 对应一组 JS 注入规则。

**理由**：
- Create 接口的入参是 `ZoneId` + `JSInjectionRules`，返回 `JSInjectionRuleIds`
- Read 接口以 `ZoneId` 查询该站点下所有规则
- Modify 接口以 `ZoneId` + `JSInjectionRules` 修改
- Delete 接口以 `ZoneId` + `JSInjectionRuleIds` 删除
- 资源生命周期绑定在 `zone_id` 上，`zone_id` 变更意味着全新的资源

**替代方案**：使用 `zone_id + rule_id` 复合 ID —— 但这不符合 API 设计，因为 API 操作粒度是站点级别（一次创建/修改多条规则），而非单条规则级别。

### 2. Schema 设计：`js_injection_rules` 使用 TypeList 嵌套

**决策**：`js_injection_rules` 使用 `schema.TypeList` + `schema.Resource` 嵌套结构，包含 `rule_id`（Computed）、`name`、`priority`、`condition`、`inject_js` 字段。

**理由**：
- 云 API 的 `JSInjectionRules` 是一个数组，每项包含多个字段
- `rule_id` 由 API 在 Create 时返回，设为 Computed
- `name`、`condition` 为业务必需字段（Required）
- `priority`、`inject_js` 有默认值，设为 Optional + Computed

### 3. Create 流程：创建后调用 Read 刷新状态

**决策**：Create 流程为：构造请求 → 调用 CreateSecurityJSInjectionRule → 设置 ID → 调用 Read 刷新完整状态。

**理由**：Create 响应仅返回 `JSInjectionRuleIds`，不返回完整的规则详情。需要调用 Read 获取完整数据以填充 state。

### 4. Update 流程：直接调用 ModifySecurityJSInjectionRule

**决策**：当 `js_injection_rules` 发生变化时，直接调用 `ModifySecurityJSInjectionRule` 接口，传入最新的完整规则列表。

**理由**：Modify 接口接受完整的 `JSInjectionRules` 列表，采用全量覆盖模式。

### 5. Delete 流程：使用 Read 获取 rule_ids 后调用 Delete

**决策**：Delete 时先通过 Read 获取当前所有 `rule_id`，然后调用 `DeleteSecurityJSInjectionRule` 一次性删除所有规则。

**理由**：Delete 接口需要 `JSInjectionRuleIds` 参数，而用户在 TF 配置中不直接设置 `rule_id`（它是 Computed 字段），因此需要从 state 中读取。

### 6. 分页处理：DescribeSecurityJSInjectionRule 使用最大 Limit=100

**决策**：在 Read 函数中，设置 `Limit=100`（API 标注的最大值），循环分页获取所有规则。

**理由**：遵循项目规范，分页字段取 API 注释中的最大值。

## Risks / Trade-offs

- **[全量覆盖更新]** → Modify 接口采用全量覆盖模式，用户修改任意一条规则时需传入完整规则列表。如果并发修改，可能覆盖其他客户端的变更。缓解：Terraform 的 state 机制本身会锁定资源，避免并发问题。
- **[规则 ID 依赖]** → `rule_id` 是 API 返回的 Computed 字段，删除时依赖 state 中的 `rule_id`。如果 state 损坏或被外部修改，可能无法正确删除。缓解：Read 操作会刷新 state，确保 rule_id 与实际一致。
- **[规则数量限制]** → API 对单次请求的规则数量可能有限制，但当前 API 文档未明确说明。缓解：在实现中不做特殊处理，如遇限制再行调整。
