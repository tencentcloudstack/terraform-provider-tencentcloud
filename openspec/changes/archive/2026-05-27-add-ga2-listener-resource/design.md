## Context

TencentCloud GA2 (Global Accelerator 2.0) 已有 `tencentcloud_ga2_endpoint_group` 资源实现，位于 `tencentcloud/services/ga2/`。现在需要新增 `tencentcloud_ga2_listener` 资源来管理监听器的生命周期。

GA2 SDK 包路径为 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`，已在 vendor 中。服务层 `Ga2Service` 和异步任务轮询方法 `WaitForGa2TaskFinish` 已存在于 `service_tencentcloud_ga2.go` 中。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_ga2_listener` 资源的完整 CRUD 操作
- 支持 TCP/UDP/HTTP/HTTPS 协议监听器
- 支持异步操作（Create/Modify/Delete 均返回 TaskId，需轮询 DescribeTaskResult）
- 支持 Import 功能（使用 `global_accelerator_id#listener_id` 联合 ID）
- 提供单元测试（使用 gomonkey mock 方式）

**Non-Goals:**
- 不实现监听器的转发规则管理（属于其他资源）
- 不实现监听器的数据源查询资源

## Decisions

1. **资源 ID 格式**: 使用 `global_accelerator_id#listener_id` 联合 ID（tccommon.FILED_SP 分隔），与现有 ga2_endpoint_group 资源风格一致。
   - 理由: DescribeListeners/ModifyListener/DeleteListener 都需要 GlobalAcceleratorId 和 ListenerId 两个参数。

2. **ForceNew 字段**: `global_accelerator_id`、`port_ranges`、`listener_type`、`protocol` 设为 ForceNew。
   - 理由: 这些字段在 ModifyListener 接口中不存在（除了 global_accelerator_id 是标识字段），创建后不可修改。

3. **异步操作处理**: 复用现有 `Ga2Service.WaitForGa2TaskFinish` 方法。
   - 理由: 该方法已在 endpoint_group 资源中验证可用，避免重复代码。

4. **Read 实现**: 在服务层新增 `DescribeGa2ListenerById` 方法，使用 Filters 按 listener-id 过滤。
   - 理由: 与 DescribeGa2EndpointGroupById 模式一致，支持分页查询。

5. **client_affinity_time 字段处理**: 该字段仅在 ModifyListener 中可传入，CreateListener 不支持。Schema 中设为 Optional+Computed，Create 时不传，Update 时可传。
   - 理由: 云 API 接口参数不一致，需要适配。

6. **port_ranges 字段**: 使用 TypeList + MaxItems:1 嵌套结构体，包含 from_port 和 to_port。
   - 理由: SDK 中 PortRanges 是单个结构体指针（非数组），包含 FromPort 和 ToPort 两个 uint64 字段。

## Risks / Trade-offs

- [Risk] CreateListener 接口返回 TaskId 但任务可能失败 → 通过 WaitForGa2TaskFinish 轮询直到 SUCCESS，超时则返回错误
- [Risk] DescribeListeners 的 Filters 可能不严格过滤 → 在代码中对返回结果做精确匹配
- [Trade-off] port_ranges 使用 MaxItems:1 的 TypeList 而非直接展开字段 → 保持与 SDK 结构体的对应关系，便于维护
