## Why

需要为 `tencentcloud_teo_function` 资源新增 Functions 字段及其子字段，以便支持从 DescribeFunctions API 读取完整的函数信息。当前资源缺少这些关键字段，导致无法在 Terraform 中完整管理和查询 TEO (Tencent Edge One) 函数配置。

## What Changes

在 `tencentcloud_teo_function` 资源中新增以下字段定义：

- 新增 `functions` 列表字段（类型：list），包含以下子字段：
  - `function_id` (string): 函数 ID
  - `zone_id` (string): 站点 ID
  - `name` (string): 函数名字
  - `remark` (string): 函数描述
  - `content` (string): 函数内容
  - `domain` (string): 函数默认域名
  - `create_time` (string): 创建时间（UTC, ISO 8601 格式）
  - `update_time` (string): 修改时间（UTC, ISO 8601 格式）

更新 CRUD 操作逻辑：
- Read 函数：从 DescribeFunctions API 响应中读取 Functions 列表并映射到资源字段
- Create/Update/Delete 函数：根据新字段的属性（Computed/Optional）进行相应处理

更新测试代码：
- 更新单元测试以验证新字段的读写
- 更新验收测试以确保实际 API 调用正确处理新字段

## Capabilities

### New Capabilities
- `teo-function-params`: 为 tencentcloud_teo_function 资源新增 Functions 参数支持，包括函数 ID、站点 ID、函数名、描述、内容、域名、创建时间和修改时间等字段

### Modified Capabilities
(无 - 本次变更为新增字段，不修改现有功能的需求)

## Impact

- 代码影响：
  - 修改 `tencentcloud/services/teo/resource_tc_teo_function.go` 中的 Schema 定义
  - 更新 Read 函数以从 DescribeFunctions API 读取新字段
  - 更新 Create/Update/Delete 函数以支持新字段
  - 更新 `resource_tc_teo_function_test.go` 单元测试
  - 更新 `resource_tc_teo_function.md` 文档示例

- API 依赖：
  - 使用 DescribeFunctions API（Read 操作）读取函数信息

- 兼容性：
  - 新增字段为 Optional/Computed 类型，保持向后兼容
  - 不破坏现有 Terraform 配置和 state
