# Implementation Tasks: Add Private DNS Account Resource

## Overview

本文档列出实现 `tencentcloud_private_dns_account` 资源所需的所有任务。

**总任务数**: 15  
**预估工作量**: 2.5 天

---

## Phase 1: Service Layer Implementation (5 tasks)

### Task 1.1: 在 service_tencentcloud_private_dns.go 中添加 DescribePrivateDnsAccountByUin 方法
- **位置**: `tencentcloud/services/privatedns/service_tencentcloud_private_dns.go`
- **操作**: 新增方法实现按 Uin 查询单个账号
- **详情**:
  ```go
  func (me *PrivateDnsService) DescribePrivateDnsAccountByUin(ctx context.Context, uin string) (
      account *privatedns.PrivateDNSAccount, errRet error)
  ```
  - 构建 `DescribePrivateDNSAccountListRequest`
  - 设置 Filter: `Name="AccountUin", Values=[uin]`
  - 实现分页逻辑 (limit=100, offset 递增)
  - 遍历所有页查找匹配的 Uin
  - 包含重试逻辑和错误处理
- **验证**: 方法编译通过，逻辑正确

### Task 1.2: 在 service_tencentcloud_private_dns.go 中添加 CreatePrivateDnsAccount 方法
- **位置**: 同上文件
- **操作**: 新增方法实现创建关联账号
- **详情**:
  ```go
  func (me *PrivateDnsService) CreatePrivateDnsAccount(ctx context.Context, uin string) (errRet error)
  ```
  - 构建 `CreatePrivateDNSAccountRequest`
  - 设置 `Account.Uin`
  - 调用 `CreatePrivateDNSAccount` API
  - 包含重试逻辑（处理 InternalError）
  - 处理 `InvalidParameter.AccountExist` 错误（视为成功）
- **验证**: 方法编译通过，API 调用正确

### Task 1.3: 在 service_tencentcloud_private_dns.go 中添加 DeletePrivateDnsAccount 方法
- **位置**: 同上文件
- **操作**: 新增方法实现删除关联账号
- **详情**:
  ```go
  func (me *PrivateDnsService) DeletePrivateDnsAccount(ctx context.Context, uin string) (errRet error)
  ```
  - 构建 `DeletePrivateDNSAccountRequest`
  - 设置 `Account.Uin`
  - 调用 `DeletePrivateDNSAccount` API
  - 包含重试逻辑
  - 特殊处理 `UnsupportedOperation.ExistBoundVpc` 错误
- **验证**: 方法编译通过，错误处理正确

### Task 1.4: 添加服务层错误常量
- **位置**: `tencentcloud/services/privatedns/service_tencentcloud_private_dns.go`
- **操作**: 如果需要，添加重试错误码常量
- **详情**: 参考现有 `PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR` 常量
- **验证**: 常量定义正确

### Task 1.5: 服务层代码格式化和验证
- **命令**: `gofmt -w tencentcloud/services/privatedns/service_tencentcloud_private_dns.go`
- **验证**: 代码格式正确，编译通过

---

## Phase 2: Resource Implementation (5 tasks)

### Task 2.1: 创建资源文件和 Schema 定义
- **位置**: `tencentcloud/services/privatedns/resource_tc_private_dns_account.go`
- **操作**: 创建新文件并定义资源 Schema
- **详情**:
  ```go
  func ResourceTencentCloudPrivateDnsAccount() *schema.Resource {
      return &schema.Resource{
          Create: resourceTencentCloudPrivateDnsAccountCreate,
          Read:   resourceTencentCloudPrivateDnsAccountRead,
          Delete: resourceTencentCloudPrivateDnsAccountDelete,
          Importer: &schema.ResourceImporter{
              State: schema.ImportStatePassthrough,
          },
          Schema: map[string]*schema.Schema{
              "account_uin": {
                  Type:        schema.TypeString,
                  Required:    true,
                  ForceNew:    true,
                  Description: "Uin of the associated account.",
              },
              "account": {
                  Type:        schema.TypeString,
                  Computed:    true,
                  Description: "Email of the associated account.",
              },
              "nickname": {
                  Type:        schema.TypeString,
                  Computed:    true,
                  Description: "Nickname of the associated account.",
              },
          },
      }
  }
  ```
- **验证**: Schema 定义符合规范

### Task 2.2: 实现 Create 函数
- **位置**: 同上文件
- **操作**: 实现 `resourceTencentCloudPrivateDnsAccountCreate`
- **详情**:
  - 添加日志和性能追踪: `defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.create")()`
  - 提取 `account_uin` 参数
  - 调用服务层 `CreatePrivateDnsAccount()`
  - 设置资源 ID 为 Uin: `d.SetId(uin)`
  - 调用 Read 函数获取完整信息
  - 错误处理和日志记录
- **验证**: Create 逻辑正确，错误处理完善

### Task 2.3: 实现 Read 函数
- **位置**: 同上文件
- **操作**: 实现 `resourceTencentCloudPrivateDnsAccountRead`
- **详情**:
  - 添加日志和性能追踪
  - 从资源 ID 获取 Uin
  - 调用服务层 `DescribePrivateDnsAccountByUin()`
  - 如果账号不存在，清空资源 ID: `d.SetId("")`
  - 设置 `account_uin`, `account`, `nickname` 到状态
  - 错误处理和日志记录
- **验证**: Read 逻辑正确，状态设置无误

### Task 2.4: 实现 Delete 函数
- **位置**: 同上文件
- **操作**: 实现 `resourceTencentCloudPrivateDnsAccountDelete`
- **详情**:
  - 添加日志和性能追踪
  - 从资源 ID 获取 Uin
  - 调用服务层 `DeletePrivateDnsAccount()`
  - 特殊处理 VPC 绑定错误，返回明确提示
  - 错误处理和日志记录
- **验证**: Delete 逻辑正确，错误提示清晰

### Task 2.5: 资源代码格式化和验证
- **命令**: `gofmt -w tencentcloud/services/privatedns/resource_tc_private_dns_account.go`
- **验证**: 代码格式正确，编译通过

---

## Phase 3: Provider Registration (1 task)

### Task 3.1: 在 Provider 中注册新资源
- **位置**: `tencentcloud/provider.go`
- **操作**: 在资源映射中添加新资源
- **详情**:
  ```go
  "tencentcloud_private_dns_account": privatedns.ResourceTencentCloudPrivateDnsAccount(),
  ```
- **位置**: 在 Private DNS 相关资源附近
- **验证**: Provider 编译通过，资源可用

---

## Phase 4: Testing (4 tasks)

### Task 4.1: 创建测试文件
- **位置**: `tencentcloud/services/privatedns/resource_tc_private_dns_account_test.go`
- **操作**: 创建测试文件并设置基本结构
- **详情**:
  - 导入必要的包
  - 定义测试常量（测试用的 Uin）
  - 创建测试辅助函数

### Task 4.2: 编写基础 CRUD 测试
- **位置**: 同上文件
- **操作**: 实现 `TestAccTencentCloudPrivateDnsAccountResource_Basic`
- **测试步骤**:
  1. 创建关联账号
  2. 验证账号属性（account_uin, account, nickname）
  3. 执行删除
- **验证**: 测试通过

### Task 4.3: 编写导入测试
- **位置**: 同上文件
- **操作**: 在基础测试中添加 ImportState 步骤
- **详情**:
  ```go
  {
      ResourceName:      "tencentcloud_private_dns_account.test",
      ImportState:       true,
      ImportStateVerify: true,
  }
  ```
- **验证**: 导入功能正常

### Task 4.4: 编写 ForceNew 测试
- **位置**: 同上文件
- **操作**: 验证修改 Uin 触发资源重建
- **测试步骤**:
  1. 创建关联账号（Uin1）
  2. 修改 Uin 为 Uin2
  3. 验证旧资源被删除，新资源被创建
- **验证**: ForceNew 行为正确

---

## Phase 5: Documentation (3 tasks)

### Task 5.1: 创建资源文档
- **位置**: `tencentcloud/services/privatedns/resource_tc_private_dns_account.md`
- **操作**: 创建 Markdown 文档
- **内容**:
  - 资源描述
  - 使用示例（基础场景）
  - 参数说明（Argument Reference）
  - 属性说明（Attributes Reference）
  - 导入说明（Import）
- **验证**: 文档内容完整准确

### Task 5.2: 生成网站文档
- **命令**: `make doc`
- **操作**: 运行文档生成工具
- **验证**: `website/docs/r/private_dns_account.html.markdown` 生成成功

### Task 5.3: 更新 provider.md（如需要）
- **位置**: `tencentcloud/provider.md`
- **操作**: 在 Private DNS 资源列表中添加新资源
- **验证**: Provider 文档包含新资源

---

## Phase 6: Code Quality & Integration (2 tasks)

### Task 6.1: 运行代码检查
- **命令**: 
  - `make fmt` - 格式化所有代码
  - `make lint` - 运行 linter
- **操作**: 修复所有 lint 错误和警告
- **验证**: 无 linter 错误

### Task 6.2: 运行验收测试
- **命令**: `TF_ACC=1 go test -v ./tencentcloud/services/privatedns -run TestAccTencentCloudPrivateDnsAccountResource`
- **前置条件**: 
  - 设置 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
  - 准备测试用的关联账号 Uin
- **验证**: 所有测试通过

---

## Validation Checklist

完成所有任务后，验证以下内容：

- [x] Schema 定义正确，包含所有必需字段
- [x] Create 操作成功创建关联账号
- [x] Read 操作正确获取账号信息
- [x] Delete 操作成功删除关联账号
- [x] Uin 修改触发 ForceNew（资源重建）
- [x] 支持资源导入（通过 Uin）
- [x] 分页逻辑正确处理大量账号
- [x] UIN 过滤逻辑正确
- [x] 错误处理完善（账号已存在、VPC 绑定等）
- [ ] 所有验收测试通过（需要实际测试环境）
- [x] 代码格式符合规范
- [x] 无 linter 错误或警告（仅有预存在的废弃警告）
- [x] 文档完整且准确
- [x] Provider 编译成功

---

## Dependencies

任务依赖关系：
1. Phase 1 (Service Layer) → Phase 2 (Resource)
2. Phase 2 (Resource) → Phase 3 (Registration)
3. Phase 3 (Registration) → Phase 4 (Testing)
4. Phase 4 (Testing) → Phase 5 (Documentation)
5. 所有阶段 → Phase 6 (Quality)

阶段内任务可以按顺序执行。

---

## Rollback Plan

如果出现问题：
1. **合并前**: 简单回滚代码更改
2. **合并后**: 
   - 资源是纯新增功能，不影响现有资源
   - 可以通过移除 Provider 注册禁用
   - 完全回滚不会影响用户

---

## Testing Notes

### 测试前准备

1. **准备测试账号**: 需要至少一个测试用的腾讯云账号 Uin
2. **权限要求**: 测试账号需要 Private DNS 服务权限
3. **清理资源**: 测试前确保没有遗留的测试关联账号

### 测试配置示例

```hcl
resource "tencentcloud_private_dns_account" "test" {
  account_uin = "100123456789"  # 替换为实际测试 Uin
}
```

### 手动测试步骤

1. 使用 Terraform 创建关联账号
2. 在腾讯云控制台验证账号已添加
3. 修改 Uin 验证 ForceNew 行为
4. 使用 `terraform import` 导入现有账号
5. 删除资源并验证清理

---

## Implementation Notes

### 设计要点

1. **资源 ID**: 使用 Uin 作为资源唯一标识符
2. **只读字段**: Account 和 Nickname 由系统管理，不可配置
3. **ForceNew**: Uin 是唯一可配置字段，修改需要重建资源
4. **分页处理**: Read 操作需要遍历所有页查找匹配账号
5. **过滤优化**: 使用 API Filter 参数减少数据传输

### 代码质量要求

- 遵循项目代码规范
- 使用 `tccommon` 包的通用函数
- 正确使用日志记录（logId）
- 实现性能追踪（LogElapsed）
- 完善的错误处理和上下文信息

---

**任务状态**: 待开始  
**下一步**: 执行 Phase 1 - Service Layer Implementation
