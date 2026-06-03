## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_function_component_binding` resource SHALL define the following schema:
- `zone_id` (Required, ForceNew, String): 站点 ID
- `function_id` (Required, ForceNew, String): 函数 ID
- `function_component_bindings` (Required, List): 函数组件绑定列表，每个元素包含：
  - `type` (Required, String): 绑定的组件类型，取值为 `kv_namespace`
  - `variable_name` (Required, String): 用于绑定的变量名
  - `kv_namespace_parameters` (Optional, List, MaxItems 1): KV 命名空间配置参数，当 type 为 kv_namespace 时必填，包含：
    - `zone_id` (Required, String): KV 命名空间所属的站点 ID
    - `namespace` (Required, String): KV 命名空间名称

资源 ID 格式为 `zone_id#function_id`（使用 tccommon.FILED_SP 分隔）。

#### Scenario: Resource schema is correctly defined
- **WHEN** user defines a `tencentcloud_teo_function_component_binding` resource in HCL
- **THEN** the resource SHALL accept `zone_id`, `function_id`, and `function_component_bindings` as parameters

#### Scenario: Import with composite ID
- **WHEN** user imports the resource using `zone_id#function_id` format
- **THEN** the resource SHALL correctly parse the composite ID and read the configuration

### Requirement: Read configuration via DescribeFunctionComponentBindings
The resource Read method SHALL call `DescribeFunctionComponentBindings` API to query the current binding list. It SHALL handle pagination by setting Limit to 1000 (maximum value) and iterating with Offset until all bindings are retrieved. The Read method SHALL set `function_component_bindings` in state from the response.

#### Scenario: Read bindings successfully
- **WHEN** the Read method is called
- **THEN** it SHALL call `DescribeFunctionComponentBindings` with `ZoneId` and `FunctionId` from the resource ID, handle pagination, and set the binding list in Terraform state

#### Scenario: Read with empty bindings
- **WHEN** the API returns an empty binding list
- **THEN** the Read method SHALL set `function_component_bindings` to an empty list in state

### Requirement: Create configuration via ModifyFunctionComponentBindings with rebind
The resource Create method SHALL call `ModifyFunctionComponentBindings` API with `Operation` set to `rebind` and the user-specified `FunctionComponentBindings` list. After successful creation, it SHALL set the resource ID as `zone_id#function_id` and call Read to refresh state.

#### Scenario: Create bindings successfully
- **WHEN** user applies a new `tencentcloud_teo_function_component_binding` resource
- **THEN** the Create method SHALL call `ModifyFunctionComponentBindings` with Operation=`rebind`, set the composite ID, and read back the configuration

### Requirement: Update configuration via ModifyFunctionComponentBindings with rebind
The resource Update method SHALL call `ModifyFunctionComponentBindings` API with `Operation` set to `rebind` and the updated `FunctionComponentBindings` list when `function_component_bindings` has changed.

#### Scenario: Update bindings successfully
- **WHEN** user modifies the `function_component_bindings` in HCL and applies
- **THEN** the Update method SHALL call `ModifyFunctionComponentBindings` with Operation=`rebind` and the new binding list

### Requirement: Delete configuration by clearing all bindings
The resource Delete method SHALL call `ModifyFunctionComponentBindings` API with `Operation` set to `rebind` and an empty `FunctionComponentBindings` list to clear all bindings.

#### Scenario: Delete clears all bindings
- **WHEN** user destroys the `tencentcloud_teo_function_component_binding` resource
- **THEN** the Delete method SHALL call `ModifyFunctionComponentBindings` with Operation=`rebind` and an empty list

### Requirement: Retry on API errors
All API calls SHALL be wrapped with `resource.Retry` using `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations. Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Retry on transient error
- **WHEN** an API call fails with a transient error
- **THEN** the operation SHALL be retried according to the configured timeout

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** user references `tencentcloud_teo_function_component_binding` in HCL
- **THEN** the provider SHALL recognize and handle the resource
