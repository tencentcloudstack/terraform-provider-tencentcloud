## Context

`tencentcloud_teo_l7_acc_rule_v2` 资源封装了 TEO 七层加速规则（RuleEngine）的 CRUD，覆盖 `CreateL7AccRules` / `ModifyL7AccRule` / `DescribeL7AccRules` 三个接口。现有资源已支持顶层 `zone_id` / `status` / `rule_name` / `description` / `branches`（含子分支 actions 全套参数与 sub_rules 递归结构）。

当前状态：
- `resource_tc_teo_l7_acc_rule_v2.go` 中 schema 已声明 `rule_id`（Computed）、`rule_priority`（Computed），Create/Update 仅处理 `status` / `rule_name` / `description` / `branches`。
- `resource_tc_teo_l7_acc_rule_extension.go` 中 `TencentTeoL7RuleBranchBasicInfo(depth)` 定义 actions 的 schema（含 cache/cache_key/.../advanced_origin_routing/shield/site_failover 等全套 Action 参数），以及 `sub_rules` 递归块；`resourceTencentCloudTeoL7AccRuleGetBranchs` / `resourceTencentCloudTeoL7AccRuleSetBranchs` 负责 actions/sub_rules 的扁平化与还原。
- SDK `teov20220901.RuleEngineItem` 已新增 `CustomVariables`（`[]*CustomVariable`）、`CustomVariableOperations`（`[]*CustomVariableOperation`）字段；`RuleEngineAction` 已新增 `CustomActionParameters`（`*CustomActionParameters`）字段。三个云 API 接口均支持这些字段的读写。

约束：
- Terraform Provider 向后兼容：现有 TF 配置（未配置新字段）必须无 plan drift。
- `custom_variables` / `custom_variable_operations` 属于顶层 `RuleEngineItem` 字段，与 `branches` 同级，不可嵌套进 branches。
- `custom_action_parameters` 属于 `RuleEngineAction`（actions 内）字段，与 cache_parameters 等同级。
- 所有新增字段均为 Optional，未配置时不向 API 传值。

## Goals / Non-Goals

**Goals:**
- 支持 Terraform 声明式管理规则级自定义变量（custom_variables）及其运算（custom_variable_operations）。
- 支持 actions 内通过 `CustomAction` 操作名调用定制配置（custom_action_parameters）。
- 保持向后兼容：未配置新字段时资源行为与之前完全一致。
- Create/Update/Read 三个接口的字段映射完整且符合云 API 的实际能力。

**Non-Goals:**
- 不修改 `rule_id` / `rule_priority` 的 schema 属性（保持 Computed 不变）。
- 不改变 branches / sub_rules 现有递归结构与现有 actions 参数的行为。
- 不新增独立的 custom_variable 资源类型（变量定义与规则生命周期绑定，复用现有资源）。
- 不引入新的 Timeouts 子块结构变化。

## Decisions

### Decision 1: 顶层新增 custom_variables / custom_variable_operations，而非嵌套进 branches

**选择**：在 `ResourceTencentCloudTeoL7AccRuleV2()` 的 schema map 顶层新增 `custom_variables` 与 `custom_variable_operations`（与 `branches` 同级），映射到 `RuleEngineItem.CustomVariables` / `CustomVariableOperations`。

**备选**：嵌套进 branches 结构。

**理由**：
- 云 API 的 `CustomVariables` / `CustomVariableOperations` 是 `RuleEngineItem` 的直属字段，与 `Branches` 平级，嵌套进 branches 会与 API 结构不一致导致映射错误。
- SDK 注释明确指出"当 Branches 为空时 CustomVariable / CustomVariableOperations 不允许填写"，说明三者是并列关系。

### Decision 2: custom_variables 与 custom_variable_operations 使用列表类型（非单元素块）

**选择**：`custom_variables` 为 `TypeList`（元素为 Resource，含 name/initial_value/description）；`custom_variable_operations` 为 `TypeList`（元素为 Resource，含 branches/description）。不设 MaxItems=1（custom_variables 支持多变量；custom_variable_operations 列表只支持填一项但 schema 不强制）。

**理由**：
- `CustomVariables` 是数组（`[]*CustomVariable`），支持多变量按数组顺序初始化。
- `CustomVariableOperations` 数组只支持填一项（SDK 注释"此列表当前只支持填写一项规则，多填无效"），但保持 TypeList 与 API 一致，服务端会校验数量。

### Decision 3: custom_variable_operations 的 branches 结构独立定义，不复用 TencentTeoL7RuleBranchBasicInfo

**选择**：为 `custom_variable_operations.branches` 定义独立的 schema（含 `condition` / `actions`（含 `name` / `set_parameters`（含 `name` / `value`）） / `sub_rules`（含 `branches`（递归同结构） / `description`）），映射到 SDK 的 `CustomVariableOperationRuleBranch` / `CustomVariableOperationRuleAction` / `SetParameters` / `CustomVariableOperationSubRule`。

**备选**：复用 `TencentTeoL7RuleBranchBasicInfo`。

**理由**：
- `CustomVariableOperationRuleBranch` 的 `Actions` 类型是 `[]*CustomVariableOperationRuleAction`（仅含 Name + SetParameters），与 `RuleBranch.Actions`（`[]*RuleEngineAction` 含全套参数）是不同的 SDK 类型，不可复用同一个 schema 与扁平化函数。
- `CustomVariableOperationRuleBranch.SubRules` 类型是 `[]*CustomVariableOperationSubRule`（含 Branches []*CustomVariableOperationRuleBranch + Description），与 `RuleBranch.SubRules`（`[]*RuleEngineSubRule`）也是不同类型。
- 独立定义避免类型混淆与不必要的 nil 字段。

### Decision 4: custom_action_parameters 嵌套进 TencentTeoL7RuleBranchBasicInfo 的 actions schema

**选择**：在 `TencentTeoL7RuleBranchBasicInfo` 的 actions Elem.Schema 中新增 `custom_action_parameters`（TypeList，MaxItems=1，Optional），含 `custom_actions` 列表（每元素含 `action_id` TypeString、`parameters` TypeList 含 `name`/`value_type`/`string_value`/`integer_value`/`float_value`/`boolean_value`/`string_array_value`/`integer_array_value`/`float_array_value`），映射到 `RuleEngineAction.CustomActionParameters` → `CustomActionParameters.CustomActions` → `CustomAction`。

**理由**：
- `CustomActionParameters` 是 `RuleEngineAction` 的直属字段，与 cache_parameters 等同级，属于 actions 的合法参数。
- `CustomAction.Parameters` 是 `[]*CustomActionParameter`，每个参数根据 `ValueType` 取不同类型的 Value 字段，schema 需覆盖全部 Value 类型字段以支持所有定制配置场景。
- `CustomActionParameters` 为单元素结构（指针），使用 MaxItems=1 与现有 advanced_origin_routing_parameters / shield_parameters 等保持一致。

### Decision 5: 扁平化/还原逻辑拆分为独立辅助函数

**选择**：
- 新增 `resourceTencentCloudTeoL7AccRuleGetCustomVariables` / `resourceTencentCloudTeoL7AccRuleSetCustomVariables` 处理 `custom_variables` 映射。
- 新增 `resourceTencentCloudTeoL7AccRuleGetCustomVariableOperations` / `resourceTencentCloudTeoL7AccRuleSetCustomVariableOperations` 处理 `custom_variable_operations` 映射（含递归 branches/sub_rules）。
- 在现有 `resourceTencentCloudTeoL7AccRuleGetBranchs` / `resourceTencentCloudTeoL7AccRuleSetBranchs` 中追加 `custom_action_parameters` 分支。

**理由**：
- custom_variables / custom_variable_operations 是顶层字段，逻辑独立，拆分函数保持单一职责。
- custom_action_parameters 属于 actions 内字段，复用现有 branches 遍历框架追加分支即可。

### Decision 6: Read 回填遵循 nil 安全规范

**选择**：在 Read 中对 `CustomVariables` / `CustomVariableOperations` / `CustomActionParameters` 回填前先判断 nil，nil 时不调用 setXX；在 SetBranchs 中对每个字段判断 nil 后再写入 map，与现有实现风格一致。

### Decision 7: Update 的 d.HasChange 检测纳入新字段

**选择**：在 `ResourceTencentCloudTeoL7AccRuleV2Update` 的 `if d.HasChange(...)` 条件中加入 `custom_variables` / `custom_variable_operations`，并在 rule 构造中追加这两个字段的映射。`custom_action_parameters` 属于 branches 子字段，已被 `d.HasChange("branches")` 覆盖。

## SDK Struct 字段定义参考（从 vendor 摘录）

以下为本次新增参数涉及的云 API struct 字段定义，供实施阶段直接参考，无需回读 vendor 大文件。

### RuleEngineItem（models.go:25625）
```go
type RuleEngineItem struct {
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`
	Description []*string `json:"Description,omitnil,omitempty" name:"Description"`
	// 规则级自定义变量列表。CustomVariable.Name 需要使用 user.rule. 作为前缀。
	// 注意：此字段可能返回 null，表示取不到有效值。
	CustomVariables []*CustomVariable `json:"CustomVariables,omitnil,omitempty" name:"CustomVariables"`
	// 规则级自定义变量运算详情。此列表当前只支持填写一项规则，多填无效。
	// 注意：此字段可能返回 null，表示取不到有效值。
	CustomVariableOperations []*CustomVariableOperation `json:"CustomVariableOperations,omitnil,omitempty" name:"CustomVariableOperations"`
	// 子规则分支。此列表当前只支持填写一项规则，多填无效。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Branches []*RuleBranch `json:"Branches,omitnil,omitempty" name:"Branches"`
	// 规则优先级。仅作为出参使用。
	RulePriority *int64 `json:"RulePriority,omitnil,omitempty" name:"RulePriority"`
}
```

### CustomVariable（models.go:6211）
```go
type CustomVariable struct {
	// 变量名称。需填写完整前缀：user.zone.* 站点级，user.rule.* 规则级。前缀后仅支持大小写字母、数字和下划线。长度≤50。创建成功后名称不可修改。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 变量初始值。支持常量字符串、变量以及公式。长度≤255。
	InitialValue *string `json:"InitialValue,omitnil,omitempty" name:"InitialValue"`
	// 变量描述。长度≤60。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`
}
```

### CustomVariableOperation（models.go:6222）
```go
type CustomVariableOperation struct {
	// 子规则分支。此列表当前只支持填写一项规则，多填无效。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Branches []*CustomVariableOperationRuleBranch `json:"Branches,omitnil,omitempty" name:"Branches"`
	// 规则注释。可以填写多个注释。
	Description []*string `json:"Description,omitnil,omitempty" name:"Description"`
}
```

### CustomVariableOperationRuleBranch（models.go:6240）
```go
type CustomVariableOperationRuleBranch struct {
	// 匹配条件。
	Condition *string `json:"Condition,omitnil,omitempty" name:"Condition"`
	// 操作。注意：Actions 和 SubRules 不可同时为空。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Actions []*CustomVariableOperationRuleAction `json:"Actions,omitnil,omitempty" name:"Actions"`
	// 子规则列表。当前只支持填写一层 SubRules。
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubRules []*CustomVariableOperationSubRule `json:"SubRules,omitnil,omitempty" name:"SubRules"`
}
```

### CustomVariableOperationRuleAction（models.go:6231）
```go
type CustomVariableOperationRuleAction struct {
	// 操作名称。当前仅支持填写 Set（自定义变量设置）。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 自定义变量设置参数。当 Name=Set 时必填。若存在多条运算，按数组顺序依次执行。
	// 注意：此字段可能返回 null，表示取不到有效值。
	SetParameters *SetParameters `json:"SetParameters,omitnil,omitempty" name:"SetParameters"`
}
```

### SetParameters（models.go:26037）
```go
type SetParameters struct {
	// 自定义变量名称。自定义变量必须先被定义才可进行运算。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 自定义变量赋值。支持常量字符串、变量以及公式，不支持中文。长度≤1000。
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`
}
```

### CustomVariableOperationSubRule（models.go:6253）
```go
type CustomVariableOperationSubRule struct {
	// 子规则分支
	// 注意：此字段可能返回 null，表示取不到有效值。
	Branches []*CustomVariableOperationRuleBranch `json:"Branches,omitnil,omitempty" name:"Branches"`
	// 规则注释。
	Description []*string `json:"Description,omitnil,omitempty" name:"Description"`
}
```

### RuleEngineAction.CustomActionParameters（models.go:25554）
```go
// 定制配置操作参数，当 Name 取值为 CustomAction 时，该参数必填。
CustomActionParameters *CustomActionParameters `json:"CustomActionParameters,omitnil,omitempty" name:"CustomActionParameters"`
```

### CustomActionParameters（models.go:6090）
```go
type CustomActionParameters struct {
	// 需要配置的定制配置列表。
	CustomActions []*CustomAction `json:"CustomActions,omitnil,omitempty" name:"CustomActions"`
}
```

### CustomAction（models.go:6053）
```go
type CustomAction struct {
	// 定制配置的配置项 Id。可通过 DescribeAvailableCustomActionsForRuleEngine 获取。
	ActionId *string `json:"ActionId,omitnil,omitempty" name:"ActionId"`
	// 该定制配置项下各参数字段的取值。可通过 DescribeAvailableCustomActionsForRuleEngine 获取。
	Parameters []*CustomActionParameter `json:"Parameters,omitnil,omitempty" name:"Parameters"`
}
```

### CustomActionParameter（models.go:6061）
```go
type CustomActionParameter struct {
	// 定制配置项下各参数字段名称。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 参数字段值的类型。枚举：String/Integer/Float/Boolean/ArrayOfString/ArrayOfInteger/ArrayOfFloat。
	ValueType *string `json:"ValueType,omitnil,omitempty" name:"ValueType"`
	// 字符串类型参数值。当 ValueType=String 时必填。
	StringValue *string `json:"StringValue,omitnil,omitempty" name:"StringValue"`
	// 整型类型参数值。当 ValueType=Integer 时必填。
	IntegerValue *int64 `json:"IntegerValue,omitnil,omitempty" name:"IntegerValue"`
	// 浮点类型参数值。当 ValueType=Float 时必填。
	FloatValue *float64 `json:"FloatValue,omitnil,omitempty" name:"FloatValue"`
	// 布尔类型参数值。当 ValueType=Boolean 时必填。
	BooleanValue *bool `json:"BooleanValue,omitnil,omitempty" name:"BooleanValue"`
	// 字符串数组类型参数值。当 ValueType=ArrayOfString 时必填。
	StringArrayValue []*string `json:"StringArrayValue,omitnil,omitempty" name:"StringArrayValue"`
	// 整型数组类型参数值。当 ValueType=ArrayOfInteger 时必填。
	IntegerArrayValue []*int64 `json:"IntegerArrayValue,omitnil,omitempty" name:"IntegerArrayValue"`
	// 浮点数组类型参数值。当 ValueType=ArrayOfFloat 时必填。
	FloatArrayValue []*float64 `json:"FloatArrayValue,omitnil,omitempty" name:"FloatArrayValue"`
}
```

### 接口请求/响应 struct 参考

- `CreateL7AccRulesRequest`（models.go:4145）：`ZoneId *string`、`Rules []*RuleEngineItem`
- `ModifyL7AccRuleRequest`（models.go:21411）：`ZoneId *string`、`Rule *RuleEngineItem`
- `DescribeL7AccRulesResponseParams`（models.go:11966）：`TotalCount *int64`、`Rules []*RuleEngineItem`

三个接口共用 `RuleEngineItem`，新增字段在 Create/Update 作为入参、在 Describe 作为出参，CRUD 映射对称。

## 现有 schema 与新参数的结构关系说明

目标资源 `ResourceTencentCloudTeoL7AccRuleV2()` schema 结构（新增部分以 `// 新增` 标注）：

```
zone_id (Required, ForceNew)
status (Optional)
rule_name (Optional)
description (Optional, list)
custom_variables (Optional, list)            // 新增 - 顶层
  └ name (Optional, string)
  └ initial_value (Optional, string)
  └ description (Optional, string)
custom_variable_operations (Optional, list)  // 新增 - 顶层
  └ branches (Optional, list)
  │   └ condition (Optional, string)
  │   └ actions (Optional, list)
  │   │   └ name (Optional, string)
  │   │   └ set_parameters (Optional, list, MaxItems=1)
  │   │       └ name (Optional, string)
  │   │       └ value (Optional, string)
  │   └ sub_rules (Optional, list)           // 递归一层
  │       └ branches (同上 branches 结构)
  │       └ description (Optional, list)
  └ description (Optional, list)
branches (Optional, list)
  └ condition (Optional, string)
  └ actions (Optional, list)
  │   └ name (Required, string)
  │   └ ... 全套现有参数 ...
  │   └ custom_action_parameters (Optional, list, MaxItems=1)  // 新增 - actions 内
  │       └ custom_actions (Optional, list)
  │           └ action_id (Optional, string)
  │           └ parameters (Optional, list)
  │               └ name (Optional, string)
  │               └ value_type (Optional, string)
  │               └ string_value (Optional, string)
  │               └ integer_value (Optional, int)
  │               └ float_value (Optional, float)
  │               └ boolean_value (Optional, bool)
  │               └ string_array_value (Optional, list)
  │               └ integer_array_value (Optional, list)
  │               └ float_array_value (Optional, list)
  └ sub_rules (Optional, list, 递归)
rule_id (Computed)
rule_priority (Computed)
```

## Risks / Trade-offs

- [CustomVariableOperation 的 branches/sub_rules 递归深度] → SDK 注释指出"当前只支持填写一层 SubRules"，schema 中 sub_rules 内不再嵌套 sub_rules（与 SDK 一致），避免无限递归。
- [CustomActionParameter 多类型 Value 字段] → schema 暴露所有 Value 类型字段（string_value/integer_value/...），由用户根据 value_type 填写对应字段，未填写的字段不传值；服务端按 value_type 取值校验，实现简单且覆盖全部场景。
- [custom_variables 名称前缀校验] → schema 层面不使用 ValidateFunc 强制 `user.rule.` 前缀（与现有字段风格一致，服务端返回错误），仅在 Description 中说明约束。
- [向后兼容] → 所有新增字段为 Optional，未配置时不进入 request，现有 TF 配置无 plan drift。