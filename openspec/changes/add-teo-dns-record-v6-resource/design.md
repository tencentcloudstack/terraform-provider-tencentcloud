## Context

TEO（边缘安全加速平台）是腾讯云提供的边缘计算服务，用户需要通过 Terraform 管理 TEO 的 DNS 记录。当前 provider 缺少对 TEO DNS 记录 v6 版本的管理能力，用户无法实现基础设施即代码的统一管理。

TEO 服务提供了四个 CAPI 接口：
- CreateDnsRecord: 创建 DNS 记录
- DescribeDnsRecords: 查询 DNS 记录列表
- ModifyDnsRecords: 批量修改 DNS 记录
- DeleteDnsRecords: 批量删除 DNS 记录

这些接口位于 tencentcloud-sdk-go 的 teo/v20220901 包中。

## Goals / Non-Goals

**Goals:**
- 实现完整的 Terraform 资源 `tencentcloud_teo_dns_record_v6`，支持 CRUD 操作
- 支持所有 TEO DNS 记录类型（A、AAAA、CNAME、TXT、MX、NS、CAA、SRV）
- 实现参数验证（TTL 范围 60-86400，权重范围 -1~100，MX 优先级范围 0~50）
- 支持可选参数：Location（仅限 A/AAAA/CNAME）、Weight（仅限 A/AAAA/CNAME）、Priority（仅限 MX）、TTL
- 生成完整的单元测试和验收测试代码
- 添加完整的文档

**Non-Goals:**
- 不实现 DNS 记录的批量导入/导出功能
- 不实现 DNS 记录的版本控制或变更历史
- 不实现 DNS 记录的高级查询或分析功能

## Decisions

### 1. Resource 命名
**Decision:** 使用 `tencentcloud_teo_dns_record_v6` 作为资源名称

**Rationale:**
- 遵循 provider 的命名规范：`tencentcloud_teo_<resource>`
- 使用 `_v6` 后缀以区别于可能存在的其他版本 DNS 记录资源
- 清晰表达资源所属服务（TEO）和类型（DNS record）

### 2. Resource ID 格式
**Decision:** 使用复合 ID 格式 `zoneId#recordId`

**Rationale:**
- 遵循 provider 的复合 ID 模式（如 "instanceId#userId"）
- `zoneId` 是必需的查询参数，用于标识站点
- `recordId` 是 DNS 记录的唯一标识
- 使用 `#` 作为分隔符，便于解析和识别

### 3. Schema 参数映射
**Decision:** 直接映射 CAPI 接口参数到 Terraform Schema

| Terraform 参数 | CAPI 参数 | 类型 | 必需 | 说明 |
|---------------|----------|------|------|------|
| zone_id | ZoneId | string | 是 | 站点 ID |
| name | Name | string | 是 | DNS 记录名 |
| type | Type | string | 是 | DNS 记录类型 |
| content | Content | string | 是 | DNS 记录内容 |
| location | Location | string | 否 | 解析线路（仅限 A/AAAA/CNAME） |
| ttl | TTL | int | 否 | 缓存时间（60-86400，默认 300） |
| weight | Weight | int | 否 | 权重（-1~100，仅限 A/AAAA/CNAME） |
| priority | Priority | int | 否 | MX 优先级（0~50，仅限 MX） |
| record_id | RecordId | string | 是（Computed）| DNS 记录 ID |

**Rationale:**
- 保持与 CAPI 接口的一致性，降低理解和维护成本
- 使用 snake_case 符合 Terraform 命名规范
- ZoneId 作为复合 ID 的一部分，在 Schema 中仍然声明为 Required，便于用户理解
- RecordId 作为 Computed 字段，由 Create 操作返回，用于后续的 Read、Update、Delete 操作

### 4. CRUD 操作实现策略

**Create:**
- 调用 `CreateDnsRecord` 接口
- 解析返回的 `RecordId`，设置到 state 中
- 使用复合 ID 格式 `zoneId#recordId` 设置 resource ID

**Read:**
- 从 resource ID 解析出 `zoneId` 和 `recordId`
- 调用 `DescribeDnsRecords` 接口，使用 Filters 参数过滤 `id` 字段
- 如果找到记录，更新 state；否则返回 "not found"

**Update:**
- 构造 `DnsRecord` 结构体，包含 `recordId` 和需要更新的字段
- 调用 `ModifyDnsRecords` 接口（支持批量修改，但这里只修改单条记录）
- 更新 state

**Delete:**
- 从 resource ID 解析出 `zoneId` 和 `recordId`
- 调用 `DeleteDnsRecords` 接口，传入 `recordId` 列表（只包含单条记录）
- 清除 state

**Rationale:**
- CreateDnsRecord 和 ModifyDnsRecords 使用不同的数据结构，需要区分处理
- DescribeDnsRecords 是列表查询接口，需要使用 Filters 参数定位单条记录
- DeleteDnsRecords 支持批量删除，但这里只删除单条记录
- 所有操作都需要 zoneId 参数，通过解析 resource ID 获取

### 5. 参数验证
**Decision:** 在 Schema 中使用 `ValidateFunc` 进行参数验证

**Rationale:**
- 在 terraform apply 阶段就能发现参数错误，而不是等到调用 API
- 提供清晰的错误提示，改善用户体验
- 减少无效的 API 调用

验证规则：
- TTL: 60 <= ttl <= 86400
- Weight: -1 <= weight <= 100（仅限 A/AAAA/CNAME 类型）
- Priority: 0 <= priority <= 50（仅限 MX 类型）

### 6. 错误处理
**Decision:** 使用 provider 标准错误处理模式

**Rationale:**
- 使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck()` 处理最终一致性问题
- 使用 `helper.Retry()` 实现重试逻辑（针对最终一致性）
- 对于 ResourceNotFound 错误，在 Delete 操作中返回成功（幂等性）

### 7. 测试策略
**Decision:** 实现单元测试和验收测试

**单元测试:**
- 使用 mock 方法将云 API 调用 mock 掉
- 测试业务代码逻辑，不依赖实际云 API
- 测试参数验证逻辑
- 测试 ID 解析和构造逻辑

**验收测试:**
- 使用 TF_ACC=1 运行实际的 API 调用
- 测试完整的 CRUD 流程
- 测试所有支持的记录类型
- 测试参数验证边界情况
- 需要 TENCENTCLOUD_SECRET_ID/KEY 环境变量

**Rationale:**
- 单元测试快速执行，不依赖外部环境，适合 CI/CD
- 验收测试验证与实际 API 的集成，确保代码在生产环境可用
- 两者结合提供全面的测试覆盖

## Risks / Trade-offs

### Risk 1: DescribeDnsRecords 接口返回列表而非单条记录
**Risk:** DescribeDnsRecords 是列表查询接口，可能返回多条记录或没有匹配记录

**Mitigation:**
- 使用 Filters 参数精确匹配 `id` 字段
- 检查返回结果数量，如果为 0 则返回 ResourceNotFound，如果 >1 则返回错误（理论上不应该发生）
- 优先匹配 `zoneId` 和 `recordId` 都正确的记录

### Risk 2: ModifyDnsRecords 接口是批量操作
**Risk:** ModifyDnsRecords 接口设计为批量操作，单条记录修改可能效率不高

**Mitigation:**
- 目前只实现单条记录修改，符合 Terraform 资源模型
- 如果未来需要批量更新，可以考虑添加独立的资源或数据源
- 接口本身支持单条记录修改，不影响功能

### Risk 3: DNS 记录类型相关的参数限制
**Risk:** Location 和 Weight 仅适用于 A/AAAA/CNAME 类型，Priority 仅适用于 MX 类型

**Mitigation:**
- 在 Schema 文档中明确说明参数适用范围
- 可以考虑使用 `DiffSuppressFunc` 或 `CustomizeDiff` 来实现动态的参数验证
- 当前选择文档说明的方式，保持简单和一致性

### Risk 4: 中文域名需要 punycode 转换
**Risk:** CAPI 接口要求中文、韩文、日文域名转换为 punycode

**Mitigation:**
- 在设计文档中明确说明这个要求
- 用户需要确保传入正确的 punycode 格式
- 不在代码中实现自动转换，避免增加复杂度和潜在的转换错误

## Migration Plan

### 部署步骤
1. 在开发环境完成代码开发和测试
2. 提交代码 review，确保符合 provider 规范
3. 合并到主分支后，发布新版本 provider
4. 用户升级到新版本后即可使用新资源

### 回滚策略
- 如果发现问题，可以在 provider 新版本中标记资源为废弃
- 不会影响现有资源的正常使用
- 可以通过禁用资源定义来移除新功能

## Open Questions

无。设计已经明确，所有技术决策已经做出。
