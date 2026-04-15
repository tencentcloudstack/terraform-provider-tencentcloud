## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘加速产品，提供全球分布的边缘节点和智能加速服务。DNS 记录管理是 TEO 的核心功能之一，允许用户将域名解析到不同的目标地址，实现流量调度和负载均衡。

当前 Terraform Provider for TencentCloud 已支持 TEO 部分资源（如站点、域名配置等），但缺少 DNS 记录的资源支持。用户无法通过 Terraform 以基础设施即代码的方式管理 TEO 的 DNS 记录，这导致用户需要手动操作控制台或使用其他工具，无法实现完整的 IaC 自动化。

Terraform Provider 基于 Terraform Plugin SDK v2 构建，使用 Go 语言开发，并通过 tencentcloud-sdk-go 调用腾讯云 API。资源代码组织在 `tencentcloud/services/teo/` 目录下，每个资源对应一个 Go 文件，遵循命名规范 `resource_tc_teo_<name>.go`。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_v13` 资源的完整 CRUD 操作
- 支持所有 DNS 记录类型（A、AAAA、CNAME、TXT、MX、NS、CAA、SRV）
- 支持记录的配置参数（TTL、Location、Weight、Priority）
- 实现异步操作的轮询机制，确保操作完成后返回
- 提供完整的单元测试和文档
- 确保与现有 Terraform Provider 架构和编码规范一致

**Non-Goals:**
- 实现批量导入/导出 DNS 记录的功能（可通过 terraform import 单独导入）
- 提供 DNS 记录的健康检查或监控功能（属于数据源范畴）
- 实现高级的 DNS 策略配置（如基于地理位置的智能调度）

## Decisions

### 1. Resource ID Design

**Decision:** 使用复合 ID 格式 `zone_id#record_id` 作为 Terraform 资源 ID。

**Rationale:**
- `zone_id` 是站点标识，`record_id` 是 DNS 记录的唯一标识
- 使用 `#` 分隔符遵循现有 Terraform Provider 的复合 ID 惯例
- 复合 ID 在全局范围内唯一，避免冲突
- 便于实现资源导入功能（`terraform import tencentcloud_teo_dns_record_v13.example zone_id#record_id`）

**Alternatives Considered:**
- 仅使用 `record_id`：不同站点可能存在相同的 record_id，会导致冲突
- 使用 `zone_id#name#type`：DNS 记录的名称和类型可能被修改，不适合作为标识

### 2. Read Operation Implementation

**Decision:** 通过 `DescribeDnsRecords` API 查询单个 DNS 记录，使用 `record_id` 过滤器。

**Rationale:**
- DescribeDnsRecords API 支持通过 `Filters` 参数过滤，其中 `id` 过滤器支持精确匹配
- API 一次请求可以返回匹配的 DNS 记录，默认返回所有字段
- 相比先查询列表再遍历查找，使用过滤器更高效且减少网络请求

**Alternatives Considered:**
- 先查询所有记录再在本地过滤：效率低，当记录数量多时性能差
- 使用专门的 Get API：TEO 云 API 没有提供单个记录的查询接口

### 3. Update Operation Strategy

**Decision:** 使用 `ModifyDnsRecords` API 批量修改 DNS 记录，每次更新仅修改当前资源。

**Rationale:**
- ModifyDnsRecords API 支持批量修改，但 Terraform 更新是逐个资源进行的
- API 接受 `DnsRecord` 数组，每个记录包含完整的配置参数
- 通过在请求中只包含需要修改的记录，避免批量操作带来的复杂性和事务问题
- 遵循 Terraform Provider 的标准更新模式，每个资源独立更新

**Alternatives Considered:**
- 批量修改多个资源：Terraform 的资源更新是异步的，批量操作难以协调状态
- 使用专用的 Update API：TEO 云 API 没有提供单个记录的更新接口

### 4. Delete Operation Implementation

**Decision:** 使用 `DeleteDnsRecords` API 删除 DNS 记录，传入单个 record_id。

**Rationale:**
- DeleteDnsRecords API 支持批量删除，但接受 `[]*string` 类型参数
- 传入单个 record_id 不会增加复杂度，保持 API 调用的一致性
- 资源删除是幂等操作，删除不存在的记录不会返回错误

**Alternatives Considered:**
- 使用专用的 Delete API：TEO 云 API 没有提供单个记录的删除接口

### 5. Async Operation Handling

**Decision:** 在 Create、Update、Delete 操作后调用 DescribeDnsRecords 进行轮询，直到操作生效。

**Rationale:**
- TEO DNS 记录操作是异步的，API 返回后记录可能尚未生效
- 通过轮询 DescribeDnsRecords 可以确保记录状态与预期一致
- 使用 helper.Retry 实现带超时的重试机制
- Create 操作轮询直到记录存在且状态为 "enable"
- Update 操作轮询直到记录参数与更新值一致
- Delete 操作轮询直到记录不存在

**Alternatives Considered:**
- 不等待操作直接返回：可能导致后续操作（如 Import）失败，因为记录尚未生效
- 使用 API 返回的 RequestId 查询状态：TEO 云 API 没有提供查询任务状态的接口

### 6. Validation Strategy

**Decision:** 在 Schema 层面进行基本验证（类型、必填、范围），在 Cloud API 层面进行业务验证。

**Rationale:**
- Schema 验证可以快速失败，减少不必要的 API 调用
- Terraform Plugin SDK 提供了丰富的 Schema 验证机制
- 业务规则验证（如记录类型的特殊要求）由 Cloud API 处理，避免 Provider 逻辑与 API 规则不一致
- TTL、Weight、Priority 的范围验证在 Schema 中实现，提供更友好的错误提示

**Alternatives Considered:**
- 所有验证都在 Provider 层实现：需要与 Cloud API 保持一致，维护成本高
- 所有验证都依赖 Cloud API：错误提示不够友好，需要额外的网络请求

### 7. Status Management

**Decision:** 将 DNS 记录的 status（enable/disable）作为 Computed 属性，不支持在创建时指定。

**Rationale:**
- CreateDnsRecord API 不接受 status 参数，创建后默认启用
- status 是系统维护的属性，不应通过 Terraform 直接修改
- 用户如需禁用记录，可以使用 ModifyDnsRecordsStatus API（但不在本次变更范围）
- 将 status 作为 Computed 属性，确保 Terraform state 与实际状态一致

**Alternatives Considered:**
- 将 status 作为可选参数：Cloud API 不支持，强行实现会增加复杂度
- 忽略 status 属性：用户无法了解记录的实际状态，不符合最佳实践

### 8. Timeouts Configuration

**Decision:** 在 Schema 中声明 Timeouts 块，支持用户自定义 Create、Update、Delete 超时时间。

**Rationale:**
- 网络环境或云 API 响应时间不确定，用户可能需要调整超时配置
- 默认超时设置为 Create=10m、Update=5m、Delete=5m，覆盖大多数场景
- 使用 resource.Timeout 协议实现，与现有资源保持一致
- 超时后返回清晰的错误信息，提示用户调整配置或检查网络

**Alternatives Considered:**
- 使用固定的超时时间：无法适应不同的网络环境和用户需求
- 不设置超时：可能导致资源操作无限等待，用户体验差

### 9. Unit Testing Strategy

**Decision:** 使用 mock 方式实现单元测试，不依赖真实的 Cloud API。

**Rationale:**
- 单元测试应快速、可重复，不应依赖外部网络和服务
- 使用 mock 可以模拟各种 API 响应场景（成功、失败、超时）
- 符合 TDD 最佳实践，便于 CI/CD 集成
- 验收测试（TF_ACC=1）负责验证与真实 Cloud API 的集成

**Alternatives Considered:**
- 使用真实 Cloud API：测试速度慢，不稳定，需要认证凭据
- 不编写测试：不符合质量要求，难以保证代码正确性

### 10. Helper Functions

**Decision:** 如需要，在独立的 `tea_dns_record_v13_helper.go` 文件中实现辅助函数。

**Rationale:**
- 将复杂的逻辑（如参数转换、过滤构建）抽离到辅助函数
- 保持资源文件的简洁性，提高代码可读性和可维护性
- 辅助函数可以被其他资源或数据源复用（如未来实现 DNS 记录数据源）
- 遵循单一职责原则

**Alternatives Considered:**
- 所有逻辑都在资源文件中实现：代码冗长，难以阅读和维护
- 不抽取辅助函数：重复代码多，不符合 DRY 原则

## Risks / Trade-offs

**Risk 1: Async operation polling timeout**
[Risk] 网络不稳定或云 API 响应慢导致轮询超时，资源操作失败。
→ Mitigation: 提供合理的默认超时时间，支持用户自定义超时配置。在错误提示中说明可能的原因和解决方法。

**Risk 2: State drift detection**
[Risk] DNS 记录在 Terraform 外部被修改（如通过控制台），导致 Terraform state 与实际状态不一致。
→ Mitigation: 实现完整的 Read 操作，用户可以运行 `terraform refresh` 同步状态。在文档中说明最佳实践（避免跨工具管理）。

**Risk 3: Complex validation logic**
[Risk] 不同 DNS 记录类型有不同的格式要求（如 SRV 记录的特殊格式），验证逻辑复杂。
→ Mitigation: 将大部分验证委托给 Cloud API，Provider 只做基本验证。在文档中说明各记录类型的格式要求。

**Risk 4: Backward compatibility**
[Risk] 未来 TEO Cloud API 可能变更，导致现有资源无法正常工作。
→ Mitigation: 使用 vendor 模式管理依赖，确保 API 版本固定。定期检查 API 变更公告，及时更新资源。

**Risk 5: Performance impact of polling**
[Risk] 频繁轮询 DescribeDnsRecords API 可能增加云 API 调用量和成本。
→ Mitigation: 使用指数退避策略，设置合理的轮询间隔。在文档中说明 API 调用频率限制。

**Trade-off 1: Batch vs Individual operations**
[Trade-off] ModifyDnsRecords 和 DeleteDnsRecords 支持批量操作，但 Terraform 资源是逐个更新的。
→ Rationale: 单资源更新更符合 Terraform 的资源模型，简化状态管理和错误处理。批量优化可以在未来通过数据源实现。

**Trade-off 2: Validation complexity**
[Trade-off] 在 Schema 中实现完整验证 vs 依赖 Cloud API 验证。
→ Rationale: Schema 验证提供快速失败和友好错误，Cloud API 验证确保业务逻辑正确性。两者结合在成本和收益之间取得平衡。

**Trade-off 3: Helper function abstraction**
[Trade-off] 抽取辅助函数 vs 所有逻辑在资源文件中。
→ Rationale: 辅助函数提高代码可读性和可维护性，但增加了文件数量和抽象层级。对于相对简单的资源，可以将逻辑保留在资源文件中，根据实际需求决定是否抽取。

## Migration Plan

本次变更为新增资源，不涉及现有资源的修改，无需迁移计划。

用户升级到新版本后，即可使用 `tencentcloud_teo_dns_record_v13` 资源。对于已存在的 DNS 记录，用户可以使用 `terraform import` 命令导入到 Terraform state 中：

```bash
terraform import tencentcloud_teo_dns_record_v13.example zone_id#record_id
```

## Open Questions

无。所有技术决策和实现细节已在设计阶段明确。
