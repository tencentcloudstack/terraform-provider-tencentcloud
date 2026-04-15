## Context

当前 Terraform Provider for TencentCloud 已经支持多个 TEO 相关资源，但缺少 DNS 记录管理能力。用户需要能够通过 Terraform 代码来管理 TEO 站点的 DNS 记录，包括创建、读取、更新和删除操作。

TEO SDK (v20220901) 已经提供了完整的 DNS 记录管理 API：
- CreateDnsRecord: 创建单个 DNS 记录
- DescribeDnsRecords: 查询 DNS 记录列表，支持过滤和排序
- ModifyDnsRecords: 批量修改 DNS 记录
- DeleteDnsRecords: 批量删除 DNS 记录

现有的 TEO 资源实现模式可以参考，文件结构位于 `tencentcloud/services/teo/` 目录。

## Goals / Non-Goals

**Goals:**
- 实现完整的 Terraform 资源 `tencentcloud_teo_dns_record_v11`，支持 CRUD 操作
- 支持所有 DNS 记录类型：A、AAAA、MX、CNAME、TXT、NS、CAA、SRV
- 支持所有可配置字段：ZoneId、Name、Type、Content、TTL、Weight、Priority、Location
- 确保资源的幂等性和一致性处理
- 提供完整的文档和示例
- 提供单元测试和集成测试

**Non-Goals:**
- 不支持 DNS 记录的批量操作（虽然云 API 支持，但 Terraform 资源模型是单个资源）
- 不支持 DNS 记录状态切换（enable/disable），这需要使用 ModifyDnsRecordsStatus 接口，但不在本次实现范围
- 不实现数据源（仅实现资源）

## Decisions

### 1. 资源标识符 (ID) 设计

**决策：** 使用复合 ID `zoneId#recordId` 格式

**理由：**
- DNS 记录在云 API 中由 ZoneId 和 RecordId 唯一标识
- 复合 ID 格式遵循 Provider 的现有模式，便于解析和维护
- RecordId 由云 API 创建后返回，ZoneId 需要在用户配置中指定

**替代方案考虑：**
- 使用单一 RecordId：不合适，因为 RecordId 在不同 Zone 下可能重复
- 使用 ZoneId 作为独立字段：不符合 Terraform 资源的最佳实践

### 2. 读取策略 (Read) 设计

**决策：** 使用 DescribeDnsRecords 接口，通过 RecordId 过滤获取单个记录

**理由：**
- DescribeDnsRecords 支持通过 Filters 参数按 RecordId 过滤
- 可以获取完整的 DNS 记录信息，包括只读字段（Status、CreatedOn、ModifiedOn）
- 保持与其他 TEO 资源的一致性

**实现细节：**
- 使用 Filters 参数：`[{"name": "id", "values": ["<recordId>"]}]`
- 期望返回的 TotalCount 为 1，且 DnsRecords[0] 的 RecordId 匹配

### 3. 更新策略 (Update) 设计

**决策：** 使用 ModifyDnsRecords 接口，传递完整的 DnsRecord 对象

**理由：**
- ModifyDnsRecords 是云 API 提供的标准更新接口
- 需要传递完整的 DnsRecord 对象，包括 RecordId、ZoneId 以及所有可更新字段
- 注意 ZoneId 参数在 ModifyDnsRecords 请求中会被忽略，但仍然需要传递

**实现细节：**
- 从 state 中读取 RecordId
- 构建 DnsRecord 对象，包含所有可更新字段（排除只读字段）
- 调用 ModifyDnsRecords 接口
- 更新后调用 Read 接口刷新状态

### 4. 删除策略 (Delete) 设计

**决策：** 使用 DeleteDnsRecords 接口，传递 RecordId 数组

**理由：**
- DeleteDnsRecords 是云 API 提供的标准删除接口
- 支持批量删除，但我们只删除单个记录

**实现细节：**
- 从复合 ID 中解析出 RecordId
- 构建参数：ZoneId（从 state 读取）和 RecordIds（包含单个 RecordId）
- 调用 DeleteDnsRecords 接口

### 5. 幂等性和错误处理

**决策：** 遵循 Provider 标准的幂等性模式和错误处理策略

**理由：**
- 确保重复操作不会导致资源不一致
- 提供良好的用户体验，明确区分临时错误和永久错误

**实现细节：**
- 使用 `helper.Retry()` 实现最终一致性重试
- 使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck()` 检查不一致状态
- 对于资源不存在的错误，返回 ResourceNotFound 以正确处理

### 6. 异步操作处理

**决策：** 如果云 API 标记为异步接口，则在操作后调用 Read 接口轮询直到生效

**理由：**
- 确保 Terraform state 与云资源状态保持一致
- 遵循 Terraform Provider 的最佳实践

**实现细节：**
- 使用 `schema.Timeout` 配置默认超时时间
- 在 Create/Update/Delete 后调用 Read 接口
- 使用 `resource.Retry()` 实现轮询逻辑

## Risks / Trade-offs

### 风险 1: DNS 记录修改的限制
**风险：** ModifyDnsRecords 接口要求传递完整的 DnsRecord 对象，某些字段（如 Name、Type）可能不允许修改

**缓解措施：**
- 在 Update 函数中实现字段变更检测，只允许修改可更新的字段
- 如果用户尝试修改不允许修改的字段，返回明确的错误信息
- 在文档中明确说明哪些字段可以修改，哪些字段创建后不可修改

### 风险 2: 分页查询的边界情况
**风险：** DescribeDnsRecords 接口返回分页数据，可能存在边界情况

**缓解措施：**
- 在 Read 函数中，如果使用 RecordId 过滤，理论上只会返回一条记录
- 设置合理的 Limit 参数（如 1 或 10）避免不必要的查询
- 验证返回的 TotalCount 和 DnsRecords 数量

### 风险 3: 复合 ID 解析错误
**风险：** 复合 ID `zoneId#recordId` 可能包含特殊字符或格式不正确

**缓解措施：**
- 在 ID 解析时使用正则表达式验证格式
- 提供明确的错误信息，指导用户正确使用 ID
- 在文档中说明 ID 的格式和组成

### 权衡 1: 不支持批量操作
**权衡：** 虽然 ModifyDnsRecords 和 DeleteDnsRecords 支持批量操作，但 Terraform 资源模型是单个资源

**理由：**
- Terraform 的资源模型是声明式的，每个资源代表一个云资源实例
- 批量操作可以通过 Terraform 的 `for_each` 或 `count` 机制实现
- 保持单个资源的简单性和一致性

### 权衡 2: ZoneId 作为配置字段
**权衡：** ZoneId 既出现在配置中，也出现在复合 ID 中

**理由：**
- ZoneId 是创建资源时的必需参数
- 从复合 ID 中解析 ZoneId 可以避免用户重复配置
- 但为了向后兼容和清晰性，我们选择将 ZoneId 作为配置字段，并在 ID 中重复包含

## Open Questions

无
