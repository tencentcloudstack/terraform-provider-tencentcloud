## 1. Schema Changes

- [x] 1.1 Add `order_type` sub-field (TypeString, Optional, Description: "Instance billing method. Valid values: `PayByCcnOwner` (CCN owner pays), `PayByInstanceOwner` (instance owner pays).") to the `instances` block schema in `ResourceTencentCloudCcnInstancesAcceptAttach()` in `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach.go`

## 2. Create Function Changes

- [x] 2.1 In `resourceTencentCloudCcnInstancesAcceptAttachCreate`, read `order_type` from the `dMap` for each instance in the `instances` list
- [x] 2.2 Set `CcnInstance.OrderType` via `helper.String(v.(string))` only when `order_type` is non-empty (consistent with the existing `route_table_id` pattern)

## 3. Unit Tests

- [x] 3.1 Create `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach_test.go` with unit tests using gomonkey mock for the `AcceptAttachCcnInstances` API
- [x] 3.2 Add test case: Create CCN instances accept attach with `order_type` specified, verify `CcnInstance.OrderType` is set in the request
- [x] 3.3 Add test case: Create CCN instances accept attach without `order_type`, verify `CcnInstance.OrderType` is NOT set in the request

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach.md` to add `order_type` field in the example usage within the `instances` block
