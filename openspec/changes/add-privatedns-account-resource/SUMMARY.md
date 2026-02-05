# Summary: Add Private DNS Account Resource

**Change ID**: `add-privatedns-account-resource`  
**Status**: ✅ Proposal Complete - Ready for Review  
**Validation**: ✅ Passed `openspec validate --strict`

---

## Overview

新增 `tencentcloud_private_dns_account` Terraform 资源，用于管理 Private DNS 跨账号关联场景下的关联账号。

---

## Problem

当前 Terraform Provider 缺少对 Private DNS 关联账号的管理能力。在跨账号绑定 VPC 的场景下，用户需要：
- 添加关联账号以获取对应账号的 VPC 资源访问权限
- 查询已关联的账号列表
- 移除不再需要的账号关联

目前只能通过控制台或 API 手动管理，无法使用 Terraform 进行自动化。

---

## Solution

### Resource Definition

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"  # 必填，关联账号的 Uin
}

# 输出计算属性
output "account_email" {
  value = tencentcloud_private_dns_account.example.account
}

output "account_nickname" {
  value = tencentcloud_private_dns_account.example.nickname
}
```

### Schema

| 字段 | 类型 | 必填 | ForceNew | 描述 |
|------|------|------|----------|------|
| `account_uin` | String | ✅ | ✅ | 关联账号的 Uin |
| `account` | String | ❌ (Computed) | - | 关联账号的邮箱 |
| `nickname` | String | ❌ (Computed) | - | 关联账号的昵称 |

### API Mapping

| 操作 | API | 说明 |
|------|-----|------|
| **Create** | CreatePrivateDNSAccount | 添加关联账号 |
| **Read** | DescribePrivateDNSAccountList | 查询账号列表（需分页和过滤） |
| **Update** | - | 不支持（Uin 是 ForceNew） |
| **Delete** | DeletePrivateDNSAccount | 移除关联账号 |
| **Import** | DescribePrivateDNSAccountList | 通过 Uin 导入 |

---

## Key Features

### ✅ 完整的 CRUD 支持
- ✅ 创建关联账号
- ✅ 读取账号信息（分页 + 过滤）
- ✅ 删除关联账号
- ✅ 导入现有账号

### ✅ 智能 Read 实现
- **分页逻辑**: 自动处理超过 100 个账号的场景
- **UIN 过滤**: 使用 API Filter 参数精确查询
- **高效查找**: 遍历分页结果直到找到目标账号

```go
// Read 实现伪代码
func DescribePrivateDnsAccountByUin(uin string) {
    filters = [{"Name": "AccountUin", "Values": [uin]}]
    
    for offset = 0; offset < totalCount; offset += 100 {
        response = DescribePrivateDNSAccountList(offset, 100, filters)
        
        for _, account := range response.AccountSet {
            if account.Uin == uin {
                return account  // 找到了！
            }
        }
    }
    
    return nil  // 未找到
}
```

### ✅ 错误处理
| 错误场景 | 处理策略 |
|----------|----------|
| 账号已存在 | 视为幂等操作，调用 Read 获取信息 |
| 存在 VPC 绑定 | 返回明确错误，提示先解绑 VPC |
| 账号不存在 | Read 时清空资源 ID，标记为需重建 |
| 服务未开通 | 返回错误提示用户开通服务 |

### ✅ ForceNew 行为
- Uin 是唯一可配置字段
- 修改 Uin 触发资源重建
- 先删除旧资源，再创建新资源

---

## Implementation Details

### File Structure

```
tencentcloud/services/privatedns/
├── service_tencentcloud_private_dns.go        # 扩展，新增 3 个方法
├── resource_tc_private_dns_account.go         # 新建，资源实现
├── resource_tc_private_dns_account.md         # 新建，源文档
└── resource_tc_private_dns_account_test.go    # 新建，测试
```

### Service Layer (新增方法)

```go
type PrivateDnsService struct {
    client *connectivity.TencentCloudClient
}

// 1. 按 Uin 查询账号（实现分页 + 过滤）
func (me *PrivateDnsService) DescribePrivateDnsAccountByUin(
    ctx context.Context, 
    uin string,
) (*privatedns.PrivateDNSAccount, error)

// 2. 创建关联账号
func (me *PrivateDnsService) CreatePrivateDnsAccount(
    ctx context.Context,
    uin string,
) error

// 3. 删除关联账号
func (me *PrivateDnsService) DeletePrivateDnsAccount(
    ctx context.Context,
    uin string,
) error
```

### Resource Layer

```go
func ResourceTencentCloudPrivateDnsAccount() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudPrivateDnsAccountCreate,
        Read:   resourceTencentCloudPrivateDnsAccountRead,
        Delete: resourceTencentCloudPrivateDnsAccountDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: {...},
    }
}
```

---

## Requirements & Scenarios

**9 个需求，15 个测试场景**：

1. **PDNS-ACCT-001**: Resource Schema Definition (2 scenarios)
   - 定义必需的 account_uin 字段
   - 访问计算属性 account, nickname

2. **PDNS-ACCT-002**: Create Account Association (3 scenarios)
   - 创建新关联账号
   - 处理账号已存在
   - 处理服务未开通

3. **PDNS-ACCT-003**: Read Account Information (3 scenarios)
   - 读取现有账号
   - 处理账号不存在
   - 正确处理分页

4. **PDNS-ACCT-004**: Delete Account Association (3 scenarios)
   - 成功删除账号
   - 处理 VPC 绑定存在
   - 处理账号已删除

5. **PDNS-ACCT-005**: Import Existing Account (1 scenario)
   - 导入现有账号关联

6. **PDNS-ACCT-006**: ForceNew on Uin Change (1 scenario)
   - Uin 修改触发替换

7. **PDNS-ACCT-007**: Error Handling and Retry Logic (3 scenarios)
   - 临时错误重试
   - 业务错误立即失败
   - 记录所有 API 交互

8. **PDNS-ACCT-008**: Service Layer Abstraction (2 scenarios)
   - 定义服务层方法
   - 封装分页逻辑

9. **PDNS-ACCT-009**: Documentation Completeness (3 scenarios)
   - 提供全面的使用示例
   - 文档化所有参数和属性
   - 提供导入说明

---

## Tasks Breakdown

**15 个任务，6 个阶段**：

### Phase 1: Service Layer (5 tasks)
- 实现 `DescribePrivateDnsAccountByUin`（分页 + 过滤）
- 实现 `CreatePrivateDnsAccount`
- 实现 `DeletePrivateDnsAccount`
- 添加错误常量
- 代码格式化

### Phase 2: Resource Implementation (5 tasks)
- 创建文件和 Schema 定义
- 实现 Create 函数
- 实现 Read 函数
- 实现 Delete 函数
- 代码格式化

### Phase 3: Provider Registration (1 task)
- 在 Provider 中注册新资源

### Phase 4: Testing (4 tasks)
- 创建测试文件
- 编写基础 CRUD 测试
- 编写导入测试
- 编写 ForceNew 测试

### Phase 5: Documentation (3 tasks)
- 创建资源文档
- 生成网站文档
- 更新 provider.md

### Phase 6: Code Quality (2 tasks)
- 运行代码检查
- 运行验收测试

---

## Testing Strategy

### Acceptance Tests

```go
func TestAccTencentCloudPrivateDnsAccountResource_Basic(t *testing.T) {
    // 1. 创建关联账号
    // 2. 验证账号属性
    // 3. 导入测试
    // 4. 删除账号
}

func TestAccTencentCloudPrivateDnsAccountResource_ForceNew(t *testing.T) {
    // 1. 创建账号（Uin1）
    // 2. 修改 Uin 为 Uin2
    // 3. 验证资源被替换
}
```

### Manual Tests
- 控制台验证账号添加
- 测试 VPC 绑定场景
- 验证导入功能

---

## Documentation

### Example Usage

```hcl
# Basic usage
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

# Access computed attributes
output "account_info" {
  value = {
    uin      = tencentcloud_private_dns_account.example.account_uin
    email    = tencentcloud_private_dns_account.example.account
    nickname = tencentcloud_private_dns_account.example.nickname
  }
}
```

### Import

```bash
$ terraform import tencentcloud_private_dns_account.example 100123456789
```

---

## Benefits

1. **自动化管理**: Terraform 自动化管理关联账号
2. **一致性**: 与其他 Private DNS 资源保持一致
3. **可追溯**: Terraform 状态跟踪变更历史
4. **可导入**: 支持导入现有关联账号
5. **高效查询**: 智能分页和过滤减少 API 调用

---

## Risks & Mitigations

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| API 限频 | 中 | 使用现有重试逻辑和速率限制 |
| Read 性能 | 低 | 使用 Filter 参数 + 高效分页 |
| 删除失败（VPC 绑定） | 中 | 返回清晰错误消息 |
| 账号已存在 | 低 | 幂等操作，Create 时检测已存在 |

---

## Timeline

- ✅ **Proposal**: 0.5 day (完成)
- ⏳ **Implementation**: 1 day
  - Service layer: 0.3 day
  - Resource implementation: 0.4 day
  - Tests: 0.3 day
- ⏳ **Documentation**: 0.5 day
- ⏳ **Review & Testing**: 0.5 day
- **Total**: ~2.5 days

---

## Dependencies

- **SDK**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028` ✅ 已存在
- **Service**: `PrivateDnsService` ✅ 已存在，需扩展
- **Breaking Changes**: ❌ 无，纯新增功能

---

## API References

- [CreatePrivateDNSAccount](https://cloud.tencent.com/document/api/1338/64976) - 添加关联账号
- [DeletePrivateDNSAccount](https://cloud.tencent.com/document/api/1338/64975) - 移除账号关联
- [DescribePrivateDNSAccountList](https://cloud.tencent.com/document/api/1338/61417) - 获取关联账号列表

---

## Files to Create/Modify

### New Files (3)
```
✨ tencentcloud/services/privatedns/resource_tc_private_dns_account.go
✨ tencentcloud/services/privatedns/resource_tc_private_dns_account.md
✨ tencentcloud/services/privatedns/resource_tc_private_dns_account_test.go
```

### Modified Files (2)
```
📝 tencentcloud/services/privatedns/service_tencentcloud_private_dns.go
📝 tencentcloud/provider.go
```

### Generated Files (1)
```
🤖 website/docs/r/private_dns_account.html.markdown (via make doc)
```

---

## Validation Status

```bash
$ openspec validate add-privatedns-account-resource --strict
✅ Change 'add-privatedns-account-resource' is valid
```

---

## Next Steps

1. **Review**: 团队审查提案
2. **Approval**: 获取利益相关者批准
3. **Implementation**: 按照 tasks.md 执行实施
4. **Testing**: 运行完整测试套件
5. **Documentation**: 生成用户文档
6. **Merge**: 合并到主分支
7. **Release**: 包含在下一个 provider 版本中

---

**提案状态**: ✅ **完整且已验证，准备审查和实施！**

---

## Related Resources

该资源与以下现有资源配合使用：
- `tencentcloud_private_dns_zone` - 私有域管理
- `tencentcloud_private_dns_zone_vpc_attachment` - VPC 绑定管理
- `tencentcloud_private_dns_record` - 解析记录管理

**使用场景示例**：

```hcl
# 1. 添加关联账号
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

# 2. 使用关联账号的 VPC 绑定私有域
resource "tencentcloud_private_dns_zone" "example" {
  domain = "example.com"
  
  account_vpc_set {
    uniq_vpc_id = "vpc-xxxxx"
    region      = "ap-guangzhou"
    uin         = tencentcloud_private_dns_account.example.account_uin
  }
}
```
