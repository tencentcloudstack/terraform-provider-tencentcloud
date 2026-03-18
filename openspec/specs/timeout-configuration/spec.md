# timeout-configuration Specification

## Purpose
TBD - created by archiving change add-cfs-file-system-timeout-block. Update Purpose after archive.
## Requirements
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

