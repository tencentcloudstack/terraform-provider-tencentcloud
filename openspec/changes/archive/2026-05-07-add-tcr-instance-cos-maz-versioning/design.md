## Context

TCR (Tencent Container Registry) 的 CreateInstance 接口在最新的 SDK 中已经支持 `EnableCosMAZ` 和 `EnableCosVersioning` 两个字段。当前 `resource_tc_tcr_instance.go` 的实现中：
- Schema 定义缺少这两个字段
- Create 函数中未处理这两个字段的传递
- Read 函数中未从 API 响应中读取并 set 这两个字段

需要补充这两个字段的完整支持，包括 schema 定义、Create 逻辑和 Read 逻辑。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_tcr_instance` 资源 schema 中添加 `enable_cos_maz` 和 `enable_cos_versioning` 字段
- 在 Create 函数中支持从 schema 读取并传递给 API
- 在 Read 函数中从 API 响应读取并 set 回 state
- 更新资源文档
- 完全向后兼容，不影响现有配置

**Non-Goals:**
- 不支持这两个字段的 Update 操作（CreateInstance 接口的字段都是 ForceNew）
- 不涉及其他 TCR 资源（如 namespace、repository）

## Decisions

### 1. Schema 字段设计

**决策:** 两个字段都设置为 `Optional: true` + `Computed: true`

**理由:**
- `Optional`: 用户可以显式指定是否启用这些特性
- `Computed`: 如果用户不指定，API 会使用默认值（false），Read 时从 API 读取实际值
- 不设置 `ForceNew`: 虽然 CreateInstance 接口的字段理论上都是 ForceNew，但由于这是新增字段，现有资源的 state 中没有这些字段。如果设置 ForceNew，会导致所有现有资源在下次 apply 时被标记为需要重建。通过不设置 ForceNew 并使用 Computed，可以让现有资源平滑地在 Read 时补充这些字段。

**类型选择:** `TypeBool`
- SDK 中这两个字段都是 `*bool` 类型
- Terraform schema 中使用 `schema.TypeBool`

### 2. Create 逻辑实现

**决策:** 使用 `GetOkExists` 读取 bool 字段值

**理由:**
- Bool 字段的 false 值与未设置状态需要区分
- `d.GetOk("enable_cos_maz")` 在值为 false 时会返回 false, false，无法区分
- 使用 `d.GetOkExists("enable_cos_maz")` 可以准确判断用户是否显式设置了该字段
- 只有在用户显式设置时才传递给 API，否则让 API 使用默认值

**实现位置:** 在现有的 params 构建逻辑之后，调用 CreateTCRInstance 之前添加

### 3. Read 逻辑实现

**决策:** 在 Read 函数中从 `instance` 对象读取并 set

**理由:**
- SDK 的 `Registry` 结构体包含 `EnableCosMAZ` 和 `EnableCosVersioning` 字段
- `DescribeTCRInstanceById` 返回的 instance 对象包含这些字段
- 在现有的字段 set 逻辑之后添加，保持代码结构一致性

### 4. 文档更新

**决策:** 在 Arguments Reference 部分添加这两个字段的说明

**理由:**
- 新增字段必须在文档中说明用途和默认值
- 参考现有文档格式，提供清晰的字段描述

## Risks / Trade-offs

### 风险 1: 现有资源 state 升级

**风险:** 现有的 TCR 实例资源在下次 `terraform refresh` 或 `terraform plan` 时，state 会新增这两个字段

**缓解措施:** 
- 字段设置为 `Computed: true`，state 变化不会触发重建
- 这是正常的 state schema 演进，不会产生实际资源变更
- 用户可以选择性地在配置中显式指定这些字段

### 风险 2: SDK 版本依赖

**风险:** 如果 vendor 中的 tencentcloud-sdk-go 版本过旧，可能不包含这两个字段

**缓解措施:**
- 已确认 vendor 中的 SDK models.go (739-745 行, 8636-8639 行) 包含这两个字段
- 代码中使用指针类型 `*bool`，兼容 SDK 结构

### 权衡: 不设置 ForceNew

**权衡:** 理论上 CreateInstance 的所有参数都应该是 ForceNew，因为创建后无法修改

**决策理由:**
- 如果设置 ForceNew，所有现有资源会在 plan 时显示需要重建（因为 state 中缺少这些字段）
- 实际上这两个字段在创建后也无法通过 API 修改，但 Terraform Provider 层面不强制 ForceNew 不会造成实际问题
- 如果用户尝试修改这些字段，Terraform 会 apply 但 API 调用会失败（因为没有 Update 接口支持），这是可接受的行为
- 优先保证现有资源不被误标记为需要重建
