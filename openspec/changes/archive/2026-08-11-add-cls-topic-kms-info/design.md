## Context

The `tencentcloud_cls_topic` resource manages CLS (Cloud Log Service) log topics. It already supports an `encryption` parameter (TypeInt, Optional, Computed) that enables KMS-CLS cloud product key encryption when set to `1`. However, when encryption is enabled, users currently always use the **system default** KMS key (alias `KMS-CLS`) and have no way to specify their own custom KMS key.

The CLS cloud API supports custom KMS keys via a nested `CustomKmsInfo` struct on both the `CreateTopicRequest` and `ModifyTopicRequest`. The `DescribeTopics` response also returns `CustomKmsInfo` inside each `TopicInfo`. The vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`) already defines the `CustomKmsInfo` struct with two string fields:

```go
type CustomKmsInfo struct {
    KmsRegion *string `json:"KmsRegion,omitnil,omitempty"`
    KmsKeyId  *string `json:"KmsKeyId,omitnil,omitempty"`
}
```

No SDK upgrade is required — the structs and fields are already present in the vendored code.

Current resource file: `tencentcloud/services/cls/resource_tc_cls_topic.go`

**API behavior analysis (from vendor):**

| API | `CustomKmsInfo.KmsRegion` in Request | `CustomKmsInfo.KmsKeyId` in Request | `CustomKmsInfo.KmsRegion` in Response | `CustomKmsInfo.KmsKeyId` in Response |
|-----|-------------------------------------|-------------------------------------|---------------------------------------|-------------------------------------|
| `CreateTopic` | Yes (nested in `CustomKmsInfo`) | Yes (nested in `CustomKmsInfo`) | No (create returns only `TopicId`) | No |
| `DescribeTopics` | N/A | N/A | Yes (in `TopicInfo.CustomKmsInfo`) | Yes (in `TopicInfo.CustomKmsInfo`) |
| `ModifyTopic` | Yes (nested in `CustomKmsInfo`) | Yes (nested in `CustomKmsInfo`) | No | No |
| `DeleteTopic` | No | No | N/A | N/A |

## Goals / Non-Goals

**Goals:**
- Add `kms_region` (TypeString, Optional, Computed) and `kms_key_id` (TypeString, Optional, Computed) parameters to the `tencentcloud_cls_topic` resource schema
- Construct and populate a `cls.CustomKmsInfo{}` struct in the Create method when `encryption = 1` and either field is set, then assign it to `request.CustomKmsInfo`
- Construct and populate a `cls.CustomKmsInfo{}` struct in the Update method when encryption-related changes occur, then assign it to `request.CustomKmsInfo` so the `ModifyTopic` call carries the custom key info
- Read `KmsRegion` and `KmsKeyId` from `topic.CustomKmsInfo` (with nil checks) in the Read method and set them in state
- Maintain full backward compatibility — existing configurations without these fields continue to work (CLS uses the default key)

**Non-Goals:**
- This change does NOT make `encryption` itself mutable after creation beyond what it already is
- This change does NOT add `kms_region` / `kms_key_id` to any CLS data source (e.g., `tencentcloud_cls_topics`)
- This change does NOT modify any other CLS resources
- This change does NOT add new resources

## Decisions

### Decision 1: Flatten `CustomKmsInfo` into two top-level schema fields

**Rationale:** The cloud API nests `KmsRegion` and `KmsKeyId` inside a `CustomKmsInfo` struct, but the Terraform schema convention (per the codebase rules) is to flatten nested single-object structs into top-level fields. This matches how the existing `encryption` field is already a top-level scalar. Flattening gives users a cleaner HCL experience (`kms_region = "..."` and `kms_key_id = "..."`) rather than requiring a nested block.

### Decision 2: Construct `CustomKmsInfo` struct inline in Create/Update rather than in the service layer

**Rationale:** The existing Create and Update methods in `resource_tc_cls_topic.go` already build the `cls.NewCreateTopicRequest()` / `cls.NewModifyTopicRequest()` inline and assign fields directly. Following this existing pattern, the `CustomKmsInfo` struct is constructed inline and assigned to `request.CustomKmsInfo`. This avoids adding a new service-layer function for a simple nested-struct assignment, keeping the change minimal.

### Decision 3: Gate `CustomKmsInfo` construction on `encryption = 1`

**Rationale:** The cloud API documentation states `CustomKmsInfo` only takes effect when `Encryption = 1`. To avoid sending meaningless payloads in other scenarios (and to respect the API's "other scenarios cannot pass this parameter" note), the provider only constructs and assigns `CustomKmsInfo` when `encryption` is set to `1` and at least one of `kms_region` / `kms_key_id` is provided by the user.

### Decision 4: Both parameters are Optional + Computed (not ForceNew)

**Rationale:** The `ModifyTopic` API accepts `CustomKmsInfo`, so the parameters are updatable, not immutable. `Computed: true` ensures that resources created without these fields (or imported resources) still populate state from the `DescribeTopics` response. This is consistent with the existing `encryption` field which is also Optional + Computed.

### Decision 5: Nil-check `CustomKmsInfo` and sub-fields before reading in Read

**Rationale:** The `DescribeTopics` response may return `TopicInfo.CustomKmsInfo` as `nil` when encryption is not enabled or the default key is used. Per the codebase rules, the Read method must check `topic.CustomKmsInfo != nil` before accessing `.KmsRegion` / `.KmsKeyId`, and only call `d.Set(...)` when the sub-field is non-nil.

## Risks / Trade-offs

- **[Risk] `CustomKmsInfo` is nil in Describe response for default-key topics** → Mitigation: nil-check `topic.CustomKmsInfo` before reading sub-fields; do not call `d.Set` when nil, so state keeps the user's configured value or Computed zero-value.
- **[Risk] Users set `kms_region` / `kms_key_id` without `encryption = 1`** → Mitigation: the provider only constructs `CustomKmsInfo` when `encryption = 1`; otherwise the values are ignored in the API request. No hard error is returned (the API itself would ignore them), keeping behavior lenient and backward compatible.
- **[Backward compatibility]** → Both new parameters are Optional + Computed. Existing configurations without `kms_region` / `kms_key_id` continue to work unchanged.
