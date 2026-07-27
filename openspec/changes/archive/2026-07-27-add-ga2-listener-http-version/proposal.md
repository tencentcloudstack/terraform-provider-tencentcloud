## Why

The `tencentcloud_ga2_listener` resource currently exposes `http_version` as a Computed-only field, meaning users cannot specify the HTTP version when creating an HTTPS listener. The cloud API `CreateListener` now supports `HttpVersion` as an input parameter (allowing users to choose `HTTP/1.1` or `HTTP/2`), but the Terraform resource does not expose this capability. Users who need to control the HTTP version for their HTTPS listeners are forced to fall back to the console or API directly.

## What Changes

- Promote `http_version` from a Computed-only field to an Optional+Computed field in the `tencentcloud_ga2_listener` schema.
  - The field is only meaningful for HTTPS listeners (the SDK comment states: "HTTPS监听器支持选择版本").
- Update the `Create` function to forward `http_version` to `CreateListenerRequest.HttpVersion` when set.
- The `Read` function already reads `HttpVersion` from `ListenerSet`; no change needed there.
- The `Update` function adds an immutable args check: `http_version` is listed in an `immutableArgs` slice, and if any immutable field changes during update, the function returns an error instructing the user to recreate the resource.
- Update the resource markdown documentation to reflect the new input parameter.

## Capabilities

### New Capabilities
<!-- None: this change modifies an existing resource, not introducing a new capability. -->

### Modified Capabilities
- `ga2-listener-resource`: `http_version` is promoted from Computed-only to Optional+Computed, allowing users to specify the HTTP version at creation time for HTTPS listeners. The `Update` function checks for changes to `http_version` and returns an error if the user attempts to modify it.

## Impact

- **Modified code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.go`: Change `http_version` schema from Computed to Optional+Computed; add `http_version` forwarding in the Create function; add immutable args check in the Update function so that changes to `http_version` return an error.
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.md`: Update example HCL to show the `http_version` parameter.
- **No breaking change**: The field was previously Computed-only; promoting it to Optional+Computed is backward-compatible. Existing Terraform configurations that do not set `http_version` will continue to work unchanged (the API default will apply, and the value will be read back as before).
- **APIs consumed**: No new API calls required; the existing `CreateListener` call gains one additional optional field.
- **SDK upgrade**: Not required — `CreateListenerRequest.HttpVersion` already exists in the vendored SDK.
