## Why

The `tencentcloud_ccn_instances_accept_attach` resource allows users to accept cross-region CCN attachment instances, but it does not expose the `OrderType` parameter that controls which account pays for the instances. The Tencent Cloud VPC SDK's `CcnInstance` struct already supports the `OrderType` field (`PayByCcnOwner` or `PayByInstanceOwner`), but the Terraform provider does not pass it through. Users who need to specify the billing account for accepted CCN attachments are forced to use the console or API directly.

## What Changes

- Add an optional `order_type` sub-field to the `instances` block of the `tencentcloud_ccn_instances_accept_attach` resource schema. Since `instances` is `ForceNew`, `order_type` inherits the same immutability behavior.
- Pass `OrderType` to the `AcceptAttachCcnInstances` API request by setting `CcnInstance.OrderType` for each instance in the Create function.
- The `OrderType` field maps to the cloud API path `request.Instances.OrderType` (a field on the `CcnInstance` struct, not the top-level request).

## Capabilities

### New Capabilities
- `ccn-instances-accept-attach-order-type`: Add the `order_type` optional sub-field to the `instances` block of the `tencentcloud_ccn_instances_accept_attach` resource, allowing users to specify the billing account (`PayByCcnOwner` or `PayByInstanceOwner`) for accepted CCN attachment instances.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach.go` — add `order_type` schema field under `instances` block, wire through Create function
  - `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach_test.go` — add unit test coverage using gomonkey mock
  - `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach.md` — update documentation with `order_type` usage example
- **SDK dependency:** No SDK update needed — the `CcnInstance` struct in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312` already includes the `OrderType *string` field.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional; existing configurations continue to work unchanged.
- **API constraints:** `OrderType` is only part of the `AcceptAttachCcnInstances` Create request (via the `CcnInstance` struct). Since the resource has no Update operation and `instances` is `ForceNew`, `order_type` is immutable after creation.
