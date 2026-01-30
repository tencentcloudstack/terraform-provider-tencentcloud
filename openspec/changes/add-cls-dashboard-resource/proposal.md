# Proposal: Add CLS Dashboard Resource

## Change ID
`add-cls-dashboard-resource`

## Summary
新增 `tencentcloud_cls_dashboard` Terraform 资源，用于管理腾讯云日志服务（CLS）的仪表盘。该资源支持创建、查询、更新和删除仪表盘，以及标签管理功能。

## Motivation
日志服务的仪表盘是用户进行日志数据可视化分析的重要工具。当前 terraform-provider-tencentcloud 尚未提供 CLS Dashboard 资源的管理能力，用户无法通过 IaC 方式管理仪表盘的生命周期。通过添加此资源，用户可以：

1. 使用基础设施即代码（IaC）方式管理日志仪表盘
2. 实现仪表盘配置的版本控制和自动化部署
3. 统一管理仪表盘及其关联的标签
4. 简化多环境（开发、测试、生产）的仪表盘配置同步

## Background
腾讯云日志服务（CLS）提供了完整的 Dashboard API，包括：
- **CreateDashboard**：创建仪表盘（https://cloud.tencent.com/document/product/614/127511）
- **DescribeDashboards**：查询仪表盘列表（https://cloud.tencent.com/document/product/614/95636）
- **ModifyDashboard**：修改仪表盘（https://cloud.tencent.com/document/product/614/127509）
- **DeleteDashboard**：删除仪表盘（https://cloud.tencent.com/document/product/614/127510）

SDK 支持：腾讯云 Go SDK（tencentcloud-sdk-go/tencentcloud/cls/v20201016）已包含所有所需 API。

## Proposed Solution

### Resource Definition
```hcl
resource "tencentcloud_cls_dashboard" "example" {
  dashboard_name = "my-dashboard"
  data           = jsonencode({
    timezone = "browser"
    subType  = "CLS_Host"
  })
  
  tags = {
    "team"        = "ops"
    "environment" = "production"
  }
}
```

### API Mapping

| Terraform Operation | CLS API | Notes |
|---------------------|---------|-------|
| Create | CreateDashboard | 创建仪表盘 |
| Read | DescribeDashboards | 查询仪表盘详情（支持分页和重试） |
| Update | ModifyDashboard | 更新名称、配置数据和标签 |
| Delete | DeleteDashboard | 删除仪表盘 |
| Import | DescribeDashboards | 通过 dashboard_id 导入 |

### Schema Design

#### Input Arguments
- `dashboard_name` (Required, String) - 仪表盘名称，账户内唯一
- `data` (Optional, String) - 仪表盘配置数据（JSON 字符串），默认为空
- `tags` (Optional, Map) - 标签键值对，最多 10 个

#### Computed Attributes
- `dashboard_id` (String) - 仪表盘 ID（全局唯一标识符）
- `create_time` (String) - 创建时间
- `update_time` (String) - 更新时间

#### Resource ID Format
使用 `dashboard_id` 作为资源 ID（简单格式，不使用复合 ID）。

### Key Implementation Details

#### 1. Create Function
```go
func resourceTencentCloudClsDashboardCreate(d *schema.ResourceData, meta interface{}) error {
    // 1. 构造 CreateDashboardRequest
    // 2. 设置 DashboardName 和 Data
    // 3. 处理 Tags（如果提供）
    // 4. 调用 API（带重试机制）
    // 5. 设置 resource ID 为 dashboard_id
    // 6. 调用 Read 刷新状态
}
```

**重点**：
- 名称唯一性：API 会返回 `InvalidParameter.DashboardNameConflict` 错误
- Data 字段为可选，默认创建空仪表盘
- 使用 `tccommon.WriteRetryTimeout`（5 分钟）

#### 2. Read Function with Pagination and Retry
```go
func resourceTencentCloudClsDashboardRead(d *schema.ResourceData, meta interface{}) error {
    // 1. 获取 dashboard_id
    // 2. 调用 Service 层方法查询
    // 3. 如果不存在，清空 ID
    // 4. 设置 schema 字段
}
```

**Service 层实现（关键）**：
```go
func (me *ClsService) DescribeClsDashboardById(ctx context.Context, dashboardId string) (dashboard *cls.DashboardInfo, errRet error) {
    // 实现分页查询和重试逻辑
    var (
        offset int64 = 0
        limit  int64 = 20
    )
    
    for {
        request.Offset = &offset
        request.Limit = &limit
        
        // 使用 resource.Retry 包装 API 调用
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            response, e := me.client.UseClsClient().DescribeDashboards(request)
            // 处理错误和响应
        })
        
        // 遍历当前页结果
        for _, item := range response.DashboardInfos {
            if *item.DashboardId == dashboardId {
                return item, nil
            }
        }
        
        // 判断是否有下一页
        if len(response.DashboardInfos) < int(limit) {
            break
        }
        offset += limit
    }
}
```

**重点**：
- **分页查询**：使用 Offset 和 Limit 参数遍历所有仪表盘
- **重试机制**：使用 `resource.Retry` + `tccommon.ReadRetryTimeout`（1 分钟）
- **错误处理**：使用 `tccommon.RetryError(e, tccommon.InternalError)`
- 资源不存在时返回 nil（不报错）

#### 3. Update Function
```go
func resourceTencentCloudClsDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
    // 1. 检查哪些字段发生变更
    // 2. 构造 ModifyDashboardRequest
    // 3. 设置 DashboardId
    // 4. 更新变更的字段（dashboard_name, data, tags）
    // 5. 调用 API（带重试）
    // 6. 调用 Read 刷新状态
}
```

**支持修改的字段**：
- `dashboard_name`
- `data`
- `tags`

#### 4. Delete Function
```go
func resourceTencentCloudClsDashboardDelete(d *schema.ResourceData, meta interface{}) error {
    // 1. 获取 dashboard_id
    // 2. 构造 DeleteDashboardRequest
    // 3. 调用 API（带重试和幂等性处理）
}
```

**幂等性处理**：
- 如果仪表盘已不存在，视为删除成功
- 容忍 `ResourceNotFound` 类错误

### Error Handling Strategy

| 错误码 | 场景 | 处理策略 |
|--------|------|----------|
| `InvalidParameter.DashboardNameConflict` | 创建/更新时名称冲突 | 立即返回错误，提示用户更换名称 |
| `LimitExceeded` | 仪表盘数量超限 | 立即返回错误，提示用户检查配额 |
| `LimitExceeded.Tag` | 标签数量超过 10 个 | 立即返回错误，提示减少标签数量 |
| `InternalError` | 内部错误 | 重试处理 |
| `OperationDenied.AccountIsolate` | 账户欠费 | 立即返回错误，提示充值 |
| 资源不存在 | 查询或删除不存在的仪表盘 | Read 返回 nil，Delete 视为成功 |

### Testing Strategy

#### Acceptance Tests
```go
func TestAccTencentCloudClsDashboard_basic(t *testing.T) {
    // 测试基本 CRUD 操作
    // 1. 创建仪表盘（空配置）
    // 2. 验证字段正确性
    // 3. 更新名称和配置
    // 4. 验证更新生效
    // 5. 删除仪表盘
}

func TestAccTencentCloudClsDashboard_tags(t *testing.T) {
    // 测试标签管理
    // 1. 创建带标签的仪表盘
    // 2. 验证标签正确绑定
    // 3. 更新标签
    // 4. 删除标签
}

func TestAccTencentCloudClsDashboard_withData(t *testing.T) {
    // 测试完整配置数据
    // 1. 创建包含复杂配置的仪表盘
    // 2. 验证配置正确保存
    // 3. 更新配置
}
```

### Documentation

#### Example Usage
```hcl
# 基础示例
resource "tencentcloud_cls_dashboard" "basic" {
  dashboard_name = "basic-dashboard"
}

# 完整示例
resource "tencentcloud_cls_dashboard" "full" {
  dashboard_name = "production-dashboard"
  data           = jsonencode({
    timezone = "browser"
    subType  = "CLS_Host"
    charts   = [
      {
        id    = "chart1"
        type  = "line"
        title = "Request Count"
      }
    ]
  })
  
  tags = {
    "team"        = "platform"
    "environment" = "production"
    "owner"       = "admin"
  }
}
```

#### Import
```bash
terraform import tencentcloud_cls_dashboard.example dashboard-xxxx-xxxx-xxxx-xxxx
```

## Impact Assessment

### Compatibility
- **向后兼容**：新增资源，不影响现有功能
- **Provider 版本**：无最低版本要求
- **SDK 版本**：需要 tencentcloud-sdk-go/tencentcloud/cls/v20201016 包

### Dependencies
- 依赖 CLS 服务的其他资源（如 logset, topic）可能与 dashboard 关联
- 但 dashboard 可独立管理，无强制依赖

### Breaking Changes
无

## Alternatives Considered

### Alternative 1: 将 Data 字段设计为结构化对象
**拒绝理由**：
- API 使用 JSON 字符串格式
- 仪表盘配置结构复杂且可能频繁变化
- 使用字符串格式给用户更大灵活性
- 用户可使用 `jsonencode()` 函数构造配置

### Alternative 2: 支持通过 DashboardName 导入
**拒绝理由**：
- 名称可能重复（虽然 API 限制唯一，但历史数据可能存在）
- DashboardId 是全局唯一且不可变的标识符
- 与其他 CLS 资源保持一致（使用 ID 导入）

## Implementation Plan

### Phase 1: Core Implementation (3 days)
- [x] 创建 proposal.md
- [ ] 实现 Resource Schema
- [ ] 实现 Create 函数
- [ ] 实现 Service 层（包含分页和重试逻辑）
- [ ] 实现 Read 函数
- [ ] 实现 Update 函数
- [ ] 实现 Delete 函数
- [ ] 注册资源到 provider

### Phase 2: Testing (2 days)
- [ ] 编写验收测试
- [ ] 测试基本 CRUD 操作
- [ ] 测试标签管理
- [ ] 测试错误场景
- [ ] 测试导入功能

### Phase 3: Documentation (1 day)
- [ ] 编写资源文档
- [ ] 添加示例配置
- [ ] 更新 provider.md
- [ ] 运行 `make doc` 生成最终文档

### Total Estimated Time: 6 days

## Success Criteria
- [ ] 资源成功创建、读取、更新、删除
- [ ] 分页查询和重试逻辑正确实现
- [ ] 所有验收测试通过
- [ ] 代码通过 lint 检查
- [ ] 文档完整且格式正确
- [ ] 支持资源导入
- [ ] 错误处理覆盖所有已知场景

## References
- [CreateDashboard API](https://cloud.tencent.com/document/product/614/127511)
- [DescribeDashboards API](https://cloud.tencent.com/document/product/614/95636)
- [ModifyDashboard API](https://cloud.tencent.com/document/product/614/127509)
- [DeleteDashboard API](https://cloud.tencent.com/document/product/614/127510)
- [CLS SDK Go](https://github.com/TencentCloud/tencentcloud-sdk-go/tree/master/tencentcloud/cls/v20201016)
- [Project Standards](../../../openspec/project.md)
