## Why

tencentcloud_teo_function 资源需要支持通过 CreateFunction API 创建函数时传入 FunctionId 参数，以满足腾讯云边缘函数（TEO）服务的新需求，允许用户直接引用已存在的函数 ID 而非总是创建新函数。

## What Changes

- 为 tencentcloud_teo_function 资源新增 FunctionId 参数（可选类型）
- 支持通过 CreateFunction API 接入 FunctionId 参数
- 保持向后兼容，不破坏现有配置和 state

## Capabilities

### New Capabilities
- `teo-function-function-id`: 支持在 tencentcloud_teo_function 资源中使用 FunctionId 参数，允许通过 CreateFunction API 引用已存在的函数 ID

### Modified Capabilities
- 无（不修改现有能力的规格要求，仅新增可选参数）

## Impact

- **Affected Resources**: tencentcloud/services/teo/resource_tencentcloud_teo_function.go
- **Affected API**: CreateFunction API (用于 create 操作)
- **Schema Changes**: 新增可选字段 FunctionId
- **Test Changes**: 需要更新或新增测试用例以验证 FunctionId 参数的行为
- **Documentation**: 需要更新 resource_tc_teo_function.md 文档
