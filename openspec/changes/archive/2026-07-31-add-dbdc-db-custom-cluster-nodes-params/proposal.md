## Why

`tencentcloud_dbdc_db_custom_cluster_nodes` 数据源当前仅返回节点的部分字段（node_id、node_name、lan_ip、ssh_endpoint、status、zone、node_type），但 `DescribeDBCustomClusterNodes` API 实际返回的 `DBCustomClusterNode` 结构体中还包含网络模式（`NetworkMode`）和 ENI 访问 IP（`EniIP`）字段。用户在配置三层网络联通（cross_tenant_eni）模式时，需要知道节点的网络模式以及可访问的 ENI IP 地址，当前数据源无法提供这些信息。

## What Changes

- 在 `tencentcloud_dbdc_db_custom_cluster_nodes` 数据源的 `node_set` schema 中新增 `network_mode` 计算字段（TypeString, Computed），映射 API 的 `response.NodeSet.NetworkMode`
- 在 `tencentcloud_dbdc_db_custom_cluster_nodes` 数据源的 `node_set` schema 中新增 `eni_ip` 计算字段（TypeString, Computed），映射 API 的 `response.NodeSet.EniIP`
- 在数据源 Read 方法中新增对应字段的 nil 检查与赋值逻辑
- 更新 `.md` 文档示例

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `dbdc-db-custom-cluster-nodes-datasource`: 在"完整的节点信息映射"需求中新增 `network_mode` 和 `eni_ip` 两个计算字段，使用户能够获取节点的网络模式和 ENI 访问 IP 地址

## Impact

### Affected Code
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.go` - 在 `node_set` schema 中新增 `network_mode`、`eni_ip` 字段，并在 Read 方法中新增赋值逻辑
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.md` - 更新文档（由 `make doc` 自动生成）
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes_test.go` - 在单元测试中补充新增字段的验证

### Affected APIs
- `DescribeDBCustomClusterNodes`（dbdc v20201029）- 读取 `Response.NodeSet[].NetworkMode` 和 `Response.NodeSet[].EniIP`，均为已有出参，无需修改请求参数

### Breaking Changes
无 - 所有变更均为新增 Computed 字段，向后兼容

### Dependencies
无 - 使用已有的 `DescribeDBCustomClusterNodes` API 调用，`NetworkMode` 和 `EniIP` 字段已存在于云 API SDK 的 `DBCustomClusterNode` 结构体中
