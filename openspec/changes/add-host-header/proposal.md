## Why

为了支持在 TEO（TencentCloud EdgeOne）源站组配置中设置回源 Host Header，需要在 tencentcloud_teo_origin_group 资源中新增 HostHeader 字段。该字段允许用户指定回源时使用的 Host Header，仅当源站类型为 HTTP 时生效。这是基于 CreateOriginGroup API 的新增参数需求。

## What Changes

- **新增字段**: 在 tencentcloud_teo_origin_group 资源的 Schema 中新增 `host_header` 字段（string 类型，Optional）
- **更新 CRUD 逻辑**: 在 Create 和 Update 操作中传递 HostHeader 参数到 CreateOriginGroup API
- **读取逻辑**: 在 Read 操作中从 API 响应中读取 HostHeader 值
- **测试更新**: 为新字段添加单元测试和验收测试用例

## Capabilities

### New Capabilities
- `origin-group-host-header`: 支持在 TEO 源站组资源中配置回源 Host Header 参数

### Modified Capabilities
- 无（仅实现层面的改动，不涉及现有能力的需求变更）

## Impact

- 受影响的文件:
  - `tencentcloud/services/teo/resource_tc_teo_origin_group.go` (Schema 和 CRUD 函数)
  - `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go` (单元测试)
  - 可能的验收测试文件
- 无破坏性变更（仅新增 Optional 字段）
- 依赖: 需要确保 tencentcloud-sdk-go 已包含 CreateOriginGroup API 的 HostHeader 参数支持