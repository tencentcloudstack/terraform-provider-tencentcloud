# tke-addon-config-json-diff-fix Specification

## Purpose

Fix false-positive diff detection in `tencentcloud_kubernetes_addon_config` resource when the TKE API returns JSON responses with elements in a different order than user input. This ensures Terraform only shows diffs when the JSON content actually changes, not when only the element ordering differs.

此规范修复了 `tencentcloud_kubernetes_addon_config` 资源中因 TKE API 返回的 JSON 元素顺序与用户输入不同而导致的误报 diff 问题。确保 Terraform 仅在 JSON 内容真正改变时显示 diff，而不是仅元素顺序不同时。

## Requirements

### Requirement: REQ-TKE-ADDON-001 - raw_values 字段必须忽略 JSON 元素顺序差异

The `raw_values` field in `tencentcloud_kubernetes_addon_config` resource MUST perform semantic JSON comparison instead of string comparison to suppress diffs caused by element ordering differences.

`tencentcloud_kubernetes_addon_config` 资源的 `raw_values` 字段必须执行语义化的 JSON 比较而非字符串比较，以忽略元素顺序差异导致的 diff。

**理由**:
- TKE API 返回的 JSON 字符串可能与用户输入的顺序不同
- JSON 对象的键顺序在语义上是无关紧要的
- 字符串比较会将顺序差异误判为内容变更
- 导致用户每次 `terraform plan` 都看到不必要的 diff

**影响范围**:
- `tencentcloud_kubernetes_addon_config` 资源的 `raw_values` 字段
- 资源的 diff 检测逻辑

#### Scenario: 用户配置 JSON 与 API 返回 JSON 顺序不同但内容相同

**Given**: 用户在 Terraform 配置中指定 `raw_values = '{"replicas": 2, "image": "nginx:latest", "port": 80}'`  
**And**: TKE API 返回 `raw_values = '{"image": "nginx:latest", "port": 80, "replicas": 2}'`  
**When**: 用户执行 `terraform plan`  
**Then**: 
- Terraform 不应显示 `raw_values` 字段的 diff
- Plan 输出: "No changes. Your infrastructure matches the configuration."
- 不会提示用户执行 `terraform apply`

**验收标准**:
- ✅ 相同 JSON 内容不同顺序不产生 diff
- ✅ 真正的内容变更仍然产生 diff
- ✅ 接受测试通过

---

#### Scenario: raw_values 内容确实发生变化

**Given**: 当前 state 中 `raw_values = '{"replicas": 2, "image": "nginx:latest"}'`  
**And**: 用户修改配置为 `raw_values = '{"replicas": 3, "image": "nginx:latest"}'`  
**When**: 用户执行 `terraform plan`  
**Then**: 
- Terraform 必须显示 `raw_values` 字段的 diff
- Diff 应明确显示 `replicas: 2 → 3`
- 提示用户执行 `terraform apply` 以应用变更

**验收标准**:
- ✅ 真实的内容变更被正确识别
- ✅ Diff 输出清晰准确
- ✅ 更新操作正常执行

---

#### Scenario: raw_values 为空或不存在

**Given**: `raw_values` 字段为空字符串或未设置  
**When**: 与另一个空值或非空值比较  
**Then**: 
- 两个空值比较: 不产生 diff
- 空值与非空值比较: 产生 diff
- 行为与标准字段一致

**验收标准**:
- ✅ 空值处理正确
- ✅ 不会因为空值导致错误

---

#### Scenario: raw_values 包含无效 JSON

**Given**: `raw_values` 包含非 JSON 格式的字符串  
**When**: 执行 diff 检测  
**Then**: 
- 回退到字符串比较
- 记录警告日志但不中断操作
- 确保向后兼容性

**验收标准**:
- ✅ 不会因为无效 JSON 而崩溃
- ✅ 日志中记录警告信息
- ✅ 功能降级为字符串比较

---

### Requirement: REQ-TKE-ADDON-002 - 实现必须使用 Terraform 标准模式

The diff suppression implementation MUST use Terraform Plugin SDK v2's `DiffSuppressFunc` mechanism.

Diff 抑制实现必须使用 Terraform Plugin SDK v2 的 `DiffSuppressFunc` 机制。

**理由**:
- `DiffSuppressFunc` 是 Terraform 官方推荐的标准模式
- 非侵入式，不影响 CRUD 操作
- 在 plan/diff 阶段自动调用
- 易于测试和维护

**实现要求**:
- 在 schema 定义中添加 `DiffSuppressFunc` 字段
- 函数签名: `func(k, old, new string, d *schema.ResourceData) bool`
- 返回 `true` 表示抑制 diff，`false` 表示显示 diff

#### Scenario: DiffSuppressFunc 正确集成到 Schema

**Given**: `raw_values` 字段的 schema 定义  
**When**: 添加 `DiffSuppressFunc: suppressJSONOrderDiff`  
**Then**: 
- Schema 编译无错误
- Provider 初始化成功
- Diff 检测时自动调用该函数

**验收标准**:
- ✅ Schema 定义正确
- ✅ 函数签名匹配 Terraform 要求
- ✅ 集成测试通过

---

### Requirement: REQ-TKE-ADDON-003 - JSON 比较必须处理嵌套结构

The JSON comparison MUST correctly handle nested objects and preserve array ordering semantics.

JSON 比较必须正确处理嵌套对象，并保留数组顺序语义。

**理由**:
- Addon 配置可能包含复杂的嵌套结构
- 对象键顺序无关紧要，但数组元素顺序有意义
- 必须使用深度相等比较

**实现要求**:
- 使用 `json.Unmarshal` 解析 JSON
- 使用 `reflect.DeepEqual` 进行深度比较
- 不需要手动处理嵌套，标准库已处理

#### Scenario: 嵌套对象的键顺序不同

**Given**: `raw_values = '{"outer": {"a": 1, "b": 2}}'`  
**And**: API 返回 `'{"outer": {"b": 2, "a": 1}}'`  
**When**: 执行 diff 检测  
**Then**: 
- 不产生 diff
- 深度比较识别两者相同

**验收标准**:
- ✅ 嵌套对象正确处理
- ✅ 单元测试覆盖此场景

---

#### Scenario: 数组元素顺序不同

**Given**: `raw_values = '{"list": [1, 2, 3]}'`  
**And**: API 返回 `'{"list": [3, 2, 1]}'`  
**When**: 执行 diff 检测  
**Then**: 
- 必须产生 diff
- 数组顺序是有意义的，不应忽略

**验收标准**:
- ✅ 数组顺序差异被识别
- ✅ 单元测试验证此行为

---

### Requirement: REQ-TKE-ADDON-004 - 实现必须保持向后兼容

The implementation MUST be fully backward compatible with existing configurations and state.

实现必须与现有配置和状态完全向后兼容。

**理由**:
- 不能破坏现有的 Terraform 配置
- 不能要求用户迁移 state
- 必须是纯优化，无破坏性变更

**兼容性要求**:
- 所有现有配置继续工作
- 不需要 state 迁移
- API 调用逻辑不变
- 只影响 diff 检测

#### Scenario: 现有配置无需修改

**Given**: 用户有现有的 `tencentcloud_kubernetes_addon_config` 资源  
**When**: 升级到包含此修复的 Provider 版本  
**Then**: 
- 现有配置继续正常工作
- 不需要修改 HCL 配置
- `terraform plan` 正常运行
- 行为改进但不破坏兼容性

**验收标准**:
- ✅ 现有配置无需修改
- ✅ 现有 state 无需迁移
- ✅ 升级平滑无影响

---

### Requirement: REQ-TKE-ADDON-005 - 代码必须遵循项目规范

The code MUST follow the project's coding conventions and organization.

代码必须遵循项目的编码规范和组织结构。

**项目规范**:
- 使用 `gofmt` 格式化代码
- 辅助函数放在文件末尾（所有 CRUD 函数之后）
- 添加清晰的注释
- 导入标准库按字母顺序排列

**代码位置**:
- 文件: `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
- 函数位置: 文件末尾，`resourceTencentCloudKubernetesAddonConfigDelete` 之后

#### Scenario: 代码组织符合项目规范

**Given**: 新增的 `suppressJSONOrderDiff` 函数  
**When**: 检查代码位置和格式  
**Then**: 
- 函数位于文件末尾
- 代码已使用 `go fmt` 格式化
- 导入按字母顺序排列
- 注释清晰描述功能

**验收标准**:
- ✅ 代码格式化检查通过
- ✅ golangci-lint 通过
- ✅ 代码审查批准

---

## Implementation Details

### Technical Design

**Function Implementation:**
```go
// suppressJSONOrderDiff compares two JSON strings and ignores ordering differences
func suppressJSONOrderDiff(k, old, new string, d *schema.ResourceData) bool {
    // Handle empty strings
    if old == "" && new == "" {
        return true
    }
    if old == "" || new == "" {
        return false
    }

    // Parse both JSON strings
    var oldJSON, newJSON interface{}
    if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
        log.Printf("[WARN] Failed to unmarshal old value as JSON: %v", err)
        return old == new
    }
    if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
        log.Printf("[WARN] Failed to unmarshal new value as JSON: %v", err)
        return old == new
    }

    // Deep equality comparison
    return reflect.DeepEqual(oldJSON, newJSON)
}
```

**Schema Update:**
```go
"raw_values": {
    Type:             schema.TypeString,
    Optional:         true,
    Computed:         true,
    Description:      "Params of addon, base64 encoded json format.",
    DiffSuppressFunc: suppressJSONOrderDiff,
},
```

### Test Coverage

**Unit Tests:**
- Empty string handling
- Invalid JSON fallback
- Same JSON different order
- Different JSON content
- Nested objects
- Array ordering

**Acceptance Tests:**
- Create addon with JSON values
- Verify no spurious diff
- Verify real changes detected
- Update and delete operations

## Status

- **Status**: ✅ Implemented
- **Implementation Date**: 2026-03-24
- **Version**: Pending release
- **Related PR**: TBD

## References

- Original Change: `openspec/changes/fix-tke-addon-config-raw-values-json-diff/`
- Resource File: `tencentcloud/services/tke/resource_tc_kubernetes_addon_config.go`
- Terraform Plugin SDK: https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas
