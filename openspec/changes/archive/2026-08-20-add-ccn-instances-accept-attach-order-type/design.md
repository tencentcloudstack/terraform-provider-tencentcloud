## Context

The `tencentcloud_ccn_instances_accept_attach` resource is located at `tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach.go`. It wraps the VPC `AcceptAttachCcnInstances` API to accept cross-region CCN attachment instances.

**Current state:**
- The resource schema has a top-level `ccn_id` field and an `instances` block (TypeList, ForceNew).
- The `instances` block contains sub-fields: `instance_id`, `instance_region`, `instance_type`, `description`, and `route_table_id`.
- The Create function builds `vpc.CcnInstance` structs from the schema and appends them to `request.Instances`.
- The resource has no Update operation; `instances` is `ForceNew`, so any change triggers recreation.
- The Read function is a no-op (the API has no direct query endpoint for this resource).
- The Delete function is a no-op.

**SDK analysis:**
The `CcnInstance` struct in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312/models.go` already includes an `OrderType *string` field with the following documentation:
> 实例付费方式。枚举值：PayByCcnOwner（CCN所在账号付费）、PayByInstanceOwner（关联实例所在账号付费）

No SDK update is required — the field is already present in the vendored SDK.

## Goals / Non-Goals

**Goals:**
- Add an optional `order_type` sub-field (TypeString) to the `instances` block of the `tencentcloud_ccn_instances_accept_attach` resource schema.
- Pass `OrderType` to the `AcceptAttachCcnInstances` API by setting `CcnInstance.OrderType` in the Create function when the user specifies it.
- Maintain full backward compatibility — existing configurations continue to work unchanged.
- Add unit test coverage using gomonkey mock (no terraform test suite, per project guidelines for new resources).

**Non-Goals:**
- Adding `order_type` to any other CCN resource or datasource.
- Adding an Update operation (the resource only supports Create and Delete).
- Adding Read support for `order_type` (the Read function is a no-op because the `AcceptAttachCcnInstances` API has no corresponding query endpoint).

## Decisions

### Decision 1: `order_type` as a sub-field of `instances` block

**Rationale:** The cloud API path is `request.Instances.OrderType`, meaning `OrderType` is a field on the `CcnInstance` struct (each element of the `Instances` array), not a top-level field on `AcceptAttachCcnInstancesRequest`. Therefore, in the Terraform schema, `order_type` must be a sub-field within the `instances` block, consistent with the existing `instance_id`, `instance_region`, `instance_type`, `description`, and `route_table_id` fields.

### Decision 2: `order_type` is Optional and inherits ForceNew from parent block

**Rationale:** The `instances` block is `ForceNew: true` at the block level. Individual sub-fields within a `ForceNew` TypeList block do not need their own `ForceNew` flag — any change to any sub-field triggers recreation of the entire resource. The `order_type` field is Optional (TypeString) so users can omit it when they want the API default behavior.

### Decision 3: Only set `OrderType` when the value is non-empty

**Rationale:** Consistent with the existing pattern for `route_table_id` in the Create function, `OrderType` should only be set on the `CcnInstance` struct when the user provides a non-empty value. This avoids overriding API defaults with an empty string.

### Decision 4: No validation of enum values

**Rationale:** The API accepts `PayByCcnOwner` and `PayByInstanceOwner`. Rather than adding provider-side `ValidateFunc` validation (which could become stale if the API adds new values), we rely on the API to reject invalid values. This is consistent with the existing `instance_type` field which also has no `ValidateFunc`.

## Risks / Trade-offs

- **[Risk] Read function does not refresh `order_type`**: The Read function is a no-op, so `order_type` will not be refreshed from the cloud after creation.
  - **Mitigation:** This is consistent with the existing behavior of all fields in this resource. Terraform preserves the configured value in state. No additional handling needed.
- **[Risk] No Update support**: Users cannot change `order_type` without recreating the resource.
  - **Mitigation:** This is an inherent limitation of the `AcceptAttachCcnInstances` API, which only has a Create endpoint. The `instances` block is already `ForceNew`, so this behavior is consistent.
