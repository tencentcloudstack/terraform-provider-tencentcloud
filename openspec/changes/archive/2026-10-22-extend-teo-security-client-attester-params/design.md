## Context

`tencentcloud_teo_security_client_attester` 是一个 RESOURCE_KIND_GENERAL 资源，封装了 TEO 客户端认证选项（Client Attester）的完整生命周期：Create 调用 `CreateSecurityClientAttester`，Read 调用 `DescribeSecurityClientAttester`（经 `DescribeTeoSecurityClientAttesterById` 分页定位），Update 调用 `ModifySecurityClientAttester`，Delete 调用 `DeleteSecurityClientAttester`。资源 ID 为 `zoneId#clientAttesterId`（`tccommon.FILED_SP` 分隔）。

当前资源的 `client_attesters` 嵌套 schema 仅暴露了 `name`、`attester_source`、`attester_duration`、`tc_rce_option`、`tc_captcha_option`、`tc_eo_captcha_option`、`id`(computed)、`type`(computed)。但云 API `ClientAttester` 结构体已支持更多字段未在 Terraform 中暴露：

- `UsageLimit *string`：是否开启认证结果使用限制（`on`/`off`）
- `MaxUsageCount *int64`：认证结果最大使用次数（`UsageLimit=on` 时必填，范围 1-100000）
- `TCEOCaptchaOption.MobileCaptchaMode *string`：移动端原生 App 人机校验交互模式（`Invisible`/`Visible`）

约束：
- 必须保持向后兼容：现有 TF 配置（未写新字段）升级后无 plan drift
- 所有新增字段均为 Optional，未配置时不传给云 API（与现有 `attester_duration` 行为一致）
- Create/Update 调用同一 `ClientAttester` 结构体，新增字段在两个接口中行为一致

### 云 API 字段定义（摘录自 vendor，实施阶段无需回读 models.go）

`ClientAttester` 结构体（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go:1884`）：

```go
type ClientAttester struct {
	// 认证选项 ID。
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// 认证选项名称。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// 认证规则类型。仅出参返回，取值有：PRESET/CUSTOM
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// 认证方法。取值有：TC-RCE / TC-CAPTCHA / TC-EO-CAPTCHA
	AttesterSource *string `json:"AttesterSource,omitnil,omitempty" name:"AttesterSource"`

	// 认证有效时间。默认 60s，单位 s(60~43200)/m(1~720)/h(1~12)
	AttesterDuration *string `json:"AttesterDuration,omitnil,omitempty" name:"AttesterDuration"`

	// 是否开启认证结果使用限制。枚举值：on/off。默认值：off
	UsageLimit *string `json:"UsageLimit,omitnil,omitempty" name:"UsageLimit"`

	// 认证结果最大使用次数。当 UsageLimit 为 on 时，此字段必填。取值范围：[1, 100000]
	MaxUsageCount *int64 `json:"MaxUsageCount,omitnil,omitempty" name:"MaxUsageCount"`

	// TC-RCE 认证的配置信息。
	TCRCEOption *TCRCEOption `json:"TCRCEOption,omitnil,omitempty" name:"TCRCEOption"`

	// TC-CAPTCHA 认证的配置信息。
	TCCaptchaOption *TCCaptchaOption `json:"TCCaptchaOption,omitnil,omitempty" name:"TCCaptchaOption"`

	// TC-EO-CAPTCHA 认证的配置信息。
	TCEOCaptchaOption *TCEOCaptchaOption `json:"TCEOCaptchaOption,omitnil,omitempty" name:"TCEOCaptchaOption"`
}
```

`TCEOCaptchaOption` 结构体（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go:26249`）：

```go
type TCEOCaptchaOption struct {
	// 浏览器端 WebView 人机校验交互模式。枚举值：Invisible/Adaptive。默认值：Invisible
	CaptchaMode *string `json:"CaptchaMode,omitnil,omitempty" name:"CaptchaMode"`

	// 移动端原生 App 人机校验交互模式。枚举值：Invisible/Visible。默认值：Invisible
	MobileCaptchaMode *string `json:"MobileCaptchaMode,omitnil,omitempty" name:"MobileCaptchaMode"`
}
```

`CreateSecurityClientAttesterRequest`（`models.go:5556`）与 `ModifySecurityClientAttesterRequest`（`models.go:22603`）字段一致：

```go
type CreateSecurityClientAttesterRequest struct {
	*tchttp.BaseRequest
	// 站点 ID。
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
	// 认证选项列表。
	ClientAttesters []*ClientAttester `json:"ClientAttesters,omitnil,omitempty" name:"ClientAttesters"`
}
```

`CreateSecurityClientAttesterResponse`（`models.go:5595`）：

```go
type CreateSecurityClientAttesterResponseParams struct {
	// 认证选项 ID 列表。
	ClientAttesterIds []*string `json:"ClientAttesterIds,omitnil,omitempty" name:"ClientAttesterIds"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_teo_security_client_attester` 资源新增 `usage_limit`、`max_usage_count`、`mobile_captcha_mode` 三个 Optional 字段
- Create/Read/Update 完整支持新字段的传递与回填
- 保持向后兼容：未配置新字段时不触发 plan diff，不向云 API 传递多余参数
- 通过单元测试（gomonkey mock 云 API）覆盖新字段

**Non-Goals:**
- 不修改 `CreateSecurityClientAttesterRequest`/`ModifySecurityClientAttesterRequest` 的请求构造方式（仍复用 `buildClientAttesterFromMap`）
- 不改变 `id`、`type` 字段的 Computed 语义
- 不新增独立的认证选项资源
- 不引入 Timeouts 块变化（资源为同步接口）

## Decisions

### Decision 1: 新增字段全部为 Optional，与现有字段保持一致

**选择**：`usage_limit`、`max_usage_count` 为 `client_attesters` 顶层 Optional 字段；`mobile_captcha_mode` 为 `tc_eo_captcha_option` 子块的 Optional 字段。

**备选**：将 `usage_limit`/`max_usage_count` 独立为顶层字段或单独子块。

**理由**：
- 云 API 中这些字段本就属于 `ClientAttesters[]` 元素属性，应平铺在 `client_attesters` 嵌套 schema 内，与 `name`/`attester_source` 同级
- `MobileCaptchaMode` 属于 `TCEOCaptchaOption`，应放进 `tc_eo_captcha_option` 子块，与 `captcha_mode` 并列
- Optional 保证向后兼容，未配置时 `buildClientAttesterFromMap` 不设置对应指针

### Decision 2: `max_usage_count` 使用 TypeInt

**选择**：`max_usage_count` schema 类型为 `schema.TypeInt`，对应云 API `MaxUsageCount *int64`。

**理由**：取值范围 1-100000 为整数，TypeInt 语义最准确，且 `helper.IntInt64()` 转换匹配 SDK 字段类型。

### Decision 3: 未配置字段不传云 API（omitnil）

**选择**：在 `buildClientAttesterFromMap` 中仅当字段非零值时设置结构体指针。

```go
if val, ok := m["usage_limit"].(string); ok && val != "" {
    attester.UsageLimit = helper.String(val)
}
if val, ok := m["max_usage_count"].(int); ok && val != 0 {
    attester.MaxUsageCount = helper.IntInt64(val)
}
// tc_eo_captcha_option 子块内：
if v, ok := opt["mobile_captcha_mode"].(string); ok && v != "" {
    eoOpt.MobileCaptchaMode = helper.String(v)
}
```

**理由**：与现有 `attester_duration`、`captcha_mode` 处理方式一致；SDK 字段为 `omitnil`，不传即不更新，避免误覆盖云端已有配置。

### Decision 4: Read 回填使用 nil 判断后再 set

**选择**：在 Read 的 `attesterMap` 构造中，对 `respData.UsageLimit`、`respData.MaxUsageCount`、`respData.TCEOCaptchaOption.MobileCaptchaMode` 做非 nil 判断后再回填。

**理由**：遵循代码生成要求，`Response` 中字段为 nil 时不调用 set，避免 panic。

## Risks / Trade-offs

- **Risk**：`MaxUsageCount` 为 `int64`，Terraform `d.Get()` 返回 `int`，在 32 位平台可能溢出 → **Mitigation**：取值范围上限 100000 远小于 int32 范围，无实际风险
- **Risk**：用户设置 `usage_limit = "off"` 但同时设置 `max_usage_count`，云 API 可能报错 → **Mitigation**：由云 API 返回错误，Provider 不做前置校验（与现有字段风格一致）
- **Trade-off**：`max_usage_count` 为 Optional 但仅在 `usage_limit=on` 时有意义，未做条件校验 → 可接受，服务端会返回明确错误