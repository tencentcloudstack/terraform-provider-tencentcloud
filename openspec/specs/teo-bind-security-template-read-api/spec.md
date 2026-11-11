## Requirements

### Requirement: Read locates binding via DescribeZones and DescribeWebSecurityTemplates
The `DescribeTeoBindSecurityTemplateById` service function SHALL locate a template-to-domain binding by first paging through the `DescribeZones` API (with `Limit` set to the documented maximum of 100) to collect all zone IDs under the account, then calling `DescribeWebSecurityTemplates` in batches of at most 100 zone IDs per request. It SHALL NOT use the `DescribeSecurityTemplateBindings` API.

#### Scenario: Binding found in the first batch
- **WHEN** the account has fewer than 100 zones and one of the returned `SecurityPolicyTemplates` matches the requested `templateId` and contains a `BindDomains` entry whose `Domain` equals the requested `entity`
- **THEN** the function SHALL return an `EntityStatus` with `Entity` set to the requested entity and `Status` set to the matched `BindDomainInfo.Status`

#### Scenario: Binding found in a later batch
- **WHEN** the account has more than 100 zones and the matching binding is only returned for a zone ID in the second (or later) `DescribeWebSecurityTemplates` batch
- **THEN** the function SHALL issue at least two `DescribeWebSecurityTemplates` requests and return the matching `EntityStatus`

#### Scenario: API call retry on failure
- **WHEN** a `DescribeZones` or `DescribeWebSecurityTemplates` API call fails with a retryable error
- **THEN** the function SHALL retry the call with `tccommon.ReadRetryTimeout` duration using `resource.Retry`, and wrap failures with `tccommon.RetryError`

### Requirement: Read returns nil when no binding matches
When no `SecurityPolicyTemplateInfo` matches the requested `templateId`, or the matched template has no `BindDomains` entry whose `Domain` equals the requested `entity`, the function SHALL return a nil `EntityStatus` (and nil error) so the resource Read method clears the Terraform state id.

#### Scenario: No zone available
- **WHEN** `DescribeZones` returns no zones
- **THEN** the function SHALL log `[DEBUG] ... no zone found when reading teo bind_security_template` and return nil without calling `DescribeWebSecurityTemplates`

#### Scenario: Template or entity not found
- **WHEN** the returned `SecurityPolicyTemplates` do not contain the requested `templateId` or none of the `BindDomains` matches the `entity`
- **THEN** the function SHALL return nil without error

### Requirement: Resource Read preserves id in logs before clearing state
The `resourceTencentCloudTeoBindSecurityTemplateRead` function SHALL, when the service lookup returns nil (binding not found), print a `[CRUD] teo_bind_security_template id=%s` log line containing `d.Id()` BEFORE calling `d.SetId("")`, so the cleared id remains traceable in logs.

#### Scenario: Binding not found during read
- **WHEN** the service function returns a nil `EntityStatus`
- **THEN** the read function SHALL log the current id, then clear the id, then log a warning that the resource was not found

### Requirement: Create state refresh tolerates nil status
The `resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0` state refresh function SHALL return the response object with an empty state string when `resp.Status` is nil, instead of dereferencing a nil pointer.

#### Scenario: Status nil during create polling
- **WHEN** the service function returns a non-nil `EntityStatus` whose `Status` field is nil
- **THEN** the state refresh function SHALL return `(resp, "", nil)` without error
