## Context

当前 `tencentcloud_teo_function` 资源在 Read 方法中，直接将 DescribeFunctions API 返回的 name 字段设置到 terraform state 中。然而，云API在创建函数后，会将 name 字段拼接为 `原始name + "-" + zone_id + "-" + app_id` 的格式返回。这导致 terraform 在后续 plan/apply 时，对比用户配置中的原始 name 与 state 中的拼接 name，产生不一致，错误地提示需要修改 name。而实际上云API不支持修改 name 字段（name 已在 immutableArgs 中）。

当前代码中：
- `resourceTencentCloudTeoFunctionRead` 方法第185-187行直接将 `respData.Name` 设置到 state
- `resourceTencentCloudTeoFunctionUpdate` 方法第220行已将 name 列入 immutableArgs，但仍然会在 plan 阶段误报变更
- 现有测试用例中使用了拼接后的完整函数名 `"aaa-zone-2qtuhspy7cr6-1310708577"` 作为 name，而非原始 name

关键约束：
- zone_id 格式为 `zone-xxxxxxx`，其中 xxxxxxx 为数字加字母的组合
- app_id 为一串纯数字字符串
- 原始 name 中可能也包含 `-`，因此不能简单地按 `-` 分割
- Function 结构体中包含 ZoneId 字段，可用于辅助拆分

## Goals / Non-Goals

**Goals:**
- 在 Read 方法中，对 DescribeFunctions 返回的拼接 name 进行拆分，仅将原始 name 部分设置到 terraform state
- 确保 terraform plan/apply 不再因 name 字段不一致而误报变更
- 保持向后兼容：已存在的 state 中的 name 值（拼接后的完整函数名）在下次 Read 后会被自动修正为原始 name

**Non-Goals:**
- 不修改云API的 CreateFunction 接口行为
- 不修改资源的 schema 定义
- 不修改 Create/Update/Delete 方法的核心逻辑

## Decisions

### Decision 1: name 拆分算法 —— 基于 zone_id 参数的精确后缀匹配

**选择**: 从拼接后的 name 中，查找 `-zoneId`（即 `-` + zone_id 参数值，如 `-zone-2qtuhspy7cr6`）子串的位置，以该位置为分割点，左侧即为原始 name。

**理由**: 拼接后的 name 格式为 `original_name + "-" + zone_id + "-" + app_id`，其中 zone_id 在 Read 方法中是已知的（从 resource ID 中解析得到）。使用完整的 `-zoneId` 作为匹配后缀，无论原始 name 中是否包含 `-zone-` 子串，都能准确找到分割点。例如，当原始 name 为 `my-zone-func`、zone_id 为 `zone-2qtuhspy7cr6` 时，拼接后的 name 为 `my-zone-func-zone-2qtuhspy7cr6-1310708577`，查找 `-zone-2qtuhspy7cr6` 可正确定位，避免将 `my-zone-func` 错误截断为 `my`。

**算法**:
1. 获取 DescribeFunctions 返回的 name（如 `my-zone-func-zone-2qtuhspy7cr6-1310708577`）
2. 构造匹配后缀 `-` + zoneId（如 `-zone-2qtuhspy7cr6`）
3. 查找该后缀在 name 中的位置
4. 取该位置之前的部分作为原始 name（如 `my-zone-func`）
5. 如果未找到匹配后缀，则保留原始返回值（兼容性保护）

**替代方案**:
- 查找 `-zone-` 子串：简单但在原始 name 包含 `-zone-` 时（如 `my-zone-func`）会错误截断
- 使用正则表达式匹配：可行但不如字符串查找直观，且同样存在边界问题

### Decision 2: 拆分函数的放置位置

**选择**: 在 `resource_tc_teo_function_extension.go` 中添加 `parseTeoFunctionOriginalName` 辅助函数。

**理由**: 遵循项目惯例，辅助函数放在 `_extension.go` 文件中，主资源文件由代码生成器管理。

### Decision 3: 测试策略

**选择**: 在 `resource_tc_teo_function_test.go` 中添加 `parseTeoFunctionOriginalName` 函数的单元测试，覆盖各种边界情况。

**理由**: name 拆分是一个纯函数逻辑，适合用单元测试验证。同时更新现有验收测试中的 name 值。

## Risks / Trade-offs

- [Risk] 如果云API未来改变 name 的拼接格式，拆分逻辑可能失效 → Mitigation: 在未找到 `-zoneId` 后缀时保留原始返回值，避免破坏性影响
- [Risk] 已存在的 state 中 name 为拼接后的完整函数名，下次 Read 后会变为原始 name → Mitigation: 这是预期行为，state 会被自动修正，不会影响资源本身
