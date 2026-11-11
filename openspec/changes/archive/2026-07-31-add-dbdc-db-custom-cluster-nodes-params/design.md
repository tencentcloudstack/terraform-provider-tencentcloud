## Context

`tencentcloud_dbdc_db_custom_cluster_nodes` 数据源（RESOURCE_KIND_DATASOURCE）当前通过 `DescribeDBCustomClusterNodes` API 查询 DB Custom 集群节点列表，并在 `node_set` schema 中返回节点信息。现有 schema 已映射 7 个字段（node_id、node_name、lan_ip、ssh_endpoint、status、zone、node_type），但云 API 的 `DBCustomClusterNode` 结构体还包含 `NetworkMode` 和 `EniIP` 两个出参字段尚未映射到 Terraform schema 中。

经核实 vendor 目录下 SDK（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029/models.go`），`DBCustomClusterNode` 结构体已包含：
- `NetworkMode *string`：网络模式，枚举值为 `privatelink`（四层网络联通）和 `cross_tenant_eni`（三层网络联通，双网卡模式）
- `EniIP *string`：当网络模式为 `cross_tenant_eni` 时，节点的可访问 IP 地址

两个字段均标注"注意：此字段可能返回 null，表示取不到有效值"。

## Goals / Non-Goals

**Goals:**
- 在 `node_set` schema 中新增 `network_mode` 和 `eni_ip` 两个 Computed 字段
- 在 Read 方法中新增对应字段的 nil 检查与赋值逻辑，遵循项目规则 #8（设置字段前先判断是否为 nil）
- 保持向后兼容：所有变更为新增 Computed 字段，不影响现有配置和 state

**Non-Goals:**
- 不修改请求参数或 API 调用逻辑（字段已存在于 API 出参中）
- 不修改 service 层分页逻辑
- 不新增独立的 spec capability（仅修改已有 spec）
- 不修改 provider.go / provider.md（数据源已注册）

## Decisions

### 决策 1：字段类型为 TypeString, Computed
**理由**：`NetworkMode` 和 `EniIP` 在 SDK 中均为 `*string` 类型，且为 API 出参，用户不可设置，因此使用 `schema.TypeString` + `Computed: true`。与现有 `node_set` 内其他字段（如 node_id、node_name）保持一致。

**备选方案**：将 `NetworkMode` 定义为枚举校验类型。**否决**：项目中数据源字段通常不做枚举校验，且 API 枚举值可能扩展，保持简单 TypeString 更稳妥。

### 决策 2：使用 MODIFIED Requirements 更新已有 spec
**理由**：`dbdc-db-custom-cluster-nodes-datasource` spec 已存在"完整的节点信息映射"需求，本次在该需求中新增两个字段。根据 openspec 规则，修改已有需求应使用 `## MODIFIED Requirements` 并包含完整的更新后内容。

### 决策 3：遵循 nil 检查模式
**理由**：两个字段在 SDK 中标注"可能返回 null"。根据项目规则 #8，在 Read 方法中设置字段前必须先判断 `node.NetworkMode != nil` 和 `node.EniIP != nil`，为 nil 时不调用 set 方法，与现有字段处理方式一致。

### 决策 4：字段命名使用 snake_case
- `NetworkMode` → `network_mode`
- `EniIP` → `eni_ip`（保持大小写缩写的 snake_case 形式）

## Risks / Trade-offs

- **[风险] 字段可能返回 null** → 通过 nil 检查缓解，nil 时不设置字段，Terraform state 中该字段为零值（空字符串），与现有字段处理一致。
- **[风险] 向后兼容性** → 所有变更为新增 Computed 字段，现有 TF 配置无需修改，无 breaking change。
- **[权衡] 文档生成** → `.md` 文件由 `make doc` 在收尾阶段自动生成，本阶段不手动编写 website/docs/。
