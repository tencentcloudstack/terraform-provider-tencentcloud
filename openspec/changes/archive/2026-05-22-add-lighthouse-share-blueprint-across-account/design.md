## Context

Terraform Provider for TencentCloud 当前不支持 Lighthouse 产品的跨账号共享镜像操作。用户需要通过 Terraform 自动化执行 `ShareBlueprintAcrossAccounts` API，将自定义镜像共享给其他账号。

该 API 是同步接口，响应仅包含 RequestId，无需异步轮询。API 参数包括：
- `BlueprintId`（string）：镜像ID
- `AccountIds`（[]*string）：接收共享镜像的账号ID列表，最多10个

当前 Lighthouse 服务在 provider 中已有相关资源注册（如 `tencentcloud_lighthouse_blueprint`），需要遵循相同的代码组织模式。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_lighthouse_share_blueprint_across_account` OPERATION 类型资源
- 实现 Create 接口，调用 `ShareBlueprintAcrossAccounts` API
- Read/Update/Delete 接口为空（一次性操作，无需记录状态）
- 在 provider.go 中注册新资源
- 编写单元测试（使用 gomonkey mock）
- 生成资源样例文档

**Non-Goals:**
- 不实现 CancelShareBlueprintAcrossAccounts 取消共享操作（可作为后续需求）
- 不实现 Read 接口查询共享状态
- 不支持异步轮询（该接口为同步接口）

## Decisions

1. **资源ID生成策略**：使用 `helper.BuildToken()` 生成随机 token 作为资源 ID。因为一次性操作没有持久化状态，也没有可唯一标识的返回值，与 `advisor_authorization_operation` 等现有 operation 资源保持一致。

2. **Schema 字段设计**：
   - `blueprint_id`：Required, ForceNew, TypeString - 镜像ID
   - `account_ids`：Required, ForceNew, TypeList(of TypeString) - 接收账号ID列表
   - 所有字段设置 ForceNew 因为 OPERATION 类型资源无 Update 操作

3. **CRUD 实现**：
   - Create: 调用 `ShareBlueprintAcrossAccounts` API，使用 `tccommon.WriteRetryTimeout` 重试
   - Read: 空实现，直接返回 nil
   - Update: 不实现（OPERATION 类型）
   - Delete: 空实现，直接返回 nil

4. **重试策略**：在 Create 中使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，失败时使用 `tccommon.RetryError(e)` 包装错误。

5. **代码文件组织**：遵循现有 Lighthouse 服务目录结构，资源文件命名为 `resource_tc_lighthouse_share_blueprint_across_account_operation.go`。

## Risks / Trade-offs

- [API 限制] AccountIds 最多支持10个账号 → 在 schema 描述中说明限制，不做强制校验，由云 API 返回错误
- [无状态追踪] 一次性操作不记录状态，无法通过 Terraform 查询共享结果 → 用户可通过控制台或 API 查询确认
- [幂等性] 重复调用相同参数的 API 不会报错（云 API 行为），但 Terraform 会重新创建资源 → 可接受，与现有 operation 资源行为一致
