## Context

The `tencentcloud_scf_function` resource currently manages SCF (Serverless Cloud Function) lifecycle operations but lacks the `qualifier` parameter. The `qualifier` is a SCF concept representing a function version number or alias name. Several SCF APIs accept a `Qualifier` parameter to target specific versions/aliases:

- `GetFunction`: accepts `Qualifier` as input and returns it in the response
- `DeleteFunction`: accepts `Qualifier` to delete a specific version
- `CreateTrigger`: accepts `Qualifier` to specify which version/alias the trigger binds to
- `DeleteTrigger`: accepts `Qualifier` to identify which trigger to delete

The SDK (`v20180416`) already defines the `Qualifier` field in these request/response structs. The task is to expose this as a Terraform schema parameter.

### SDK Verification (vendor directory)

Based on `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416/models.go`:

| API | Request Qualifier | Response Qualifier |
|-----|-------------------|---------------------|
| `CreateFunction` | ❌ Not supported | N/A |
| `UpdateFunctionCode` | ❌ Not supported | N/A |
| `UpdateFunctionConfiguration` | ❌ Not supported | N/A |
| `GetFunction` | ✅ `GetFunctionRequest.Qualifier` (line 2347) | ✅ `GetFunctionResponseParams.Qualifier` (line 2496) |
| `DeleteFunction` | ✅ `DeleteFunctionRequest.Qualifier` (line 1113) | N/A |
| `CreateTrigger` | ✅ `CreateTriggerRequest.Qualifier` (line 894) | N/A |
| `DeleteTrigger` | ✅ `DeleteTriggerRequest.Qualifier` (line 1518) | N/A |

Additionally, `Trigger.Qualifier` (line 4949) in the `GetFunction` response's `Triggers` list carries per-trigger qualifier info.

## Goals / Non-Goals

**Goals:**
- Add an optional `qualifier` parameter (TypeString, Optional + Computed) to the `tencentcloud_scf_function` resource schema
- Pass the qualifier value when calling `GetFunction`, `DeleteFunction`, `CreateTrigger`, and `DeleteTrigger` APIs
- Read back the qualifier from `GetFunction` response into Terraform state
- Maintain full backward compatibility (existing configurations without qualifier continue to work)

**Non-Goals:**
- Do NOT attempt to pass qualifier to `CreateFunction` (SDK does not support it)
- Do NOT attempt to pass qualifier to `UpdateFunctionCode` or `UpdateFunctionConfiguration` (SDK does not support it)
- Do NOT expose qualifier as anything other than a simple string (no enum validation, no custom types)
- Do NOT change the resource ID format (remains `namespace+name`)

## Decisions

### Decision 1: Schema Type — Optional + Computed

**Choice**: `Optional: true, Computed: true`

**Rationale**: The qualifier defaults to `$LATEST` on the cloud side when not specified. Making it `Computed` allows the provider to read back the actual value from the API response without requiring the user to always specify it. Making it `Optional` means users can specify it when they need version-aware operations.

### Decision 2: Service Layer Changes — Add qualifier parameter to relevant methods

**Choice**: Add a `qualifier` string parameter to `DescribeFunction`, `DeleteFunction`, `CreateTriggers`, and `DeleteTriggers` methods.

**Rationale**: These are the only service methods whose corresponding APIs support the qualifier parameter. Passing an empty string will result in the SDK not setting the field (nil), preserving default behavior.

### Decision 3: Resource ID Format — Keep unchanged

**Choice**: Keep `namespace+name` as the resource ID. Do NOT include qualifier in the ID.

**Rationale**: The qualifier is about which version of a function to operate on, not about which function instance. Including qualifier in the ID would:
- Break backward compatibility with existing state
- Create separate Terraform resources for different versions of the same function
- Not align with the resource being a "function" rather than a "function version"

### Decision 4: Schema Field Path

**Choice**: Use `"qualifier"` (lowercase) as the Terraform schema field name, matching the lowercase convention used by most existing fields in this resource.

**Rationale**: Consistency with existing naming patterns in the resource (e.g., `handler`, `runtime`, `timeout`).

## Risks / Trade-offs

- **Risk**: A user sets `qualifier` to a specific version and then deletes the resource — only that version will be deleted, not the entire function. This is the expected SCF behavior but may surprise users.
  - **Mitigation**: Document this behavior clearly in the parameter description.

- **Risk**: The `CreateTrigger` and `CreateFunction` operations may have implicit qualifier dependencies. If a user creates a function then triggers with a qualifier that doesn't exist yet, the trigger creation will fail.
  - **Mitigation**: This is an inherent SCF API behavior. The Terraform provider should surface the API error clearly.

- **Risk**: The `DeleteFunction` with qualifier only deletes the specified version; if the qualifier is `$LATEST` or empty, it deletes the entire function. Users may not realize they need to set qualifier to a specific version to do version-level deletion.
  - **Mitigation**: Document this clearly in the description.

## Migration Plan

- No migration needed. The new parameter is purely additive (Optional + Computed).
- Existing state files do not contain `qualifier` and will continue to work.
- On the next `terraform apply`, if the user doesn't specify `qualifier`, the API default (`$LATEST`) will apply.
- On the next `terraform refresh`, the qualifier value returned by the API will be stored in state.

## Open Questions

- None. All API capabilities have been verified against the vendor SDK.