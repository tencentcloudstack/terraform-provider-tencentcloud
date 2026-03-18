# Spec: CFS File System Timeout 配置能力

**能力 ID**: `timeout-configuration`  
**关联变更**: `add-cfs-file-system-timeout-block`  
**类型**: Capability (能力规范)  
**优先级**: Medium

---

## ADDED Requirements

### Requirement: CFS File System 资源 MUST 支持 timeouts 块配置

**优先级**: Medium  
**类型**: Feature Enhancement  

**描述**:
`tencentcloud_cfs_file_system` 资源 MUST 支持 Terraform 标准的 `timeouts` 块,允许用户为 Create、Update、Delete 操作配置自定义超时时间。资源 SHALL 提供合理的默认超时值,确保所有类型的文件系统(包括 Turbo 系列)都能在默认超时内完成操作。

**适用范围**:
- 所有使用 `tencentcloud_cfs_file_system` 资源的 Terraform 配置
- 所有类型的文件系统: Standard NFS (SD), High-Performance NFS (HP), Standard Turbo (TB), High-Performance Turbo (TP)

**业务价值**:
- 用户可以根据实际情况优化 Terraform 执行时间
- 特别适用于 Turbo 系列文件系统,避免默认超时不足
- 提升用户体验和灵活性,与其他资源保持一致

---

#### Scenario: 用户不指定 timeouts 块,使用默认超时

**Given** 用户在 Terraform 配置中定义 `tencentcloud_cfs_file_system` 资源  
**And** 用户没有指定 `timeouts` 块  

**When** 用户执行 `terraform apply` 创建资源  

**Then** 资源 MUST 使用默认超时值:
- Create: 20 分钟
- Update: 10 分钟  
- Delete: 10 分钟

**And** 资源在默认超时内成功创建  
**And** 创建过程包含 API 调用和状态轮询,都在同一超时范围内  

**验证方法**:
- 创建 Standard NFS 文件系统,验证在 20 分钟内完成
- 创建 Turbo 文件系统,验证在 20 分钟内完成
- 查看 Terraform 日志,确认使用默认超时值

---

#### Scenario: 用户为 Turbo 系列配置更长的 create timeout

**Given** 用户需要创建 Standard Turbo (TB) 文件系统  
**And** 用户知道 Turbo 系列创建时间较长  

**When** 用户在配置中指定 `timeouts { create = "30m" }`  
**And** 用户执行 `terraform apply`  

**Then** 资源 MUST 使用用户配置的 30 分钟超时  
**And** 资源最多等待 30 分钟直到状态变为 available  
**And** 如果 30 分钟内完成,操作成功返回  
**And** 如果超过 30 分钟,返回超时错误  

**验证方法**:
- 配置 `timeouts { create = "30m" }`
- 创建 TB 类型文件系统
- 验证 Terraform 最多等待 30 分钟

---

#### Scenario: 用户为快速 CI/CD 流程优化超时时间

**Given** 用户在 CI/CD 流程中使用 Standard NFS 文件系统  
**And** 用户希望减少 Terraform 执行时间  

**When** 用户配置 `timeouts { create = "5m", update = "3m", delete = "3m" }`  
**And** 用户执行 `terraform apply`  

**Then** 资源 MUST 使用用户配置的短超时值  
**And** Standard NFS 通常在 2-3 分钟内完成,在超时前成功  
**And** 如果操作超过配置的超时,立即返回错误而不是等待默认的 20 分钟  

**验证方法**:
- 配置短超时值
- 创建 Standard NFS 文件系统
- 验证在配置的超时内完成或失败

---

#### Scenario: 用户只配置 create timeout,其他使用默认值

**Given** 用户只需要调整 create 操作的超时  
**And** 用户对 update 和 delete 的默认超时满意  

**When** 用户配置 `timeouts { create = "25m" }`(只指定 create)  
**And** 用户执行 create/update/delete 操作  

**Then** Create 操作 MUST 使用 25 分钟超时  
**And** Update 操作 MUST 使用默认 10 分钟超时  
**And** Delete 操作 MUST 使用默认 10 分钟超时  

**验证方法**:
- 配置只包含 create timeout
- 执行各种操作,验证各自使用正确的超时

---

#### Scenario: 超时错误时,用户可以调整配置重试

**Given** 用户首次创建 Turbo 文件系统使用默认 20 分钟超时  
**And** 由于网络或云端原因,创建超过 20 分钟未完成  

**When** Terraform 返回超时错误  

**Then** 用户 SHOULD 能够:
1. 检查云控制台确认资源状态
2. 如果资源仍在创建中,增加 timeout 配置(如 `create = "30m"`)
3. 重新执行 `terraform apply`
4. 或使用 `terraform import` 导入已创建的资源

**And** 资源状态 MUST 保持一致,不会因超时产生僵尸资源  

**验证方法**:
- 模拟超时场景(配置很短的超时)
- 验证错误消息清晰
- 验证用户可以调整配置重试

---

#### Scenario: Update 操作支持自定义超时

**Given** 用户有一个已存在的 CFS 文件系统  
**And** 用户需要修改名称或访问组  

**When** 用户配置 `timeouts { update = "15m" }`  
**And** 用户修改资源配置并执行 `terraform apply`  

**Then** Update 操作 MUST 使用 15 分钟超时  
**And** 修改操作在超时内完成  
**And** 如果多个字段同时修改,每个修改操作都有完整的超时时间  

**验证方法**:
- 配置 update timeout
- 修改文件系统名称
- 修改访问组
- 验证各操作使用正确超时

---

#### Scenario: Delete 操作支持自定义超时

**Given** 用户需要删除 CFS 文件系统  
**And** 用户担心删除操作可能需要较长时间  

**When** 用户配置 `timeouts { delete = "15m" }`  
**And** 用户执行 `terraform destroy`  

**Then** Delete 操作 MUST 使用 15 分钟超时  
**And** 删除操作在超时内完成  
**And** 如果超时,返回明确的错误信息  

**验证方法**:
- 配置 delete timeout
- 删除文件系统
- 验证使用正确超时

---

## 技术规范

### 资源定义

**MUST** 在 `ResourceTencentCloudCfsFileSystem()` 函数中添加 `Timeouts` 字段:

```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(20 * time.Minute),
    Update: schema.DefaultTimeout(10 * time.Minute),
    Delete: schema.DefaultTimeout(10 * time.Minute),
},
```

### CRUD 函数

**MUST** 使用 `d.Timeout()` 替代硬编码的超时值:

- **Create**: `resource.Retry(d.Timeout(schema.TimeoutCreate), ...)`
- **Update**: `resource.Retry(d.Timeout(schema.TimeoutUpdate), ...)`
- **Delete**: `resource.Retry(d.Timeout(schema.TimeoutDelete), ...)`

**MUST NOT** 修改 Read 函数的超时(Terraform SDK 不支持 Read timeout)

### 默认超时值

**MUST** 提供以下默认超时值:

| 操作 | 默认值 | 说明 |
|------|-------|------|
| Create | 20 分钟 | 足够覆盖所有类型文件系统(包括 Turbo) |
| Update | 10 分钟 | 足够完成配置修改 |
| Delete | 10 分钟 | 足够完成删除操作 |

**SHOULD** 在文档中说明不同类型文件系统的典型创建时间,帮助用户选择合适的超时值。

### 向后兼容性

**MUST** 保持完全向后兼容:

- 使用 `schema.DefaultTimeout()` 确保现有配置无需修改
- 不修改资源 ID 格式或状态结构
- 不引入破坏性变更

### 错误处理

**MUST** 在超时时返回清晰的错误信息:

```
Error: timeout while waiting for state to become 'available'
```

**SHOULD** 在错误信息中提示用户可以增加 timeout 配置或检查云控制台。

### 文档要求

**MUST** 在资源文档中添加 Timeouts 章节,包含:

1. 支持的 timeout 类型及默认值
2. 使用示例(默认和自定义)
3. 不同类型文件系统的建议超时值

**SHOULD** 提供针对 Turbo 系列的 timeout 配置最佳实践。

---

## 测试要求

### 功能测试

**MUST** 验证:
- [ ] 不指定 timeouts,使用默认值
- [ ] 指定自定义 timeouts,使用用户配置
- [ ] 只指定部分 timeout,其他使用默认值
- [ ] Create/Update/Delete 操作都正确使用配置的超时

### 兼容性测试

**MUST** 验证:
- [ ] 现有配置(不含 timeouts 块)继续正常工作
- [ ] 所有现有验收测试通过
- [ ] 不同类型文件系统(SD/HP/TB/TP)都正常工作

### 边界测试

**SHOULD** 验证:
- [ ] 超时触发时的错误处理
- [ ] 极短超时(如 1秒)的行为
- [ ] 极长超时(如 2小时)的行为

---

## 影响范围

### 代码文件

**修改的文件**:
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go`
  - 添加 Timeouts 字段
  - 修改 Create 函数(合并两个 retry 块)
  - 修改 Update 函数(2处)
  - 修改 Delete 函数(1处)

**修改的文档**:
- `website/docs/r/cfs_file_system.html.markdown`
  - 添加 Timeouts 章节
  - 添加使用示例

### 影响的用户

**受益用户**:
- 使用 Turbo 系列文件系统的用户
- 在 CI/CD 中使用 CFS 的用户
- 需要精确控制 Terraform 执行时间的用户

**无影响用户**:
- 不指定 timeouts 块的现有用户(使用默认值,行为基本一致)

---

## 非功能需求

### 性能

**MUST** 确保:
- 默认超时值足够覆盖所有正常场景
- 不引入额外的性能开销
- 状态轮询频率合理(由 resource.Retry 退避策略控制)

### 可维护性

**SHOULD** 确保:
- 代码清晰易懂,符合项目编码规范
- 减少重复代码(Create 函数合并 retry 块)
- 文档完整准确

### 可扩展性

**SHOULD** 为未来扩展留空间:
- 如果需要为不同类型配置不同默认值,易于扩展
- 如果需要添加更多超时类型,易于扩展

---

## 依赖关系

### 前置条件
- Terraform Plugin SDK v2 已安装
- Go 1.17+ 环境

### 相关能力
- 无特定依赖其他能力
- 与项目中其他资源的 timeout 实现保持一致

---

## 参考实现

**类似实现**:
- `tencentcloud_instance`: Create timeout
- `tencentcloud_mysql_instance`: Create/Delete timeout
- `tencentcloud_cynosdb_cluster_slave_zone`: 完整的 CRUD timeout

**Terraform 标准**:
- [ResourceTimeout](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ResourceTimeout)
- [Timeouts Documentation](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts)

---

**审查人**: _待定_  
**审查日期**: _待定_  
**批准状态**: 待审批
