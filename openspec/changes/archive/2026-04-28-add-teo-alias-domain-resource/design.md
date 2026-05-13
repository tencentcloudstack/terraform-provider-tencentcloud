## Context

TEO (EdgeOne) 是腾讯云的边缘安全加速平台，当前 Terraform Provider 已有 41 个 TEO 资源，但缺少对别称域名（Alias Domain）的管理能力。别称域名允许用户将一个域名作为另一个域名的别名，配合证书配置实现灵活的域名管理。

### 云 API 分析

经过对 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` 的详细分析，AliasDomain 相关 API 如下：

| API | 用途 | 关键参数 |
|-----|------|---------|
| `CreateAliasDomain` | 创建别称域名 | ZoneId, AliasName, TargetName, CertType, CertId([]*string) |
| `DescribeAliasDomains` | 查询别称域名列表 | ZoneId, Offset, Limit, Filters |
| `ModifyAliasDomain` | 修改别称域名配置 | ZoneId, AliasName, TargetName, CertType, CertId([]*string) |
| `ModifyAliasDomainStatus` | 修改别称域名状态 | ZoneId, Paused, AliasNames |
| `DeleteAliasDomain` | 删除别称域名 | ZoneId, AliasNames([]*string) |

**关键发现**：`DescribeAliasDomains` 返回的 `AliasDomain` 结构体仅包含 AliasName、ZoneId、TargetName、Status、ForbidMode、CreatedOn、ModifiedOn，**不包含** CertType 和 CertId 字段。这意味着 CertType 和 CertId 是只写字段（Write-Only），在 Read 操作中无法从 API 获取。

### 现有代码参考

- TEO 资源模式参考：`resource_tc_teo_origin_group.go`（完整 CRUD）
- 服务层模式参考：`service_tencentcloud_teo.go`
- 注册模式参考：`provider.go` 中 TEO 资源注册区域

## Goals / Non-Goals

**Goals:**
- 实现完整的 `tencentcloud_teo_alias_domain` 资源 CRUD 操作
- 支持别称域名的创建、读取、更新和删除
- 支持资源导入（Import）
- 正确处理 CertType/CertId 只写字段
- 在服务层实现 DescribeTeoAliasDomainById 查询方法
- 添加单元测试和文档

**Non-Goals:**
- 不实现 ModifyAliasDomainStatus 的独立调用（enable/disable 通过状态管理，当前 Terraform 资源仅管理配置，不管理启用/停用状态）
- 不实现别称域名的批量操作（遵循 Terraform 单资源原则）
- 不创建 `_extension.go` 文件（除非有异步等待需求）

## Decisions

### Decision 1: 复合 ID 使用 `zone_id#alias_name`

**选择**: `d.SetId(zoneId + tccommon.FILED_SP + aliasName)`

**理由**:
- AliasName 在同一 ZoneId 下唯一，但不同 ZoneId 可能有相同的 AliasName
- 与其他 TEO 资源（如 origin_group 使用 `zoneId#groupId`）保持一致
- CreateAliasDomain API 不返回唯一 ID，使用 zone_id + alias_name 组合作为标识

**替代方案**: 使用 `zoneId#aliasName#targetName` — 拒绝，因为 TargetName 可变，不适合作为 ID 的一部分

### Decision 2: CertType 和 CertId 作为 Optional 字段，不在 Read 中设置

**选择**: CertType 和 CertId 在 Schema 中定义为 Optional 字段，在 Read 操作中不从 API 设置（因为 API 不返回），依赖 Terraform 状态保留

**理由**:
- DescribeAliasDomains 返回的 AliasDomain 结构体不包含 CertType/CertId
- 在 Create/Update 后调用 Read 时，Terraform 状态中已保留这些值
- 这与项目中其他类似只写字段的处理方式一致

**替代方案**: 将 CertType/CertId 标记为 Computed — 拒绝，因为 Computed 意味着由服务端计算，但 API 不返回这些值

### Decision 3: CertId 使用 TypeList 而非 TypeString

**选择**: `cert_id` 在 Schema 中使用 `TypeList` of `TypeString`，对应 API 的 `[]*string` 类型

**理由**:
- CreateAliasDomain 和 ModifyAliasDomain API 中 CertId 类型均为 `[]*string`
- 遵循 API 实际类型定义，确保参数映射正确

**替代方案**: 使用 TypeString 并用逗号分隔 — 拒绝，不符合 Terraform 最佳实践

### Decision 4: 使用 DescribeAliasDomains 的 Filters 按 alias-name 精确查询

**选择**: 在服务层 `DescribeTeoAliasDomainById` 方法中使用 `AdvancedFilter`，设置 `Name="alias-name"`, `Fuzzy=false` 进行精确查询

**理由**:
- API 支持按 alias-name 过滤，支持精确匹配和模糊匹配
- 使用服务端过滤减少网络开销和数据量
- 设置 Limit 为 API 最大值 1000 以减少分页请求

### Decision 5: zone_id 和 alias_name 设为 ForceNew

**选择**: `zone_id` 和 `alias_name` 在 Schema 中标记为 `ForceNew: true`

**理由**:
- 别称域名创建后 zone_id 和 alias_name 不可更改（API 不支持修改这两个字段）
- 修改 zone_id 或 alias_name 需要删除旧资源并创建新资源
- target_name、cert_type、cert_id 为可更新字段

## Risks / Trade-offs

- **CertType/CertId 只写字段** → Read 操作后 Terraform 状态中这些字段值来自用户配置而非 API 返回值。在导入（Import）场景下，这些字段不会被填充，用户需要在配置中手动指定
- **CreateAliasDomain 不返回 ID** → 使用 zone_id + alias_name 组合作为资源 ID，创建后立即可用
- **DeleteAliasDomain 接受批量参数** → 本资源仅传递单个 alias_name，遵循 Terraform 单资源原则
- **API 标注为企业版内测功能** → 功能可能变化，但当前 SDK 已支持所有必要接口
