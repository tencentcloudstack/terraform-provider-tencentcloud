## Why

TencentCloud EdgeOne (TEO) 边缘函数支持绑定组件（如 KV 命名空间），用户需要通过 Terraform 管理边缘函数的组件绑定配置，实现基础设施即代码的自动化管理。当前 provider 缺少对函数组件绑定配置的支持。

## What Changes

- 新增 `tencentcloud_teo_function_component_binding` 资源（RESOURCE_KIND_CONFIG 类型），用于管理 TEO 边缘函数的组件绑定配置
- 资源仅包含 Read 和 Update 操作（配置型资源，资源存在则配置存在）
- Read 使用 `DescribeFunctionComponentBindings` 接口查询绑定列表
- Update 使用 `ModifyFunctionComponentBindings` 接口修改绑定关系，支持 bind/bind-override/unbind/rebind 四种操作模式
- 资源 ID 使用 `zone_id` 和 `function_id` 的联合 ID

## Capabilities

### New Capabilities
- `teo-function-component-binding-config`: 管理 TEO 边缘函数的组件绑定配置，支持查询和修改函数与 KV 命名空间等组件的绑定关系

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config.md`
- 修改文件: `tencentcloud/provider.go`（注册资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中）
