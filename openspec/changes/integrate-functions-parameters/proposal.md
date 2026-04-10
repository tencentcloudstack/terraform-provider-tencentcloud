## Why

TEO（EdgeOne）服务的 Function 资源当前仅实现了基本的 CRUD 操作，但在实际使用中，用户需要更丰富的函数配置能力，包括环境变量、规则绑定、区域选择等高级参数。这些功能对于复杂的边缘计算场景至关重要，能够提升函数管理的灵活性和可配置性。

## What Changes

本次变更将在 `tencentcloud_teo_function` 资源中新增以下功能：

- **环境变量支持**：允许用户配置函数的环境变量（Environment Variables），支持字符串和 JSON 类型
- **规则绑定配置**：添加函数规则（Function Rule）的绑定能力，实现更灵活的函数触发规则
- **区域选择参数**：支持为函数指定特定的部署区域（Region Selection）

## Capabilities

### New Capabilities

- `teo-function-env-vars`: 在 teo_function 资源中添加环境变量（environment_variables）字段，支持通过环境变量配置函数运行时参数
- `teo-function-rules`: 添加函数规则（rules）字段，允许配置函数触发规则和执行条件
- `teo-function-regions`: 添加区域选择（region_selection）字段，支持指定函数部署的区域范围

### Modified Capabilities

- `teo-function`: 更新现有的 teo_function 资源，添加新参数字段以支持高级配置功能

## Impact

- **代码修改**：`tencentcloud/services/teo/resource_tc_teo_function.go` - 修改 Schema 定义和 CRUD 函数
- **测试代码**：`tencentcloud/services/teo/resource_tc_teo_function_test.go` - 添加新字段的测试用例
- **API 依赖**：需要确认 TEO SDK 中 CreateFunction 和 ModifyFunction API 是否支持这些新参数
- **文档更新**：需要更新 `website/docs/r/teo_function.html` 文档，说明新增字段的使用方法
