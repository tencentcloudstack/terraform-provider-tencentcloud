## Context

tencentcloud_vpc_end_point 资源是 Terraform Provider for TencentCloud 中的一个 VPC 服务资源，用于管理私有网络终端节点的创建和配置。当前资源实现基于 v1.82.73 版本，通过调用腾讯云的 CreateVpcEndPoint API 创建资源。

根据用户反馈和实际使用场景，当前资源缺少三个重要参数：SecurityGroupId（安全组配置）、Tags（资源标签管理）和 IpAddressType（IP 地址类型配置）。这些参数在腾讯云 API 中已经支持，但在 Terraform Provider 中尚未实现。

由于这些参数都是可选的（Optional），可以安全地添加到现有 Schema 中而不破坏向后兼容性。

## Goals / Non-Goals

**Goals:**
1. 在 tencentcloud_vpc_end_point 资源 Schema 中新增 SecurityGroupId、Tags 和 IpAddressType 三个可选字段
2. 更新 Create 函数以支持这三个新参数的传递到 CreateVpcEndPoint API
3. 更新 Read 函数以从 DescribeVpcEndPoints API 读取并同步这些参数到 Terraform state
4. 更新 Update 函数以支持这些参数的修改（通过 ModifyVpcEndPointAttribute API 或资源重建）
5. 确保新字段的验证逻辑与 API 定义一致
6. 更新单元测试和验收测试以覆盖新字段的功能

**Non-Goals:**
1. 不修改现有字段的任何行为（保持向后兼容）
2. 不涉及其他 VPC Endpoint 相关资源的修改
3. 不实现自定义的标签管理逻辑（使用 Terraform 标签机制）
4. 不实现自定义的 IP 地址类型转换逻辑

## Decisions

### 1. Schema 定义方式
**Decision:** 使用 Terraform Plugin SDK v2 的标准 Schema 定义方式，新字段声明为 Optional。

**Rationale:**
- SecurityGroupId: 使用 `schema.TypeString` 类型，标记为 Optional
- Tags: 使用 `schema.TypeList` + `schema.TypeMap` 类型嵌套，Key 为必填，Value 为可选
- IpAddressType: 使用 `schema.TypeString` 类型，标记为 Optional，默认值为 "Ipv4"

**Alternatives considered:**
- 将字段标记为 Computed: 考虑过，但这些参数是用户可配置的，应该使用 Optional
- 使用 Set 类型替代 List: 考虑过，但 Tags 是有序的键值对列表，List 更合适

### 2. Create 操作实现
**Decision:** 在 resourceTencentCloudVpcEndPointCreate 函数中，从 Terraform configuration 读取新字段并映射到 API 请求参数。

**Rationale:**
- CreateVpcEndPoint API 已经支持这三个参数
- 通过 `d.Get()` 方法获取值，然后设置到 API 请求结构体中
- 使用 SDK 的标准错误处理和日志记录机制

### 3. Read 操作实现
**Decision:** 在 resourceTencentCloudVpcEndPointRead 函数中，从 DescribeVpcEndPoints API 响应中读取这些字段并更新到 Terraform state。

**Rationale:**
- 保持 state 与云资源状态一致
- 使用 `d.Set()` 方法更新 state
- 处理 API 返回的 null/empty 情况

### 4. Update 操作实现
**Decision:** 优先使用 ModifyVpcEndPointAttribute API 更新可变参数，如果 API 不支持则标记为 ForceNew。

**Rationale:**
- SecurityGroupId: 通常可以通过 Modify API 更新
- Tags: 可以通过标签管理 API 或 Modify API 更新
- IpAddressType: 如果不支持修改，标记为 ForceNew 触发资源重建

**Alternatives considered:**
- 所有参数都标记为 ForceNew: 会增加资源重建成本
- 所有参数都尝试通过 Modify API 更新: 需要先测试 API 支持情况

### 5. 验证逻辑
**Decision:** 使用 Terraform 的 Schema validators 或在 Create/Update 函数中进行验证。

**Rationale:**
- IpAddressType: 验证值为 "Ipv4" 或 "Ipv6"
- Tags: 验证 Key 字段不为空
- SecurityGroupId: 验证格式（如果需要）

### 6. 测试策略
**Decision:** 为每个新字段添加单元测试和验收测试场景。

**Rationale:**
- 单元测试覆盖 Create/Read/Update 逻辑
- 验收测试验证与云 API 的集成
- 测试边界情况和错误处理

## Risks / Trade-offs

**Risk 1:** ModifyVpcEndPointAttribute API 可能不支持某些字段的更新
→ Mitigation: 提前测试 API 支持情况，不支持的参数标记为 ForceNew

**Risk 2:** Tags 字段的数据结构可能与其他资源不一致
→ Mitigation: 遵循 Terraform Provider 的标准 Tags 实现模式，参考其他资源的实现

**Risk 3:** IpAddressType 的默认值可能导致意外行为
→ Mitigation: 在文档中明确说明默认值，并在测试中验证默认值行为

**Risk 4:** 新增字段可能影响现有资源的 state 迁移
→ Mitigation: 使用 Optional 字段，现有 state 会自动兼容（新字段为空）

**Trade-off:** 对于不支持更新的字段，使用 ForceNew 会导致资源重建，但可以保证配置的一致性
→ Acceptable: 这种做法在 Terraform 中是标准的，用户可以通过合理规划变更来最小化影响

## Migration Plan

### 部署步骤
1. 修改 resource_tencentcloud_vpc_end_point.go 中的 Schema 定义
2. 更新 Create 函数以传递新字段到 API
3. 更新 Read 函数以从 API 读取新字段
4. 更新 Update 函数以支持新字段的修改
5. 编写和运行单元测试
6. 编写和运行验收测试（TF_ACC=1）
7. 更新资源文档和示例

### 回滚策略
如果新实现引入问题，可以通过以下方式回滚：
1. Git revert 代码变更
2. 发布补丁版本
3. 对于已经使用新字段的用户，需要调整配置文件
4. 提供迁移指南帮助用户处理 state 更新

### 兼容性
- **向后兼容**: 是（新字段都是 Optional）
- **State 迁移**: 不需要（自动兼容）
- **配置迁移**: 不需要（新字段可选，不影响现有配置）
