## Why

腾讯云 TEO (TencentCloud EdgeOne) 边缘函数资源需要在读取操作中支持通过 `FunctionIds` 参数进行查询过滤。当前 `tencentcloud_teo_function` 资源虽然在内部使用了 `DescribeFunctions` API 的 `FunctionIds` 参数，但未将其作为资源参数暴露给用户，限制了用户在特定场景下对函数资源的灵活查询和管理能力。

## What Changes

- 在 `tencentcloud_teo_function` 资源 Schema 中新增 `function_ids` 参数字段
- 更新 Read 函数，支持使用 `function_ids` 参数调用 `DescribeFunctions` API
- 调整 Read 函数逻辑，支持单个函数 ID 查询（向后兼容）和多个函数 ID 查询（通过 `function_ids` 参数）
- 确保 `function_ids` 参数为 Optional 属性，与 CAPI 接口定义一致
- 更新单元测试和验收测试代码，覆盖新增的 `function_ids` 参数场景

**Breaking Change**: 无。此为新增字段，不破坏现有功能。

## Capabilities

### New Capabilities

- `teo-function-ids-query`: 新增在 TEO Function 资源中通过 `function_ids` 参数进行多函数查询的能力

### Modified Capabilities

- 无。现有资源的基本 CRUD 行为保持不变，仅在 Read 操作中扩展了查询能力

## Impact

- **Affected Code**:
  - `/repo/tencentcloud/services/teo/resource_tc_teo_function.go` - 修改 Schema 定义和 Read 函数逻辑
  - `/repo/tencentcloud/services/teo/resource_tc_teo_function_test.go` - 更新单元测试
  - `/repo/tencentcloud/services/teo/service_tencentcloud_teo.go` - 可能需要调整服务层方法以支持多函数查询

- **Affected APIs**:
  - `DescribeFunctions` API 的 `FunctionIds` 参数将被更充分利用

- **Affected Dependencies**: 无新增依赖

- **Systems**: 仅影响 Terraform Provider 的 TEO Function 资源功能，不影响其他系统
