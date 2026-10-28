## Context

Terraform Provider for TencentCloud 需要新增 KMS 托管实例的购买开通能力。云 API `CreateKmsInstance`（位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118`）支持购买预付费 KMS 实例，但当前 provider 未封装该接口。

本实现为 terraform-plugin-framework **action** 类型：一次性操作，操作完成后不需要记录任何状态。仅实现 Invoke 接口，无 Read/Update/Delete。资源代码文件命名为 `action_tc_kms_instance_operation.go`，符合 framework action 命名规范（`action_tc_<product>_<name>.go`）。

### 平台参考

- terraform-plugin-framework 的 `action.Action` 接口，参考现有 action 实现：`tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.go`、`tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`、`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`、`tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.go`。
- action 工厂在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片中注册。
- action 文档文件命名 `action_tc_kms_instance_operation.md`，由 `make doc` 生成到 `website/docs/actions/`。

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
KMS client 获取方法定义于 `tencentcloud/connectivity/client.go` `UseKmsClient()`。

## Goals / Non-Goals

**Goals:**
- 新增 framework action `tencentcloud_kms_instance_operation`，封装 `CreateKmsInstance` 云 API。
- 支持 instance_charge_type、product_version、client_token 顶层属性。
- 支持 instance_charge_prepaid SingleNestedBlock（period、renew_flag）。
- 支持 tags ListNestedBlock（tag_key、tag_value）。
- 在 framework/registry.go 中注册 action 工厂。
- 在 provider.md 中更新文档索引。
- 提供 gomonkey mock 云 API 的单元测试。

**Non-Goals:**
- 不实现 KMS 实例的读取/更新/删除（action 一次性操作，无对应 RUD 云 API）。
- 不支持后付费实例创建（云 API 入参 InstanceChargeType 仅支持 PREPAID）。
- 不对外暴露输出属性（framework action 限制，instance_id 等仅通过日志输出）。

## Decisions

### 1. framework action vs SDKv2 resource

采用 terraform-plugin-framework **action** 而非 SDKv2 resource。理由：
- action 是 Terraform 1.14+ 提供的语义化一次性操作原语，正好匹配"购买开通实例"的场景——不需要管理状态、不需要 destroy。
- provider 已有空白的 framework action 基础设施（framework/registry.go 的 actionFactories、fw.ActionWithConfigure），现有 6 个 action 已验证该模式可行。
- action 在 `framework/registry.go` 中注册，不污染 SDKv2 的 `provider.go` ResourcesMap。

### 2. Schema 设计

- `instance_charge_type`：StringAttribute，Required。云 API 必填。
- `product_version`：StringAttribute，Optional。
- `client_token`：StringAttribute，Optional。
- `instance_charge_prepaid`：SingleNestedBlock，含 period（Int64Attribute Required）、renew_flag（Int64Attribute Required）。对应云 API `InstanceChargePrepaid *InstanceChargePrepaid`。
- `tags`：ListNestedBlock，含 tag_key（StringAttribute Required）、tag_value（StringAttribute Required）。对应云 API `Tags []*Tag`。

framework action schema 使用 `action/schema` 包（非 resource/schema），SingleNestedBlock 在 model 中映射为 `types.Object`，ListNestedBlock 映射为 `types.List`。

### 3. 输出属性处理

framework action 无法对外暴露输出属性（与 SDKv2 resource 的 Computed 字段不同）。CreateKmsInstance 返回的 InstanceId、DealName、BillId 通过 `log.Printf("[DEBUG]...")` 记录到日志，便于用户排障。

### 4. Retry 与错误处理

在 service 方法 `KmsService.CreateKmsInstance` 中使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装云 API 调用，失败时用 `tccommon.RetryError(e)` 包装返回。retry 块内只调用接口。空响应在 service 方法中返回 error。

### 5. 文件命名

遵循 framework README 的命名规范：`action_tc_<product>_<name>.go`。由于该 action 明确属于 KMS 产品，文件命名为 `action_tc_kms_instance_operation.go`，Go package 为 `kms`。

### 6. 参考代码风格

- action 骨架参考 `tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources.go` 和 `tencentcloud/services/dlc/action_tc_dlc_initialize_tc_lake.go`。
- 嵌套 block 处理参考 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`（SingleNestedBlock 用 types.Object + As，ListNestedBlock 用 types.List + ElementsAs）。
- service 方法 retry 参考现有 kms service 和 dlc service。
- 测试参考 `tencentcloud/services/antiddos/action_tc_antiddos_unblock_resources_test.go`（tftypes 原始值构建 InvokeRequest）。

## Risks / Trade-offs

- **[约束] 需 Terraform 1.14+**：action 是较新的 Terraform 特性，低版本用户无法使用。在 .md 文档中已用 NOTE 提示。
- **[风险] 重复购买**：CreateKmsInstance 为下单操作，若用户重复执行 action 可能重复下单。action 无状态，不会阻止重复调用。提供 client_token 支持幂等重入。
- **[权衡] 无输出属性**：instance_id 等仅日志输出，用户无法在 HCL 中引用。符合 framework action 设计限制。
- **[风险] InstanceChargePrepaid 字段类型为 uint64**：period、renew_flag 在云 API 中为 `*uint64`。Schema 中使用 `schema.Int64Attribute`，读取时用 `helper.IntUint64(v)` 转换。