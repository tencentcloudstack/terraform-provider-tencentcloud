# dbdc-db-custom-cluster-node-config-list-datasource Specification

## Purpose
Defines the `tencentcloud_dbdc_db_custom_cluster_node_config_list` data source, which queries DB Custom cluster node scheduling configuration (labels and taints) via the TencentCloud DBDC `DescribeDBCustomClusterNodeConfig` API, enabling Terraform configurations to dynamically discover and reference node Kubernetes scheduling settings.

## Requirements

### Requirement: 支持查询 DB Custom 集群节点配置列表
数据源 `tencentcloud_dbdc_db_custom_cluster_node_config_list` 必须能够查询指定 DB Custom 集群内节点的配置信息（labels 和 taints）。

**Rationale**: 用户需要在 Terraform 配置中动态发现和引用 DB Custom 集群内节点的 Kubernetes 调度配置（标签和污点），用于基础设施规划和 Pod 调度决策。

#### Scenario: 按 Cluster ID 查询节点配置列表
- **WHEN** 用户指定了 `cluster_id` 参数（必填）
- **THEN** 数据源返回该集群下所有节点的配置信息

**Acceptance Criteria**:
- `cluster_id` 为 Required 参数，类型 TypeString，映射到云 API 的 `request.ClusterId`
- 调用 `DescribeDBCustomClusterNodeConfig` API 获取节点配置列表
- 该接口为同步读接口，无需异步轮询

#### Scenario: 按 Node IDs 过滤查询节点配置
- **WHEN** 用户指定了 `node_ids` 参数（可选）
- **THEN** 仅返回指定节点的配置信息

**Acceptance Criteria**:
- `node_ids` 为 Optional 参数，类型 TypeList，元素类型 TypeString
- 映射到云 API 的 `request.NodeIds`（`[]*string`）
- 每次请求的数量上限为 100（由云 API 限制）
- 在 Read 函数中将 `[]interface{}` 转换为 `[]*string`（使用 `helper.String`）

### Requirement: 完整的节点配置信息映射
数据源必须返回节点的完整调度配置信息，涵盖 `DBCustomClusterNodeConfig` 结构体中的所有字段，并将列表展开（flatten），使每个字段都可被 Terraform 单独 set/read。

**Rationale**: 用户需要完整的节点配置信息（标签和污点）用于调度决策，不能将所有字段再嵌套一层 `xxx_set` 结构。

#### Scenario: 返回节点配置基础字段
- **WHEN** 查询到节点配置列表
- **THEN** `node_set` 列表中每个元素包含以下字段

**Acceptance Criteria**:
- `node_set` 为 Computed 参数，类型 TypeList
- `node_set` 的 Elem 为 `schema.Resource`，其 Schema 包含 `node_id`、`labels`、`taints` 字段（直接平铺，不再嵌套一层）
- `node_id` - 节点ID (TypeString, Computed)，映射到 API 的 `DBCustomClusterNodeConfig.NodeId`

#### Scenario: 返回节点标签信息
- **WHEN** 查询到节点配置列表
- **THEN** 每个元素包含 `labels` 字段

**Acceptance Criteria**:
- `labels` - 节点标签列表 (TypeList, Computed)，映射到 API 的 `DBCustomClusterNodeConfig.Labels`
- 该字段可能返回 null，需进行 nil 检查，为 null 时设为空列表
- 每个 label 元素为 `schema.Resource`，包含：
  - `key` - 标签键 (TypeString, Computed)，映射到 `Label.Key`
  - `value` - 标签值 (TypeString, Computed)，映射到 `Label.Value`

#### Scenario: 返回节点污点信息
- **WHEN** 查询到节点配置列表
- **THEN** 每个元素包含 `taints` 字段

**Acceptance Criteria**:
- `taints` - 节点污点列表 (TypeList, Computed)，映射到 API 的 `DBCustomClusterNodeConfig.Taints`
- 该字段可能返回 null，需进行 nil 检查，为 null 时设为空列表
- 每个 taint 元素为 `schema.Resource`，包含：
  - `key` - 污点键 (TypeString, Computed)，映射到 `Taint.Key`
  - `value` - 污点值 (TypeString, Computed)，映射到 `Taint.Value`，可为空
  - `effect` - 污点效果 (TypeString, Computed)，映射到 `Taint.Effect`，枚举值：NoSchedule / PreferNoSchedule / NoExecute

### Requirement: 结果输出文件支持
数据源必须支持将查询结果输出到文件，与同包其他数据源保持一致。

**Rationale**: 用户可能需要将查询结果保存到文件供后续处理或集成使用。

#### Scenario: 输出结果到文件
- **WHEN** 用户指定了 `result_output_file` 参数
- **THEN** 查询结果被写入指定文件

**Acceptance Criteria**:
- `result_output_file` 为 Optional 参数，类型 TypeString
- 使用 `tccommon.WriteToFile` 将 `d`（ResourceData）写入指定路径
- 未指定时不执行文件写入

### Requirement: 数据源 ID 与重试处理
数据源 Read 函数必须遵循 `RESOURCE_KIND_DATASOURCE` 的空响应与重试约定。

**Rationale**: 腾讯云 API 可能存在短暂波动，直接清空本地 state 的 id 会造成数据丢失；应让外层 retry 继续尝试并以"重试耗尽"形式失败，便于人工介入。

#### Scenario: 成功读取后设置 ID
- **WHEN** API 返回非空 `NodeSet`
- **THEN** 数据源设置一个生成的 token 作为 id

**Acceptance Criteria**:
- 读取成功后调用 `d.SetId(helper.BuildToken())`

#### Scenario: API 返回空响应时返回 NonRetryableError
- **WHEN** `result == nil || result.Response == nil || len(result.Response.NodeSet) == 0`
- **THEN** 返回 `resource.NonRetryableError`，不调用 `d.SetId("")`

**Acceptance Criteria**:
- 在 retry 块内检查空响应，返回 `resource.NonRetryableError(fmt.Errorf("..."))`
- 不直接清空 id
- 在外层 retry 失败路径上保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

#### Scenario: API 调用失败时重试
- **WHEN** API 调用返回 error
- **THEN** 使用 `tccommon.RetryError` 包装错误并重试

**Acceptance Criteria**:
- 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 API 调用
- 错误通过 `tccommon.RetryError(e)` 返回
- 设置 id 等成功操作放在 retry 块外、retry 错误处理后

### Requirement: 服务层封装
必须在 `service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomClusterNodeConfigByFilter` 方法封装云 API 调用。

**Rationale**: 将云 API 调用、重试、限流和错误日志集中到服务层，使数据源 Read 函数保持为薄映射层。

#### Scenario: 封装 DescribeDBCustomClusterNodeConfig 调用
- **WHEN** 数据源 Read 调用服务层方法
- **THEN** 服务层方法封装 API 调用并返回 `[]*dbdcv20201029.DBCustomClusterNodeConfig`

**Acceptance Criteria**:
- 方法签名：`func (me *DbdcService) DescribeDBCustomClusterNodeConfigByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.DBCustomClusterNodeConfig, errRet error)`
- 使用 `dbdcv20201029.NewDescribeDBCustomClusterNodeConfigRequest()` 构造请求
- 从 `param` map 读取 `ClusterId`（`*string`）和 `NodeIds`（`[]*string`）设置到 request
- 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `ratelimit.Check` 和 `me.client.UseDbdcV20201029Client().DescribeDBCustomClusterNodeConfig(request)`
- defer 中当 `errRet != nil` 时打印 `[CRITAL]` 级别日志
- 成功时打印 `[DEBUG]` 级别日志

### Requirement: Provider 注册与文档
数据源必须在 `provider.go` 中注册，并提供对应的文档示例文件。

**Rationale**: 未在 provider 中注册的数据源无法被 Terraform CLI 识别和使用；缺少文档示例会降低用户体验。

#### Scenario: 在 provider.go 中注册数据源
- **WHEN** 实现完成
- **THEN** `provider.go` 的 DataSourcesMap 中包含 `tencentcloud_dbdc_db_custom_cluster_node_config_list`

**Acceptance Criteria**:
- 在 `tencentcloud/provider.go` 的 `dataSources` map 中添加 `"tencentcloud_dbdc_db_custom_cluster_node_config_list": dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeConfigList()`
- 在 `tencentcloud/provider.md` 中添加对应的文档条目（通过 `make doc` 生成）

#### Scenario: 创建数据源文档示例文件
- **WHEN** 实现完成
- **THEN** 存在 `data_source_tc_dbdc_db_custom_cluster_node_config_list.md` 文件

**Acceptance Criteria**:
- 文件位于 `tencentcloud/services/dbdc/` 目录
- 文件以一句话描述开头："Use this data source to query ..." 并带上所属云产品名称 dbdc
- 包含 Example Usage 部分
- 不包含 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）
