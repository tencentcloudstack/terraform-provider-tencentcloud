## Why

TencentCloud EdgeOne (TEO) 边缘函数支持绑定组件（如 KV 命名空间），但当前 Terraform Provider 缺少对函数组件绑定配置的管理能力。用户无法通过 Terraform 声明式地管理边缘函数与 KV 命名空间之间的绑定关系，需要新增 `tencentcloud_teo_function_component_binding` 资源来填补这一空白。

## What Changes

- 新增 `tencentcloud_teo_function_component_binding` 资源（RESOURCE_KIND_CONFIG 类型），用于管理 TEO 边缘函数的组件绑定配置
- 资源支持 Read（通过 `DescribeFunctionComponentBindings` 接口查询绑定列表）和 Update（通过 `ModifyFunctionComponentBindings` 接口修改绑定关系）操作
- 支持的绑定组件类型为 KV 命名空间（kv_namespace）
- 支持 bind-override 操作模式，实现声明式的绑定管理
- 资源 ID 使用 zone_id 和 function_id 的联合 ID

## Capabilities

### New Capabilities

- `teo-function-component-binding-config`: 管理 TEO 边缘函数的组件绑定配置，支持读取和更新函数与 KV 命名空间的绑定关系

### Modified Capabilities

（无）

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config.go`
- 新增测试文件: `tencentcloud/services/teo/resource_tc_teo_function_component_binding_config_test.go`
- 新增文档: `tencentcloud/services/teo/resource_tc_teo_function_component_binding.md`
- 修改 `tencentcloud/provider.go`: 注册新资源
- 修改 `tencentcloud/provider.md`: 添加资源文档引用
- 依赖 SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
