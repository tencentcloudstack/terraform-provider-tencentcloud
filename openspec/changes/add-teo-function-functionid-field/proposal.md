## Why

TEO Function 资源需要支持通过 CreateFunction API 的响应参数 FunctionId 来获取和展示函数的唯一标识符。这对于用户识别和管理已创建的边缘函数至关重要，特别是在后续的函数更新、删除或查询操作中需要通过 FunctionId 进行精确操作的场景。

## What Changes

- 在 tencentcloud_teo_function 资源的 Schema 中新增 FunctionId 字段（string 类型，Optional）
- 更新 Create 函数，在创建函数后从 API 响应中读取并设置 FunctionId 到状态中
- 更新 Read 函数，从 DescribeFunction 或相关 API 响应中读取 FunctionId 并更新到状态中
- 更新 Update 和 Delete 函数中涉及 FunctionId 字段的逻辑（如需要）
- 确保新增字段的 Optional 属性与 CAPI 接口定义一致（FunctionId 为可选字段）
- 更新相关的单元测试和验收测试代码，验证 FunctionId 字段的正确读写

## Capabilities

### New Capabilities

- `teo-function-functionid`: 为 tencentcloud_teo_function 资源添加 FunctionId 字段支持，实现函数 ID 的读取和展示功能

### Modified Capabilities

无（仅添加新字段，不修改现有行为要求）

## Impact

- **受影响的代码**: tencentcloud/services/teo/resource_tencentcloud_teo_function.go
- **受影响的测试**: resource_tencentcloud_teo_function_test.go
- **API 依赖**: TencentCloud TEO CreateFunction API（Read 响应参数）
- **兼容性**: 向后兼容，仅新增 Optional 字段，不影响现有配置和状态
- **系统影响**: 无（仅字段级别的修改）
