## Why

腾讯云 Elasticsearch Service (ES) 的云 API `NodeInfo.Type` 字段已支持 `dedicatedCoordinating`（专用协调节点）类型，但 Terraform Provider 的 `tencentcloud_elasticsearch_instance` 资源当前仅允许 `node_info_list` 的 `type` 字段取值 `hotData`、`warmData`、`dedicatedMaster`（见 `extension_elasticsearch.go` 中的 `ES_NODE_TYPE` 白名单）。这导致用户无法通过 Terraform 创建或管理带专用协调节点的 ES 集群，必须手动到控制台操作，破坏了基础设施即代码的一致性。此外，资源 `Update` 流程中的 `typeList` 硬编码了上述三种类型，新增协调节点类型时必须同步更新 CRUD 逻辑，否则协调节点的增删改不会被正确处理。

## What Changes

- 在 `extension_elasticsearch.go` 中新增常量 `ES_NODE_TYPE_DEDICATED_COORDINATING = "dedicatedCoordinating"`，并将其加入 `ES_NODE_TYPE` 校验白名单。
- 更新 `tencentcloud_elasticsearch_instance` 资源 schema 中 `node_info_list[].type` 字段的 `Description`，将 `dedicatedCoordinating` 加入有效取值说明。
- 更新 `resourceTencentCloudElasticsearchInstanceUpdate` 中 `node_info_list` 变更分支的 `typeList`（新增 `dedicatedCoordinating`），使其在协调节点的数量、规格、磁盘大小变更以及增删时与现有 `dedicatedMaster` 非数据节点走相同的处理路径（非数据节点，不参与多可用区 node_num 校验）。
- 在 `resource_tc_elasticsearch_instance.md` 文档示例中补充 `dedicatedCoordinating` 节点的用法说明。
- 在 `_test.go` 中补充单元测试用例，覆盖协调节点类型的校验与 update diff 逻辑。

## Capabilities

### New Capabilities
- `es-instance-coordinating-node`: 为 `tencentcloud_elasticsearch_instance` 资源的 `node_info_list[].type` 字段增加对 `dedicatedCoordinating`（专用协调节点）类型的支持，覆盖 schema 校验、创建、读取、更新（增删改协调节点）全链路。

### Modified Capabilities

（无。现有 `es-instance-destroy-protection` 规范的需求不发生变更。）

## Impact

- **受影响的代码**:
  - `tencentcloud/services/es/extension_elasticsearch.go`（新增常量与校验白名单条目）
  - `tencentcloud/services/es/resource_tc_elasticsearch_instance.go`（schema 描述、Update 中 `typeList`/数据处理逻辑）
  - `tencentcloud/services/es/resource_tc_elasticsearch_instance_test.go`（新增单元测试）
- **受影响的文档**:
  - `tencentcloud/services/es/resource_tc_elasticsearch_instance.md`（补充示例与取值说明，最终通过 `make doc` 生成 website 文档）
- **API 依赖**:
  - ES API v20180416: `CreateInstance`、`UpdateInstance`、`DescribeInstances` 中的 `NodeInfo.Type` 字段，已原生支持 `dedicatedCoordinating`（见 https://cloud.tencent.com/document/api/845/30634 数据结构 NodeInfo）。
- **兼容性**: 无破坏性变更。`dedicatedCoordinating` 是对 `type` 字段可选值的纯新增，现有 `hotData`/`warmData`/`dedicatedMaster` 配置和 state 完全不受影响。
