### Requirement: 拆分 DescribeFunctions 返回的拼接 name
系统 SHALL 在 `resourceTencentCloudTeoFunctionRead` 方法中，对 DescribeFunctions API 返回的 name 字段进行拆分处理，仅将原始 name 部分设置到 terraform state 中，以避免 terraform plan/apply 时因 name 字段不一致而误报资源变更。

#### Scenario: 标准拼接 name 拆分
- **WHEN** DescribeFunctions 返回的 name 为 `my-func-zone-2qtuhspy7cr6-1310708577`，其中原始 name 为 `my-func`，zone_id 为 `zone-2qtuhspy7cr6`，app_id 为 `1310708577`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `my-func`

#### Scenario: 原始 name 不包含连字符
- **WHEN** DescribeFunctions 返回的 name 为 `myfunc-zone-2qtuhspy7cr6-1310708577`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `myfunc`

#### Scenario: 原始 name 包含多个连字符
- **WHEN** DescribeFunctions 返回的 name 为 `my-test-func-v2-zone-2qtuhspy7cr6-1310708577`
- **THEN** 设置到 terraform state 中的 name 值 SHALL 为 `my-test-func-v2`

#### Scenario: 拼接 name 中未找到 -zone- 子串
- **WHEN** DescribeFunctions 返回的 name 不包含 `-zone-` 子串（异常情况或 API 格式变更）
- **THEN** 系统 SHALL 保留原始返回值设置到 terraform state，不做拆分处理

### Requirement: name 拆分辅助函数
系统 SHALL 提供 `parseTeoFunctionOriginalName` 辅助函数，用于从 DescribeFunctions 返回的拼接 name 中提取原始 name 部分。该函数通过查找 `-zone-` 子串的位置来确定原始 name 的边界。

#### Scenario: 辅助函数正确提取原始 name
- **WHEN** 调用 `parseTeoFunctionOriginalName("my-func-zone-2qtuhspy7cr6-1310708577")`
- **THEN** 返回 `my-func`

#### Scenario: 辅助函数处理不含 -zone- 的情况
- **WHEN** 调用 `parseTeoFunctionOriginalName("myfunc")`
- **THEN** 返回 `myfunc`

#### Scenario: 辅助函数处理空字符串
- **WHEN** 调用 `parseTeoFunctionOriginalName("")`
- **THEN** 返回空字符串 `""`
