## Context

TEO（TencentCloud EdgeOne）是腾讯云的边缘加速产品，用户在接入域名时需要配置 CNAME 记录。为了验证配置是否正确生效，用户需要检查 CNAME 状态。当前 Terraform Provider for TencentCloud 缺少相应的操作资源支持，用户需要通过其他方式（如控制台或直接调用 API）来验证 CNAME 配置状态，这增加了运维复杂度。

Terraform Provider for TencentCloud 使用 Go 1.17+ 和 Terraform Plugin SDK v2 开发，通过腾讯云 Go SDK 调用云 API。资源文件遵循命名规范 `resource_tc_<service>_<name>.go`，操作资源使用 `_operation.go` 后缀。测试文件使用 mock（gomonkey）方式对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试。

## Goals / Non-Goals

**Goals:**
- 新增操作资源 `tencentcloud_teo_check_cname_status`，支持批量检查多个域名的 CNAME 配置状态
- 返回每个域名的记录名称、CNAME 地址和状态信息（active/moved）
- 使用 mock 测试方式验证业务逻辑的正确性
- 生成相应的文档文件供用户参考

**Non-Goals:**
- 不支持修改 CNAME 配置（只读操作）
- 不支持异步轮询检查（一次性操作资源）
- 不持久化状态（RESOURCE_KIND_OPERATION 类型）

## Decisions

**资源类型选择：RESOURCE_KIND_OPERATION**
- **决策原因**：CNAME 状态检查是一个一次性操作，不需要在 Terraform state 中持久化状态。用户只需要获取当前状态信息，不需要后续的更新或删除操作。
- **替代方案**：如果需要持续监控 CNAME 状态变化，可以考虑使用数据源（datasource），但这会增加不必要的复杂度。

**测试策略：使用 mock（gomonkey）**
- **决策原因**：作为一次性操作资源，不需要进行集成测试，使用 mock 可以快速验证业务逻辑的正确性，同时避免需要真实的云环境配置。
- **替代方案**：可以使用 Terraform 测试套件进行集成测试，但这需要真实的环境和配置，增加测试成本。

**参数命名：使用 snake_case**
- **决策原因**：遵循 Terraform 资源参数命名规范，使用 snake_case 命名（如 `zone_id`, `record_names`, `cname_status`）。
- **考虑**：云 API 使用 PascalCase（如 `ZoneId`, `RecordNames`），在代码中需要进行映射转换。

**文件命名：resource_tc_teo_check_cname_status_operation.go**
- **决策原因**：遵循现有操作资源的命名规范，使用 `_operation.go` 后缀明确标识这是一个操作资源。
- **考虑**：如果与其他资源命名冲突，可以考虑使用更具体的名称，但当前没有冲突风险。

## Risks / Trade-offs

**风险 1：CNAME 状态可能存在延迟**
- **缓解措施**：在文档中说明 CNAME 状态可能存在传播延迟，用户需要在配置 CNAME 后等待一段时间再检查。

**风险 2：批量检查可能存在性能问题**
- **缓解措施**：云 API 已经支持批量检查，单个请求可以检查多个域名，减少 API 调用次数。

**风险 3：mock 测试可能无法完全覆盖所有边界情况**
- **缓解措施**：在单元测试中覆盖正常情况、空列表、单个域名、多个域名等场景，确保代码的健壮性。

**权衡：操作资源 vs 数据源**
- 操作资源适合一次性检查，数据源适合持续查询。CNAME 状态检查通常是配置验证阶段的一次性操作，因此选择操作资源更合适。如果未来需要持续监控，可以考虑添加数据源。
