## Why

用户需要通过 Terraform Provider 查询 TEO（TencentCloud EdgeOne）服务的边缘 KV 数据，目前缺少对应的 Terraform Resource 来支持这一功能。该资源允许用户声明式地查询和获取边缘 KV 键值对数据，便于在基础设施即代码的流程中使用和集成。

## What Changes

- 新增 Terraform Resource: `tencentcloud_teo_edge_k_v_get`
  - 支持指定站点 ID (ZoneId)、命名空间 (Namespace) 和键名列表 (Keys) 来查询 KV 数据
  - 返回包含键名、键值和过期时间的键值对数据列表
  - 实现 Create、Read、Update、Delete 四个标准 CRUD 操作
  - 添加对应的单元测试和验收测试代码
  - 添加资源文档示例

## Capabilities

### New Capabilities
- `teo-edge-kv-get`: 支持 Edge KV 数据查询的 Terraform Resource，提供基于 ZoneId、Namespace 和 Keys 的键值对数据查询能力

### Modified Capabilities
(无)

## Impact

- 新增文件位置：`tencentcloud/services/teo/resource_tencentcloud_teo_edge_k_v_get.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tencentcloud_teo_edge_k_v_get_test.go`
- 新增文档文件：`website/docs/r/teo_edge_k_v_get.html.markdown`
- 依赖 TEO 服务的 EdgeKVGet CAPI 接口
- 不影响现有资源和数据源的功能，保持完全向后兼容
