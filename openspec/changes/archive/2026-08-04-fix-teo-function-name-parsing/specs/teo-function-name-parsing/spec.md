## ADDED Requirements

### Requirement: 拆分 DescribeFunctions 返回的拼接 name
系统 SHALL 在 `resourceTencentCloudTeoFunctionRead` 方法中，对 DescribeFunctions API 返回的 name 字段进行拆分处理，仅将原始 name 部分设置到 terraform state 中，以避免 terraform plan/apply 时因 name 字段不一致而误报资源变更。拆分算法使用已知的 zone_id 参数（从 resource ID 解析得到）构造匹配后缀 `-zoneId`，在拼接 name 中精确定位分割点。

#### Scenario: 标准拼接 name 拆分
- **WHEN** DescribeFunctions 返回的 name 为 `my-func-zone-2qtuhspy7cr6-1310708577`，zone_id 为 `zone-2qtuhspy7cr6`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `my-func`

#### Scenario: 原始 name 不包含连字符
- **WHEN** DescribeFunctions 返回的 name 为 `myfunc-zone-2qtuhspy7cr6-1310708577`，zone_id 为 `zone-2qtuhspy7cr6`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `myfunc`

#### Scenario: 原始 name 包含多个连字符
- **WHEN** DescribeFunctions 返回的 name 为 `my-test-func-v2-zone-2qtuhspy7cr6-1310708577`，zone_id 为 `zone-2qtuhspy7cr6`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `my-test-func-v2`

#### Scenario: 原始 name 包含 -zone- 子串
- **WHEN** DescribeFunctions 返回的 name 为 `my-zone-func-zone-2qtuhspy7cr6-1310708577`，zone_id 为 `zone-2qtuhspy7cr6`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `my-zone-func`（而非错误截断为 `my`）

#### Scenario: 拼接 name 中未找到 -zoneId 后缀
- **WHEN** DescribeFunctions 返回的 name 不包含 `-zoneId` 后缀（异常情况或 API 格式变更）
- **THEN** 系统 SHALL 保留原始返回值设置到 terraform state，不做拆分处理

### Requirement: name 拆分辅助函数
系统 SHALL 提供 `ParseTeoFunctionOriginalName` 辅助函数，接收 `name` 和 `zoneId` 两个字符串参数，用于从 DescribeFunctions 返回的拼接 name 中提取原始 name 部分。该函数通过查找 `-zoneId` 后缀（即 `-` + zone_id 参数值）的位置来确定原始 name 的边界。

#### Scenario: 辅助函数正确提取原始 name
- **WHEN** 调用 `ParseTeoFunctionOriginalName("my-zone-func-zone-2qtuhspy7cr6-1310708577", "zone-2qtuhspy7cr6")`
- **THEN** 返回 `my-zone-func`

#### Scenario: 辅助函数处理不含 -zoneId 后缀的情况
- **WHEN** 调用 `ParseTeoFunctionOriginalName("myfunc", "zone-2qtuhspy7cr6")`
- **THEN** 返回 `myfunc`

#### Scenario: 辅助函数处理空字符串
- **WHEN** 调用 `ParseTeoFunctionOriginalName("", "zone-2qtuhspy7cr6")`
- **THEN** 返回空字符串 `""`
