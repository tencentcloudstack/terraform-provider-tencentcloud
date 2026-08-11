## Why

The TencentCloud provider currently lacks a dedicated resource for managing the binding of a single tag key/value to a single cloud resource (resource six-segment). While `tencentcloud_tag` manages the tag (key/value pair) lifecycle and `tencentcloud_tag_attachment` attaches a tag to a set of resources in batch, there is no first-class resource to manage a single tag-on-resource binding through its full CRUD lifecycle (create/read/update/delete). Users who need to bind a specific tag key and value to a specific resource, and later update only the value (or remove the binding), must rely on imperative API calls outside Terraform. This change introduces a `tencentcloud_tag_resource_tag` resource to fill that gap.

## What Changes

- Add a new Terraform resource `tencentcloud_tag_resource_tag` (RESOURCE_KIND_GENERAL) under `tencentcloud/services/tag/resource_tc_tag_resource_tag.go`, with full Create/Read/Update/Delete operations:
  - **Create**: calls `AddResourceTag` to bind a tag key/value to a resource six-segment.
  - **Read**: calls `GetResources` to query the tags currently bound to the resource and flatten the matched tag.
  - **Update**: calls `UpdateResourceTagValue` to modify the tag value bound to the resource (tag_key and resource are immutable).
  - **Delete**: calls `DeleteResourceTag` to unbind the tag from the resource.
- Add a composite resource ID (`tag_key` + `resource`, joined by `tccommon.FILED_SP`) so each binding is uniquely addressable.
- Register the new resource in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`.
- Add a markdown doc `tencentcloud/services/tag/resource_tc_tag_resource_tag.md` (will be mirrored to `website/docs/` via `make doc`).
- Add unit tests using gomonkey mocks in `resource_tc_tag_resource_tag_test.go`.

## Capabilities

### New Capabilities
- `tag-resource-tag`: Manage the lifecycle of a single tag key/value binding to a single cloud resource (resource six-segment), including create, read, update tag value, and delete (unbind).

### Modified Capabilities
<!-- None. No existing spec requirements are changing. -->

## Impact

- **New files**:
  - `tencentcloud/services/tag/resource_tc_tag_resource_tag.go` (resource CRUD implementation)
  - `tencentcloud/services/tag/resource_tc_tag_resource_tag.md` (resource doc)
  - `tencentcloud/services/tag/resource_tc_tag_resource_tag_test.go` (unit tests, gomonkey mocks)
- **Modified files**:
  - `tencentcloud/provider.go` (register `tencentcloud_tag_resource_tag`)
  - `tencentcloud/provider.md` (register resource entry, generated via `make doc`)
- **Cloud APIs** (all in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813`):
  - `AddResourceTag` (create/binding)
  - `GetResources` (read/flatten)
  - `UpdateResourceTagValue` (update tag value)
  - `DeleteResourceTag` (delete/unbinding)
- **Dependencies**: No new third-party dependencies; the tag SDK package is already vendored.
- **Backward compatibility**: Fully additive; no existing resource schema or behavior is changed.
