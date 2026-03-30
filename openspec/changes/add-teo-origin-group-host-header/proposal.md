## Why

腾讯云 TEO (Tencent EdgeOne) 服务提供了源站组（OriginGroup）管理功能，CreateOriginGroup API 中包含 HostHeader 参数，用于设置请求源站时携带的自定义 Host 头信息。该参数在用户需要将请求路由到特定的虚拟主机或需要自定义 Host 头的场景下非常重要。当前 tencentcloud_teo_origin_group 资源未支持该参数，无法满足用户配置自定义 Host 头的需求。

## What Changes

- 为 tencentcloud_teo_origin_group 资源新增 `host_header` 参数
- 该参数对应 CreateOriginGroup API 的 HostHeader 字段
- 支持用户在创建和更新源站组时配置自定义 Host 头信息
- 确保参数正确传递到底层 API 调用

## Capabilities

### New Capabilities
- `teo-origin-group-host-header`: 支持 tencentcloud_teo_origin_group 资源的 host_header 参数配置，允许用户为源站请求设置自定义 Host 头

### Modified Capabilities
- （无，仅新增参数，不修改现有行为）

## Impact

- **Affected Code**:
  - `tencentcloud/services/teo/resource_tencentcloud_teo_origin_group.go`: 资源定义和 CRUD 函数
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`: 服务层 API 调用（如需要）
- **API Calls**: CreateOriginGroup API 的 HostHeader 参数映射
- **Schema Changes**: 新增 Optional 类型的 `host_header` 字段（TypeString）
- **Documentation**: 更新资源文档（resource_tc_teo_origin_group.md）和示例
- **Tests**: 新增或更新相关测试用例以验证 host_header 参数功能
