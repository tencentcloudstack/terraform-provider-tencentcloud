## Context

The `tencentcloud_ga2_listener` resource already exists and is fully functional, managing the lifecycle of a Global Accelerator V2 listener. Its current schema exposes `http_version` as a **Computed-only** field — it can only be read from the API response, never set by the user.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` now supports `HttpVersion` as an input parameter on `CreateListenerRequest`:
- `CreateListenerRequestParams.HttpVersion *string` — comment: "HTTPS监听器支持选择版本. 枚举值: HTTP/1.1, HTTP/2"
- `ModifyListenerRequestParams` does **not** carry `HttpVersion` — it is immutable after creation.
- `ListenerSet.HttpVersion *string` — already read by the existing Read function.

This means users can now choose the HTTP version at listener creation time, but the Terraform resource does not expose this ability.

## Goals / Non-Goals

**Goals:**
- Promote `http_version` from Computed-only to Optional+Computed+ForceNew, allowing users to specify `HTTP/1.1` or `HTTP/2` when creating an HTTPS listener.
- Forward the user-supplied value to `CreateListenerRequest.HttpVersion` in the Create function.
- Maintain full backward compatibility: existing configurations that do not set `http_version` continue to work unchanged (the API default applies, and the value is read back as before).

**Non-Goals:**
- Adding support for modifying `http_version` after creation — the API does not support it, and ForceNew handles this correctly.
- Changing any other schema field or behavior of the `tencentcloud_ga2_listener` resource.
- Adding a new data source or separate resource.

## Decisions

### D1. Schema change: Computed → Optional+Computed+ForceNew

**Why ForceNew:** `ModifyListenerRequest` has no `HttpVersion` field. The API does not allow changing the HTTP version after creation, so any change to this field requires destroying and recreating the listener. This matches the existing ForceNew pattern used for `port_ranges`, `listener_type`, `protocol`, and `global_accelerator_id`.

**Why Optional+Computed (not just Optional):** The field is only meaningful for HTTPS listeners. For other protocols (TCP, UDP, HTTP), the API ignores this value. Using Optional+Computed ensures that:
- Users who don't set it get the API's default (computed from the response).
- The field always has a value in state after Read, even if the user never explicitly set it.
- No state inconsistency arises when importing existing listeners that were created with a specific HTTP version.

### D2. Create function: forward `http_version` to `CreateListenerRequest.HttpVersion`

The Create function currently does not set `HttpVersion` on the request because the schema field was Computed-only. After the schema change, we add a standard `if v, ok := d.GetOk("http_version"); ok` block that forwards the value to `request.HttpVersion = helper.String(v.(string))`.

This follows the same pattern used for all other Optional+Computed fields in the existing Create function.

### D3. Read function: no structural change needed

The Read function already reads `HttpVersion` from `ListenerSet` and sets it via `_ = d.Set("http_version", respData.HttpVersion)`. No change required — the field will continue to be populated from the API response regardless of whether the user set it or not.

### D4. Update function: no change needed

Since `http_version` is ForceNew, Terraform automatically triggers a destroy/create cycle if the user changes it. The Update function never sees a `HasChange("http_version")` event for in-place updates.

## Risks / Trade-offs

- **[Risk]** Users may set `http_version` on a TCP/UDP listener, where the API ignores it. → **Mitigation**: The field is Optional, and the API performs its own validation. If the API rejects the value, the existing retry and error handling will surface the error. We do not add client-side validation to keep the schema simple (matching the project's established pattern for protocol-dependent fields).

- **[Risk]** Promoting a Computed field to Optional+Computed could cause existing state to show a diff if the stored computed value differs from what the user might want to set. → **Mitigation**: Since the field was Computed-only before, existing states already have the correct API-returned value stored. No diff will appear for existing configurations because the user hasn't set the field explicitly, and Computed fields are not compared for changes unless the user adds an explicit value.

- **[Trade-off]** ForceNew means changing `http_version` destroys and recreates the entire listener (including its endpoint groups). This is unavoidable because the API does not support in-place modification.
