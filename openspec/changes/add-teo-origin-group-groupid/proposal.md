## Why

当前 tencentcloud_teo_origin_group 资源的 DeleteOriginGroup API 操作需要使用 GroupId 参数，但该字段在资源定义中缺失，导致删除操作无法正确执行。为了确保资源的完整 CRUD 功能，需要添加 GroupId 字段支持。

## What Changes

- 在 tencentcloud_teo_origin_group 资源的 Schema 中新增 `group_id` 字段（类型：string，必填）
- 更新 Create、Read、Update、Delete 函数中涉及 group_id 字段的逻辑
- 确保新字段与 DeleteOriginGroup API 的 GroupId 参数正确映射
- 更新相关的单元测试和验收测试代码

## Capabilities

### New Capabilities

- `teo-origin-group-groupid`: 支持 tencentcloud_teo_origin_group 资源的 group_id 字段，用于标识源站组 ID 并支持完整的 CRUD 操作

### Modified Capabilities

- 无（这是新增字段，不修改现有功能需求）

## Impact

- 受影响的代码：`tencentcloud/services/teo/resource_tencentcloud_teo_origin_group.go`
- 受影响的测试：`tencentcloud/services/teo/resource_tencentcloud_teo_origin_group_test.go`
- 依赖的 API：TencentCloud TEO DeleteOriginGroup API
- 向后兼容性：新增字段，不影响现有配置，保持向后兼容
