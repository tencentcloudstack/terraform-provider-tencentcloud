# Spec: ES Instance Coordinating Node

## ADDED Requirements

### Requirement: node_info_list type 字段放行 dedicatedCoordinating 取值

`tencentcloud_elasticsearch_instance` 资源 SHALL 允许 `node_info_list[].type` 字段取值 `dedicatedCoordinating`，并与现有 `hotData`、`warmData`、`dedicatedMaster` 一同构成 `type` 字段的有效取值集合。校验 SHALL 在 plan 阶段通过 `ValidateFunc` 白名单完成，默认值 SHALL 保持为 `hotData`。

**Rationale**: 腾讯云 ES 云 API `NodeInfo.Type` 已支持 `dedicatedCoordinating`（专用协调节点，见 https://cloud.tencent.com/document/api/845/30634 ）。白名单校验能在 plan 阶段拦截拼写错误，避免错误取值透传到云 API 才暴露。

#### Scenario: 配置 dedicatedCoordinating 类型节点通过 plan 校验

- **WHEN** 用户在 `tencentcloud_elasticsearch_instance` 的 `node_info_list` 中声明一个 `type = "dedicatedCoordinating"` 的节点项
- **THEN** `terraform plan` SHALL 成功
- **AND** 不会出现 "expected type to be one of ..." 校验错误

#### Scenario: 未声明 type 时默认值为 hotData

- **WHEN** 用户在 `node_info_list` 中声明节点项但不指定 `type`
- **THEN** 该节点项的 `type` SHALL 取默认值 `hotData`
- **AND** `dedicatedCoordinating` 不会被意外设为默认值

#### Scenario: 现有配置不受影响

- **WHEN** 用户应用一个仅含 `hotData`/`warmData`/`dedicatedMaster` 的既有 `tencentcloud_elasticsearch_instance` 配置
- **THEN** plan SHALL 成功
- **AND** 不会要求用户新增 `dedicatedCoordinating` 节点

### Requirement: 创建含专用协调节点的集群

资源 Create SHALL 将 `node_info_list` 中 `type = "dedicatedCoordinating"` 的节点项原样映射到 `CreateInstance` 请求的 `NodeInfoList` 元素的 `Type` 字段，完成协调节点的创建。

**Rationale**: `CreateInstance` 的 `NodeInfoList` 原生支持 `dedicatedCoordinating`，无需额外接口或后置步骤。

#### Scenario: Create 携带协调节点

- **WHEN** 用户创建 `tencentcloud_elasticsearch_instance` 且 `node_info_list` 包含一项 `type = "dedicatedCoordinating"`
- **THEN** 资源 SHALL 在 `CreateInstance` 请求的 `NodeInfoList` 中包含一个 `Type` 为 `dedicatedCoordinating` 的 `NodeInfo`
- **AND** 该 `NodeInfo` 的 `NodeNum`、`NodeType`、`DiskType`、`DiskSize`、`DiskEncrypt` SHALL 与用户配置一致

#### Scenario: Create 后 Read 回填协调节点

- **WHEN** 创建完成并执行 Read
- **THEN** Read SHALL 将云 API 返回的 `Type == "dedicatedCoordinating"` 节点写回 `node_info_list` state
- **AND** SHALL NOT 将其当作 `kibana` 节点过滤掉

### Requirement: 更新流程识别并处理 dedicatedCoordinating 节点变更

资源 Update SHALL 在 `node_info_list` 发生变更时，将 `dedicatedCoordinating` 纳入逐类型 diff 的 `typeList`，并按非数据节点路径处理其新增、删除、数量修改、规格修改与磁盘大小修改。`dedicatedCoordinating` SHALL NOT 被加入 `dataTypeList`，即不参与多可用区 node_num 倍数校验。

**Rationale**: Update 中 `typeList` 原硬编码为 `["hotData","warmData","dedicatedMaster"]`，若不扩展，协调节点的增删改会被 diff 逻辑跳过，导致 state 与云上不一致。协调节点不承载数据分片，语义与 `dedicatedMaster` 一致（非数据节点），故归类为非数据节点。

#### Scenario: 新增 dedicatedCoordinating 节点

- **WHEN** 用户在一个已有 `tencentcloud_elasticsearch_instance` 的 `node_info_list` 中新增一项 `type = "dedicatedCoordinating"`
- **THEN** 资源 SHALL 将该新增节点追加到 `baseNodeList` 后调用 `UpdateInstance` 下发
- **AND** 下发的 `NodeInfoList` SHALL 包含该协调节点
- **AND** 资源 SHALL 等待实例升级完成

#### Scenario: 删除 dedicatedCoordinating 节点

- **WHEN** 用户从 `node_info_list` 中移除一项 `type = "dedicatedCoordinating"`
- **THEN** 资源 SHALL 以不含该协调节点的 `baseNodeList` 调用 `UpdateInstance`
- **AND** 资源 SHALL 等待实例升级完成

#### Scenario: 修改 dedicatedCoordinating 节点数量

- **WHEN** 用户修改 `dedicatedCoordinating` 节点项的 `node_num`
- **THEN** 资源 SHALL 调用 `UpdateInstance`，其 `NodeInfoList` 中协调节点的 `NodeNum` 为新值
- **AND** 该变更 SHALL NOT 触发多可用区 node_num 倍数校验（协调节点为非数据节点）

#### Scenario: 修改 dedicatedCoordinating 节点规格或磁盘大小

- **WHEN** 用户修改 `dedicatedCoordinating` 节点项的 `node_type` 或 `disk_size`
- **THEN** 资源 SHALL 调用 `UpdateInstance`，其 `NodeInfoList` 中协调节点对应字段为新值
- **AND** 资源 SHALL 等待实例升级完成

#### Scenario: 修改 dedicatedCoordinating 节点的不可变字段报错

- **WHEN** 用户修改 `dedicatedCoordinating` 节点项的 `disk_type`、`encrypt` 或 `type`
- **THEN** 资源 SHALL 返回 "xxx not support change" 错误
- **AND** SHALL NOT 调用 `UpdateInstance`

### Requirement: 文档补充 dedicatedCoordinating 用法

`resource_tc_elasticsearch_instance.md` SHALL 在示例中补充含 `type = "dedicatedCoordinating"` 节点的用法，并在 `node_info_list[].type` 的取值说明中列出 `dedicatedCoordinating`。website 文档 SHALL 由收尾阶段 `make doc` 生成，不在本阶段直接修改 `website/` 目录。

**Rationale**: 用户需要明确的用法示例与取值清单，与既有节点类型文档保持一致。

#### Scenario: 文档包含协调节点示例

- **WHEN** 查看生成的 `website/docs/r/elasticsearch_instance.html.markdown`
- **THEN** 示例 SHALL 包含一个 `node_info_list` 项其 `type = "dedicatedCoordinating"`
- **AND** `type` 参数说明 SHALL 列出 `dedicatedCoordinating` 为有效取值之一
