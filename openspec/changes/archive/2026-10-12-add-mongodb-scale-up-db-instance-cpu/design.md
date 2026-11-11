## Context

腾讯云 MongoDB 提供 CPU 弹性扩容能力，可在业务高峰期为实例临时增加 CPU 核数，低峰时回缩。相关云 API（位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725`）：

- `ScaleUpDBInstanceCpu`：手动开启 CPU 弹性扩容，入参 `InstanceId`(*string)、`ExtraCpu`(*int64)，出参 `FlowId`(*int64)。**异步接口**。
- `ScaleDownDBInstanceCpu`：关闭弹性 CPU 扩容并回缩，入参 `InstanceId`(*string)，出参 `FlowId`(*int64)。**异步接口**。
- `CheckDBInstanceElasticCpuScalable`：查询 cpu 弹性扩容信息，入参 `InstanceId`(*string)，出参 `Scalable`(*bool)、`IsScaled`(*bool)、`IsLocked`(*bool)、`Reason`(*string)、`MaxExtraCpu`(*int64)、`ExtraCpu`(*int64)、`TriggerType`(*int64)、`ScaleUpTime`(*string)。

这三个接口已在 vendor 中存在（models.go / client.go）。当前 provider 中没有资源管理 MongoDB 实例的弹性 CPU 扩容生命周期。

代码风格参照 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（通用资源 CRUD 骨架）以及 `tencentcloud/services/mongodb/resource_tc_mongodb_instance_ssl.go`（mongodb 包内同款客户端访问方式 `MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`，客户端取用 `me.client.UseMongodbClient()`）。

### 云 API request/response 关键字段定义（vendor 摘录）

**ScaleUpDBInstanceCpuRequest**（models.go:8766）
```
InstanceId *string  // 实例ID
ExtraCpu    *int64   // 要扩容的CPU核数（增量），每个节点都会增加这么多核数
```
**ScaleUpDBInstanceCpuResponse**（models.go:8805）→ `Response.FlowId *int64`（任务流ID）

**ScaleDownDBInstanceCpuRequest**（models.go:8706）
```
InstanceId *string  // 实例ID
```
**ScaleDownDBInstanceCpuResponse**（models.go:8741）→ `Response.FlowId *int64`（任务流ID）

**CheckDBInstanceElasticCpuScalableRequest**（models.go:384）
```
InstanceId *string  // 实例ID
```
**CheckDBInstanceElasticCpuScalableResponse**（models.go:440）→ `Response` 为 `CheckDBInstanceElasticCpuScalableResponseParams`：
```
Scalable     *bool    // 是否可以进行弹性CPU扩容
IsScaled     *bool    // 当前是否处于扩容状态
IsLocked     *bool    // 实例是否被锁定（有流程在执行）
Reason       *string  // 不可扩容的原因
MaxExtraCpu  *int64   // 最大可扩容的CPU核数，MIN(最小分片CPU核数, 24)
ExtraCpu     *int64   // 当前扩容的CPU核数（如果处于扩容状态）
TriggerType  *int64   // 扩容触发类型: 1-手动, 2-周期, 3-一次性时间段, 4-监控
ScaleUpTime  *string  // 扩容时间
RequestId    *string
```

客户端调用方式（client.go）：
```go
meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ScaleUpDBInstanceCpu(request)
```
辅助构造：`mongodb.NewScaleUpDBInstanceCpuRequest()` / `mongodb.NewScaleUpDBInstanceCpuResponse()`（ScaleDown、Check 同理）。

### 参照资源结构说明（tencentcloud_igtm_strategy）
- `ResourceTencentCloudIgtmStrategy()` 返回 `*schema.Resource`，含 Create/Read/Update/Delete + Importer（`schema.ImportStatePassthrough`）。
- Create：`defer tccommon.LogElapsed` + `defer tccommon.InconsistentCheck`；`logId := tccommon.GetLogId(tccommon.ContextNil)`；`ctx = tccommon.NewResourceLifeCycleHandleFuncContext(...)`；用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用云 API，失败用 `tccommon.RetryError(e)`；返回 nil 前检查 `result == nil || result.Response == nil`。
- 资源 id 用 `strings.Join([]string{...}, tccommon.FILED_SP)`（复合 id）。本资源仅用实例 id，无需拼接。
- Read：拆 id → 调 service 描述 → 空则 `log.Printf("[WARN]...")` + `d.SetId("")`；逐字段判 nil 后 `d.Set`。
- Update：检测 `d.HasChange` → 构造请求 → retry 调用 → 调 Read 回填。
- Delete：拆 id → 构造请求 → retry 调用。

### mongodb 服务层既有约定（service_tencentcloud_mongodb.go）
- `MongodbService{client *connectivity.TencentCloudClient}`，构造：`MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`。
- 异步等待已有 `DescribeAsyncRequestInfo(ctx, asyncId string, timeout time.Duration)`（行 446）：内部 `resource.Retry` 调 `DescribeAsyncRequestInfo`，状态 `MONGODB_TASK_SUCCESS`/`MONGODB_TASK_FAILED`（常量在 extension_mongodb.go）。**但本次三个接口的 FlowId 不能直接用 DescribeAsyncRequestInfo 轮询**（弹性 CPU 扩容的 FlowId 与异步任务查询体系不一定匹配），采用轮询 `CheckDBInstanceElasticCpuScalable` 的 `IsScaled` 字段判断生效更可靠（见 Decisions）。

## Goals / Non-Goals

**Goals:**
- 以声明式方式管理 MongoDB 实例的弹性 CPU 扩容：开启（Create）、查询（Read）、调整 extra_cpu（Update）、关闭并回缩（Delete）。
- 正确处理异步接口：ScaleUp/ScaleDown 返回 FlowId 后，轮询 Check 接口直至 `IsScaled` 达到期望值；且在拿到资源 id 后、异步轮询之前先 `d.SetId()`，保证失败时仍可 destroy。
- 提供 gomonkey 单元测试（不走 Terraform 验收测试套件）。

**Non-Goals:**
- 不管理周期性/监控触发型扩容策略配置（仅支持手动触发场景）。
- 不改动现有 `tencentcloud_mongodb_instance` 资源 schema。
- 不封装除上述三接口外的其它 mongodb 接口。

## Decisions

### 1. 资源 id 用实例 id（单字段）
**Decision**：`d.SetId(instanceId)`，不拼接 FlowId。
**Rationale**：一个实例同一时刻只有一个弹性 CPU 扩容状态，FlowId 是一次性任务标识、每次扩缩容都会变，不适合作为持久 id。实例 id 唯一标识该实例的弹性 CPU 管理。Import 时直接传实例 id。
**Alternatives**：`instanceId#flowId` —— FlowId 会被覆盖，Read 时无法用它反查，故否决。

### 2. 异步生效判定：轮询 CheckDBInstanceElasticCpuScalable 的 IsScaled
**Decision**：ScaleUp 之后轮询 `CheckDBInstanceElasticCpuScalable`，直到 `IsScaled == true` 视为扩容生效；ScaleDown 之后轮询直到 `IsScaled == false` 视为回缩完成。使用 `resource.Retry(3*tccommon.ReadRetryTimeout, ...)` 作为超时上限。
**Rationale**：Check 接口直接返回当前是否处于扩容状态，是判定弹性 CPU 生效与否的最直接语义来源，比依赖通用异步任务查询更贴合业务语义。
**Alternatives**：用既有 `DescribeAsyncRequestInfo(flowId)` 轮询 —— 弹性 CPU FlowId 是否在异步任务查询体系中可查未经验证，且语义不如 `IsScaled` 直接；故不采用。
**注意**：Check 接口在扩缩容进行中可能短暂返回旧状态或 `IsLocked == true`，轮询时对 API 错误用 `tccommon.RetryError(e)`，对 `IsLocked == true` 或未达期望状态返回 `resource.RetryableError` 继续重试。

### 3. extra_cpu 为可选但 Create 需校验
**Decision**：`extra_cpu` schema 标 Optional（非 ForceNew）。Create 时若用户未填则取 `d.GetOk` 缺省传 0 或不传，依赖云 API 默认行为；Update 时仅 `extra_cpu` 变化触发重新 ScaleUp。`instance_id` 为 ForceNew（变更实例等同于重建资源）。
**Rationale**：参数映射表标注两者可选；ForceNew 仅在 instance_id 上以避免资源与实例解绑。

### 4. Delete 调用 ScaleDown
**Decision**：Delete 调用 `ScaleDownDBInstanceCpu` 关闭弹性扩容并回缩，随后轮询 Check 直到 `IsScaled == false`。
**Rationale**：声明式资源销毁语义即"取消弹性 CPU 扩容"，与 ScaleDown 行为一致。若实例本就未扩容，ScaleDown 可能返回错误，此时按"资源已不存在"处理（记录日志后返回 nil，避免 destroy 卡死）。

### 5. Read 空响应处理
**Decision**：Read 调 Check 接口，若 `response == nil || response.Response == nil`，先 `log.Printf("[CRUD] ... id=%s", d.Id())` 保留现场，再 `d.SetId("")`。
**Rationale**：遵循全局规则——禁止直接清空 id 丢失排障信息。

### 6. FlowId 回填
**Decision**：Create/Update/Delete 中拿到 ScaleUp/ScaleDown 返回的 `FlowId` 后写入 state 的 `flow_id`（Computed）。FlowId 为 `*int64`，需判 nil。
**Rationale**：参数映射要求 `flow_id` 输出；便于用户追踪任务。

### 7. 测试方式：gomonkey
**Decision**：新增资源，单元测试使用 gomonkey mock 云 API（`UseMongodbClient().ScaleUpDBInstanceCpu` / `ScaleDownDBInstanceCpu` / `CheckDBInstanceElasticCpuScalable`），不使用 Terraform 验收测试套件。
**Rationale**：遵循全局规则——新增 terraform 资源用 gomonkey 做业务逻辑单测。

## Risks / Trade-offs

- [Check 接口短暂返回 IsLocked=true 或旧状态] → 轮询重试，超时上限 3*ReadRetryTimeout，耗尽后返回错误便于人工介入。
- [Delete 时实例本未扩容，ScaleDown 报错] → 捕获错误并按"已不存在"处理，记录日志后返回 nil，避免 destroy 卡死。
- [FlowId 为异步任务标识，无法持久反查] → 不用作资源 id；仅作为 Computed 输出供追踪。
- [轮询用 Check 接口而非 DescribeAsyncRequestInfo] → 语义更直接但依赖 Check 接口 `IsScaled` 字段准确性；若云 API 行为异常，超时会暴露问题，可人工介入。