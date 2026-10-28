## Context

Terraform Provider for TencentCloud 需要新增 KMS 托管实例的购买开通能力。云 API `CreateKmsInstance`（位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118`）支持购买预付费 KMS 实例，但当前 provider 未封装该接口。

本资源为 RESOURCE_KIND_OPERATION 类型：一次性操作，操作完成后不需要记录任何状态。只需实现 Create 接口，Read/Update/Delete 为空实现。资源代码文件命名为 `resource_tc_kms_kms_instance_operation.go`。

### vendor 云 API 关键 struct 定义（已从 models.go 摘录，实施阶段无需回读大文件）

**CreateKmsInstanceRequest**（models.go:624）:
```go
type CreateKmsInstanceRequest struct {
    *tchttp.BaseRequest
    // 实例付费类型(PREPAID)
    InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty"`
    // 版本类型，枚举值：professional(专业版)
    ProductVersion *string `json:"ProductVersion,omitnil,omitempty"`
    // 预付费实例信息
    InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty"`
    // 可重入客户端ID
    ClientToken *string `json:"ClientToken,omitnil,omitempty"`
    // 实例标签
    Tags []*Tag `json:"Tags,omitnil,omitempty"`
}
```

**InstanceChargePrepaid**（models.go:3533）:
```go
type InstanceChargePrepaid struct {
    // 按月购买时长（1～100）月
    Period *uint64 `json:"Period,omitnil,omitempty"`
    // 续费标记（0:用户未设置，1:自动续费 2:不自动续费）
    RenewFlag *uint64 `json:"RenewFlag,omitnil,omitempty"`
}
```

**Tag**（models.go:4855）:
```go
type Tag struct {
    // 标签键
    TagKey *string `json:"TagKey,omitnil,omitempty"`
    // 标签值
    TagValue *string `json:"TagValue,omitnil,omitempty"`
}
```

**CreateKmsInstanceResponse**（models.go:681，Response 字段类型为 CreateKmsInstanceResponseParams，定义于 models.go:667）:
```go
type CreateKmsInstanceResponseParams struct {
    // 实例ID
    InstanceId *string `json:"InstanceId,omitnil,omitempty"`
    // 计费订单号
    DealName *string `json:"DealName,omitnil,omitempty"`
    // 计费冻结流水号
    BillId *string `json:"BillId,omitnil,omitempty"`
    // 唯一请求 ID
    RequestId *string `json:"RequestId,omitnil,omitempty"`
}
type CreateKmsInstanceResponse struct {
    *tchttp.BaseResponse
    Response *CreateKmsInstanceResponseParams `json:"Response"`
}
```

**Client 调用方式**（client.go:637）:
```go
meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().CreateKmsInstanceWithContext(ctx, request)
```
KMS client 获取方法定义于 `tencentcloud/connectivity/client.go:1097` `UseKmsClient()`。

**错误码**（client.go:605-616）: FailedOperation.OrderCheckFailed / InternalError / InvalidParameter / InvalidParameter.ConflictParameter / InvalidParameter.MissingProductVersion / InvalidParameterValue.InvalidInstanceChargeType / InvalidParameterValue.InvalidPeriod / InvalidParameterValue.InvalidProductVersion / LimitExceeded.InstanceQuotaExceeded / ResourceInUse.OrderInProgress / ResourceInUse.PostpaidInstanceAlreadyExists / UnauthorizedOperation。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_kms_kms_instance_operation` 资源，封装 `CreateKmsInstance` 云 API。
- 支持 instance_charge_type、product_version、client_token 顶层参数。
- 支持 instance_charge_prepaid 嵌套结构（period、renew_flag）。
- 支持 tags 标签列表（tag_key、tag_value）。
- 返回 instance_id、deal_name、bill_id 作为 computed 属性。
- 在 provider.go / provider.md 中完成注册与文档。
- 提供 gomonkey mock 云 API 的单元测试。

**Non-Goals:**
- 不实现 KMS 实例的读取/更新/删除（RESOURCE_KIND_OPERATION 一次性操作，无对应 RUD 云 API）。
- 不支持后付费实例创建（云 API 入参 InstanceChargeType 仅支持 PREPAID）。
- 不实现异步轮询（CreateKmsInstance 同步返回 InstanceId 等结果，非异步接口）。

## Decisions

### 1. 参数映射：嵌套结构展平 vs 嵌套 block

需求中 `instance_charge_prepaid` 同时标注为可选顶层块，且其子字段 `period`、`renew_flag` 也单独映射为顶层参数（必填）。

**决策**：采用嵌套 block 方式定义 `instance_charge_prepaid`（schema.TypeList，maxItems=1，elem 为 Resource，包含 period、renew_flag 两个必填字段）。同时保留 `period`、`renew_flag` 作为顶层必填参数的映射关系——即用户可以通过顶层 `period` / `renew_flag` 配置，也可通过 `instance_charge_prepaid` block 配置。

经权衡，为避免同一参数两处配置造成歧义，最终方案：**`instance_charge_prepaid` 作为可选嵌套 block（TypeList, maxItems=1），内部含 period（Required）、renew_flag（Required）；不再单独在顶层重复暴露 period/renew_flag**。Create 方法读取 block 内字段构造 `InstanceChargePrepaid` struct。

理由：符合 Terraform 惯例的分组参数用 block 表达；避免重复字段导致 set/diff 混乱。`tags` 作为 TypeList，每个元素含 tag_key（Required）、tag_value（Required），对应云 API `Tags []*Tag`。

### 2. 资源 ID 策略

OPERATION 类型资源无后续 RUD 操作。参考 `tencentcloud_tcss_refresh_task_operation`（使用 `helper.BuildToken()` 生成随机 token 作为 id）和 `tencentcloud_kms_overwrite_white_box_device_fingerprints`（使用业务 key_id 作为 id）。

**决策**：由于 CreateKmsInstance 返回 InstanceId，使用返回的 `InstanceId` 作为 `d.SetId()`。这样 tfstate 中保留真实实例 id，便于用户后续追踪。但需注意：此资源无 Read/Delete 实现，id 仅作占位记录，不影响一次性操作语义。

若返回的 InstanceId 为空（含 Response 为 nil、InstanceId 为 nil 或空字符串），返回 `resource.NonRetryableError`，避免写入空 id 触发后续状态混乱（符合代码生成要求第 9 条）。

### 3. Retry 与错误处理

Create 方法中使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装云 API 调用，失败时用 `tccommon.RetryError(e)` 包装返回。retry 块内只调用接口，不做 set id 等成功操作（set id 放到 retry 块外）。符合代码生成要求第 4 条。

### 4. 单元测试策略

新增资源属新增 terraform 资源，单元测试使用 gomonkey mock 云 API（`UseKmsClient` 与 `CreateKmsInstanceWithContext`），仅做业务代码逻辑单元测试，不使用 terraform 测试套件，不实际调用云 API。符合 go 代码生成要求第 1 条。

### 5. 文档与注册

- `provider.go`：在 KMS 资源区块（约第 2090 行 `tencentcloud_kms_overwrite_white_box_device_fingerprints` 之后）追加 `"tencentcloud_kms_instance_operation": kms.ResourceTencentCloudKmsInstanceOperation(),`。
- `provider.md`：在 KMS Resource 区块（约第 711 行之后）追加 `tencentcloud_kms_instance_operation`。
- 资源 .md 文件 `resource_tc_kms_kms_instance_operation.md`：一句话描述带云产品名 KMS，Example Usage，无 Import 部分（OPERATION 类型无 Import）。不手写 Argument/Attribute Reference（由工具自动生成）。

### 6. 参考代码风格

- OPERATION 骨架参考 `tencentcloud/services/tcss/resource_tc_tcss_refresh_task_operation.go`：Schema 定义、Create/Read/Delete 空实现、`helper.BuildToken()` 占位 id、defer LogElapsed/InconsistentCheck。
- 嵌套 list/block 处理参考 `tencentcloud/services/kms/resource_tc_kms_overwrite_white_box_device_fingerprints.go`（device_fingerprints 列表遍历构造 SDK struct 的模式）。
- retry + 错误处理参考 `tencentcloud/services/crs/resource_tc_redis_renew_instance_operation.go`。
- kms 包名、import、client 获取参考现有 kms 资源（`import kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"`，`UseKmsClient()`）。

## Risks / Trade-offs

- **[风险] 重复购买**：CreateKmsInstance 为下单操作，若用户重复 apply 可能重复下单。→ 缓解：资源为 OPERATION 类型，apply 成功后 id 已设置，再次 apply 不会触发 create（无 update 实现）。用户需自行管理资源生命周期。提供 client_token 支持幂等重入。
- **[风险] 无回滚能力**：OPERATION 资源 Delete 为空，无法通过 terraform destroy 销毁已购买的 KMS 实例。→ 缓解：在 .md 文档中明确说明此资源为一次性购买操作，删除资源不会销毁 KMS 实例，需用户手动在控制台处理。
- **[权衡] 使用 InstanceId 作为 id vs BuildToken**：使用真实 InstanceId 便于追踪，但该资源无 Read，state 刷新不会校验实例是否存在。权衡后选择 InstanceId 以保留可追溯性。
- **[风险] InstanceChargePrepaid 字段类型为 uint64**：period、renew_flag 在云 API 中为 `*uint64`。Schema 中使用 `schema.TypeInt`，读取时用 `helper.IntUint64(v.(int))` 转换。