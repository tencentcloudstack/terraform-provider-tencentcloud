## Context

Terraform Provider for TencentCloud 目前已经支持多种云服务的管理功能，但对 TEO（边缘安全加速平台）产品的 DNS 记录管理功能尚不完善。用户需要通过 Terraform 来管理 TEO 的 DNS 记录，以实现基础设施即代码（IaC）的管理方式。

TEO 产品提供了完整的 DNS 记录管理 API：
- `CreateDnsRecord`：创建 DNS 记录
- `DescribeDnsRecords`：查询 DNS 记录列表
- `ModifyDnsRecords`：批量修改 DNS 记录
- `DeleteDnsRecords`：批量删除 DNS 记录

当前项目使用 Terraform Plugin SDK v2，采用标准的资源文件组织结构。资源文件命名格式为 `resource_tc_<service>_<name>.go`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_dns_record_v2` 资源，提供完整的 CRUD 功能
- 实现资源创建、读取、更新和删除的 Terraform 资源接口
- 提供单元测试，使用 mock 方式测试业务逻辑
- 遵循项目现有的编码规范和模式

**Non-Goals:**
- 不实现数据源（DataSource）功能
- 不提供 DNS 记录的批量操作功能（单次操作仅针对单个记录）
- 不实现复杂的权限管理和状态管理

## Decisions

### 资源标识符设计
使用复合 ID 格式 `zone_id#record_id` 作为 Terraform 资源的 ID，因为：
- `zone_id` 是站点 ID，是资源所属的站点标识
- `record_id` 是 DNS 记录的唯一标识符
- 使用 `#` 分隔符符合项目的命名规范
- 便于解析和状态管理

### Schema 参数设计
根据云 API 的参数定义，设计 Terraform Schema：
- **Required 参数**：`zone_id`、`name`、`type`、`content`（创建记录时必须提供）
- **Optional 参数**：`location`、`ttl`、`weight`、`priority`（根据实际需要可选）
- **Computed 参数**：`record_id`（创建后由云 API 返回）

### CRUD 接口实现方案
1. **Create**：调用 `CreateDnsRecord` API，获取返回的 `record_id` 并设置到状态中
2. **Read**：调用 `DescribeDnsRecords` API，通过 `Filters` 过滤 `record_id` 来查询单个记录
3. **Update**：调用 `ModifyDnsRecords` API，构造包含更新字段的 `DnsRecord` 对象
4. **Delete**：调用 `DeleteDnsRecords` API，传入 `record_id` 列表（单个记录）

### 错误处理和重试机制
- 使用 `defer tccommon.LogElapsed(ctx)` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck(d, &resource)` 检查数据一致性
- 对于可能的临时错误，使用 `helper.Retry()` 进行最终一致性重试

### 测试策略
由于这是新增资源，采用 mock 方式进行单元测试：
- 使用 gomonkey 库对云 API 进行 mock
- 测试重点在业务逻辑和状态管理，不依赖实际云 API
- 不使用 Terraform acceptance test（TF_ACC），避免需要真实凭证

### 超时配置
虽然云 API 没有明确标注为异步接口，但考虑到 DNS 记录的全球生效特性：
- 在 Schema 中声明 `Timeouts` 块
- 默认设置合理的超时时间
- 在 CRUD 操作中使用 context 支持超时取消

## Risks / Trade-offs

### 风险 1：批量修改接口的限制
**风险**：`ModifyDnsRecords` API 是批量修改接口，需要构造完整的 `DnsRecord` 对象，包括所有字段。

**缓解措施**：
- 在 Update 函数中，从 state 中读取当前记录的完整数据
- 仅更新用户指定的字段，其他字段保持不变
- 确保不丢失任何未修改的字段

### 风险 2：查询接口的性能
**风险**：`DescribeDnsRecords` API 返回记录列表，需要通过 Filters 过滤，可能在记录较多时性能较差。

**缓解措施**：
- 使用精确的 `record_id` 过滤条件
- 设置合理的 `Limit` 参数（如 20）
- 如果性能问题严重，可以考虑后续优化

### 风险 3：参数校验
**风险**：云 API 对某些参数有特定的格式要求（如 MX 记录的优先级、权重值的范围），Terraform 层面需要提供友好的错误提示。

**缓解措施**：
- 在 Schema 中定义合理的验证规则
- 在 Resource 函数中进行参数校验
- 提供清晰的错误信息

### 权衡 1：复合 ID vs 单一 ID
**权衡**：
- 复合 ID（`zone_id#record_id`）：便于理解和调试，但解析逻辑稍复杂
- 单一 ID（仅 `record_id`）：实现简单，但缺少上下文信息

**决策**：选择复合 ID，因为：
- 符合项目规范
- 便于在日志和错误信息中定位问题
- 可以通过辅助函数简化解析逻辑

### 权衡 2：是否支持批量操作
**权衡**：
- 支持批量操作：可以一次性管理多个记录，但增加复杂度
- 仅支持单个操作：实现简单，但效率较低

**决策**：仅支持单个操作，原因：
- Terraform 资源通常管理单个实体
- 用户可以通过 Terraform 的 count/for_each 功能实现批量管理
- 降低实现复杂度和维护成本

## Migration Plan

### 部署步骤
1. 创建资源文件 `resource_tc_teo_dns_record_v2.go`
2. 创建测试文件 `resource_tc_teo_dns_record_v2_test.go`
3. 在 `service_tencentcloud_teo.go` 中注册新资源
4. 运行单元测试确保功能正确
5. 提交代码到版本控制系统

### 回滚策略
如果出现问题，可以通过以下方式回滚：
- 从主分支回退到上一个稳定版本
- 从资源注册文件中移除新资源的注册代码
- 如果已经发布到用户，需要发布公告建议用户暂停使用该资源

## Open Questions

1. **问**：是否需要支持 `status` 字段（启用/停用 DNS 记录）？
   **答**：当前云 API 的 `ModifyDnsRecords` 接口不支持修改 `status` 字段，因此本次实现不包含该功能。如需要，可以考虑后续使用 `ModifyDnsRecordsStatus` API。

2. **问**：TTL 的默认值应该是多少？
   **答**：根据云 API 文档，TTL 的默认值是 300 秒，我们将在 Schema 中将 TTL 设置为 Optional，并在用户未指定时使用云 API 的默认值。

3. **问**：是否需要在 Terraform 中支持 DNS 记录的校验（如 MX 记录的优先级范围）？
   **答**：云 API 会进行参数校验，因此 Terraform 层面只需要进行基本的格式校验（如数值范围），详细的业务校验由云 API 负责即可。
