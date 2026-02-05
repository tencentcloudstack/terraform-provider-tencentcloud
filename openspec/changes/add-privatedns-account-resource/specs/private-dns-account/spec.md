# Specification: Private DNS Account Resource

**Capability ID**: `private-dns-account`  
**Resource**: `tencentcloud_private_dns_account`  
**Status**: Draft

---

## ADDED Requirements

### Requirement: PDNS-ACCT-001 - Resource Schema Definition
**Priority**: High  
**Type**: Functional

资源 MUST 支持管理 Private DNS 关联账号的生命周期。

**Acceptance Criteria**:
- 资源名称为 `tencentcloud_private_dns_account`
- Schema 包含必需的 `account_uin` 字段
- Schema 包含计算字段 `account` 和 `nickname`
- 支持资源导入功能

#### Scenario: Define resource with required account_uin
**Given** 用户需要关联一个账号到 Private DNS  
**When** 在 Terraform 配置中定义资源  
**Then** 用户可以指定 `account_uin` 参数  
**And** 参数必须是字符串类型  
**And** 参数是必填的  
**And** 参数修改触发资源重建（ForceNew）

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}
```

#### Scenario: Access computed account information
**Given** 关联账号已创建  
**When** 读取资源状态  
**Then** `account` 字段显示账号邮箱  
**And** `nickname` 字段显示账号昵称  
**And** 这些字段是只读的（Computed）

```hcl
output "account_email" {
  value = tencentcloud_private_dns_account.example.account
}

output "account_nickname" {
  value = tencentcloud_private_dns_account.example.nickname
}
```

---

### Requirement: PDNS-ACCT-002 - Create Account Association
**Priority**: High  
**Type**: Functional

创建资源时 MUST 调用 CreatePrivateDNSAccount API 建立账号关联。

**Acceptance Criteria**:
- 提取 `account_uin` 参数
- 构建正确的 API 请求
- 调用 CreatePrivateDNSAccount API
- 设置资源 ID 为 Uin
- 处理账号已存在错误（幂等性）
- 创建后立即读取完整信息

#### Scenario: Create new account association
**Given** 账号未关联  
**When** 执行 `terraform apply`  
**Then** CreatePrivateDNSAccount API 被调用  
**And** Request.Account.Uin 设置为用户提供的 Uin  
**And** 资源 ID 设置为 Uin  
**And** Read 操作被调用获取完整信息  
**And** 状态包含 account 和 nickname

```hcl
resource "tencentcloud_private_dns_account" "new" {
  account_uin = "100123456789"
}
```

#### Scenario: Handle account already exists
**Given** 账号已在其他地方关联  
**When** 尝试创建相同 Uin 的资源  
**And** API 返回 `InvalidParameter.AccountExist` 错误  
**Then** 创建操作视为成功  
**And** 调用 Read 获取现有账号信息  
**And** 资源状态正确反映现有账号

#### Scenario: Handle service not subscribed
**Given** Private DNS 服务未开通  
**When** 尝试创建关联账号  
**And** API 返回 `ResourceNotFound.ServiceNotSubscribed` 错误  
**Then** 创建失败  
**And** 返回明确的错误消息提示用户开通服务

---

### Requirement: PDNS-ACCT-003 - Read Account Information
**Priority**: High  
**Type**: Functional

读取资源时 MUST 通过 DescribePrivateDNSAccountList API 查询账号信息。

**Acceptance Criteria**:
- 使用资源 ID（Uin）作为查询条件
- 通过 Filter 参数精确匹配 Uin
- 实现分页逻辑处理大量账号
- 如果账号不存在，清空资源 ID
- 正确设置所有状态字段

#### Scenario: Read existing account association
**Given** 关联账号存在  
**When** Terraform 读取资源状态  
**Then** DescribePrivateDNSAccountList API 被调用  
**And** Filters 包含 AccountUin 过滤器  
**And** 遍历分页结果查找匹配 Uin  
**And** `account_uin` 设置为 Uin  
**And** `account` 设置为账号邮箱  
**And** `nickname` 设置为账号昵称

#### Scenario: Handle account not found
**Given** 关联账号已被删除（控制台或其他方式）  
**When** Terraform 读取资源状态  
**And** DescribePrivateDNSAccountList 未返回匹配账号  
**Then** 资源 ID 被清空 `d.SetId("")`  
**And** Terraform 将资源标记为需要重建

#### Scenario: Handle pagination correctly
**Given** 系统中有超过 100 个关联账号  
**When** 读取特定 Uin 的账号  
**Then** 实现分页逻辑（limit=100）  
**And** 遍历所有页直到找到匹配账号或结束  
**And** 正确返回匹配的账号信息

```go
// Pagination logic example
for {
    request.Offset = &offset
    request.Limit = &limit
    
    response := callAPI(request)
    
    for _, account := range response.AccountSet {
        if *account.Uin == targetUin {
            return account
        }
    }
    
    if offset + limit >= response.TotalCount {
        break
    }
    offset += limit
}
```

---

### Requirement: PDNS-ACCT-004 - Delete Account Association
**Priority**: High  
**Type**: Functional

删除资源时 MUST 调用 DeletePrivateDNSAccount API 移除账号关联。

**Acceptance Criteria**:
- 使用资源 ID（Uin）构建删除请求
- 调用 DeletePrivateDNSAccount API
- 处理 VPC 绑定错误
- 实现重试逻辑
- 正确清理资源状态

#### Scenario: Delete account association successfully
**Given** 关联账号存在  
**And** 账号没有绑定任何 VPC  
**When** 执行 `terraform destroy`  
**Then** DeletePrivateDNSAccount API 被调用  
**And** Request.Account.Uin 设置为资源 ID  
**And** 删除成功  
**And** 资源从 Terraform 状态中移除

#### Scenario: Handle VPC binding exists
**Given** 关联账号存在  
**And** 账号绑定了 VPC 资源  
**When** 尝试删除关联账号  
**And** API 返回 `UnsupportedOperation.ExistBoundVpc` 错误  
**Then** 删除失败  
**And** 返回清晰的错误消息  
**And** 错误消息提示用户先解绑 VPC

```
Error: Cannot delete Private DNS account association

The account 100123456789 has VPC resources bound to it.
Please unbind all VPCs from this account before deleting the association.

Use the tencentcloud_private_dns_zone_vpc_attachment resource to manage VPC bindings.
```

#### Scenario: Handle account already deleted
**Given** 关联账号已被删除（控制台或其他方式）  
**When** 执行 `terraform destroy`  
**And** API 返回资源不存在错误  
**Then** 删除操作视为成功  
**And** 资源从 Terraform 状态中移除

---

### Requirement: PDNS-ACCT-005 - Import Existing Account
**Priority**: Medium  
**Type**: Functional

资源 MUST 支持导入现有的关联账号到 Terraform 管理。

**Acceptance Criteria**:
- 支持通过 Uin 导入
- 导入后正确读取所有属性
- 导入后状态与新创建一致
- 后续 plan 不显示变更

#### Scenario: Import existing account association
**Given** 关联账号已存在（通过控制台创建）  
**And** Uin 为 "100123456789"  
**When** 执行 `terraform import tencentcloud_private_dns_account.example 100123456789`  
**Then** 资源成功导入  
**And** 资源 ID 设置为 Uin  
**And** Read 操作被调用  
**And** 所有属性正确设置  
**And** `terraform show` 显示完整资源信息  
**And** `terraform plan` 显示无变更

```bash
$ terraform import tencentcloud_private_dns_account.example 100123456789

tencentcloud_private_dns_account.example: Importing from ID "100123456789"...
tencentcloud_private_dns_account.example: Import prepared!
  Prepared tencentcloud_private_dns_account for import
tencentcloud_private_dns_account.example: Refreshing state... [id=100123456789]

Import successful!

$ terraform show
resource "tencentcloud_private_dns_account" "example" {
    id          = "100123456789"
    account_uin = "100123456789"
    account     = "test@example.com"
    nickname    = "Test Account"
}

$ terraform plan
No changes. Your infrastructure matches the configuration.
```

---

### Requirement: PDNS-ACCT-006 - ForceNew on Uin Change
**Priority**: High  
**Type**: Functional

修改 `account_uin` 字段 MUST 触发资源重建（ForceNew）。

**Acceptance Criteria**:
- `account_uin` 字段标记为 ForceNew
- 修改 Uin 时 Terraform plan 显示重建
- 旧资源先删除，再创建新资源
- 无 Update 操作

#### Scenario: Change account uin triggers replacement
**Given** 关联账号已存在  
**And** 当前 Uin 为 "100123456789"  
**When** 修改配置将 Uin 改为 "200987654321"  
**And** 执行 `terraform plan`  
**Then** Plan 显示资源将被替换  
**And** 显示 `-/+ (forces replacement)`  
**And** 执行 apply 时旧账号先删除  
**And** 然后创建新账号关联

```hcl
# Before
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

# After
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "200987654321"  # Changed
}
```

```
$ terraform plan

Terraform will perform the following actions:

  # tencentcloud_private_dns_account.example must be replaced
-/+ resource "tencentcloud_private_dns_account" "example" {
      ~ account_uin = "100123456789" -> "200987654321" # forces replacement
      ~ account     = "old@example.com" -> (known after apply)
      ~ id          = "100123456789" -> (known after apply)
      ~ nickname    = "Old Account" -> (known after apply)
    }

Plan: 1 to add, 0 to change, 1 to destroy.
```

---

### Requirement: PDNS-ACCT-007 - Error Handling and Retry Logic
**Priority**: High  
**Type**: Non-Functional

所有 API 调用 MUST 包含适当的错误处理和重试逻辑。

**Acceptance Criteria**:
- 使用 `resource.Retry` 实现重试
- 内部错误自动重试
- 业务错误立即返回
- 错误消息包含上下文信息
- 日志记录 API 请求和响应

#### Scenario: Retry on transient errors
**Given** API 调用遇到临时错误  
**When** 错误类型为 `InternalError`  
**Then** 自动重试操作  
**And** 使用指数退避策略  
**And** 最多重试到 ReadRetryTimeout 或 WriteRetryTimeout  
**And** 重试失败后返回最后一个错误

#### Scenario: Immediate failure on business errors
**Given** API 调用遇到业务逻辑错误  
**When** 错误类型为 `InvalidParameter` 或 `UnsupportedOperation`  
**Then** 不进行重试  
**And** 立即返回错误  
**And** 错误消息包含原始 API 错误信息

#### Scenario: Log all API interactions
**Given** 任何 API 调用  
**When** 发起请求  
**Then** 记录 logId, action, request body  
**When** 收到响应  
**Then** 记录 response body (如果成功)  
**When** 发生错误  
**Then** 记录错误详情和原因

```go
log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
    logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
    logId, request.GetAction(), request.ToJsonString(), err.Error())
```

---

### Requirement: PDNS-ACCT-008 - Service Layer Abstraction
**Priority**: High  
**Type**: Technical

服务层 MUST 提供清晰的抽象接口封装 API 调用。

**Acceptance Criteria**:
- 在 `PrivateDnsService` 中定义三个方法
- 方法签名清晰，参数最小化
- 包含完整的错误处理
- 实现重试和分页逻辑
- 方法可复用于其他资源或数据源

#### Scenario: Define service layer methods
**Given** 需要管理 Private DNS 账号  
**When** 定义服务层接口  
**Then** 提供以下方法:

```go
type PrivateDnsService struct {
    client *connectivity.TencentCloudClient
}

// Query account by Uin with pagination and filtering
func (me *PrivateDnsService) DescribePrivateDnsAccountByUin(
    ctx context.Context, 
    uin string,
) (account *privatedns.PrivateDNSAccount, errRet error)

// Create account association
func (me *PrivateDnsService) CreatePrivateDnsAccount(
    ctx context.Context,
    uin string,
) (errRet error)

// Delete account association
func (me *PrivateDnsService) DeletePrivateDnsAccount(
    ctx context.Context,
    uin string,
) (errRet error)
```

#### Scenario: Encapsulate pagination logic in service layer
**Given** DescribePrivateDNSAccountList API 支持分页  
**When** 实现 `DescribePrivateDnsAccountByUin` 方法  
**Then** 方法内部处理所有分页逻辑  
**And** 调用者无需关心分页细节  
**And** 方法返回单个账号或 nil

---

### Requirement: PDNS-ACCT-009 - Documentation Completeness
**Priority**: Medium  
**Type**: Non-Functional

资源 MUST 有完整的用户文档。

**Acceptance Criteria**:
- 源文档文件存在
- 包含资源描述
- 包含使用示例
- 包含参数说明
- 包含属性说明
- 包含导入说明
- 网站文档自动生成

#### Scenario: Provide comprehensive usage examples
**Given** 用户查看资源文档  
**When** 阅读示例部分  
**Then** 文档包含以下示例:
1. 基础使用 - 创建关联账号
2. 访问计算属性 - 输出账号信息
3. 导入现有账号

#### Scenario: Document all parameters and attributes
**Given** 用户查看参数参考  
**When** 阅读 Argument Reference 部分  
**Then** 列出 `account_uin` 参数  
**And** 说明类型、是否必填、ForceNew 行为  
**When** 阅读 Attributes Reference 部分  
**Then** 列出 `id`, `account`, `nickname` 属性  
**And** 说明它们是计算字段

#### Scenario: Provide import instructions
**Given** 用户需要导入现有账号  
**When** 阅读 Import 部分  
**Then** 提供清晰的导入命令  
**And** 说明需要提供 Uin 作为导入 ID  
**And** 给出具体示例

```markdown
## Import

Private DNS account association can be imported using the account Uin, e.g.

```bash
$ terraform import tencentcloud_private_dns_account.example 100123456789
```
```

---

## Success Metrics

- **功能完整性**: 所有 9 个需求完全实现
- **测试覆盖**: 100% 场景覆盖率（15 个场景全部测试）
- **代码质量**: 通过所有 linter 检查
- **文档质量**: 用户文档完整且准确
- **性能**: Read 操作在大量账号下仍保持高效（分页优化）

---

## Related Changes

无（独立功能）

---

## Notes

- **API 限制**: 
  - CreatePrivateDNSAccount 限频 10次/秒
  - DeletePrivateDNSAccount 和 DescribePrivateDNSAccountList 限频 20次/秒
- **依赖关系**: 与 `tencentcloud_private_dns_zone_vpc_attachment` 资源相关（VPC 绑定场景）
- **实现参考**: 参考现有 Private DNS 资源实现模式
- **SDK 版本**: 使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028`
