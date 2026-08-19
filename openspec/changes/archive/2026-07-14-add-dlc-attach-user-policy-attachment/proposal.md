## Why

The DLC (Data Lake Compute) product currently lacks a Terraform resource to manage the binding of authorization policies to users. Users need to bind (attach) and unbind (detach) permission policies to/from DLC sub-users through Terraform, enabling infrastructure-as-code management of DLC data access permissions. Today this binding relationship can only be managed via the console or API calls, making it difficult to keep permission grants consistent with other infrastructure definitions.

## What Changes

- Add a new Terraform resource `tencentcloud_dlc_attach_user_policyr_attachment` of kind `RESOURCE_KIND_ATTACHMENT` to manage the binding between a DLC user and authorization policies.
- The resource supports Create (bind via `AttachUserPolicy`), Read (query via `DescribeUserInfo`), and Delete (unbind via `DetachUserPolicy`). There is no dedicated Update API; the resource is immutable and uses `ForceNew` for changing any parameter.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add the resource documentation file `resource_tc_dlc_attach_user_policyr_attachment.md`.
- Add unit tests using gomonkey mocks (no Terraform test suite) in `resource_tc_dlc_attach_user_policyr_attachment_test.go`.

## Capabilities

### New Capabilities
- `dlc-attach-user-policy-attachment`: Manage the binding relationship between a DLC user and authorization policies, including binding (create), reading the bound policies, and unbinding (delete).

### Modified Capabilities
<!-- No existing capabilities are modified. -->

## Impact

- **New files**:
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment.go` (resource CRUD implementation)
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment_test.go` (unit tests with gomonkey mocks)
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment.md` (documentation)
- **Modified files**:
  - `tencentcloud/provider.go` (register the new resource)
  - `tencentcloud/provider.md` (documentation index entry)
- **Cloud APIs used** (from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`):
  - `AttachUserPolicy` — bind policies to a user (Create)
  - `DescribeUserInfo` — read user detail info including bound policies (Read)
  - `DetachUserPolicy` — unbind policies from a user (Delete)
- **Resource ID**: composite id composed of `user_id` + `account_type` (joined by `tccommon.FILED_SP`), used to uniquely identify the binding relationship. Since the attach API does not return a single identifier, the composite key is derived from request parameters.
- **No breaking changes**: this is a purely additive change.
