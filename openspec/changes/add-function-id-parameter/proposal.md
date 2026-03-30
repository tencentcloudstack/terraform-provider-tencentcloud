## Why

在当前的 tencentcloud_teo_function 资源实现中，function_id 参数仅作为 Computed 字段存在，由服务端自动生成。这限制了用户在某些场景下的灵活性，例如需要复用已存在的 FunctionId 或需要特定 ID 进行资源关联的场景。通过接入 FunctionId 参数，允许用户在创建函数时指定 FunctionId，提升了资源的灵活性和可控性。

## What Changes

- 在 tencentcloud_teo_function 资源的 schema 中，将 `function_id` 参数从 Computed 改为 Computed+Optional
- 修改 `resourceTencentCloudTeoFunctionCreate` 函数，支持在创建时传入用户指定的 FunctionId
- 更新 CreateFunction API 调用逻辑，当用户指定 FunctionId 时将其包含在请求中
- 保持向后兼容性，当用户未指定 FunctionId 时仍由服务端自动生成

## Capabilities

### New Capabilities
- `teo-function-id-parameter`: 支持 tencentcloud_teo_function 资源在创建时指定 FunctionId 参数的能力

### Modified Capabilities
- 无

## Impact

- 受影响的资源文件：`tencentcloud/services/teo/resource_tc_teo_function.go`
- 受影响的测试文件：`tencentcloud/services/teo/resource_tc_teo_function_test.go`
- 受影响的文档文件：`website/docs/r/teo_function.html.markdown`
- 需要验证 CreateFunction API 是否支持传入 FunctionId 参数
- 保持向后兼容，不影响现有 Terraform 配置
