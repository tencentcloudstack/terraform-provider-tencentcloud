## Context

腾讯云集团账号管理的共享单元(Share Unit)支持将组织部门添加到共享单元中,从而实现资源共享的精细化控制。目前 Provider 已有 `tencentcloud_organization_org_share_unit_member_v2` 资源用于管理共享单元成员,现在需要添加共享单元部门管理的支持。

**Current State:**
- 已有共享单元成员管理资源(`resource_tc_organization_org_share_unit_member_v2.go`)
- 已有共享单元查询数据源(`data_source_tc_organization_org_share_units.go`)
- 服务层已实现组织相关的 API 调用模式
- SDK 已支持 AddShareUnitNode, DeleteShareUnitNode, DescribeShareUnitNodes API

**Constraints:**
- 必须遵循现有的资源命名规范: `tencentcloud_organization_org_share_unit_node`
- 必须保持与现有组织资源的一致性(ID 复合模式、错误处理等)
- API 请求频率限制: 20次/秒
- 共享单元部门数量有上限限制(LimitExceeded.ShareUnitNodeOverLimit 错误)

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_organization_org_share_unit_node` 资源,支持添加/删除共享单元部门
- 实现 `tencentcloud_organization_org_share_unit_nodes` 数据源,支持查询共享单元的部门列表
- 提供完整的 CRUD 操作和 Import 支持
- 添加验收测试和完整文档
- 遵循现有代码模式和最佳实践

**Non-Goals:**
- 不修改现有共享单元成员资源的实现
- 不支持批量添加部门(API 不支持批量操作)
- 不支持共享单元本身的创建/删除(已有其他资源)

## Decisions

### Decision 1: 资源 ID 设计
**Choice:** 使用复合 ID: `{UnitId}#{NodeId}`

**Rationale:**
- 符合现有资源的 ID 模式(如 `organization_org_share_unit_member_v2` 使用 `{UnitId}#{Area}`)
- 唯一标识共享单元中的部门
- 便于 Import 操作
- NodeId 是 Integer 类型,但在复合 ID 中转换为 string

**Alternatives Considered:**
- 使用 UnitId 作为 ID,NodeId 作为属性: 不符合 Terraform 资源单实例原则
- 使用自增 ID: 不能反映实际的业务标识

### Decision 2: Schema 设计
**Resource Schema:**
```go
{
  "unit_id": {
    Type: schema.TypeString,
    Required: true,
    ForceNew: true,  // 修改需要重建
  },
  "node_id": {
    Type: schema.TypeInt,
    Required: true,
    ForceNew: true,  // 修改需要重建
  },
}
```

**Data Source Schema:**
```go
{
  "unit_id": {
    Type: schema.TypeString,
    Required: true,
  },
  "offset": {
    Type: schema.TypeInt,
    Optional: true,
    Default: 0,
  },
  "limit": {
    Type: schema.TypeInt,
    Optional: true,
    Default: 10,
  },
  "search_key": {
    Type: schema.TypeString,
    Optional: true,
  },
  "items": {
    Type: schema.TypeList,
    Computed: true,
    Elem: &schema.Resource{
      Schema: map[string]*schema.Schema{
        "share_node_id": schema.TypeInt,
        "create_time": schema.TypeString,
      },
    },
  },
}
```

**Rationale:**
- 所有字段都是 ForceNew,因为 API 不支持更新操作
- 数据源支持分页和搜索,与其他组织数据源保持一致
- 返回字段与 API 响应字段对应

### Decision 3: 服务层方法设计
在 `service_tencentcloud_organization.go` 中添加:

```go
// 添加共享单元部门
func (me *OrganizationService) AddOrganizationOrgShareUnitNodeById(ctx context.Context, unitId string, nodeId int64) (errRet error)

// 删除共享单元部门
func (me *OrganizationService) DeleteOrganizationOrgShareUnitNodeById(ctx context.Context, unitId string, nodeId int64) (errRet error)

// 查询共享单元部门
func (me *OrganizationService) DescribeOrganizationOrgShareUnitNodeById(ctx context.Context, unitId string, nodeId int64) (shareUnitNode *organization.ShareUnitNode, errRet error)

// 查询共享单元部门列表(数据源用)
func (me *OrganizationService) DescribeOrganizationOrgShareUnitNodesByFilter(ctx context.Context, param map[string]interface{}) (shareUnitNodes []*organization.ShareUnitNode, errRet error)
```

**Rationale:**
- 符合现有服务层方法命名规范
- `ById` 方法用于资源的 Read 操作
- `ByFilter` 方法用于数据源查询,支持灵活的参数组合

### Decision 4: 错误处理策略
**Approach:**
- 使用 `helper.Retry()` 包装 API 调用,处理最终一致性
- 在 Read 操作中,如果资源不存在,设置 `d.SetId("")` 而不是返回错误
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行日志和一致性检查

**Special Error Handling:**
- `LimitExceeded.ShareUnitNodeOverLimit`: 共享单元部门超过上限,应该在 Create 时返回明确错误
- `FailedOperation.ShareNodeNotExist`: 部门不存在,在 Read/Delete 时应优雅处理
- `ResourceNotFound.OrganizationNodeNotExist`: 组织节点不存在,Create 时应返回错误

### Decision 5: Import 支持
**Format:** `terraform import tencentcloud_organization_org_share_unit_node.foo {unit_id}#{node_id}`

**Implementation:**
```go
Importer: &schema.ResourceImporter{
  State: schema.ImportStatePassthrough,
},
```

在 Read 函数中解析复合 ID:
```go
idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
if len(idSplit) != 2 {
  return fmt.Errorf("id is broken,%s", d.Id())
}
unitId := idSplit[0]
nodeId, _ := strconv.ParseInt(idSplit[1], 10, 64)
```

## Risks / Trade-offs

### Risk 1: API 限制共享单元部门数量
**Mitigation:**
- 在文档中明确说明限制
- 捕获 `LimitExceeded.ShareUnitNodeOverLimit` 错误并返回清晰的错误消息
- 用户可以通过数据源查询当前部门数量

### Risk 2: NodeId 类型转换
**Issue:** NodeId 在 API 中是 Integer,但在复合 ID 中需要转换为 string
**Mitigation:**
- 在 Read 函数中使用 `strconv.ParseInt` 进行转换
- 添加错误处理确保转换安全
- 在测试中验证 ID 解析逻辑

### Risk 3: 部门可能在外部被删除
**Issue:** 部门可能通过控制台或其他方式被删除,导致 Terraform state 不一致
**Mitigation:**
- Read 操作中检查部门是否存在,不存在则清空 state
- 使用 `defer tccommon.InconsistentCheck()` 检测不一致性
- 文档中说明 drift detection 的行为

### Trade-off: 不支持批量操作
**Reason:** 腾讯云 API 不支持批量添加/删除部门
**Impact:** 用户需要为每个部门创建单独的资源实例
**Benefit:** 每个部门是独立的 Terraform 资源,状态管理更清晰

## Migration Plan

### Deployment Steps:
1. 添加服务层方法到 `service_tencentcloud_organization.go`
2. 实现资源文件 `resource_tc_organization_org_share_unit_node.go`
3. 实现数据源文件 `data_source_tc_organization_org_share_unit_nodes.go`
4. 在 `provider.go` 中注册新资源和数据源
5. 添加测试文件和文档
6. 运行验收测试确保功能正确
7. 运行 `make doc` 生成 website 文档

### Testing Strategy:
- 单元测试: 验证 ID 解析和复合逻辑
- 验收测试:
  - 创建共享单元部门
  - 查询数据源验证创建成功
  - 删除共享单元部门
  - Import 资源验证 ID 解析
  - 测试错误场景(不存在的部门、超过限制等)

### Rollback:
- 如果发现问题,可以在 `provider.go` 中注释掉资源注册
- 不影响现有资源,因为是全新添加

## Open Questions

1. ~~是否需要支持批量操作?~~ → 不需要,API 不支持且单资源模式更符合 Terraform 最佳实践
2. ~~复合 ID 分隔符使用什么?~~ → 使用 `tccommon.FILED_SP` (`#`),与现有资源保持一致
3. ~~是否需要支持 Update 操作?~~ → 不需要,API 只支持添加和删除,所有字段都是 ForceNew
