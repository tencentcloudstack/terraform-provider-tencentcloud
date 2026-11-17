## Why

The `tencentcloud_teo_bind_security_template` resource currently exposes only `status` from the template binding query result. The cloud API `EntityStatus` struct returned by the read path includes an additional `Message` field that provides human-readable delivery status information (e.g., failure reasons). Exposing this `Message` field as a computed attribute gives users visibility into configuration delivery details without having to check the cloud console.

## What Changes

- Add a new `message` computed attribute to the `tencentcloud_teo_bind_security_template` resource schema, sourced from the `Message` field of the `EntityStatus` returned by the read path.

## Capabilities

### New Capabilities
- `teo-bind-security-template-message`: Expose the `Message` field of `EntityStatus` as a computed attribute on `tencentcloud_teo_bind_security_template`.

## Impact

- **Affected files**:
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template.go` — add `message` schema field and set it in Read
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template.md` — add `message` to example usage
  - `tencentcloud/services/teo/resource_tc_teo_bind_security_template_test.go` — add test coverage for `message`
- **Backward compatibility**: Fully backward compatible — the new `message` field is computed-only, and the existing `status` field behavior is preserved
