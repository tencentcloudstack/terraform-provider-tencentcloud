## Context

The `tencentcloud_sms_template` resource is a RESOURCE_KIND_GENERAL resource in the `sms` service package. It manages the full lifecycle of an SMS message template through the following API calls:

- **Create**: `AddSmsTemplate` - Creates a new SMS template
- **Read**: `DescribeSmsTemplateList` - Queries template details (already used in the service layer)
- **Update**: `ModifySmsTemplate` - Updates an existing template
- **Delete**: `DeleteSmsTemplate` - Deletes a template

The existing Read function delegates to `SmsService.DescribeSmsTemplate()`, which calls `DescribeSmsTemplateList` and returns a `*sms.DescribeTemplateListStatus` struct. This struct already contains the `StatusCode *int64` field (template audit status), but the Terraform resource does not expose it in its schema or read it into state.

Current state:
- Resource schema has 5 fields: `template_name`, `template_content`, `international`, `sms_type`, `remark`
- All fields are `Required: true`
- Read function sets 3 fields: `template_name`, `template_content`, `international`
- No computed/read-only fields exist

## Goals / Non-Goals

**Goals:**
- Add a new `status_code` computed field to the `tencentcloud_sms_template` resource schema
- Update the Read function to populate `status_code` from the `DescribeTemplateListStatus.StatusCode` field
- Ensure backward compatibility - existing Terraform configurations continue to work unchanged
- Keep the implementation minimal - no new API calls, no changes to Create/Update/Delete logic

**Non-Goals:**
- Not modifying the service layer (`service_tencentcloud_sms.go`) - the existing `DescribeSmsTemplate` method already returns the `DescribeTemplateListStatus` struct with `StatusCode`
- No changes to the Create, Update, or Delete operations
- No integration with other resources or data sources
- No state migration or import behavior changes

## Decisions

### Decision 1: Schema field type - `TypeInt`

**Choice**: Use `schema.TypeInt` for the `status_code` field.

**Rationale**:
- The SDK field `StatusCode` is `*int64` (a pointer to int64)
- Terraform's `TypeInt` maps naturally to Go's int types
- The status code values are integer constants (0, 1, 2, -1 as documented in the API)
- Alternative: `TypeString` would require string conversion and lose the semantic meaning of numeric status codes
- The `international` field in the same resource already uses `TypeInt`, providing consistency

### Decision 2: Schema field properties - `Computed: true`

**Choice**: Set `Computed: true` (not `Optional`, not `Required`, `ForceNew: false`).

**Rationale**:
- `StatusCode` is an output-only field returned by the Describe API - users cannot set it during Create
- Using `Computed: true` (without `Optional`) ensures Terraform knows the value comes from the provider, not user config
- No `ForceNew` needed since this field is never provided as input
- This is the standard pattern for read-only status fields in the provider

### Decision 3: Read function placement

**Choice**: Add the `status_code` setting after the existing `international` field in the Read function.

**Rationale**:
- Follows the existing pattern in the resource
- `d.Set()` calls are already grouped by field in the Read function
- Adds the nil check pattern consistent with other fields: `if template.StatusCode != nil { _ = d.Set("status_code", template.StatusCode) }`

### Decision 4: No changes to test file (unit tests via gomonkey)

**Choice**: Add unit test coverage using gomonkey (mock) approach when applicable, or skip test changes if the existing test framework doesn't easily support computed field testing.

**Rationale**:
- The existing test uses Terraform acceptance tests (requires real cloud credentials)
- Per the project's go code generation requirements, for schema-only additions to existing resources, the existing acceptance test pattern should be followed
- A computed field cannot be set by the user, so it would need to be verified via `TestCheckResourceAttrSet` in the existing acceptance test

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| **Nil pointer dereference**: `template.StatusCode` could be nil for templates created before this field was added to the API response | Add nil check before `d.Set()` - standard pattern used by all other fields in the resource |
| **Schema drift**: If the `StatusCode` value changes on the cloud side (e.g., template transitions from pending to approved) between Read calls, the Terraform state will eventually converge | This is expected behavior for computed fields - the next Read operation will sync the latest value. No action needed |
| **Backward compatibility**: Existing users running `terraform plan` after upgrading may see a diff showing `status_code` being added to state | This is a cosmetic diff only for the first plan after upgrade. Adding `Computed: true` means Terraform knows the value will be populated by the provider, so no unexpected diffs should appear |
| **ForceNew concern**: Since this is `Computed: true` and not `Optional`, users might worry about ForceNew triggers | No `ForceNew` is set, and `Computed: true` fields alone never trigger resource replacement |