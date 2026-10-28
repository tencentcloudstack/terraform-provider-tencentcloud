## Context

The Terraform Provider for TencentCloud manages cloud resources via the `tencentcloud-sdk-go` SDK. The DLC (Data Lake Compute) product exposes APIs to bind (`AttachUserPolicy`) and unbind (`DetachUserPolicy`) authorization policies to/from sub-users, and to query a user's detailed info (`DescribeUserInfo`). Currently there is no Terraform resource representing this binding relationship, so users cannot manage DLC user-policy bindings as code.

The binding is described by two entities: a user (identified by `user_id`) and a set of policies (`policy_set`). The cloud APIs do not return a single stable identifier for a binding; instead the binding is identified by the combination of `user_id` and `account_type`. The cloud APIs are synchronous (no async polling is required per the API descriptions and the `AttachUserPolicy`/`DetachUserPolicy` responses contain only a `RequestId`).

This resource is `RESOURCE_KIND_ATTACHMENT`: it primarily manages bind/unbind operations and only needs Create/Read/Delete (CRD). Because the cloud API offers no dedicated update endpoint for an existing binding, the resource is treated as immutable — any schema change triggers recreation (`ForceNew`).

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_dlc_attach_user_policyr_attachment` that binds a set of DLC authorization policies to a user.
- Support Create (bind via `AttachUserPolicy`), Read (query via `DescribeUserInfo`), and Delete (unbind via `DetachUserPolicy`).
- Use a composite ID (`user_id` + `account_type` joined by `tccommon.FILED_SP`) since the cloud API returns no single identifier.
- Follow the established provider code style (referencing `tencentcloud_igtm_strategy` for non-DATASOURCE resources), including retry logic with `tccommon.ReadRetryTimeout`/`tccommon.WriteRetryTimeout` and `tccommon.RetryError`.
- Add unit tests using gomonkey mocks (no Terraform acceptance test suite) per project conventions for new resources.

**Non-Goals:**
- No update-in-place support: the resource is immutable; changing any top-level argument recreates the resource.
- No import support (ATTACHMENT resources do not include an Import section per documentation rules — only GENERAL/ATTACHMENT/CONFIG kinds include import, but this attachment has no uniquely re-queryable single id and the binding is reconstructed from request params; import is omitted to avoid ambiguity).
- No data source for listing user-policy bindings.

## Decisions

### Decision 1: Resource kind and lifecycle (CRD only, immutable)
The cloud API provides `AttachUserPolicy` (bind), `DescribeUserInfo` (read), and `DetachUserPolicy` (unbind), with no dedicated modify endpoint. Therefore the Terraform resource implements only Create/Read/Delete. Per the project rule for CRD-only resources, the `Id()` field is set to `ForceNew` and all other top-level fields are added to an `immutableArgs` array in the update method; if any of them changed, the update returns an error (which triggers Terraform to recreate).

**Why:** This matches the API surface and the project's CRD-only resource convention. **Alternative considered:** Implementing a custom update that diffs the policy_set and re-binds — rejected because the cloud API has no partial update endpoint and the attachment semantics are bind/unbind; recreation is simpler and correct.

### Decision 2: Composite ID composition
The resource ID is composed of `user_id` and `account_type` joined by `tccommon.FILED_SP` (`user_id#account_type`). In Read/Update/Delete, the ID is split back into its parts and each part is used as a request parameter.

**Why:** `AttachUserPolicy` returns only `PolicySet` + `RequestId` (no id). The pair `(user_id, account_type)` uniquely identifies the user whose policies are being managed. **Alternative considered:** Using `user_id` alone — rejected because `account_type` distinguishes `TencentAccount` vs `EntraAccount` users and is required by all three APIs.

### Decision 3: Read strategy using DescribeUserInfo
`DescribeUserInfo` returns a `UserDetailInfo` containing multiple `Policys` collections (`DataPolicyInfo`, `EnginePolicyInfo`, `CatalogPolicyInfo`, `ModelPolicyInfo`, `RowFilterInfo`). To verify the binding still exists, Read calls `DescribeUserInfo` with `Type=DataAuth` and checks whether the policies in `policy_set` are present among the user's bound data policies (matched by the `PolicyId` field returned by the API, or by the policy content fields when `PolicyId` is unavailable).

**Why:** There is no single "describe binding" API; `DescribeUserInfo` is the only read path that returns the user's bound policies. The `policy_id` request parameter of `DescribeUserInfo` (described as "TF 资源 ID") can be used to narrow the query when a specific policy id is known.

### Decision 4: PolicySet schema structure
`policy_set` is a `TypeList` of `Policy` objects. The `Policy` struct has many fields; the schema exposes the fields that are meaningful as inputs for binding: `database`, `catalog`, `table`, `operation`, `policy_type`, `function`, `view`, `column`, `data_engine`, `re_auth`, `engine_generation`, `model`, and `policy_id`. Fields marked "入参不填" in the API docs (e.g. `source`, `mode`, `operator`, `create_time`, `source_id`, `source_name`, `id`, `is_admin_policy`) are output-only and not set in Create, but may be read back where applicable.

**Why:** Only input-eligible fields should be `Required`/`Optional`; read-only fields are `Computed`.

### Decision 5: Delete via DetachUserPolicy
Delete calls `DetachUserPolicy` with the `user_id`, `account_type`, and the `policy_set` (and/or `policy_ids`) used at creation. Since `DetachUserPolicy` accepts both `PolicySet` and `PolicyIds`, the implementation passes the stored `policy_set` to unbind exactly the policies that were bound by this resource.

### Decision 6: Retry and error handling
- Create/Update/Delete wrap the cloud API call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and translate errors via `tccommon.RetryError`.
- Read wraps the `DescribeUserInfo` call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Setting the ID and other state mutations happen outside the retry block, after the retry error handling.
- In Create, after the API returns, the response is checked for emptiness; `AttachUserPolicy` returns no id, so the composite ID is set from request params. The `policy_set` returned in the response is set into state.
- In Read, if the response is empty or the bound policies are not found, `log.Printf("[CRUD] ...")` is emitted first to preserve the id, then `d.SetId("")`.

## Risks / Trade-offs

- **[No single resource id from API]** → Mitigated by deriving a composite ID from `user_id` + `account_type`. The user must not delete the same binding outside Terraform, otherwise Read will detect the binding is gone and remove it from state.
- **[DescribeUserInfo returns multiple policy categories]** → Read focuses on `DataAuth` type and matches bound policies; if a policy is bound but not reflected in `DataPolicyInfo`, Read may consider the binding absent. Mitigated by querying with the appropriate `Type` and, when available, the `PolicyId` filter to precisely locate the bound policy.
- **[Immutable resource recreation on change]** → Any change recreates the binding (unbind + bind). This is acceptable for attachment resources and matches the CRD-only convention; documented in the resource docs.
- **[PolicySet matching during Read]** → The API may return policies with `PolicyId` populated. Read matches by `PolicyId` when present, otherwise by content fields. Edge case: two identical policies could cause ambiguity, which is inherent to the API design.
