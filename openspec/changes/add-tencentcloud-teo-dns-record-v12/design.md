## Context

当前 Terraform TencentCloud Provider 已支持 TEO（TencentCloud EdgeOne）产品的多个资源，但缺少对 DNS 记录资源的支持。TEO 是腾讯云的边缘加速服务，提供全球边缘节点加速和智能 DNS 解析能力。用户需要通过 Terraform 自动化管理 TEO 的 DNS 记录配置，以实现基础设施即代码的完整覆盖。

TEO 产品在 provider 中已有其他资源的实现参考，现有的 TEO 服务包（v20220901）已经提供了 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords` 四个 API 接口，可以支持完整的 CRUD 操作。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_v12` 资源的完整 CRUD 操作
- 支持 A、AAAA、MX、CNAME、TXT、NS、CAA、SRV 等 DNS 记录类型
- 支持解析线路、缓存时间（TTL）、权重、优先级等高级配置
- 遵循 Terraform Provider SDK v2 的标准实现模式
- 确保与现有 TEO 资源实现的一致性
- 提供完整的文档和测试

**Non-Goals:**
- 不实现 DNS 记录的批量操作（Terraform 每次操作单条记录）
- 不实现 DNS 记录的健康检查或监控功能
- 不实现 DNS 记录的版本历史管理
- 不支持 DNS 记录的导入功能（如果 API 不支持从 record_id 查询完整信息）

## Decisions

**1. 资源标识符设计**
- 使用复合 ID 格式：`zone_id#record_id`
- 决策理由：
  - `record_id` 在 TEO 产品中是 DNS 记录的唯一标识
  - `zone_id` 是站点 ID，确保同一站点下的记录
  - 复合 ID 避免了全局唯一 ID 的限制，符合 provider 常见模式
- 替代方案：仅使用 `record_id`（缺点：无法跨站点验证记录归属）

**2. Schema 设计模式**
- 必填参数：`zone_id`, `name`, `type`, `content`
- 可选参数：`location`, `ttl`, `weight`, `priority`
- 计算属性（Computed）：`record_id`, `status`, `created_on`
- 决策理由：
  - 直接映射云 API 的参数结构，减少转换复杂度
  - 将系统返回字段标记为 Computed，提供完整的状态可见性
  - 可选参数的默认值由云 API 处理，不在 provider 层面硬编码
- 替代方案：为所有可选参数设置默认值（缺点：与云 API 默认值可能不一致）

**3. Update 接口使用策略**
- 使用 `ModifyDnsRecords` 接口进行更新操作
- 请求体构建：将单条记录信息封装为 `DnsRecord` 数组
- 决策理由：
  - `ModifyDnsRecords` 是批量修改接口，可以处理单条记录
  - 接口返回成功即认为更新生效（同步接口）
  - Update 后立即调用 Read 接口轮询直到状态一致
- 替代方案：使用删除+重建模式（缺点：无法保留 `record_id`，可能影响其他依赖）

**4. 错误处理和重试机制**
- 最终一致性重试：使用 `helper.Retry()` 包裹 API 调用
- 资源状态检查：使用 `defer tccommon.InconsistentCheck()` 确保状态一致性
- 错误日志：使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 决策理由：
  - TEO 作为边缘服务，可能出现短暂的一致性延迟
  - 统一使用 provider 已有的错误处理模式，保持一致性
- 替代方案：自定义重试逻辑（缺点：增加代码复杂度，与项目风格不一致）

**5. 测试策略**
- 单元测试：使用 mock 方式模拟云 API 响应，测试业务逻辑
- 不使用 Terraform 测试套件（TF_ACC）
- 决策理由：
  - TF_ACC 测试需要真实的云环境和凭证，不适合单元测试
  - mock 测试可以快速验证 CRUD 逻辑的正确性
  - 符合项目的测试要求："生成*_test.go时，不要使用terraform的测试套件，而是使用mock的方法把 云APImock掉"
- 替代方案：集成测试（缺点：依赖云环境，执行慢，不稳定）

## Risks / Trade-offs

**Risk 1: ModifyDnsRecords 接口限制**
- 风险：`ModifyDnsRecords` 一次最多修改 100 条记录，当前实现仅处理单条记录，未充分利用批量能力
- 缓解：Terraform 资源按实例操作，单条更新模式符合使用场景，批量操作可通过多个资源实例实现

**Risk 2: DNS 记录类型的参数校验**
- 风险：不同 DNS 记录类型（如 SRV、CAA）对记录值格式有特殊要求，provider 层面未做格式验证
- 缓解：依赖云 API 的参数校验，在 API 返回错误时将错误信息返回给用户

**Risk 3: 异步操作的轮询机制**
- 风险：虽然云 API 是同步接口，但 DNS 解析生效可能有延迟，Read 接口可能短暂返回不一致状态
- 缓解：使用 `helper.Retry()` 提供合理的重试机制，并在文档中说明可能的生效延迟

**Trade-off: 字段默认值的处理**
- 权衡：provider 不为可选参数设置默认值，直接依赖云 API 的默认行为
- 影响：用户无法在 Terraform 代码中看到明确的默认值，需要查看云 API 文档
- 优点：保持与云 API 的严格同步，避免不同步问题

## Migration Plan

此变更新增资源，不涉及现有资源的迁移：

1. **部署步骤**：
   - 在 `tencentcloud/services/teo/` 目录下创建 `resource_tc_teo_dns_record_v12.go` 和 `resource_tc_teo_dns_record_v12_test.go`
   - 在 `website/docs/r/` 目录下创建 `teo_dns_record_v12.md` 文档
   - 在 `examples/resources/tencentcloud_teo_dns_record_v12/` 目录下创建使用示例
   - 在 TEO 服务包中注册新资源（`service_tencentcloud_teo.go`）

2. **回滚策略**：
   - 如需回滚，直接删除新增的资源文件、文档和示例即可
   - 不影响现有资源和数据源的运行

3. **验证步骤**：
   - 运行单元测试：`go test -v ./tencentcloud/services/teo/ -run TestAccTencentCloudTeoDnsRecordV12`
   - 验证文档生成：检查 `website/docs/r/teo_dns_record_v12.md` 是否正确生成
   - 示例验证：使用示例配置创建资源并验证 CRUD 操作

## Open Questions

无。所有设计决策基于现有云 API 和 provider 实现模式，无需额外调研。
