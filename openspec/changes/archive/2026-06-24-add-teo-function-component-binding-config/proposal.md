## Why

TencentCloud EdgeOne (TEO) 边缘函数支持绑定组件（如 KV 命名空间），但当前 Terraform Provider 缺少对函数组件绑定配置的管理能力。用户无法通过 Terraform 声明式地管理边缘函数与 KV 命名空间等组件的绑定关系，需要新增 `tencentcloud_teo_function_component_binding` 资源来填补这一空白。

## What Changes

- 新增 RESOURCE_KIND_CONFIG 类型资源 `tencentcloud_teo_function_component_binding`，用于管理 TEO 边缘函数的组件绑定配置
- 资源支持 Read（通过 `DescribeFunctionComponentBindings` 接口查询绑定列表）和 Update（通过 `ModifyFunctionComponentBindings` 接口修改绑定关系）操作
- Update 操作使用 `rebind` 模式，实现声明式的全量绑定管理
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源

## Capabilities

### New Capabilities
- `teo-function-component-binding-config`: 管理 TEO 边缘函数的组件绑定配置，支持查询和更新函数与 KV 命名空间等组件的绑定关系

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding.md`
- 修改文件: `tencentcloud/provider.go`（注册资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中）
