## Why

当前 `tencentcloud_teo_function` 资源的 `function_id` 参数为 Computed 字段，仅支持由系统自动生成。在某些场景下，用户需要在创建函数时指定特定的 function_id（例如与现有系统保持一致、实现特定的命名规范或满足业务需求），因此需要将 function_id 改为 Optional 字段，支持用户自定义输入。

**然而，在实施过程中发现：** 腾讯云 TEO API 的 `CreateFunction` 接口**不支持**用户自定义 `function_id` 参数，`function_id` 只能由系统自动生成。这是 API 层面的硬性限制，无法通过 Terraform Provider 层面的修改来绕过。

## What Changes

~~- 将 `tencentcloud_teo_function` 资源的 `function_id` 参数从 Computed 改为 Optional + Computed~~
~~- 在 CreateFunction API 调用时支持传入用户提供的 function_id 参数~~
- ~~保持向后兼容：当用户未提供 function_id 时，仍由系统自动生成~~
- ~~更新资源创建逻辑，在调用 CreateFunction API 时传入 FunctionId 参数（如果用户提供了值）~~

**实际变更：**
- 无代码变更（因 API 限制）
- 记录 API 限制，为未来可能的 API 改进提供参考
- 建议联系腾讯云 API 团队，请求添加 `CreateFunction` API 对自定义 `function_id` 的支持

## Capabilities

### New Capabilities
~~- `teo-function-id-input`: 支持在创建 tencentcloud_teo_function 资源时通过 function_id 参数指定函数 ID~~

**说明：** 由于 API 限制，无法实现此能力。

### Modified Capabilities
- 无（现有功能的行为要求未改变）

## Impact

~~- 影响文件：`tencentcloud/services/teo/resource_tc_teo_function.go`~~
~~- 影响的 API：`CreateFunction` API 调用增加 FunctionId 参数支持~~
~~- Schema 变更：`function_id` 参数从 Computed 改为 Optional + Computed~~
- ~~向后兼容：未提供 function_id 的现有配置将继续正常工作~~
~~- 测试影响：需要新增测试用例验证 function_id 参数输入功能~~
~~- 文档影响：需要更新资源文档说明 function_id 参数的可选性~~

**实际影响：**
- 无代码变更
- 记录发现的问题和限制
- 为未来可能的 API 改进提供参考
- 建议腾讯云 API 团队添加 `CreateFunction` API 对自定义 `function_id` 的支持
