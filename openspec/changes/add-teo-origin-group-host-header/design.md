## Context

当前 tencentcloud_teo_origin_group 资源已实现基本的源站组管理功能，支持配置源站服务器列表、权重等参数。腾讯云 TEO 服务的 CreateOriginGroup API 提供了 HostHeader 参数，允许用户在请求源站时携带自定义 Host 头信息，这在以下场景中非常重要：
- 多个域名共享同一个源站 IP 时，需要通过 Host 头区分不同站点
- 源站服务需要特定的 Host 头才能正确路由请求
- 用户需要在 CDN/边缘节点和源站之间传递自定义 Host 信息

资源位于 `tencentcloud/services/teo/resource_tencentcloud_teo_origin_group.go`，使用 Terraform Plugin SDK v2 构建，通过 tencentcloud-sdk-go 调用底层 API。

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_teo_origin_group 资源新增 `host_header` 参数（TypeString，Optional）
- 在资源创建和更新时正确处理 host_header 参数
- 确保参数正确映射到 CreateOriginGroup 和 ModifyOriginGroup API 的 HostHeader 字段
- 更新资源文档和测试用例
- 保持向后兼容性，不破坏现有配置和状态

**Non-Goals:**
- 不修改现有参数的行为或语义
- 不影响其他 TEO 资源
- 不改变资源的 ID 或状态结构

## Decisions

**1. 参数命名和类型**
- 使用 `host_header` 作为参数名（遵循 Terraform 命名规范，下划线分隔）
- 参数类型为 `schema.TypeString`，使用 `schema.Optional` 和 `schema.Computed`
  - Optional: 用户可以省略该参数
  - Computed: 允许 API 返回默认值（如果有），但实际用户期望显式设置，可考虑仅 Optional

**2. Schema 位置**
- 将 host_header 参数添加到资源 Schema 的根级别
- 参数描述："Custom Host header to use when making requests to the origin"

**3. API 参数映射**
- CreateOriginGroup: 将 `d.Get("host_header")` 映射到 API 的 `HostHeader` 字段
- ModifyOriginGroup: 同样映射 HostHeader 参数
- ReadOriginGroup: 从 API 响应读取 HostHeader 并设置到 state
- DeleteOriginGroup: 无需处理（删除操作不涉及该参数）

**4. 状态管理**
- 在 Create 函数中：使用 `d.Set("host_header", params.HostHeader)` 设置初始状态
- 在 Update 函数中：对比新旧值，仅在变化时调用 ModifyOriginGroup
- 在 Read 函数中：从 API 响应读取并更新 state：`d.Set("host_header", resp.HostHeader)`

**5. 向后兼容性**
- 新增参数为 Optional，现有配置无需修改
- 不修改现有字段或删除任何功能
- 确保 state 升级平滑（新用户使用新参数，旧用户配置不受影响）

**6. 测试策略**
- 新增测试用例验证 host_header 参数：
  - 创建包含 host_header 的资源
  - 更新 host_header 值
  - 删除资源
- 确保所有现有测试仍通过

## Risks / Trade-offs

**[Risk 1] API 字段变更**
- CreateOriginGroup API 的 HostHeader 字段可能在 future 版本中废弃或行为变化
- **Mitigation**: 遵循 Terraform Provider 最佳实践，添加参数校验，监控 API 变更公告

**[Risk 2] 类型不匹配**
- API 的 HostHeader 可能接受多种类型（如字符串数组、结构体等）
- **Mitigation**: 在实现前查阅 API 文档确认字段类型，当前假设为字符串类型

**[Risk 3] 空值处理**
- 用户可能传递空字符串，API 可能不接受空值
- **Mitigation**: 在 Create/Update 函数中添加校验，空字符串不传递给 API 或传递默认值

**[Trade-off] Computed vs Optional**
- 如果仅使用 Optional，用户必须显式设置值；如果使用 Computed，可以接受 API 返回的默认值
- **Decision**: 当前需求未明确说明 API 是否有默认值，暂定使用 Optional + Computed，根据 API 文档确认后调整

**[Trade-off] 参数分组**
- host_header 可以与 origin 相关参数分组（如 origin_servers），但为简化 Schema 结构，放在根级别
- **Decision**: 根级别放置，保持与其他 TEO 资源一致的风格

## Migration Plan

**部署步骤:**
1. 更新资源 Schema 和 CRUD 函数代码
2. 更新资源文档和示例
3. 新增或更新测试用例
4. 运行完整测试套件确保现有功能不受影响
5. 发布新版本 Provider

**回滚策略:**
- 由于是纯新增 Optional 参数，回滚只需删除相关代码
- 不会影响已部署的用户配置（参数未使用时无副作用）
- 如果发现 API 行为异常，可快速发布补丁移除参数支持

## Open Questions

1. **API 字段类型确认**: HostHeader 字段的确切类型（string、[]string、或其他）需要查阅腾讯云 API 文档确认
2. **空值处理规则**: API 是否接受空字符串？如果不接受，应如何处理？
3. **默认值行为**: CreateOriginGroup API 是否返回默认的 HostHeader？如果是，Computed 是否必要？

（问题将在实现阶段通过查阅 API 文档或测试解决）
