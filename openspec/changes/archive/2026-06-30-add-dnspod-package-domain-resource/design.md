## Context

DNSPod 付费套餐支持绑定到指定域名以提供增值 DNS 服务（如更高 QPS、更多线路等）。套餐与域名之间的绑定关系可以通过 `ModifyPackageDomain` API 进行管理（绑定、解绑、换绑），绑定状态可以通过 `DescribeDomainVipList` API 查询。当前 Terraform Provider 提供了 `tencentcloud_dnspod_package_order` 资源用于购买套餐，但缺少将套餐绑定到域名的资源。

本设计基于 vendor 目录下已存在的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323` SDK，参考现有资源 `resource_tc_dnspod_custom_line.go` 和 `resource_tc_dnspod_domain_lock.go` 的实现模式。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_dnspod_package_domain` 资源的完整 CRUD 生命周期
- Create → 绑定域名到套餐（ModifyPackageDomain, Operation="bind"）
- Read → 查询套餐域名绑定信息和状态（DescribeDomainVipList）
- Update → 换绑套餐域名（ModifyPackageDomain, Operation="change"）
- Delete → 解绑套餐域名（ModifyPackageDomain, Operation="unbind"）
- 支持 Import 导入已有绑定关系

**Non-Goals:**
- 不实现套餐购买功能（已有 `tencentcloud_dnspod_package_order`）
- 不实现套餐自动续费管理（已有 `tencentcloud_dnspod_package_order`）
- 不实现批量绑定/解绑（API 一次只操作一个套餐和域名）

## Decisions

### 1. 资源 ID 设计

**决策**: 使用 `{resource_id}#{domain_id}` 作为复合 ID。

**理由**:
- `ModifyPackageDomain` 不返回独立资源 ID（仅返回 RequestId），无法使用云端生成的 ID
- `resource_id` 唯一标识套餐，`domain_id` 唯一标识域名，两者联合唯一标识一条绑定关系
- 使用 `#` 分隔符是项目标准（参考 `resource_tc_dnspod_custom_line.go` 中 `domain#name` 的模式）
- 支持 Import 时从 ID 解析参数

### 2. Schema 参数设计

**决策**: 仅暴露两个 Required 输入参数，其他字段为 Computed 输出。

**输入参数**:
- `domain_id` (TypeInt, Required): 要绑定的域名 ID
- `resource_id` (TypeString, Required, ForceNew): 套餐资源 ID，变更会触发重建（先解绑旧套餐再绑定新套餐）

**输出参数** (来自 `PackageListItem`):
- `domain` (TypeString): 域名原始格式
- `grade` (TypeString): 套餐等级代码
- `grade_title` (TypeString): 套餐名称
- `vip_start_at` (TypeString): 付费套餐开通时间
- `vip_end_at` (TypeString): 付费套餐到期时间
- `vip_auto_renew` (TypeString): 自动续费状态（YES/NO/DEFAULT）
- `remain_times` (TypeInt): 剩余换绑/绑定次数
- `grade_level` (TypeInt): 域名等级代号
- `status` (TypeString): 绑定状态
- `is_grace_period` (TypeString): 是否宽限期
- `downgrade` (TypeBool): 是否降级

**理由**: 
- API 设计使得套餐绑定关系以 `resource_id` 为核心，`domain_id` 为可变属性
- `resource_id` 标记为 ForceNew 而非 `domain_id`，因为套餐是核心资源
- 所有 Computed 字段来自 `DescribeDomainVipList` 返回的 `PackageListItem`

### 3. CRUD API 映射

**Create**: `ModifyPackageDomain(Operation="bind", ResourceId=resource_id, NewDomainId=domain_id)`
- 响应仅有 RequestId，ID 从输入参数构造

**Read**: `DescribeDomainVipList(ResourceIdList=[resource_id])`
- 通过 ResourceIdList 精确查询指定套餐
- 遍历 PackageList，匹配 DomainId 找到绑定关系
- 未找到匹配时清空 d.SetId("")

**Update**: 仅支持 `domain_id` 变更
- `ModifyPackageDomain(Operation="change", ResourceId=resource_id, DomainId=旧domain_id, NewDomainId=新domain_id)`
- 更新成功后更新资源 ID 中的 domain_id 部分

**Delete**: `ModifyPackageDomain(Operation="unbind", ResourceId=resource_id, DomainId=domain_id)`

### 4. 重试与错误处理策略

- 所有写操作使用 `tccommon.WriteRetryTimeout`
- 所有读操作使用 `tccommon.ReadRetryTimeout`（虽然 `DescribeDomainVipList` 是读接口，但作为 Read 资源函数的标准做法）
- API 错误通过 `tccommon.RetryError(e)` 包装返回，触发 SDK 层重试
- Read 中返回空时记录日志并清空 ID，不在 retry 块内清空 ID

### 5. Import 支持

**格式**: `resource_id#domain_id`

**示例**: `terraform import tencentcloud_dnspod_package_domain.example "res-xxxxx#12345"`

Import 时使用 `schema.ImportStatePassthrough`，由 Read 函数解析 ID 并填充状态。

## Risks / Trade-offs

- **[风险] ModifyPackageDomain 无返回数据**: 无法验证操作后的实际状态。→ **缓解**: Create/Update/Delete 后立即调用 Read 确认状态
- **[风险] DescribeDomainVipList 可能有延迟**: 绑定/解绑后查询可能短暂不一致。→ **缓解**: Read 使用 ReadRetryTimeout 重试
- **[风险] resource_id 标记为 ForceNew**: 用户更换套餐需要重建资源。→ **缓解**: 这是合理行为，套餐变更本质上是新建绑定
- **[风险] 并发操作**: 同一个套餐的绑定/解绑操作可能冲突。→ **缓解**: API 层面保证，SDK 重试机制兜底
