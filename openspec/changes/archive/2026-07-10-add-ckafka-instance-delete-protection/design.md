## Context

`tencentcloud_ckafka_instance` 是 CKafka 实例的通用资源（RESOURCE_KIND_GENERAL），封装了实例的创建、读取、更新、删除全生命周期。

当前状态：
- 资源 schema 中未包含删除保护相关字段
- Create 流程在实例创建完成后，会调用一次 `ModifyInstanceAttributes` 设置 `msg_retention_time`、`config`、`dynamic_retention_config`、`rebalance_time`、`public_network`、`max_message_byte` 等属性
- Read 流程先调用 `DescribeInstancesDetail`（返回 `InstanceDetail`，不含 `DeleteProtectionEnable`），再调用 `DescribeInstanceAttributes`（返回 `InstanceAttributesResponse`，含 `DeleteProtectionEnable *int64`）回写属性
- Update 流程通过 `ModifyInstanceAttributes` 修改可变属性，并通过 `InstanceScalingDown` / `ModifyInstancePre` 处理规格变更

云 API 能力确认（vendor `ckafka/v20190819`）：
- `ModifyInstanceAttributesRequest.DeleteProtectionEnable`：`*int64`，注释"实例删除保护开关: 1 开启 0 关闭" ✅
- `InstanceAttributesResponse.DeleteProtectionEnable`：`*int64`，注释同上 ✅
- `InstanceDetail`（`DescribeInstancesDetail` 响应项）**不含** `DeleteProtectionEnable` ❌

约束：
- Terraform Provider 向后兼容：未配置 `delete_protection_enable` 的既有配置不能产生 plan diff
- 该字段为可选属性，未配置时应保留云端现有值（Computed 回填）

## Goals / Non-Goals

**Goals:**
- 支持通过 Terraform 声明式开启/关闭 CKafka 实例的删除保护
- 保持向后兼容：未配置时不触发 plan diff
- Create / Update / Read 三条路径均正确处理 `delete_protection_enable`
- 通过单元测试覆盖三条路径

**Non-Goals:**
- 不改变 Delete 操作行为（删除保护开启时，云 API 会在 Delete 阶段返回错误，由云端拦截，provider 不做特殊处理）
- 不修改 `DescribeInstancesDetail` 相关逻辑（该接口不返回该字段）
- 不调整 `immutableArgs` 列表（`delete_protection_enable` 是可变属性）
- 不新增独立资源

## Decisions

### Decision 1: 字段类型使用 TypeInt 而非 TypeBool

**选择**：`delete_protection_enable` 使用 `schema.TypeInt`，取值 `1`（开启）/ `0`（关闭），与云 API `*int64` 类型一一对应。

**备选**：使用 `schema.TypeBool`，在 provider 层做 bool↔int64 转换。

**理由**：
- 与资源内同类型开关字段风格一致（如 `elastic_bandwidth_switch`、`dynamic_retention_config.enable` 等均为 TypeInt）
- 避免引入额外的类型转换逻辑，降低出错概率
- 云 API 直接接受 int64，无需转换

### Decision 2: 字段属性为 Optional + Computed

**选择**：`Optional: true, Computed: true`。

**理由**：
- `Optional` 允许用户显式配置；`Computed` 保证未配置时由 Read 回填云端值，不产生 plan diff
- 这是 provider 处理"可选的、未配置时由云端回填"字段的标准模式

### Decision 3: Create 阶段在已有 ModifyInstanceAttributes 调用中顺带设置

**选择**：在 Create 末尾构造 `modifyRequest` 的现有逻辑中，新增对 `delete_protection_enable` 的判断与填充，复用同一次 `ModifyInstanceAttributes` 调用。

**备选**：Create 后单独发起一次 `ModifyInstanceAttributes` 调用。

**理由**：
- 现有 Create 流程已经在创建后调用一次 `ModifyInstanceAttributes` 设置其他属性，复用该调用减少 API 往返
- 使用 `d.GetOk("delete_protection_enable")` 判断用户是否显式配置，仅在配置时填充，避免覆盖云端默认值

### Decision 4: Update 阶段在已有 ModifyInstanceAttributes 分支中扩展

**选择**：在 `resourceTencentCloudCkafkaInstanceUpdate` 的 `modifyInstanceAttributesFlag` 逻辑块中，新增 `if d.HasChange("delete_protection_enable")` 分支，填充 `request.DeleteProtectionEnable` 并置 `modifyInstanceAttributesFlag = true`。

**理由**：
- 与现有 `instance_name`、`msg_retention_time`、`public_network`、`max_message_byte` 等可变字段的处理模式完全一致
- `delete_protection_enable` 不加入 `immutableArgs`，允许原地更新

### Decision 5: Read 阶段从 DescribeInstanceAttributes 响应回填

**选择**：在 Read 流程调用 `DescribeInstanceAttributes` 后，判断 `attr.DeleteProtectionEnable != nil`，然后 `d.Set("delete_protection_enable", attr.DeleteProtectionEnable)`。

**理由**：
- `DescribeInstanceAttributes` 的响应 `InstanceAttributesResponse` 已包含 `DeleteProtectionEnable` 字段
- 遵循"set 前判 nil"的规范，避免 nil 指针问题
- 不在 `DescribeInstancesDetail` 分支处理（该接口不返回此字段）

### Decision 6: 使用 d.GetOkExists 处理 0 值

**选择**：在 Create 中使用 `d.GetOkExists("delete_protection_enable")` 读取用户配置，确保用户显式设置 `0`（关闭）时也能正确传入。

**理由**：
- `d.GetOk` 在值为 `0` 时会返回 `ok=false`，会导致用户显式关闭删除保护时参数不被传入
- `d.GetOkExists` 能区分"未配置"与"配置为 0"，与现有 `max_message_byte` 字段处理方式一致

## Risks / Trade-offs

- **Risk**：用户在配置中设置 `delete_protection_enable = 1` 后执行 `terraform destroy`，云 API 会因删除保护开启而拒绝删除 → **Mitigation**：这是删除保护本身的预期行为，用户需先在配置中设为 `0` 并 apply 后再 destroy；provider 不绕过该保护。
- **Risk**：旧 state 中无 `delete_protection_enable` 字段，首次 refresh 时由 Computed 回填 → **Mitigation**：Optional+Computed 模式保证无 plan diff，向后兼容。
- **Trade-off**：Create 阶段复用同一次 `ModifyInstanceAttributes` 调用，若该调用失败会导致删除保护与其他属性一起未设置 → 可接受，Create 失败会整体报错并由用户重试。

## Migration Plan

- 纯加法变更（新增 Optional 字段），无 state 迁移需求
- 存量资源：升级后 `terraform plan` 对未在 HCL 配置 `delete_protection_enable` 的资源不产生 diff
- 文档更新：在 `resource_tc_ckafka_instance.md` 中补充字段说明和示例
- 回滚：移除 schema 字段与 CRUD 分支即可，state 中多余字段会被 terraform 忽略

## Open Questions

- 无
