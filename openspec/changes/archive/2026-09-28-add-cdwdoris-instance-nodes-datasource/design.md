## Context

当前 Terraform Provider for TencentCloud 的 CDWDoris 服务已具备数据源 `tencentcloud_cdwdoris_instances`（查询集群列表）、资源 `tencentcloud_cdwdoris_instance`、`tencentcloud_cdwdoris_workload_group`、`tencentcloud_cdwdoris_user` 等，但缺少查询某集群下节点明细的数据源。

云 API `DescribeInstanceNodes`（包名 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228`）已存在于 vendor 中，可按 `InstanceId` 获取集群节点信息列表，返回每个节点的 IP、规格、CPU/内存/磁盘、角色、可用区等信息，以及集群支持的节点角色列表。本设计基于该接口新增数据源 `tencentcloud_cdwdoris_instance_nodes`。

数据源为只读（RESOURCE_KIND_DATASOURCE），仅实现 Read 操作，无状态、不修改云上资源，完全向后兼容。

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_cdwdoris_instance_nodes`，支持按 `instance_id`（必填）、`node_role`（可选）、`display_policy`（可选）查询集群节点列表。
- 输出字段对齐云 API `DescribeInstanceNodes` 响应：`instance_nodes_list`（节点列表，每项含 ip/spec/core/memory/disk_type/disk_size/role/status/rip/fe_role/uuid/zone/virtual_zone/create_time/compute_group_id）及 `node_roles`（集群支持的节点角色列表）。
- 在 service 层封装 `DescribeCdwdorisInstanceNodesByFilter`，内部自动分页拉取全部节点（不向用户暴露 limit/offset）。
- 在 provider.go 注册数据源、补充 provider.md 文档条目、生成 `.md` 说明文件。
- 数据源 Read 方法的 retry 块内对空响应返回 NonRetryableError，避免因 API 短暂波动清空 state。

**Non-Goals:**
- 不新增/修改任何资源（RESOURCE_KIND_GENERAL 等），不涉及 CUD 操作。
- 不暴露 limit/offset 分页参数给用户。
- 不修改现有 `tencentcloud_cdwdoris_instances` 数据源的 schema 或行为。

## Decisions

### 1. 接口选择：DescribeInstanceNodes

vendor 中存在三个相近接口：`DescribeInstanceNodes`、`DescribeInstanceNodesInfo`、`DescribeInstanceNodesRole`。需求明确指定使用 `DescribeInstanceNodes`，其 request 含 `InstanceId`/`NodeRole`/`DisplayPolicy` 等参数，response 含 `InstanceNodesList []*InstanceNode` 与 `NodeRoles []*string`，与需求出参映射完全一致，确认采用此接口。

### 2. 参数映射与 schema 组织

入参（request）：
- `request.InstanceId`（`*string`，必填）→ schema `instance_id`
- `request.NodeRole`（`*string`，可选）→ schema `node_role`
- `request.DisplayPolicy`（`*string`，可选）→ schema `display_policy`

出参（response）：
- `response.Response.InstanceNodesList`（`[]*InstanceNode`）→ schema `instance_nodes_list`（TypeList，Computed），elem 为 Resource，下展字段均为 Computed：
  - `Ip *string` → `ip`
  - `Spec *string` → `spec`
  - `Core *int64` → `core`
  - `Memory *int64` → `memory`
  - `DiskType *string` → `disk_type`
  - `DiskSize *int64` → `disk_size`
  - `Role *string` → `role`
  - `Status *string` → `status`
  - `Rip *string` → `rip`
  - `FeRole *string` → `fe_role`
  - `UUID *string` → `uuid`
  - `Zone *string` → `zone`
  - `VirtualZone *string` → `virtual_zone`
  - `CreateTime *string` → `create_time`
  - `ComputeGroupId *string` → `compute_group_id`
- `response.Response.NodeRoles`（`[]*string`）→ schema `node_roles`（TypeSet/List，Computed，elem TypeString）

额外提供 `result_output_file`（Optional）用于保存结果，与同服务现有数据源一致。

### 3. service 层自动分页

`DescribeInstanceNodes` request 含 `Offset`/`Limit` 分页参数。service 方法 `DescribeCdwdorisInstanceNodesByFilter` 内部循环分页，每页取上限值 `Limit=10`（云 API 注释标注默认步长为 10），累计返回全部节点列表，不在 schema 暴露分页参数，符合数据源分页规范。

### 4. Read 函数与 retry 处理

参照 `tencentcloud_igtm_instance_list` 与 `tencentcloud_cdwdoris_instances` 数据源风格：
- `defer tccommon.LogElapsed()` + `defer tccommon.InconsistentCheck()`
- `resource.Retry(tccommon.ReadRetryTimeout, ...)`，错误用 `tccommon.RetryError(e)` 包装
- retry 块内对空响应（`response == nil || response.Response == nil || len(InstanceNodesList) == 0`），**直接返回 `NonRetryableError`**，不 `d.SetId("")`，并在外层失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示
- 拿到数据后，按字段 nil 检查逐个 set 到 map，再 `d.Set("instance_nodes_list", list)`、`d.Set("node_roles", nodeRolesList)`
- `d.SetId(helper.BuildToken())`
- 支持 `result_output_file` 输出

### 5. provider 注册与文档

- `provider.go` 在 dataSources map 中注册：`"tencentcloud_cdwdoris_instance_nodes": cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes(),`
- `provider.md` 补充数据源条目（由 `make doc` 在收尾阶段统一生成，本阶段不直接改 website/）
- 新增 `data_source_tc_cdwdoris_instance_nodes.md`：一句话描述（带产品名 CDWDoris）+ Example Usage。

## 参考信息（供实施阶段直接使用，避免回读 vendor 大文件）

### 云 API struct 定义（摘自 vendor/.../cdwdoris/v20211228/models.go）

#### DescribeInstanceNodesRequestParams（models.go:2806）
```go
type DescribeInstanceNodesRequestParams struct {
	InstanceId    *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`    // 集群实例ID
	NodeRole      *string `json:"NodeRole,omitnil,omitempty" name:"NodeRole"`        // 集群角色类型，默认为 "data"数据节点
	Offset        *int64  `json:"Offset,omitnil,omitempty" name:"Offset"`            // 分页参数，第一页为0，第二页为10
	Limit         *int64  `json:"Limit,omitnil,omitempty" name:"Limit"`              // 分页参数，分页步长，默认为10
	DisplayPolicy *string `json:"DisplayPolicy,omitnil,omitempty" name:"DisplayPolicy"` // 展现策略，All时显示所有
}
// DescribeInstanceNodesRequest 嵌入 *tchttp.BaseRequest，字段同上
```

#### DescribeInstanceNodesResponseParams（models.go:2866）
```go
type DescribeInstanceNodesResponseParams struct {
	TotalCount        *int64         `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	InstanceNodesList []*InstanceNode `json:"InstanceNodesList,omitnil,omitempty" name:"InstanceNodesList"`
	NodeRoles         []*string      `json:"NodeRoles,omitnil,omitempty" name:"NodeRoles"`
	RequestId         *string        `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
// DescribeInstanceNodesResponse { *tchttp.BaseResponse; Response *DescribeInstanceNodesResponseParams }
```

#### InstanceNode（models.go:4589）—— instance_nodes_list 元素
```go
type InstanceNode struct {
	Ip             *string `json:"Ip,omitnil,omitempty" name:"Ip"`               // IP地址
	Spec           *string `json:"Spec,omitnil,omitempty" name:"Spec"`           // 机型，如 S1
	Core           *int64  `json:"Core,omitnil,omitempty" name:"Core"`           // cpu核数
	Memory         *int64  `json:"Memory,omitnil,omitempty" name:"Memory"`       // 内存大小
	DiskType       *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`   // 磁盘类型
	DiskSize       *int64  `json:"DiskSize,omitnil,omitempty" name:"DiskSize"`   // 磁盘大小
	Role           *string `json:"Role,omitnil,omitempty" name:"Role"`           // 所属clickhouse cluster名称
	Status         *string `json:"Status,omitnil,omitempty" name:"Status"`       // 状态
	Rip            *string `json:"Rip,omitnil,omitempty" name:"Rip"`             // rip
	FeRole         *string `json:"FeRole,omitnil,omitempty" name:"FeRole"`       // FE节点角色
	UUID           *string `json:"UUID,omitnil,omitempty" name:"UUID"`           // UUID
	Zone           *string `json:"Zone,omitnil,omitempty" name:"Zone"`           // 可用区
	VirtualZone    *string `json:"VirtualZone,omitnil,omitempty" name:"VirtualZone"`          // 虚拟可用区
	CreateTime     *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`            // 创建时间
	ComputeGroupId *string `json:"ComputeGroupId,omitnil,omitempty" name:"ComputeGroupId"`    // 计算组ID
}
```

#### Client 调用方式（client.go:1248）
```go
func (c *Client) DescribeInstanceNodes(request *DescribeInstanceNodesRequest) (response *DescribeInstanceNodesResponse, err error)
// provider 中通过 me.client.UseCdwdorisV20211228Client().DescribeInstanceNodes(request) 调用
```

### 代码风格参照：tencentcloud_igtm_instance_list 数据源

文件 `tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`，关键结构：
- `DataSourceTencentCloudIgtmInstanceList() *schema.Resource`：仅 `Read` 字段；`Schema` map 含入参 + Computed 列表字段 + `result_output_file`。
- `dataSourceTencentCloudIgtmInstanceListRead`：
  - `defer tccommon.LogElapsed(...)` / `defer tccommon.InconsistentCheck(d, meta)`
  - `logId := tccommon.GetLogId(nil)`；`ctx := tccommon.NewResourceLifeCycleHandleFuncContext(...)`
  - `service := IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`
  - `resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {...})`，错误 `tccommon.RetryError(e)`
  - 遍历响应构造 `[]map[string]interface{}`，逐字段 nil 检查后 `d.Set("instance_set", list)`
  - `d.SetId(helper.BuildToken())`
  - `result_output_file` 通过 `tccommon.WriteToFile` 输出

### 代码风格参照：tencentcloud_cdwdoris_instances 数据源（同服务）

文件 `tencentcloud/services/cdwdoris/data_source_tc_cdwdoris_instances.go`：
- 与 igtm 风格一致，service 为 `CdwdorisService{client: ...}`，构造 `paramMap` 传给 `service.DescribeCdwdorisInstancesByFilter(ctx, paramMap)`。
- service 文件 `service_tencentcloud_cdwdoris.go`：`func (me *CdwdorisService) DescribeCdwdorisInstancesByFilter(...)`，内部 `for k, v := range param { ... }` 赋值 request 字段，分页 `offset/limit` 循环，`ratelimit.Check(request.GetAction())`，`me.client.UseCdwdorisV20211228Client().DescribeInstances(request)`。

### 新数据源文件骨架要点

```go
package cdwdoris

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwdorisv20211228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdwdorisInstanceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdwdorisInstanceNodesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": { Type: schema.TypeString, Required: true, ... },
			"node_role":   { Type: schema.TypeString, Optional: true, ... },
			"display_policy": { Type: schema.TypeString, Optional: true, ... },
			"result_output_file": { Type: schema.TypeString, Optional: true, ... },
			"instance_nodes_list": { Type: schema.TypeList, Computed: true, Elem: &schema.Resource{ Schema: map[string]*schema.Schema{...15 fields...} } },
			"node_roles": { Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString} },
		},
	}
}
```

Read 函数 retry 块对空响应返回 `NonRetryableError`；set 完成后 `d.SetId(helper.BuildToken())`。

## Risks / Trade-offs

- [DescribeInstanceNodes 返回字段语义偏 clickhouse cluster] `InstanceNode.Role` 注释为"所属clickhouse cluster名称"，字段名与语义在 Doris 场景下沿用云 API 原样映射，不做语义改写，保证与上游一致 → 直接透传该字段值。
- [分页步长较小] 默认 Limit=10，节点数多时需多次请求；service 层循环分页已覆盖，对用户透明，可接受。
- [空响应处理] 数据源在集群无节点或 API 波动返回空时，按规范返回 NonRetryableError 而非清空 state，避免误删 → 已在 Read 决策中明确。