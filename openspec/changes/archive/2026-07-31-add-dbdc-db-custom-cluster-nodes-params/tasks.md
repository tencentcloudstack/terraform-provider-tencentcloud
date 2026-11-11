## 1. Schema 定义

- [x] 1.1 在 `data_source_tc_dbdc_db_custom_cluster_nodes.go` 的 `node_set` schema 中新增 `network_mode` 字段（TypeString, Computed），描述为网络模式，映射 API 的 `response.NodeSet.NetworkMode`
- [x] 1.2 在 `data_source_tc_dbdc_db_custom_cluster_nodes.go` 的 `node_set` schema 中新增 `eni_ip` 字段（TypeString, Computed），描述为当网络模式为 cross_tenant_eni 时节点的可访问 IP 地址，映射 API 的 `response.NodeSet.EniIP`

## 2. Read 函数实现

- [x] 2.1 在 `dataSourceTencentCloudDbdcDbCustomClusterNodesRead` 方法中，于遍历 `respData` 的循环内新增 `network_mode` 字段的赋值逻辑：判断 `node.NetworkMode != nil` 后设置 `nodeMap["network_mode"] = node.NetworkMode`
- [x] 2.2 在 `dataSourceTencentCloudDbdcDbCustomClusterNodesRead` 方法中，于遍历 `respData` 的循环内新增 `eni_ip` 字段的赋值逻辑：判断 `node.EniIP != nil` 后设置 `nodeMap["eni_ip"] = node.EniIP`

## 3. 测试

- [x] 3.1 在 `data_source_tc_dbdc_db_custom_cluster_nodes_test.go` 的 `TestDbdcDbCustomClusterNodesDS_ReadBasic` 测试中，为 mock 的 NodeSet 节点添加 `NetworkMode` 和 `EniIP` 字段值，并新增对应的 assert 断言
- [x] 3.2 在 `data_source_tc_dbdc_db_custom_cluster_nodes_test.go` 的 `TestDbdcDbCustomClusterNodesDS_Schema` 测试中，新增对 `network_mode` 和 `eni_ip` schema 字段存在的断言

## 4. 文档

- [ ] 4.1 更新 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.md` 文档（由 `make doc` 自动生成，收尾阶段执行）

## 5. 验证

- [ ] 5.1 执行 `gofmt` 格式化代码（收尾阶段执行）
- [ ] 5.2 执行 `make doc` 生成文档（收尾阶段执行）
