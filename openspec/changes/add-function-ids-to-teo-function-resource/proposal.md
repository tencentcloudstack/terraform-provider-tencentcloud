## Why

用户需要在 tencentcloud_teo_function 资源的读取操作中支持通过 FunctionIds 字段进行过滤。这个字段允许用户在查询函数时指定一个函数 ID 列表，从而提高查询效率和灵活性。DescribeFunctions API 已经支持该参数，需要在 Terraform Provider 中对接该功能。

## What Changes

- 在 tencentcloud_teo_function 资源的 Schema 中新增 FunctionIds 字段（list, 可选）
- 更新 Read 函数，使用 FunctionIds 字段作为 DescribeFunctions API 的请求参数进行过滤
- 更新相关的单元测试和验收测试代码，覆盖新增字段的功能

## Capabilities

### New Capabilities
(无 - 本次变更不引入新的能力规格)

### Modified Capabilities
(无 - 本次变更仅新增可选字段，不改变现有能力的行为要求)

## Impact

- **修改的文件**:
  - `tencentcloud/services/teo/resource_tc_teo_function.go` - 新增 FunctionIds 字段定义，更新 Read 函数逻辑
  - `tencentcloud/services/teo/resource_tc_teo_function_test.go` - 新增 FunctionIds 字段的单元测试
  - `tencentcloud/services/teo/resource_tc_teo_function.md` - 更新资源文档，说明 FunctionIds 字段的用法

- **影响的 API**:
  - DescribeFunctions API - 新增 FunctionIds 请求参数

- **兼容性**:
  - 新增字段为可选属性，保持向后兼容，不会影响现有配置
