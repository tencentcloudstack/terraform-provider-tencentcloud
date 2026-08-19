## Why

用户需要在 Terraform 中查询 EdgeOne（TEO）负载均衡实例下源站组的健康状态，以便在基础设施即代码场景下监控源站健康情况。当前 Provider 缺少对应的数据源，用户无法通过 Terraform 获取源站组健康状态信息，必须手动调用 API 或使用控制台查询，影响了自动化运维效率。

## What Changes

- 新增数据源 `tencentcloud_teo_origin_group_health_status`，对应腾讯云 TEO 服务的 `DescribeOriginGroupHealthStatus` API
- 支持通过站点 ID（ZoneId）和负载均衡实例 ID（LBInstanceId）查询源站组健康状态
- 支持可选的源站组 ID 列表（OriginGroupIds）过滤，不填写时默认获取负载均衡下所有源站组的健康状态
- 返回源站组健康状态列表，包含源站组 ID、综合决策的源站健康状态、各健康检查区域的源站健康状态

## Capabilities

### New Capabilities
- `teo-origin-group-health-status-datasource`: 实现 TEO 源站组健康状态数据源，支持按站点和负载均衡实例查询源站组健康状态

### Modified Capabilities
<!-- 无现有能力需要修改 -->

## Impact

- **新增文件**:
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group_health_status.go` (数据源实现)
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group_health_status_test.go` (单元测试)
  - `tencentcloud/services/teo/data_source_tc_teo_origin_group_health_status.md` (文档)
- **修改文件**:
  - `tencentcloud/provider.go` (注册新数据源)
  - `tencentcloud/provider.md` (文档)
- **依赖**: 使用现有的 `tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包，无需新增依赖
- **兼容性**: 纯新增功能，不影响现有资源和数据源，完全向后兼容
