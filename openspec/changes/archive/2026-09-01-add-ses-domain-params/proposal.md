## Why

The `tencentcloud_ses_domain` resource currently only supports setting the `email_identity` field when creating an SES domain. Users cannot configure the DKIM key length or associate tags with the domain at creation time. The SES cloud API (`CreateEmailIdentity`) already supports `DKIMOption` and `TagList` parameters, and `GetEmailIdentity` returns these values, but the Terraform resource does not expose them. This limits users' ability to fully manage SES domain configuration through Terraform.

## What Changes

- Add a new optional `dkim_option` parameter (uint64) to the `tencentcloud_ses_domain` resource schema, allowing users to specify the DKIM key length (0: 1024-bit, 1: 2048-bit) at creation time.
- Add a new optional `tag_list` nested block parameter to the `tencentcloud_ses_domain` resource schema, allowing users to associate tags with the domain at creation time. The `tag_list` block contains `tag_key` and `tag_value` string sub-fields and supports multiple tag entries.
- Populate `dkim_option` and `tag_list` from the `GetEmailIdentity` response in the resource Read operation so they are reflected in Terraform state.

## Capabilities

### New Capabilities

- `ses-domain-params`: Adds `dkim_option` and `tag_list` (with `tag_key` and `tag_value` sub-fields) parameters to the `tencentcloud_ses_domain` resource, enabling DKIM key length configuration and tag association for SES domains.

### Modified Capabilities

<!-- No existing spec-level requirements are being modified. -->

## Impact

- **Affected code**: `tencentcloud/services/ses/resource_tc_ses_domain.go` (schema, create, read operations), `tencentcloud/services/ses/service_tencentcloud_ses.go` (DescribeSesDomain to return additional fields), `tencentcloud/services/ses/resource_tc_ses_domain_test.go` (test updates).
- **Affected docs**: `tencentcloud/services/ses/resource_tc_ses_domain.md` (example usage update).
- **APIs**: `CreateEmailIdentity` (new request parameters `DKIMOption`, `TagList`), `GetEmailIdentity` (new response fields `DKIMOption`, `TagList`).
- **Backward compatibility**: All new parameters are optional; existing configurations remain valid. No breaking changes.