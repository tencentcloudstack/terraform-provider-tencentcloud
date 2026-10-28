## Context

Terraform Provider for TencentCloud 需要新增对云产品 VectorDB (VDB) 实例的管理能力。VDB SDK 位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616`，provider 中通过 `UseVdbV20230616Client()` 获取客户端实例。

当前 provider.go 已导入 vdb service 包并注册了 `tencentcloud_vdb_instance`（`vdb.ResourceTencentCloudVdbInstance()`），connectivity 层已有 `UseVdbV20230616Client()` 方法。vendor 中的 models.go 已包含所有所需 API 的 request/response struct。

VDB 实例的创建为**异步操作**：调用 `CreateInstance` 后返回 `InstanceIds`，但实例需要一段时间才能变为 `online` 可用状态。删除操作分两步：先 `IsolateInstance` 隔离，再 `DestroyInstances` 彻底销毁。

代码风格严格参考 `tencentcloud_igtm_monitor` 资源（`tencentcloud/services/igtm/resource_tc_igtm_monitor.go`）。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_vdb_instance` 资源的完整 CRUD 生命周期管理
- Create: 调用 `CreateInstance`，设置 `d.SetId()` 后轮询等待实例 `online`
- Read: 调用 `DescribeInstances` 查询实例并回填 state；调用 `DescribeInstanceNodes` 查询节点信息
- Update: 支持 `ScaleUpInstance` (cpu/memory/disk_size) 和 `ScaleOutInstance` (worker_node_num) 规格变更；支持 `ModifyDBInstanceSecurityGroups` 安全组更新
- Delete: 先 `IsolateInstance` 隔离，`force_delete` 为 true 时继续 `DestroyInstances` 彻底销毁
- 所有 API 调用使用 retry 重试机制
- 正确处理空指针（nil 检查）场景
- 查询接口分页字段使用最大值

**Non-Goals:**
- 不支持实例的参数模板修改（ModifyInstanceMaintenanceWindow 等）
- 不支持实例恢复（RecoverInstance）
- 不实现 datasource 类型资源

## Decisions

### 1. 资源 ID 设计
- 使用 `CreateInstance` 返回的 `InstanceIds[0]` 作为资源 ID（单个实例创建场景）。
- 不使用复合 ID，因为实例管理通过单一 `InstanceId` 标识。

### 2. Schema 字段组织
- **必填入参**: `vpc_id`, `subnet_id`, `pay_mode`, `security_group_ids`
- **可选入参**: `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags` (list of {tag_key, tag_value}), `instance_type`, `mode`, `product_type`, `node_type`, `cpu`, `memory`, `disk_size`, `worker_node_num`
- **ForceNew 字段**: `vpc_id`, `subnet_id`, `instance_type`, `mode`, `product_type`, `node_type` — 这些字段变更需要重建实例
- **可更新字段**: `cpu`, `memory`, `disk_size` (通过 ScaleUpInstance), `worker_node_num` (通过 ScaleOutInstance), `security_group_ids` (通过 ModifyDBInstanceSecurityGroups), `pay_mode`, `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags` — 标记为 immutable，变更时拒绝
- **Computed 字段**: `status`, `region`, `zone`, `product`, `shard_num`, `api_version`, `extend`, `expired_at`, `is_no_expired`, `wan_address`, `isolate_at`, `task_status`, `networks` (list of {vpc_id, subnet_id, vip, port, preserve_duration, expire_time, is_ssl}), `created_at`, `engine_name`, `engine_version`, `nodes` (list of {name, status})
- 声明 `Timeouts` 块：Create 30min, Update 30min, Delete 10min

### 3. Create 流程（异步创建）
1. 构建 `CreateInstanceRequest`，填充所有入参
2. 在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中调用 `CreateInstanceWithContext`
3. 检查 `response.Response.InstanceIds` 非空（nil 或 len==0 返回 NonRetryableError）
4. **设置 `d.SetId(instanceId)`** — 在异步等待之前
5. 轮询 `DescribeInstances` 等待 status == `online`

### 4. Read 流程
1. 通过 service 层 `DescribeVdbInstanceById` 查询单个实例
2. 若实例不存在，打印 `[CRUD]` 日志后 `d.SetId("")`
3. 回填所有 computed 字段，每个字段先判断非 nil
4. 回填 networks 子结构（循环 list）
5. 回填 resource_tags 子结构
6. 回填 security_group_ids
7. 额外调用 `DescribeInstanceNodes` 查询节点信息，回填 nodes（失败时仅 warn，不中断）

### 5. Update 流程
- **immutableFields 检查**: `pay_mode`, `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags`, `instance_type`, `mode`, `product_type`, `node_type` — 这些字段变更时返回 error
- **Step 1 ScaleUp**: 若 `cpu`/`memory`/`disk_size` 变更，调用 `ScaleUpInstance`，然后轮询等待规格生效
- **Step 2 ScaleOut**: 若 `worker_node_num` 变更，调用 `ScaleOutInstance`，然后轮询等待生效
- **Step 3 SecurityGroups**: 若 `security_group_ids` 变更，调用 `ModifyDBInstanceSecurityGroups`，然后轮询等待匹配

### 6. Delete 流程
1. 调用 `IsolateInstance` 隔离实例
2. 轮询等待 status == `isolated`
3. 若 `force_delete` == true，调用 `DestroyInstances` 彻底销毁
4. 轮询等待实例不存在（WaitForInstanceNotFound）

### 7. Service 层设计
- `VdbService` 封装 SDK 调用与轮询逻辑
- `VdbClientInterface` 接口抽象 SDK client 方法，用于单元测试 mock
- 方法: `DescribeVdbInstanceById`, `DescribeVdbInstanceNodesById`, `DescribeDBSecurityGroupsByInstanceId`, `WaitForInstanceStatus`, `WaitForInstanceScaleUp`, `WaitForInstanceScaleOut`, `WaitForInstanceNotFound`, `WaitForSecurityGroupsMatch`

### 8. 单元测试策略
- 使用 mock（testify/mock）方式，不使用 terraform 测试套件
- MockVdbClient 实现 VdbClientInterface，mock 所有 SDK 方法
- 测试覆盖: DescribeVdbInstanceById (found/notfound/nil/api error/nil response), DescribeVdbInstanceNodesById, DescribeDBSecurityGroupsByInstanceId, WaitFor* 系列, Schema 验证

## Risks / Trade-offs

- **[异步创建超时]** → 设置 Create timeout 30min，使用 `resource.Retry` 轮询，超时后 tfstate 中仍保留 id，用户可通过 `terraform destroy` 清理
- **[实例创建返回多个 InstanceIds]** → 始终取 `InstanceIds[0]`，当前资源每次创建一个实例（GoodsNum=1）
- **[ScaleUp 与 ScaleOut 串行执行]** → 先 ScaleUp 完成再 ScaleOut，避免并发操作冲突
- **[删除操作不可逆]** → 默认 `force_delete`=false 仅隔离，用户可手动恢复；`force_delete`=true 时彻底销毁不可恢复
- **[安全组更新为全量替换]** → `ModifyDBInstanceSecurityGroups` 为全量替换，需传入完整安全组列表

## Cloud API 关键字段定义（vendor 摘录）

以下为实施阶段所需的 cloud API request/response struct 关键字段定义（从 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616/models.go` 行区间摘录）。

### CreateInstanceRequest (models.go:206-301)
```go
type CreateInstanceRequest struct {
    *tchttp.BaseRequest
    VpcId            *string   // 私有网络 ID
    SubnetId         *string   // 子网 ID
    PayMode         *int64    // 0:按量付费 1:包年包月
    InstanceName     *string
    SecurityGroupIds []*string
    PayPeriod        *int64    // 包年包月续费时长(月)
    AutoRenew        *int64    // 0:不自动续费 1:自动续费
    Params            *string   // 实例额外参数,json
    ResourceTags     []*Tag
    Project           *string
    ProductType       *int64    // 0:标准版 1:容量增强版
    InstanceType      *string   // base/single/cluster
    Mode              *string   // two/three (高可用版)
    GoodsNum          *int64    // 购买数量
    NetworkType      *string   // VPC或Tcs
    TemplateId       *string
    Components        []*CreateInstancesComponent
    Zone              *string
    SlaveZones        []*string
    IsNoExpired      *bool
    EngineName       *string
    EngineVersion     *string
    Brief             *string
    Chief             *string
    DBA               *string
    NodeType          *string   // compute/normal/store
    Cpu               *uint64
    Memory            *uint64
    DiskSize         *uint64
    WorkerNodeNum    *uint64
    EnableEncryption *bool
}
```

### CreateInstancesComponent (models.go:377-401)
```go
type CreateInstancesComponent struct {
    Component    *string
    Cpu          *float64
    Memory       *float64
    ReplicaNum  *uint64
    Name         *string
    StorageSize  *uint64
    StorageType *string
    Params        *string
}
```

### CreateInstanceResponse (models.go:353-364)
```go
type CreateInstanceResponseParams struct {
    InstanceIds []*string  // 实例 ID
    RequestId   *string
}
type CreateInstanceResponse struct {
    *tchttp.BaseResponse
    Response *CreateInstanceResponseParams
}
```

### DescribeInstancesRequest (models.go:663-713)
```go
type DescribeInstancesRequest struct {
    *tchttp.BaseRequest
    InstanceIds      []*string
    InstanceNames    []*string
    InstanceKeys     []*string
    Status           []*string
    EngineNames      []*string
    EngineVersions    []*string
    ApiVersions      []*string
    CreateAt          *string
    Zones             []*string
    OrderBy          *string
    OrderDirection    *string
    Offset            *uint64
    Limit             *uint64   // 列表查询数量, 最大值 100
    ResourceTags      []*Tag
    TaskStatus        []*int64
    Networks          []*string
}
```

### DescribeInstancesResponse (models.go:749-764)
```go
type DescribeInstancesResponseParams struct {
    Items       []*InstanceInfo
    TotalCount *uint64
    RequestId   *string
}
```

### InstanceInfo (models.go:1235-1349)
```go
type InstanceInfo struct {
    InstanceId       *string
    Name             *string
    AppId            *uint64
    Region           *string
    Zone             *string
    Product          *string
    Networks         []*Network
    ShardNum         *uint64
    ReplicaNum      *uint64
    Cpu              *float64
    Memory           *float64
    Disk             *uint64
    HealthScore     *float64  // Deprecated
    Warning           *int64   // Deprecated
    Project           *string  // Deprecated
    ResourceTags      []*Tag
    CreatedAt        *string
    Status           *string
    EngineName       *string
    EngineVersion     *string
    ApiVersion        *string
    PayMode          *int64
    Extend            *string  // json
    ExpiredAt         *string
    IsNoExpired       *bool
    ProductType       *int64
    InstanceType     *string
    NodeType          *string
    WanAddress       *string
    IsolateAt        *string
    AutoRenew         *int64
    TaskStatus        *int64
    SecurityGroupIds  []*string
    UpgradeVersion    *string
    IsInternal         *bool
}
```

### Network (models.go:1540-1561)
```go
type Network struct {
    VpcId            *string
    SubnetId         *string
    Vip              *string
    Port             *int64
    PreserveDuration *int64
    ExpireTime       *string
    IsSSL             *bool
}
```

### Tag (models.go:1860-1866)
```go
type Tag struct {
    TagKey   *string
    TagValue *string
}
```

### ScaleUpInstanceRequest (models.go:1761-1787)
```go
type ScaleUpInstanceRequest struct {
    *tchttp.BaseRequest
    InstanceId   *string
    Cpu          *float64
    Memory       *float64
    StorageSize  *uint64
    RunNow       *bool
}
```

### ScaleOutInstanceRequest (models.go:1676-1689)
```go
type ScaleOutInstanceRequest struct {
    *tchttp.BaseRequest
    InstanceId  *string
    ReplicaNum  *uint64  // 扩容后的节点数量
    RunNow      *bool
}
```

### IsolateInstanceRequest (models.go:1357-1362)
```go
type IsolateInstanceRequest struct {
    *tchttp.BaseRequest
    InstanceId *string
}
```

### DestroyInstancesRequest (models.go:1100-1105)
```go
type DestroyInstancesRequest struct {
    *tchttp.BaseRequest
    InstanceIds []*string
}
```

### ModifyDBInstanceSecurityGroupsRequest (models.go:1415-1424)
```go
type ModifyDBInstanceSecurityGroupsRequest struct {
    *tchttp.BaseRequest
    SecurityGroupIds []*string  // 全量替换
    InstanceIds      []*string
}
```

### DescribeInstanceNodesRequest (models.go:542-556)
```go
type DescribeInstanceNodesRequest struct {
    *tchttp.BaseRequest
    InstanceId *string
    Limit       *int64
    Offset      *int64
    Component   *string
}
```

### NodeInfo (models.go:1563-1569)
```go
type NodeInfo struct {
    Name   *string  // Pod名称
    Status *string  // Pod状态
}
```

## 代码风格参照文件说明（tencentcloud_igtm_monitor）

参照文件: `tencentcloud/services/igtm/resource_tc_igtm_monitor.go`

### schema 组织结构
```go
func ResourceTencentCloudIgtmMonitor() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudIgtmMonitorCreate,
        Read:   resourceTencentCloudIgtmMonitorRead,
        Update: resourceTencentCloudIgtmMonitorUpdate,
        Delete: resourceTencentCloudIgtmMonitorDelete,
        Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
        Schema: map[string]*schema.Schema{ ... },
    }
}
```

### CRUD 骨架结构
- Create: `defer tccommon.LogElapsed()` + `defer tccommon.InconsistentCheck()` → 构建 request → `resource.Retry(WriteRetryTimeout)` 调用 API → 检查 nil → 提取 ID → `d.SetId()` → 调用 Read
- Read: 构建 ctx → 调用 service 层查询 → 判空 `d.SetId("")` → 逐字段 nil 检查后 `d.Set()`
- Update: 声明 immutableArgs/mutableArgs → 遍历 `d.HasChange()` → 构建 request 调用 API → 调用 Read
- Delete: 构建 request → `resource.Retry(WriteRetryTimeout)` 调用 API

### 常用 helper
- `tccommon.GetLogId(tccommon.ContextNil)`
- `tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`
- `helper.String()`, `helper.IntInt64()`, `helper.IntUint64()`
- `tccommon.RetryError(e)`