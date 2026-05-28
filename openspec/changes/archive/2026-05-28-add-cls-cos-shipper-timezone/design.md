## Context

The `tencentcloud_cls_cos_shipper` resource manages CLS (Cloud Log Service) log shipping rules to COS (Cloud Object Storage). The resource currently supports parameters like `topic_id`, `bucket`, `prefix`, `shipper_name`, `interval`, `max_size`, `filter_rules`, `partition`, `compress`, `content`, `filename_mode`, `start_time`, `end_time`, and `storage_type`.

The cloud API already supports a `TimeZone` parameter in `CreateShipperRequest`, `ModifyShipperRequest`, and returns it in `ShipperInfo` (via `DescribeShippers`). However, the Terraform resource does not yet expose this parameter.

The resource file is at `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go` and follows the standard CRUD pattern with retry logic.

## Goals / Non-Goals

**Goals:**
- Add `time_zone` as an Optional string parameter to the `tencentcloud_cls_cos_shipper` resource schema
- Wire the parameter into Create, Read, and Update methods
- Maintain backward compatibility (existing configurations without `time_zone` continue to work)

**Non-Goals:**
- Adding other missing parameters (e.g., `DSLFilter`, `RoleArn`, `ExternalId`)
- Changing existing parameter behavior
- Modifying the service layer (`service_tencentcloud_cls.go`)

## Decisions

1. **Schema field type**: Use `schema.TypeString` with `Optional: true` and `Computed: true`. The field is optional on create/update, and the API may return a default value on read.

2. **No validation function**: The API accepts a wide range of GMT/UTC timezone formats. Rather than maintaining a large validation list in Terraform, we rely on the API to validate the input and return appropriate errors.

3. **Update support**: The `time_zone` field is updatable via `ModifyShipper` API, so it will be handled in the Update method with `d.HasChange("time_zone")` pattern, consistent with other string fields like `storage_type`.

4. **Read support**: Read from `ShipperInfo.TimeZone` in the Read method with nil check before setting, consistent with existing patterns like `storage_type`.

## Risks / Trade-offs

- [Risk] API timezone format validation may change → Mitigation: Rely on API-side validation; no client-side validation needed.
- [Trade-off] No client-side validation means users get API errors instead of Terraform plan-time errors → Acceptable since the timezone list is extensive and may be updated server-side.
