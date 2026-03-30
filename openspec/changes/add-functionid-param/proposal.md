## Why

用户需要在创建 tencentcloud_teo_function 资源时能够指定自定义的 FunctionId，以满足特定的业务需求和管理规范。当前 function_id 字段为 Computed，由系统自动生成，无法支持用户自定义输入。

## What Changes

- 将 tencentcloud_teo_function 资源的 function_id 字段从 Computed 改为 Optional
- 修改 resourceTencentCloudTeoFunctionCreate 函数，支持传入用户指定的 FunctionId 参数
- 确保向后兼容性：如果用户未指定 FunctionId，则保持原有行为（由 API 自动生成）

## Capabilities

### New Capabilities
- `teo-function-functionid-input`: 支持在创建 tencentcloud_teo_function 资源时指定自定义的 FunctionId

### Modified Capabilities

(无)

## Impact

- 修改文件: `tencentcloud/services/teo/resource_tc_teo_function.go`
- 修改文件: `tencentcloud/services/teo/resource_tc_teo_function_test.go` (更新测试用例)
- 修改文件: `tencentcloud/services/teo/resource_tc_teo_function.md` (更新文档示例)
- 受影响的 API: CreateFunction (支持传入 FunctionId 参数)
- 保持向后兼容，不破坏现有 Terraform 配置和状态
