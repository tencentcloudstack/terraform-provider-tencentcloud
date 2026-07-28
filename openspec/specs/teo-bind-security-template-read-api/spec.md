## Requirements

### Requirement: Read locates binding via DescribeSecurityTemplateBindings
The `DescribeTeoBindSecurityTemplateById` service function SHALL locate a template-to-domain binding by calling the `DescribeSecurityTemplateBindings` API with the specified `zoneId` and `templateId` (as a single-element array). It SHALL NOT use the `DescribeZones` or `DescribeWebSecurityTemplates` APIs. The function SHALL return the matching `EntityStatus` from the `TemplateScope`'s `EntityStatus` list where `Entity` matches the requested `entity`.

#### Scenario: Binding found in the API response
- **WHEN** the `DescribeSecurityTemplateBindings` API returns a `SecurityTemplate` list containing a `TemplateScope` with `ZoneId` matching the requested `zoneId`, and the `EntityStatus` list contains an entry whose `Entity` matches the requested `entity`
- **THEN** the function SHALL return the matched `EntityStatus` (including `Entity`, `Status`, and `Message` fields)

#### Scenario: TemplateScope is empty
- **WHEN** the `DescribeSecurityTemplateBindings` API returns an empty `SecurityTemplate` list or an empty `TemplateScope` array
- **THEN** the function SHALL return nil without error

#### Scenario: Entity not found in EntityStatus list
- **WHEN** the `DescribeSecurityTemplateBindings` API returns a matching `TemplateScope` but no `EntityStatus` entry matches the requested `entity`
- **THEN** the function SHALL return nil without error

#### Scenario: API call retry on failure
- **WHEN** the `DescribeSecurityTemplateBindings` API call fails with a retryable error
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
