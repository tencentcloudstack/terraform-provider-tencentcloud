## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.go` 的 `node_set` 元素 schema 中新增 `network_mode`（TypeString, Computed）字段，映射 `DBCustomNode.NetworkMode`
- [x] 1.2 在同一 `node_set` 元素 schema 中新增 `eni_ip`（TypeString, Computed）字段，映射 `DBCustomNode.EniIP`

## 2. Read 函数

- [x] 2.1 在 `dataSourceTencentCloudDbdcDbCustomNodesRead` 的节点遍历循环中，新增对 `node.NetworkMode` 的 nil 判断并写入 `nodeMap["network_mode"]`（位于 `host_ip` 之后）
- [x] 2.2 在同一节点遍历循环中，新增对 `node.EniIP` 的 nil 判断并写入 `nodeMap["eni_ip"]`（位于 `network_mode` 之后）

## 3. 测试

- [x] 3.1 在 `data_source_tc_dbdc_db_custom_nodes_test.go` 的 `TestDbdcDbCustomNodesDS_ReadBasic` mock 响应中补充 `NetworkMode` 和 `EniIP` 字段值，并新增对应断言
- [x] 3.2 在 `TestDbdcDbCustomNodesDS_Schema` 中补充对 `network_mode` 和 `eni_ip` 字段存在的断言
- [x] 3.3 （可选）新增或补充一个场景测试，覆盖 `NetworkMode=NetworkModeCrossTenantENI` 且 `EniIP` 有值的情况

## 4. 文档

- [x] 4.1 同步更新 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.md`（仅在需要时补充说明，不手动编写 Argument/Attribute Reference）

## 5. 验证（由收尾阶段执行）

- [ ] 5.1 `gofmt` 格式化（由 tfpacer-finalize skill 执行）
- [ ] 5.2 `make doc` 生成 website/docs 文档（由 tfpacer-finalize skill 执行）
