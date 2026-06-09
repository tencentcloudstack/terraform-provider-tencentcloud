## Context

The `tencentcloud_teo_content_identifier` Terraform resource manages TEO (EdgeOne) content identifiers. The cloud API's `ContentIdentifier` struct returns a `Status` field indicating the resource's lifecycle state (active/deleted), but the current Terraform resource schema does not expose this field to users.

The existing resource already has computed fields `content_id`, `created_on`, and `modified_on` from the same `ContentIdentifier` struct. Adding `status` follows the same pattern — it is a read-only field populated from the `DescribeContentIdentifiers` API response.

Current schema fields in `resource_tc_teo_content_identifier.go`:
- `description` (Required, TypeString)
- `plan_id` (Required, TypeString)
- `tags` (Optional, TypeList)
- `content_id` (Computed, TypeString)
- `created_on` (Computed, TypeString)
- `modified_on` (Computed, TypeString)

Cloud API `ContentIdentifier` struct fields not yet in Terraform:
- `Status` (string) — content identifier status: `active` or `deleted`
- `ReferenceCount` (int64) — number of rule engine references
- `DeletedOn` (string) — deletion time

This change adds only `status` as the new computed parameter.

## Goals / Non-Goals

**Goals:**
- Add the `status` computed field to the `tencentcloud_teo_content_identifier` resource schema
- Populate the `status` field in the Read method from the `DescribeContentIdentifiers` API response
- Maintain full backward compatibility — computed field addition does not affect existing configurations
- Add unit test coverage for the new field

**Non-Goals:**
- Adding `reference_count` or `deleted_on` fields (out of scope for this change)
- Modifying Create, Update, or Delete methods (status is read-only)
- Changing any existing schema fields or behavior

## Decisions

1. **Add `status` as a Computed field (TypeString)**
   - Rationale: `Status` is returned by the `DescribeContentIdentifiers` API in the `ContentIdentifier` struct. It is not a user-configurable parameter. Following the existing pattern for `content_id`, `created_on`, and `modified_on`, it should be a computed field.
   - Alternative considered: Making it Optional + Computed — rejected because `status` is not an input parameter for Create or Modify API calls.

2. **Read `status` from `respData.Status` in the Read method**
   - Rationale: The existing `DescribeTeoContentIdentifierById` service method already returns the `ContentIdentifier` struct which includes `Status`. The Read method simply needs to add `d.Set("status", respData.Status)` with the nil-check pattern consistent with other computed fields.

3. **No changes to Create/Update/Delete methods**
   - Rationale: `Status` is a server-side computed field. The `CreateContentIdentifier` and `ModifyContentIdentifier` API requests do not include a `Status` parameter. Create and Update methods already call Read at the end, so `status` will be automatically populated.

4. **Unit test using gomonkey mock**
   - Rationale: Following the project's testing pattern for new resources, use gomonkey to mock the cloud API client and verify that the Read method correctly sets the `status` field.

## Risks / Trade-offs

- **[Nil pointer risk]** → Mitigation: Follow the existing nil-check pattern (`if respData.Status != nil`) before calling `d.Set("status", ...)`.
- **[State drift if status becomes "deleted"]** → Mitigation: The Read method already handles the case where the resource is not found by setting `d.SetId("")`. If the resource status is "deleted", the Terraform state will reflect this, and users can take appropriate action.
