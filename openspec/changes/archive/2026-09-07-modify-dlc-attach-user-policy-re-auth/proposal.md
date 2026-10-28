## Why

The `tencentcloud_dlc_attach_user_policy_attachment` resource currently exposes `policy_set.re_auth` as a `Computed`-only field, meaning users cannot configure whether a grantee is allowed to re-grant (further delegate) the permissions they receive. The DLC `AttachUserPolicy` cloud API accepts `ReAuth` (`*bool`) as an input parameter on each `Policy` in the `PolicySet`, so the API supports letting callers opt into delegation. Exposing this as a user-settable optional field lets users fully control the authorization semantics of the policies they attach, aligning the Terraform provider with the cloud API's capabilities.

## What Changes

- Modify the `policy_set.re_auth` schema field from `Computed: true` to `Optional: true, Computed: true` so users can set it while it is still populated from the API response on read.
- In the `resourceTencentCloudDlcAttachUserPolicyAttachmentCreate` handler, add mapping of the `re_auth` value from `policy_set` state into the `dlc.Policy.ReAuth` field before calling the `AttachUserPolicy` API.
- Update the resource documentation (`resource_tc_dlc_attach_user_policy_attachment.md`) to reflect that `re_auth` is now an optional input.
- Add/update unit test coverage in `resource_tc_dlc_attach_user_policy_attachment_test.go` for the create path that sets `re_auth`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `dlc-attach-user-policy-attachment`: The `policy_set.re_auth` field changes from read-only (`Computed`) to user-settable (`Optional` + `Computed`), and the Create handler now forwards the user-configured `ReAuth` value to the `AttachUserPolicy` API request `PolicySet[*].ReAuth`.

## Impact

- **Code**: `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment.go` (schema + create handler), `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment_test.go` (unit tests).
- **Documentation**: `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment.md`.
- **API**: No new API endpoints; the existing `AttachUserPolicy` (dlc v20210125) already accepts `ReAuth` on each `Policy` in the request `PolicySet`. Verified against the vendored SDK `Policy` struct (`ReAuth *bool`).
- **Backward compatibility**: Fully backward compatible. The field becomes optional (not required), so existing configurations without `re_auth` continue to work and the API default (`false`) applies. State continues to be populated from the API response on read.
- **Dependencies**: No new dependencies; uses the already-vendored `tencentcloud-sdk-go/tencentcloud/dlc/v20210125`.
