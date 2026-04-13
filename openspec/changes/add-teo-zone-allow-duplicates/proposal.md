## Why

为 `tencentcloud_teo_zone` 资源添加 `allow_duplicates` 参数支持，以匹配 CreateZone API 的新增字段，满足用户在创建 TEO 站点时控制是否允许重复配置的需求。

## What Changes

- 在 `tencentcloud_teo_zone` 资源的 Schema 中新增 `allow_duplicates` 字段
- 更新 Create 函数，在调用 CreateZone API 时传入 `allow_duplicates` 参数
- 更新 Read 函数，从 DescribeZone API 响应中读取 `allow_duplicates` 字段值
- 更新 Update 函数，支持通过 ModifyZone API 更新 `allow_duplicates` 参数
- 更新 Delete 函数，确保删除操作不受影响
- 更新相关的单元测试代码，覆盖新增字段的测试场景
- 更新验收测试代码，验证 `allow_duplicates` 参数的正确性

## Capabilities

### New Capabilities

- `teo-zone-allow-duplicates`: 支持 TEO 站点资源的 `allow_duplicates` 参数配置和 CRUD 操作

### Modified Capabilities

无（未修改现有能力的 REQUIREMENTS）

## Impact

- 受影响的代码文件：
  - `tencentcloud/services/teo/resource_tc_teo_zone.go`
  - `tencentcloud/services/teo/resource_tc_teo_zone_test.go`
- 依赖的云 API：
  - CreateZone API（新增参数支持）
  - DescribeZone API（读取字段）
  - ModifyZone API（更新字段）
- 向后兼容性：新字段为 Optional 属性，不影响现有配置
