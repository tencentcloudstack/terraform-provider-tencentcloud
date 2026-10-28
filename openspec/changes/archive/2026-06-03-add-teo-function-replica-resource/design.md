## Context

TencentCloud EdgeOne (TEO) 提供边缘函数副本功能，允许用户为边缘函数创建多个副本用于版本管理和灰度测试。云 API 已在 vendor 中提供完整的 CRUD 接口支持（CreateFunctionReplica、DescribeFunctionReplicas、ModifyFunctionReplica、DeleteFunctionReplica）。当前 Provider 中已有 TEO 服务目录（`tencentcloud/services/teo/`），需要在其中新增资源文件。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_function_replica` 资源的完整 CRUD 生命周期管理
- 遵循现有 Provider 代码风格（参考 tencentcloud_igtm_strategy）
- 支持 Import 功能（使用联合 ID）
- 提供单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不实现对应的 data source（本次仅新增 resource）
- 不修改现有 TEO 资源的行为
- 不处理异步操作（所有接口均为同步）

## Decisions

### 1. 资源 ID 设计
**决策**: 使用 `zone_id#function_id#replica_name` 作为联合 ID（分隔符为 `tccommon.FILED_SP`）

**理由**: Create 接口不返回独立的副本 ID，副本通过 zone_id + function_id + replica_name 三元组唯一标识。这与 Provider 中其他使用联合 ID 的资源保持一致。

### 2. Schema 字段设计
**决策**:
- `zone_id` (Required, ForceNew, String): 站点 ID
- `function_id` (Required, ForceNew, String): 函数 ID
- `replica_name` (Required, ForceNew, String): 副本名称
- `content` (Required, String): 副本内容（JavaScript 代码）
- `remark` (Optional, String): 副本描述

**理由**:
- `zone_id`、`function_id`、`replica_name` 设为 ForceNew，因为 Modify 接口使用 replica_name 定位副本，不支持改名，且 zone_id 和 function_id 是副本的归属标识
- `content` 和 `remark` 支持 Update（通过 ModifyFunctionReplica 接口）
- DescribeFunctionReplicas 的 SortBy、SortOrder、Filters 参数仅在 Read 内部使用，不暴露为 schema 字段

### 3. Read 实现方式
**决策**: 使用 DescribeFunctionReplicas 接口，通过 Filters 按 replica-name 过滤获取单个副本信息

**理由**: 没有单独的 DescribeFunctionReplica（单数）接口，只能通过列表接口加过滤条件获取。使用 Filters 中的 replica-name 过滤条件可以精确定位目标副本。

### 4. Delete 实现方式
**决策**: 将单个 replica_name 包装为 `[]*string` 列表传入 DeleteFunctionReplica 接口的 ReplicaNames 字段

**理由**: Delete 接口设计为批量删除（接受名称列表），但 Terraform 资源粒度为单个副本，因此每次只传入一个名称。

### 5. 测试方式
**决策**: 使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 验收测试套件

**理由**: 按照项目要求，新增资源使用 mock 方式进行单元测试，通过 `go test -gcflags=all=-l` 运行。

## Risks / Trade-offs

- [Risk] DescribeFunctionReplicas 的 Filters 为模糊查询 → 在 Read 方法中遍历返回结果，精确匹配 replica_name 确保正确性
- [Risk] Create 接口无返回 ID → 使用输入参数组合作为 ID，Import 时需要用户提供完整的联合 ID
- [Trade-off] replica_name 设为 ForceNew 意味着改名需要销毁重建 → 这是 API 限制决定的，Modify 接口不支持改名
