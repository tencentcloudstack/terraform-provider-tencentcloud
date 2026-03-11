# PostgreSQL 实例资源 - auto_renew_flag 字段修复

## MODIFIED Requirements

### Requirement: auto_renew_flag 字段状态读取

`tencentcloud_postgresql_instance` 资源 SHALL 在 Read 操作中正确读取并设置 `auto_renew_flag` 字段的值。

**背景**: 该字段当前在 schema 中定义但在 Read 函数中未设置,导致状态不一致。

#### Scenario: Schema 标记为 Computed

- **假设** 有一个 `tencentcloud_postgresql_instance` 资源定义
- **当** 检查 `auto_renew_flag` 字段的 schema 定义时
- **那么** 该字段 SHALL:
  - 设置 `Type: schema.TypeInt`
  - 设置 `Optional: true` (允许用户配置)
  - 设置 `Computed: true` (允许从 API 读取)
  - 包含准确的 description

**理由**: `Computed: true` 表示该字段可以从 API 读取,支持 refresh 和 import 操作。

#### Scenario: Read 函数设置 auto_renew_flag 值

- **假设** 调用 `DescribePostgresqlInstanceById` 返回了实例信息
- **并且** 实例的 `AutoRenew` 字段不为 nil
- **当** 执行 `resourceTencentCloudPostgresqlInstanceRead` 函数时
- **那么** 系统 SHALL:
  - 检查 `instance.AutoRenew` 是否为 nil
  - 如果不为 nil,将其值转换为 `int` 类型
  - 调用 `d.Set("auto_renew_flag", int(*instance.AutoRenew))` 设置到 state 中

**位置**: 在设置 `charge_type` 之后,设置 `security_groups` 之前

#### Scenario: 处理 nil 值

- **假设** 调用 `DescribePostgresqlInstanceById` 返回了实例信息
- **并且** 实例的 `AutoRenew` 字段为 nil (例如按量计费实例)
- **当** 执行 `resourceTencentCloudPostgresqlInstanceRead` 函数时
- **那么** 系统 SHALL:
  - 跳过 `d.Set("auto_renew_flag", ...)` 调用
  - 不设置该字段值(保持现有行为或默认值)

**理由**: 避免空指针异常,对于不支持自动续费的实例类型保持向后兼容。

#### Scenario: 状态一致性

- **假设** 用户创建了一个预付费实例并设置 `auto_renew_flag = 1`
- **当** 执行 `terraform refresh` 时
- **那么** Terraform state 中的 `auto_renew_flag` 值 SHALL 与云上实际的自动续费状态一致

#### Scenario: 配置漂移检测

- **假设** 用户通过 Terraform 创建了实例,设置 `auto_renew_flag = 0`
- **并且** 用户随后通过控制台将自动续费改为启用 (值为 1)
- **当** 执行 `terraform plan` 时
- **那么** Terraform SHALL 检测到配置漂移并显示变更

#### Scenario: 导入现有实例

- **假设** 云上存在一个预付费实例,自动续费已启用
- **当** 用户通过 `terraform import` 导入该实例时
- **那么** Terraform state 中 SHALL 正确包含 `auto_renew_flag = 1`

**理由**: 支持将现有资源纳入 Terraform 管理。

### Requirement: 与相关资源保持一致

该实现 SHALL 与项目中其他 PostgreSQL 资源的 `auto_renew_flag` 处理方式保持一致。

#### Scenario: 与只读实例实现对齐

- **假设** 检查 `resource_tc_postgresql_readonly_instance.go` 的实现
- **当** 比较两个资源的 `auto_renew_flag` 处理逻辑时
- **那么** 两者的实现 SHALL 保持一致:
  - Schema 定义相同 (Optional + Computed)
  - Read 函数中的设置逻辑相同
  - 类型转换方式相同

**参考**: `resource_tc_postgresql_readonly_instance.go:407`
```go
_ = d.Set("auto_renew_flag", instance.AutoRenew)
```

#### Scenario: 与数据源保持一致

- **假设** 检查 `data_source_tc_postgresql_instances.go` 的实现
- **当** 比较数据源和资源的 `auto_renew_flag` 字段处理时
- **那么** 两者应使用相同的字段名和值格式

**参考**: `data_source_tc_postgresql_instances.go:546`
```go
listItem["auto_renew_flag"] = v.AutoRenew
```
