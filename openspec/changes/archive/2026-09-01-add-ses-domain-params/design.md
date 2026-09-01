## Context

The `tencentcloud_ses_domain` resource (in `tencentcloud/services/ses/resource_tc_ses_domain.go`) currently only supports the `email_identity` parameter. The resource implements Create, Read, and Delete operations (no Update) using the SES cloud API:

- **Create**: calls `CreateEmailIdentity` with only `EmailIdentity`
- **Read**: calls `GetEmailIdentity` (via `DescribeSesDomain` service method) and reads DNS `Attributes`
- **Delete**: calls `DeleteEmailIdentity`

The cloud API `CreateEmailIdentityRequest` already supports `DKIMOption` (*uint64) and `TagList` ([]*TagList) fields. The `GetEmailIdentityResponse` returns `DKIMOption` (*uint64) and `TagList` ([]*TagList) as well. The `TagList` type contains `TagKey` and `TagValue` string fields.

Since this resource has no Update operation, all new parameters will be `ForceNew` — changing them requires recreating the resource.

## Goals / Non-Goals

**Goals:**
- Expose `DKIMOption` as an optional `dkim_option` schema field (uint64) so users can set DKIM key length at creation.
- Expose `TagList` (TagKey/TagValue) as optional `tag_key` and `tag_value` schema fields (string) so users can associate a single tag with the domain at creation.
- Read back `dkim_option`, `tag_key`, and `tag_value` from `GetEmailIdentity` so they appear in Terraform state.

**Non-Goals:**
- Supporting multiple tags (the cloud API `TagList` is a slice, but the Terraform resource will expose a single tag pair as `tag_key`/`tag_value` for simplicity, matching the requirement mapping).
- Adding an Update operation to the resource (the resource remains Create-Read-Delete only; new parameters are ForceNew).
- Modifying the `attributes` computed field behavior.

## Decisions

### Decision 1: Expose TagList as flat `tag_key` / `tag_value` fields rather than a nested list

**Rationale**: The requirement specifies `TagList.TagKey → TagKey` and `TagList.TagValue → TagValue` as flat schema names. The SES domain resource manages a single domain identity, and the typical use case is a single tag. Flat fields are simpler for users and align with the requirement mapping. The `CreateEmailIdentity` request will construct a `[]*TagList` slice with a single element from these two fields.

**Alternative considered**: A nested `tag_list` block schema (TypeList of blocks with `tag_key`/`tag_value`). Rejected because the requirement explicitly maps to flat top-level `TagKey` and `TagValue` schema names.

### Decision 2: New parameters are ForceNew

**Rationale**: The `tencentcloud_ses_domain` resource has no Update operation — only Create, Read, and Delete. The cloud API does not provide an update endpoint for email identity parameters. Therefore `dkim_option`, `tag_key`, and `tag_value` must be `ForceNew: true`, consistent with `email_identity`.

### Decision 3: Modify `DescribeSesDomain` service method to return the full `GetEmailIdentityResponseParams`

**Rationale**: Currently `DescribeSesDomain` only returns `[]*ses.DNSAttributes` (the `Attributes` field). To read `DKIMOption` and `TagList` from the response, the service method needs to return the full `GetEmailIdentityResponseParams` (or at least the additional fields). The cleanest approach is to return `*ses.GetEmailIdentityResponseParams` so the Read function can access all fields including `Attributes`, `DKIMOption`, and `TagList`.

**Alternative considered**: Add a separate service method for the new fields. Rejected because it would result in a duplicate API call to `GetEmailIdentity`.

### Decision 4: `dkim_option` schema type uses `schema.TypeInt`

**Rationale**: The cloud API uses `*uint64` for `DKIMOption`. In Terraform schema, `schema.TypeInt` is the idiomatic mapping for integer-like values. The Create function will convert the int value to uint64 when building the request.

## Risks / Trade-offs

- **[Risk] `DescribeSesDomain` signature change may affect callers** → Mitigation: `DescribeSesDomain` is only called from `resourceTencentCloudSesDomainRead`. The signature change is contained within the ses service package. The return type changes from `[]*ses.DNSAttributes` to `*ses.GetEmailIdentityResponseParams`, and the Read function will be updated accordingly.
- **[Risk] TagList may contain multiple entries in API response** → Mitigation: The Read function will read the first element of the `TagList` slice to populate `tag_key` and `tag_value`, consistent with the single-tag schema design.
- **[Trade-off] ForceNew on all new parameters** → Users must recreate the domain to change DKIM option or tags. This is acceptable because the SES API does not support updating these fields on an existing identity.
