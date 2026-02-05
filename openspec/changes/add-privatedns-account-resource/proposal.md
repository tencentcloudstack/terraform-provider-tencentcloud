# Proposal: Add Private DNS Account Resource

**Change ID**: `add-privatedns-account-resource`  
**Status**: Proposal  
**Author**: AI Agent  
**Date**: 2026-02-03

## Executive Summary

新增 `tencentcloud_private_dns_account` 资源，用于管理 Private DNS 跨账号关联场景下的关联账号。该资源支持添加、查询和删除关联账号，允许跨账号绑定 VPC 资源到私有域。

## Problem Statement

当前 Terraform Provider 缺少对 Private DNS 关联账号的管理能力。在跨账号绑定 VPC 的场景下，用户需要：
- 添加关联账号以获取对应账号的 VPC 资源访问权限
- 查询已关联的账号列表
- 移除不再需要的账号关联

目前用户只能通过控制台或 API 手动管理，无法使用 Terraform 进行自动化管理。

## Background

### API Support Analysis

1. **CreatePrivateDNSAccount** (https://cloud.tencent.com/document/api/1338/64976)
   - **功能**: 添加关联账号
   - **频率限制**: 10次/秒
   - **输入参数**: `Account.Uin` (必填，关联账号的 UIN)
   - **返回**: RequestId
   - **限制**: 关联账号不能与主账号一致，需要子账号授权

2. **DeletePrivateDNSAccount** (https://cloud.tencent.com/document/api/1338/64975)
   - **功能**: 移除账号关联
   - **频率限制**: 20次/秒
   - **输入参数**: `Account.Uin` (必填)
   - **返回**: RequestId
   - **限制**: 存在绑定的 VPC 资源时需要先解绑

3. **DescribePrivateDNSAccountList** (https://cloud.tencent.com/document/api/1338/61417)
   - **功能**: 获取关联账号列表
   - **频率限制**: 20次/秒
   - **输入参数**: 
     - `Offset` (可选，分页偏移量，从0开始)
     - `Limit` (可选，分页限制，最大100，默认20)
     - `Filters` (可选，过滤参数，支持 AccountUin 过滤)
   - **返回**: 
     - `TotalCount` (总数)
     - `AccountSet` (账号列表，包含 Uin, Account, Nickname)
   - **注意**: API 不直接支持按 UIN 精确查询单个账号，需要通过 Filters 实现

### Design Considerations

1. **资源 ID 设计**: 使用 `Uin` 作为资源 ID（唯一标识符）
2. **Read 实现**: 需要实现分页逻辑和 UIN 过滤逻辑
3. **Update 操作**: API 不支持修改操作，Uin 变化需要 ForceNew
4. **Import 支持**: 支持通过 Uin 导入现有账号关联

## Proposed Solution

### Resource Schema

```go
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"  // 必填，关联账号的 Uin
}
```

**Schema Definition**:
- `account_uin`: (Required, String, ForceNew) 关联账号的 Uin

**Computed Attributes**:
- `account`: (Computed, String) 关联账号的邮箱
- `nickname`: (Computed, String) 关联账号的昵称

### Implementation Components

#### 1. Resource File: `resource_tc_private_dns_account.go`

**Create Operation**:
```go
func resourceTencentCloudPrivateDnsAccountCreate(d *schema.ResourceData, meta interface{}) error {
    // 1. 提取 account_uin
    // 2. 构建 CreatePrivateDNSAccountRequest
    // 3. 调用 API 创建关联
    // 4. 设置资源 ID 为 Uin
    // 5. 调用 Read 获取完整信息
}
```

**Read Operation**:
```go
func resourceTencentCloudPrivateDnsAccountRead(d *schema.ResourceData, meta interface{}) error {
    // 1. 获取资源 ID (Uin)
    // 2. 调用服务层 DescribePrivateDnsAccountByUin
    // 3. 如果未找到，清空资源 ID
    // 4. 设置 account, nickname 等属性到状态
}
```

**Update Operation**:
```go
// 不需要实现 - Uin 是 ForceNew，其他字段都是 Computed
```

**Delete Operation**:
```go
func resourceTencentCloudPrivateDnsAccountDelete(d *schema.ResourceData, meta interface{}) error {
    // 1. 获取资源 ID (Uin)
    // 2. 构建 DeletePrivateDNSAccountRequest
    // 3. 调用 API 删除关联
    // 4. 处理 ExistBoundVpc 错误（需要先解绑 VPC）
}
```

#### 2. Service Layer: `service_tencentcloud_private_dns.go`

新增服务方法：

```go
func (me *PrivateDnsService) DescribePrivateDnsAccountByUin(ctx context.Context, uin string) (
    account *privatedns.PrivateDNSAccount, errRet error) {
    // 1. 构建 DescribePrivateDNSAccountListRequest
    // 2. 设置 Filters: Name="AccountUin", Values=[uin]
    // 3. 实现分页逻辑（limit=100）
    // 4. 遍历所有页查找匹配的 Uin
    // 5. 返回匹配的账号或 nil
}

func (me *PrivateDnsService) CreatePrivateDnsAccount(ctx context.Context, uin string) (errRet error) {
    // 1. 构建 CreatePrivateDNSAccountRequest
    // 2. 设置 Account.Uin
    // 3. 调用 API 并处理重试逻辑
}

func (me *PrivateDnsService) DeletePrivateDnsAccount(ctx context.Context, uin string) (errRet error) {
    // 1. 构建 DeletePrivateDNSAccountRequest
    // 2. 设置 Account.Uin
    // 3. 调用 API 并处理重试逻辑
    // 4. 处理 ExistBoundVpc 错误
}
```

#### 3. Testing: `resource_tc_private_dns_account_test.go`

**测试场景**:
1. **Basic Test**: 创建、读取、删除关联账号
2. **Import Test**: 导入现有关联账号
3. **Update Test**: 验证 Uin 修改触发 ForceNew

#### 4. Documentation: `resource_tc_private_dns_account.md`

包含：
- 资源描述
- 使用示例
- 参数说明
- 属性说明
- 导入说明

## Technical Details

### API Request/Response Mapping

**Create API**:
```go
Request:
{
    "Account": {
        "Uin": "100123456789"
    }
}

Response:
{
    "RequestId": "xxx"
}
```

**Describe API**:
```go
Request:
{
    "Offset": 0,
    "Limit": 100,
    "Filters": [
        {
            "Name": "AccountUin",
            "Values": ["100123456789"]
        }
    ]
}

Response:
{
    "TotalCount": 1,
    "AccountSet": [
        {
            "Uin": "100123456789",
            "Account": "test@example.com",
            "Nickname": "Test Account"
        }
    ]
}
```

**Delete API**:
```go
Request:
{
    "Account": {
        "Uin": "100123456789"
    }
}

Response:
{
    "RequestId": "xxx"
}
```

### Error Handling

| Error Code | Description | Handling Strategy |
|------------|-------------|-------------------|
| `InvalidParameter.AccountExist` | 账号已存在 | 在 Create 时视为成功，调用 Read 获取信息 |
| `UnsupportedOperation.ExistBoundVpc` | 存在绑定的 VPC | 在 Delete 时返回明确错误提示用户先解绑 |
| `ResourceNotFound` | 账号不存在 | 在 Read 时清空资源 ID |
| `ResourceNotFound.ServiceNotSubscribed` | 服务未开通 | 返回错误提示用户开通服务 |

### Read Operation Implementation Details

由于 API 不支持直接通过 Uin 查询单个账号，需要实现以下逻辑：

```go
func (me *PrivateDnsService) DescribePrivateDnsAccountByUin(ctx context.Context, uin string) (
    account *privatedns.PrivateDNSAccount, errRet error) {
    
    logId := tccommon.GetLogId(ctx)
    request := privatedns.NewDescribePrivateDNSAccountListRequest()
    
    // 设置过滤器
    request.Filters = []*privatedns.Filter{
        {
            Name:   helper.String("AccountUin"),
            Values: []*string{helper.String(uin)},
        },
    }
    
    // 分页参数
    var (
        limit  int64 = 100
        offset int64 = 0
    )
    
    // 循环获取所有页
    for {
        request.Limit = &limit
        request.Offset = &offset
        
        var response *privatedns.DescribePrivateDNSAccountListResponse
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            result, e := me.client.UsePrivateDnsClient().DescribePrivateDNSAccountList(request)
            if e != nil {
                return tccommon.RetryError(e, tccommon.InternalError)
            }
            response = result
            return nil
        })
        
        if err != nil {
            return nil, err
        }
        
        // 查找匹配的账号
        for _, acc := range response.Response.AccountSet {
            if *acc.Uin == uin {
                return acc, nil
            }
        }
        
        // 检查是否还有更多页
        if offset+limit >= *response.Response.TotalCount {
            break
        }
        offset += limit
    }
    
    // 未找到
    return nil, nil
}
```

## Alternatives Considered

### Alternative 1: 使用复合 ID

**考虑**: 使用 `{region}:{uin}` 作为资源 ID
**拒绝理由**: Private DNS Account 是全局资源，不绑定到特定区域，使用单一 Uin 即可

### Alternative 2: 在 Read 时不实现过滤

**考虑**: 获取所有账号列表，在客户端内存中过滤
**拒绝理由**: API 已支持 Filter 参数，应该利用服务端过滤减少网络开销

### Alternative 3: Update 支持修改 Nickname

**考虑**: 允许用户修改账号昵称
**拒绝理由**: API 不支持修改操作，Nickname 是只读属性

## Benefits

1. **自动化管理**: 用户可以通过 Terraform 自动化管理 Private DNS 关联账号
2. **一致性**: 与其他 Private DNS 资源保持一致的使用体验
3. **可追溯**: 通过 Terraform 状态跟踪账号关联的变更历史
4. **可导入**: 支持导入现有账号关联到 Terraform 管理

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| API 限频 | 中 | 使用现有重试逻辑和速率限制机制 |
| Read 性能 | 低 | 使用 Filter 参数减少数据量，实现高效分页 |
| 删除失败（VPC 绑定） | 中 | 返回清晰错误消息，指导用户先解绑 VPC |
| 账号已存在 | 低 | 视为幂等操作，Create 时检测已存在则调用 Read |

## Success Criteria

1. ✅ 用户可以通过 Terraform 创建关联账号
2. ✅ 账号信息正确读取并显示在状态中
3. ✅ Uin 修改触发资源重建
4. ✅ 账号关联可以正确删除
5. ✅ 支持导入现有关联账号
6. ✅ 所有验收测试通过
7. ✅ 文档完整准确

## Testing Strategy

1. **Unit Tests**: Schema 验证
2. **Acceptance Tests**:
   - 创建关联账号
   - 读取并验证属性
   - 导入关联账号
   - 修改 Uin 验证 ForceNew
   - 删除关联账号
3. **Manual Tests**: 
   - 在腾讯云控制台验证
   - 测试 VPC 绑定场景

## Documentation

### User Documentation

```hcl
# Create a Private DNS account association
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

# Output the account information
output "account_email" {
  value = tencentcloud_private_dns_account.example.account
}

output "account_nickname" {
  value = tencentcloud_private_dns_account.example.nickname
}
```

### Import Example

```bash
$ terraform import tencentcloud_private_dns_account.example 100123456789
```

## Timeline

- **Proposal**: 0.5 day (完成)
- **Implementation**: 1 day
  - Service layer: 0.3 day
  - Resource implementation: 0.4 day
  - Tests: 0.3 day
- **Documentation**: 0.5 day
- **Review & Testing**: 0.5 day
- **Total**: ~2.5 days

## Dependencies

- **SDK**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028` (已存在)
- **Service Layer**: `PrivateDnsService` (已存在，需要扩展)
- **No Breaking Changes**: 纯新增功能

## References

- **API Documentation**:
  - CreatePrivateDNSAccount: https://cloud.tencent.com/document/api/1338/64976
  - DeletePrivateDNSAccount: https://cloud.tencent.com/document/api/1338/64975
  - DescribePrivateDNSAccountList: https://cloud.tencent.com/document/api/1338/61417
- **Existing Implementation**: `tencentcloud/services/privatedns/resource_tc_private_dns_zone.go`
- **Service Layer**: `tencentcloud/services/privatedns/service_tencentcloud_private_dns.go`

## Open Questions

**Q1**: 是否需要在创建前验证账号是否有效？  
**A1**: 否，API 会返回相应错误，由用户处理

**Q2**: 删除时如果有 VPC 绑定，是否自动解绑？  
**A2**: 否，返回错误提示用户手动解绑，避免意外删除

**Q3**: 是否需要支持批量操作？  
**A3**: 否，遵循 Terraform 单资源原则，用户可以使用 count 或 for_each

**Q4**: Nickname 和 Account 字段是否应该支持配置？  
**A4**: 否，这些是只读属性，由腾讯云系统管理
