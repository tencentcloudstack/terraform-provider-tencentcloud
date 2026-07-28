## MODIFIED Requirements

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