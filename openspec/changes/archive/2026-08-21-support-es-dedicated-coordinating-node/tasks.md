## 1. 常量与类型枚举

- [x] 1.1 在 `tencentcloud/services/es/extension_elasticsearch.go` 新增常量 `ES_NODE_TYPE_DEDICATED_COORDINATING = "dedicatedCoordinating"`
- [x] 1.2 将 `ES_NODE_TYPE_DEDICATED_COORDINATING` 追加到 `ES_NODE_TYPE` 枚举切片（顺序置于 `ES_NODE_TYPE_DEDICATED_MATER` 之后），作为节点类型枚举的唯一数据源

## 2. 资源 Schema 与 CRUD 代码

- [x] 2.1 在 `tencentcloud/services/es/resource_tc_elasticsearch_instance.go` 的 `node_info_list[].type` schema 字段中移除 `ValidateFunc: tccommon.ValidateAllowedStringValue(ES_NODE_TYPE)`（移除白名单校验）
- [x] 2.2 在 `tencentcloud/services/es/resource_tc_elasticsearch_instance.go` 的 `node_info_list[].type` schema 字段中更新 `Description`，将有效取值列出为 `hotData`、`warmData`、`dedicatedMaster`、`dedicatedCoordinating`，默认值保持 `ES_NODE_TYPE_HOT_DATA`
- [x] 2.3 在 `resourceTencentCloudElasticsearchInstanceUpdate` 的 `node_info_list` 变更分支中，将 `typeList` 改为引用 `ES_NODE_TYPE` 切片（替代硬编码 `[]string{"hotData", "warmData", "dedicatedMaster", "dedicatedCoordinating"}`），`dataTypeList` 保持 `[]string{"hotData", "warmData"}` 不变
- [x] 2.4 确认 Create 分支（`resourceTencentCloudElasticsearchInstanceCreate`）已通过通用 `node_info_list` 处理原样下发 `Type`，无需额外修改；确认 Read 分支不会过滤 `dedicatedCoordinating`（仅过滤 `kibana`）

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/es/resource_tc_elasticsearch_instance_test.go` 中新增/补充单元测试用例，覆盖：`node_info_list[].type` 字段不再有 `ValidateFunc`（断言 `ValidateFunc == nil`）
- [x] 3.2 补充 `validateNodeInfoListUnique` 对含 `dedicatedCoordinating` 的 `node_info_list` 不报重复错误的测试
- [x] 3.3 补充 Update diff 测试：新增 `dedicatedCoordinating` 节点时 `UpdateInstance` 入参 `NodeInfoList` 包含该协调节点；删除时入参不含该节点
- [x] 3.4 补充 Update diff 测试：修改 `dedicatedCoordinating` 的 `disk_type`/`encrypt`/`type` 返回 "xxx not support change" 错误

## 4. 文档

- [x] 4.1 在 `tencentcloud/services/es/resource_tc_elasticsearch_instance.md` 示例中补充含 `type = "dedicatedCoordinating"` 节点的用法，并在 `node_info_list[].type` 取值说明中列出 `dedicatedCoordinating`

## 5. 验证（收尾阶段执行）

- [x] 5.1 运行 `gofmt` 格式化变更的 Go 文件
- [ ] 5.2 运行 `make doc` 生成 `website/docs/r/elasticsearch_instance.html.markdown` 文档
- [x] 5.3 运行 `openspec validate support-es-dedicated-coordinating-node --strict` 验证提案完整性
