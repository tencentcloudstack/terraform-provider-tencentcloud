## Context

Terraform Provider for TencentCloud 目前已支持 TEO (EdgeOne) 产品的多种资源（如 `teo_security_ip_group`、`teo_security_policy_config` 等），但尚未支持 Security API Service 资源的管理。TEO Security API Service 允许用户创建、查询、修改和删除 API 安全服务，是 TEO 安全防护能力的重要组成部分。

当前代码库中 TEO 资源遵循以下模式：
- 资源文件位于 `tencentcloud/services/teo/` 目录下
- 使用 `teov20220901` SDK 包调用云 API
- 资源 ID 使用 `tccommon.FILED_SP` 分隔符拼接多个字段
- CRUD 操作均包含 `resource.Retry` 重试逻辑
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行错误处理

云 API 接口详情（均为同步接口）：
- `CreateSecurityAPIService`：创建 API 服务，入参 ZoneId + APIServices，出参 APIServiceIds
- `DescribeSecurityAPIService`：查询 API 服务，入参 ZoneId（含分页），出参 APIServices
- `ModifySecurityAPIResource`：修改 API 资源，入参 ZoneId + APIResources
- `DeleteSecurityAPIService`：删除 API 服务，入参 ZoneId + APIServiceIds

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_security_api_service` 资源，支持完整的 CRUD 生命周期管理
- 资源 schema 设计与云 API 参数对齐，正确映射所有字段
- 支持复合 ID（zone_id + api_service_ids），确保资源唯一标识
- 支持资源导入（Import）
- 遵循现有 TEO 资源代码风格和最佳实践
- 编写单元测试验证业务逻辑
- 生成资源文档 `.md` 文件

**Non-Goals:**
- 不修改任何现有资源的 schema 或行为
- 不新增 datasource（仅新增 resource）
- 不支持异步操作轮询（所有接口均为同步）
- 不处理跨服务的依赖关系

## Decisions

### 1. 资源 ID 设计
**决策**：使用 `zone_id` + `api_service_ids`（以 `FILED_SP` 分隔）作为复合 ID

**理由**：
- `CreateSecurityAPIService` 返回 `APIServiceIds`（字符串数组），`DeleteSecurityAPIService` 需要 `ZoneId` + `APIServiceIds` 作为入参
- `DescribeSecurityAPIService` 以 `ZoneId` 查询，返回所有 API 服务列表
- 复合 ID 需包含 `zone_id` 和所有 `api_service_ids`，以便 Read 和 Delete 操作能正确还原参数
- 由于 `api_service_ids` 是数组，将其以逗号拼接后作为 ID 的一部分

### 2. Schema 字段设计
**决策**：
- `zone_id`：Required, ForceNew - 站点 ID，创建后不可变更
- `api_services`：Required - API 服务列表（TypeList），包含 `name`（Required）和 `base_path`（Required）子字段，对应 `APIService` 结构
- `api_service_ids`：Computed - API 服务 ID 列表，由 Create 接口返回
- `api_resources`：Optional - API 资源列表（TypeList），对应 `APIResource` 结构，用于 Update 操作

**理由**：
- Create 接口的 `APIServices` 入参包含 `Name` 和 `BasePath`，`Id` 由服务端生成
- Modify 接口操作的是 `APIResources`（与 `APIServices` 是不同维度的对象），属于 Update 逻辑
- `api_service_ids` 是 Create 的返回值，应设为 Computed
- `api_resources` 仅在 Update 时使用，设为 Optional

### 3. Update 逻辑设计
**决策**：Update 操作调用 `ModifySecurityAPIResource` 修改 `api_resources`，`zone_id` 不可变更

**理由**：
- `ModifySecurityAPIResource` 接口修改的是 API 资源（APIResource），而非 API 服务（APIService）本身
- API 服务的 `name` 和 `base_path` 没有对应的修改接口，因此这些字段变更时需要重建资源
- `zone_id` 变更也需要重建资源（ForceNew）

### 4. Read 逻辑设计
**决策**：Read 操作调用 `DescribeSecurityAPIService`，使用分页查询（Limit=100，最大值），并根据 ID 中的 `api_service_ids` 过滤结果

**理由**：
- `DescribeSecurityAPIService` 按 `ZoneId` 查询，返回该站点下所有 API 服务
- 需要使用分页参数确保获取所有数据（Limit=100 为接口最大值）
- 需要根据 `api_service_ids` 过滤出属于当前资源的 API 服务

### 5. Delete 逻辑设计
**决策**：Delete 操作调用 `DeleteSecurityAPIService`，传入 `ZoneId` 和 `APIServiceIds`

**理由**：
- Delete 接口需要 `ZoneId` 和 `APIServiceIds`，与复合 ID 的组成部分一致
- 从 `d.Get()` 获取参数而非直接从 `d.Id()` 解析

### 6. 测试策略
**决策**：使用 gomonkey mock 方式编写单元测试

**理由**：
- 新增资源应使用 mock（gomonkey）方法对云 API 进行 mock 处理
- 仅进行业务代码逻辑的单元测试
- 使用 `go test -gcflags=all=-l` 运行测试

## Risks / Trade-offs

- **[Risk] DescribeSecurityAPIService 返回全量数据需过滤** → Read 时需要根据 `api_service_ids` 过滤，如果 API 服务被外部删除，可能导致资源状态不一致。通过检查返回列表中是否包含所有 `api_service_ids` 来处理此情况
- **[Risk] ModifySecurityAPIResource 操作的是 APIResource 而非 APIService** → Update 只能修改 API 资源，不能修改 API 服务的基本信息（name、base_path）。如果用户需要修改这些信息，需要重建资源
- **[Risk] Create 接口批量创建返回多个 ID** → 一个资源可能对应多个 API 服务实例，ID 拼接和解析逻辑需要仔细处理
