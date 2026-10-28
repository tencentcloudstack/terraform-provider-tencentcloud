## MODIFIED Requirements

### Requirement: 完整的节点信息映射
数据源 SHALL 返回节点的完整详细信息，涵盖所有 `DBCustomClusterNode` 结构体中的字段。

**Rationale**: 用户需要完整的节点信息用于资源规划和管理决策，包括节点标识、网络配置、状态和规格。

#### Scenario: 返回节点基础信息字段
- **WHEN** 查询到节点列表
- **THEN** 每个 `node_set` 元素 SHALL 包含以下字段

**Acceptance Criteria**:
- `node_id` - 节点ID (TypeString, Computed)
- `node_name` - 节点名称 (TypeString, Computed)
- `lan_ip` - 节点内网IP地址 (TypeString, Computed)
- `ssh_endpoint` - 节点SSH访问Endpoint，格式为IP:Port (TypeString, Computed)
- `status` - 节点在集群中的实例状态 (TypeString, Computed)
- `zone` - 节点所属地域 (TypeString, Computed)
- `node_type` - 节点类型 (TypeString, Computed，枚举值包括 DB.AT5.32XLARGE512, DB.AT5.64XLARGE1152, DB.AT5.128XLARGE2304, DB.AT5.16XLARGE256, DB.AT5.8XLARGE128)
- `network_mode` - 网络模式 (TypeString, Computed，枚举值包括 `privatelink` 四层网络联通放通SSH通路、`cross_tenant_eni` 三层网络联通双网卡模式)，映射到 API 的 `response.NodeSet.NetworkMode`
- `eni_ip` - 当网络模式为 `cross_tenant_eni` 时节点的可访问 IP 地址 (TypeString, Computed)，映射到 API 的 `response.NodeSet.EniIP`

#### Scenario: 返回总数量
- **WHEN** 查询到节点列表
- **THEN** 数据源 SHALL 返回集群下总的节点数量

**Acceptance Criteria**:
- `total_count` - 集群下总的节点数量 (TypeInt, Computed)
- 映射到 API 的 `Response.TotalCount`

#### Scenario: 安全设置 network_mode 和 eni_ip 字段
- **WHEN** API 返回的 `DBCustomClusterNode` 中 `NetworkMode` 或 `EniIP` 为 nil
- **THEN** 数据源 MUST NOT 调用对应的 set 方法，字段保持零值

**Acceptance Criteria**:
- 在设置 `network_mode` 前判断 `node.NetworkMode != nil`，为 nil 时跳过
- 在设置 `eni_ip` 前判断 `node.EniIP != nil`，为 nil 时跳过
- 遵循项目规则 #8 的 nil 检查要求
- `NetworkMode` 和 `EniIP` 在 SDK 中标注"可能返回 null"，nil 处理必不可少
