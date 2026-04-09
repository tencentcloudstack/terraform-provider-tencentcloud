## Why

TEO (TencentCloud EdgeOne) 边缘 KV 存储功能需要支持查询操作，当前 Terraform Provider 缺少查询边缘 KV 数据的能力。用户需要能够查询指定的键值对信息，包括键名、键值和过期时间，以便进行数据读取和验证。

## What Changes

- 添加新的 Terraform Resource：`tencentcloud_teo_edge_k_v_get`，支持查询 TEO 边缘 KV 存储中的键值对数据
- 实现 Resource Schema 定义，包含 ZoneId、Namespace、Keys 等参数
- 实现 Create、Read、Update、Delete 四个 CRUD 操作函数
- 添加对应的单元测试和验收测试代码
- 生成相应的文档和示例文件

## Capabilities

### New Capabilities
- `teo-edge-kv-get-resource`: TEO 边缘 KV 查询资源，支持通过站点 ID、命名空间和键名列表查询键值对数据，返回包含键名、键值和过期时间的结果

### Modified Capabilities
(无现有能力的需求变更)

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_edge_k_v_get.go`
- 新增测试：`tencentcloud/services/teo/resource_tencentcloud_teo_edge_k_v_get_test.go`
- 新增文档：`website/docs/r/teo_edge_k_v_get.html.markdown`
- 新增示例：`tencentcloud/services/teo/resource_tencentcloud_teo_edge_k_v_get.md`
- 依赖 TEO SDK 的 EdgeKVGet API，无新增外部依赖
- 不影响现有资源，向后兼容
