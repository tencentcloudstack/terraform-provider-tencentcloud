## Context

`tencentcloud_elasticsearch_instance` 资源通过 `node_info_list`（`schema.TypeSet`）描述集群各节点层的规格，每项的 `type` 字段标识节点角色。当前 Provider 仅允许三种角色：`hotData`、`warmData`、`dedicatedMaster`（见 `tencentcloud/services/es/extension_elasticsearch.go` 的 `ES_NODE_TYPE` 白名单与常量定义）。

腾讯云 ES 云 API 的 `NodeInfo.Type` 字段已扩展支持 `dedicatedCoordinating`（专用协调节点）以及 `dedicatedMl`（专用机器学习节点）（见 https://cloud.tencent.com/document/api/845/30634 ）。本次变更只引入 `dedicatedCoordinating`，`dedicatedMl` 暂不在范围内。

资源的 `Update` 函数（`resourceTencentCloudElasticsearchInstanceUpdate`）在处理 `node_info_list` 变更时，使用硬编码的 `typeList := []string{"hotData", "warmData", "dedicatedMaster"}` 逐类型 diff，并用 `dataTypeList := []string{"hotData", "warmData"}` 区分数据节点（参与多可用区 node_num 倍数校验）与非数据节点。若不把 `dedicatedCoordinating` 加入 `typeList`，新增/删除/修改协调节点的配置在 `terraform apply` 时会被 diff 逻辑跳过，导致 state 与云上实际配置不一致。

`node_info_list` 的元素以 `type` 作为 `schema.Set` 哈希的一部分（`validateNodeInfoListUnique` 校验 type 唯一），且 `disk_type`、`encrypt`、`type` 在 Update 中被标记为不可修改字段。

## Goals / Non-Goals

**Goals:**
- 让用户可以通过 `node_info_list` 声明 `type = "dedicatedCoordinating"` 的专用协调节点，并在 create/read/update（含增、删、改数量、改规格、改磁盘大小）全链路生效。
- 保持对现有 `hotData`/`warmData`/`dedicatedMaster` 配置和 state 的完全向后兼容。
- 复用现有非数据节点（`dedicatedMaster`）的处理路径，最小化代码改动。

**Non-Goals:**
- 不引入 `dedicatedMl`（专用机器学习节点）类型。
- 不修改 `node_info_list` 的 schema 结构（不新增/移除字段、不改变 `TypeSet` 语义）。
- 不修改 `disk_type`/`type`/`encrypt` 不可变更的既有约束。
- 不修改多可用区（`multi_zone_infos`）变更逻辑本身，仅确保协调节点不参与数据节点的 node_num 倍数校验。
- 不修改 `provider.go` 注册（资源已注册）。

## Decisions

### D1 — 新增 `dedicatedCoordinating` 常量并加入校验白名单

在 `extension_elasticsearch.go` 新增 `ES_NODE_TYPE_DEDICATED_COORDINATING = "dedicatedCoordinating"`，并追加到 `ES_NODE_TYPE` 切片，使 schema 的 `ValidateFunc: tccommon.ValidateAllowedStringValue(ES_NODE_TYPE)` 放行该取值。

**备选方案**：放宽为不校验（去掉 `ValidateFunc`）。**否决**：白名单能在 plan 阶段拦截无效取值（如拼写错误），与现有 `hotData`/`warmData`/`dedicatedMaster` 的一致性更好，也避免把错误取值透传到云 API 才暴露。

### D2 — 将 `dedicatedCoordinating` 归类为非数据节点

在 `resourceTencentCloudElasticsearchInstanceUpdate` 的 `node_info_list` 变更分支中，把 `typeList` 扩展为 `[]string{"hotData", "warmData", "dedicatedMaster", "dedicatedCoordinating"}`，但 `dataTypeList` 保持 `[]string{"hotData", "warmData"}` 不变。

**理由**：协调节点不承载数据分片，行为与 `dedicatedMaster` 一致——不参与多可用区 node_num 倍数校验（`isDataNode && !changeMultiZone` 分支），在「新增节点」分支中走 `!isDataNode` 路径（将新节点追加到 `baseNodeList` 后整体下发 `UpdateInstance`）。这与云 API 对协调节点的语义一致。

**备选方案**：把 `dedicatedCoordinating` 加入 `dataTypeList`。**否决**：会导致多可用区变更时对协调节点强制 node_num 倍数校验，而协调节点并非按数据可用区分布，会误报错误。

### D3 — 更新 schema Description 与文档示例

将 `node_info_list[].type` 的 `Description` 更新为列举 `hotData`、`warmData`、`dedicatedMaster`、`dedicatedCoordinating` 四种取值，默认值仍为 `hotData`。在 `resource_tc_elasticsearch_instance.md` 示例中补充一个含协调节点的用法。文档最终通过收尾阶段 `make doc` 生成 website 文档，不在本阶段直接修改 `website/`。

### D4 — 单元测试覆盖校验与 update diff

在 `resource_tc_elasticsearch_instance_test.go` 中用 gomonkey mock 云 API（`ElasticsearchService.UpdateInstance` / `DescribeInstanceById`）补充单元测试：
- schema 校验白名单放行 `dedicatedCoordinating`；
- `node_info_list` 含协调节点时 `validateNodeInfoListUnique` 不报错；
- Update 中新增/删除协调节点调用 `UpdateInstance` 的入参 `NodeInfoList` 包含协调节点 `Type`。

不使用 terraform 验收测试套件，遵循对新增字段/参数用 mock 做业务逻辑单测的规范。

## Risks / Trade-offs

- **[协调节点磁盘字段语义]** → 协调节点与专用主节点一样不带数据盘，云 API 可能忽略 `disk_type`/`disk_size`。本次不改变现有 `diskTypeSizeDefault` CustomizeDiff 对非大数据/高IO机型默认填值的逻辑；若协调节点被填入默认磁盘值且云 API 忽略，Read 会回填实际值（与 `dedicatedMaster` 现有行为一致）。保持与 `dedicatedMaster` 一致以降低风险。
- **[typeList 顺序]** → `typeList` 的迭代顺序影响多步 Update 的执行顺序。把 `dedicatedCoordinating` 放在末尾，保证既有三种类型的处理顺序不变，避免对现有 state 产生行为回归。
- **[Read 过滤 kibana]** → 现有 Read 跳过 `*item.Type == "kibana"` 的节点，协调节点不在过滤名单，会被正常回填到 state，符合预期。
